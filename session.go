package hhttp

import (
	"compress/gzip"
	"golang.org/x/net/publicsuffix"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Session struct {
	Jar    *cookiejar.Jar
	Header http.Header
}

func NewSession() *Session {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	s := &Session{
		Jar:    jar,
		Header: make(http.Header),
	}
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Session) UpdateHeader(header http.Header) {
	s.Header = header
}

func (s *Session) Get(urlStr string) (resp *HResponse, err error) {
	return s.Request("GET", urlStr, nil)
}

func (s *Session) PostForm(urlStr string, data url.Values) (resp *HResponse, err error) {
	return s.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
func (s *Session) Post(urlStr string, contentType string, body io.Reader) (resp *HResponse, err error) {
	s.Header.Set("Content-Type", contentType)
	return s.Request("POST", urlStr, body)
}

func (s *Session) Request(method string, urlStr string, body io.Reader) (resp *HResponse, err error) {

	client := &http.Client{
		Jar: s.Jar,
	}
	req, err := http.NewRequest(method, urlStr, body)

	req.Header = s.Header
	res, err := client.Do(req)
	wrapped := &HResponse{*res}
	return wrapped, err
}

type HResponse struct {
	http.Response
}

func (resp *HResponse) Text() string {
	defer resp.Body.Close()
	var txt []byte
	var err error
	if resp.Header["Content-Encoding"] != nil && resp.Header["Content-Encoding"][0] == "gzip" {
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer gr.Close()
		txt, err = ioutil.ReadAll(gr)
	} else {
		txt, err = ioutil.ReadAll(resp.Body)
	}
	if err == nil {
		return string(txt)
	} else {
		return ""
	}
}
