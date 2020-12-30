package req

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"github.com/projectdiscovery/retryablehttp-go"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// Resp represents a request with it's response
type Resp struct {
	r      *Req
	req    *retryablehttp.Request
	resp   *http.Response
	client *retryablehttp.Client
	cost   time.Duration
	*multipartHelper
	reqBody          []byte
	respBody         []byte
	downloadProgress DownloadProgress
	err              error // delayed error
}

func (r *Resp) HandleDecompression(bodyOrig []byte) (bodyDec []byte, err error) {
	if r.req == nil || r.req.Request == nil {
		return bodyOrig, nil
	}

	encodingHeader := strings.TrimSpace(strings.ToLower(r.req.Request.Header.Get("Accept-Encoding")))
	if encodingHeader == "gzip" || encodingHeader == "gzip, deflate" {
		gzipreader, err := gzip.NewReader(bytes.NewReader(bodyOrig))
		if err != nil {
			return bodyDec, err
		}
		defer gzipreader.Close()

		bodyDec, err = ioutil.ReadAll(gzipreader)
		if err != nil {
			return bodyDec, err
		}

		return bodyDec, nil
	}

	return bodyOrig, nil
}

// Request returns *http.Request
func (r *Resp) Request() *retryablehttp.Request {
	return r.req
}

// Response returns *http.Response
func (r *Resp) Response() *http.Response {
	return r.resp
}

// Bytes returns response body as []byte
func (r *Resp) Bytes() []byte {
	data, _ := r.ToBytes()
	return data
}

// ToBytes returns response body as []byte,
// return error if error happend when reading
// the response body
func (r *Resp) ToBytes() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.respBody != nil {
		return r.respBody, nil
	}
	defer r.resp.Body.Close()
	respBody, err := ioutil.ReadAll(r.resp.Body)
	if err != nil {
		r.err = err
		_, copyErr := io.Copy(ioutil.Discard, r.resp.Body)
		if copyErr != nil {
			return nil, copyErr
		}
		return nil, errors.Wrap(err, "could not read http body")
	}

	respBody, err = r.HandleDecompression(respBody)
	if err != nil {
		return nil, errors.Wrap(err, "could not decompress http body")
	}

	r.respBody = respBody
	return r.respBody, nil
}

// String returns response body as string
func (r *Resp) String() string {
	data, _ := r.ToBytes()
	return string(data)
}

// ToString returns response body as string,
// return error if error happend when reading
// the response body
func (r *Resp) ToString() (string, error) {
	data, err := r.ToBytes()
	return string(data), err
}

// ToJSON convert json response body to struct or map
func (r *Resp) ToJSON(v interface{}) error {
	data, err := r.ToBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// ToXML convert xml response body to struct or map
func (r *Resp) ToXML(v interface{}) error {
	data, err := r.ToBytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// ToFile download the response body to file with optional download callback
func (r *Resp) ToFile(name string) error {
	//TODO set name to the suffix of url path if name == ""
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	if r.respBody != nil {
		_, err = file.Write(r.respBody)
		return err
	}

	if r.downloadProgress != nil && r.resp.ContentLength > 0 {
		return r.download(file)
	}

	defer r.resp.Body.Close()
	_, err = io.Copy(file, r.resp.Body)
	return err
}

func (r *Resp) download(file *os.File) error {
	p := make([]byte, 1024)
	b := r.resp.Body
	defer b.Close()
	total := r.resp.ContentLength
	var current int64
	var lastTime time.Time

	defer func() {
		r.downloadProgress(current, total)
	}()

	for {
		l, err := b.Read(p)
		if l > 0 {
			_, _err := file.Write(p[:l])
			if _err != nil {
				return _err
			}
			current += int64(l)
			if now := time.Now(); now.Sub(lastTime) > r.r.progressInterval {
				lastTime = now
				r.downloadProgress(current, total)
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

var regNewline = regexp.MustCompile(`\n|\r`)

func (r *Resp) autoFormat(s fmt.State) {
	req := r.req
	if r.r.flag&Lcost != 0 {
		fmt.Fprint(s, req.Method, " ", req.URL.String(), " ", r.cost)
	} else {
		fmt.Fprint(s, req.Method, " ", req.URL.String())
	}

	// test if it is should be outputed pretty
	var pretty bool
	var parts []string
	addPart := func(part string) {
		if part == "" {
			return
		}
		parts = append(parts, part)
		if !pretty && regNewline.MatchString(part) {
			pretty = true
		}
	}
	if r.r.flag&LreqBody != 0 { // request body
		addPart(string(r.reqBody))
	}
	if r.r.flag&LrespBody != 0 { // response body
		addPart(r.String())
	}

	for _, part := range parts {
		if pretty {
			fmt.Fprint(s, "\n")
		}
		fmt.Fprint(s, " ", part)
	}
}

func (r *Resp) miniFormat(s fmt.State) {
	req := r.req
	if r.r.flag&Lcost != 0 {
		fmt.Fprint(s, req.Method, " ", req.URL.String(), " ", r.cost)
	} else {
		fmt.Fprint(s, req.Method, " ", req.URL.String())
	}
	if r.r.flag&LreqBody != 0 && len(r.reqBody) > 0 { // request body
		str := regNewline.ReplaceAllString(string(r.reqBody), " ")
		fmt.Fprint(s, " ", str)
	}
	if r.r.flag&LrespBody != 0 && r.String() != "" { // response body
		str := regNewline.ReplaceAllString(r.String(), " ")
		fmt.Fprint(s, " ", str)
	}
}

// Format fort the response
func (r *Resp) Format(s fmt.State, verb rune) {
	if r == nil || r.req == nil {
		return
	}
	if s.Flag('+') { // include header and format pretty.
		fmt.Fprint(s, r.Dump())
	} else if s.Flag('-') { // keep all informations in one line.
		r.miniFormat(s)
	} else { // auto
		r.autoFormat(s)
	}
}
