package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"os"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var IsJapanese bool

var defaultContentType = "text/html; charset=utf-8"

// ErrMissingBoundary is returned when there is no boundary given for a multipart entity
var ErrMissingBoundary = errors.New("no boundary found for multipart entity")

// ErrMissingContentType is returned when there is no "Content-Type" header for a MIME entity
var ErrMissingContentType = errors.New("no Content-Type found for MIME entity")

func MthmlToHtml(mht string) ([]byte, error) {

	fd, err := os.Open(mht)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	tr := &trimReader{rd: fd}
	tp := textproto.NewReader(bufio.NewReader(tr))

	// Parse the main headers
	header, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	body := tp.R

	parts, err := parseMIMEParts(header, body)
	if err != nil {
		return nil, err
	}

	var html *part
	var saves = make(map[string]string)
	for _, part := range parts {
		contentType := part.header.Get("Content-Type")
		if contentType == "" {
			return nil, ErrMissingContentType
		}
		mimetype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return nil, err
		}
		if html == nil && mimetype == "text/html" {
			html = part
			continue
		}
	}

	if html == nil {
		return nil, errors.New("html not found")
	}

	var r io.Reader
	if IsJapanese {
		r = transform.NewReader(bytes.NewReader(html.body), japanese.ShiftJIS.NewDecoder())
	} else {
		r = bytes.NewReader(html.body)
	}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	doc.Find("img,link,script").Each(func(i int, e *goquery.Selection) {
		changeRef(e, saves)
	})

	data, err := doc.Html()
	if err != nil {
		return nil, err
	}

	return []byte(data), nil

}

func changeRef(e *goquery.Selection, saves map[string]string) {
	attr := "src"
	switch e.Get(0).Data {
	case "img":
		e.RemoveAttr("loading")
		e.RemoveAttr("srcset")
	case "link":
		attr = "href"
	}
	ref, _ := e.Attr(attr)
	local, exist := saves[ref]
	if exist {
		e.SetAttr(attr, local)
	}
}

// part is a copyable representation of a multipart.Part
type part struct {
	header textproto.MIMEHeader
	body   []byte
}

// trimReader is a custom io.Reader that will trim any leading
// whitespace, as this can cause email imports to fail.
type trimReader struct {
	rd      io.Reader
	trimmed bool
}

// Read trims off any unicode whitespace from the originating reader
func (tr *trimReader) Read(buf []byte) (int, error) {
	n, err := tr.rd.Read(buf)
	if err != nil {
		return n, err
	}
	if !tr.trimmed {
		t := bytes.TrimLeftFunc(buf[:n], unicode.IsSpace)
		tr.trimmed = true
		n = copy(buf, t)
	}
	return n, err
}

func parseMIMEParts(hs textproto.MIMEHeader, b io.Reader) ([]*part, error) {
	var ps []*part
	// If no content type is given, set it to the default
	if _, ok := hs["Content-Type"]; !ok {
		hs.Set("Content-Type", defaultContentType)
	}
	ct, params, err := mime.ParseMediaType(hs.Get("Content-Type"))
	if err != nil {
		return ps, err
	}
	// If it's a multipart email, recursively parse the parts
	if strings.HasPrefix(ct, "multipart/") {
		if _, ok := params["boundary"]; !ok {
			return ps, ErrMissingBoundary
		}
		mr := multipart.NewReader(b, params["boundary"])
		for {
			var buf bytes.Buffer
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return ps, err
			}
			if _, ok := p.Header["Content-Type"]; !ok {
				p.Header.Set("Content-Type", defaultContentType)
			}
			subct, _, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
			if err != nil {
				return ps, err
			}
			if strings.HasPrefix(subct, "multipart/") {
				sps, err := parseMIMEParts(p.Header, p)
				if err != nil {
					return ps, err
				}
				ps = append(ps, sps...)
			} else {
				var reader io.Reader
				reader = p
				const cte = "Content-Transfer-Encoding"
				if p.Header.Get(cte) == "base64" {
					reader = base64.NewDecoder(base64.StdEncoding, reader)
				}
				// Otherwise, just append the part to the list
				// Copy the part data into the buffer
				if _, err := io.Copy(&buf, reader); err != nil {
					return ps, err
				}
				ps = append(ps, &part{body: buf.Bytes(), header: p.Header})
			}
		}
	} else {
		// If it is not a multipart email, parse the body content as a single "part"
		switch hs.Get("Content-Transfer-Encoding") {
		case "quoted-printable":
			b = quotedprintable.NewReader(b)
		case "base64":
			b = base64.NewDecoder(base64.StdEncoding, b)
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, b); err != nil {
			return ps, err
		}
		ps = append(ps, &part{body: buf.Bytes(), header: hs})
	}
	return ps, nil
}
