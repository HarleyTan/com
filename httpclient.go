//An easier way to use http.Client
package com

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	//	"net/url"
	//	"strings"
	"compress/gzip"
	"io"
	"log"
)

type HttpClient struct {
	c       *http.Client
	cookies []*http.Cookie
	jar     *cookiejar.Jar

	Header http.Header
	ua     string

	//编码转换相关处理
	conv    bool //conv between utf-8 and charset
	charset string
	enc     mahonia.Encoder
	dec     mahonia.Decoder

	//链接转向相关处理
	redirect    bool   //是否转向了。每次Get之前置为false
	redirectUrl string //转向后的链接

	Debug bool //调试开关，打开输出日志
}

func NewHttpClient() (this *HttpClient) {

	this = &HttpClient{}

	this.cookies = nil
	this.jar, _ = cookiejar.New(nil)

	this.c = &http.Client{Jar: this.jar, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		this.redirect = true
		this.redirectUrl = req.URL.String()
		//return errors.New("Redirected!")
		return nil
	}}
	this.ua = "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36"

	this.ResetHeader()
	this.conv = false

	return this
}

func (this *HttpClient) ResetHeader() {
	this.Header = make(http.Header)
	this.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	this.Header.Add("Accept-Encoding", "gzip,deflate,sdch")
	this.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	this.Header.Add("Connection", "keep-alive")
	this.Header.Add("User-Agent", this.ua)
}

//you should set it only once!
func (this *HttpClient) SetCharSet(charset string) {
	this.conv = true
	this.charset = charset
	this.enc = mahonia.NewEncoder(charset)
	this.dec = mahonia.NewDecoder(charset)
}

func (this *HttpClient) SetUa(ua string) {
	this.ua = ua
	this.Header.Set("User-Agent", ua)
}

func (this *HttpClient) Enc(in string) string {
	if this.conv {
		return this.enc.ConvertString(in)
	}
	return in
}

func (this *HttpClient) Dec(in string) string {
	if this.conv {
		return this.dec.ConvertString(in)
	}
	return in
}

func (this *HttpClient) Get(surl string) (page string, err error) {
	this.clearRedirect()
	//自动转化为GB2312编码
	gbkUrl := this.Enc(surl)

	req, err := http.NewRequest("GET", gbkUrl, nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.Get(%s),NewRequest error:%s", surl, err.Error()))
		return
	}

	for k, v := range this.Header {
		req.Header.Add(k, v[0])
	}

	if this.Debug {
		headerStr := ""
		for k, v := range req.Header {
			headerStr += fmt.Sprintf("%s:%s\n", k, v)
		}
		log.Printf("Header======\n%sHeader======\n", headerStr)
	}

	resp, err := this.c.Do(req)

	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.Get(%s),Response error:%s", surl, err.Error()))
		return
	}

	defer resp.Body.Close()

	var body string

	if resp.StatusCode == 200 {

		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, e := gzip.NewReader(resp.Body)
			if e != nil {
				err = errors.New(fmt.Sprintf("HttpClient.Get(%s),Read gzip body error:%s", surl, e.Error()))
				return
			}
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					panic(err)
				}

				if n == 0 {
					break
				}
				body += string(buf)
			}
		default:
			bodyByte, e := ioutil.ReadAll(resp.Body)
			if e != nil {
				err = errors.New(fmt.Sprintf("HttpClient.Get(%s),Read body error:%s", surl, e.Error()))
				return
			}
			body = string(bodyByte)
		}
	}

	//log.Println(body)

	this.cookies = this.jar.Cookies(req.URL)

	page = this.Dec(body)
	if this.Debug {
		log.Printf("================\nStatusCode:%d\nContent:\n%s\n================\n", resp.StatusCode, page)
	}
	return
}

func (this *HttpClient) Post(surl, postdata string) (page string, err error) {
	this.clearRedirect()
	//自动转化为GB2312编码
	gbkUrl := this.Enc(surl)
	gbkPostdata := this.Enc(postdata)

	req, _ := http.NewRequest("POST", gbkUrl, bytes.NewReader([]byte(gbkPostdata)))
	for k, v := range this.Header {
		req.Header.Add(k, v[0])
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if this.Debug {
		headerStr := ""
		for k, v := range req.Header {
			headerStr += fmt.Sprintf("%s:%s\n", k, v)
		}
		log.Printf("Header======\n%sHeader======\n", headerStr)
	}

	resp, err := this.c.Do(req)

	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.Post(%s,%s),Response error:%s", surl, postdata, err.Error()))
		return
	}
	//	log4go.Debug("Response code:%d status:%s", resp.StatusCode, resp.Status)

	defer resp.Body.Close()

	var body string

	if resp.StatusCode == 200 {

		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, e := gzip.NewReader(resp.Body)
			if e != nil {
				err = errors.New(fmt.Sprintf("HttpClient.Post(%s),Read gzip body error:%s", surl, e.Error()))
				return
			}
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					panic(err)
				}

				if n == 0 {
					break
				}
				body += string(buf)
			}
		default:
			bodyByte, e := ioutil.ReadAll(resp.Body)
			if e != nil {
				err = errors.New(fmt.Sprintf("HttpClient.Post(%s),Read body error:%s", surl, e.Error()))
				return
			}
			body = string(bodyByte)
		}
	}

	this.cookies = this.jar.Cookies(req.URL)

	page = this.Dec(body)

	if this.Debug {
		log.Printf("================\nHttpClient.Post(%s,%s)\nStatusCode:%d\nContent:\n%s\n================\n", surl, postdata, resp.StatusCode, page)
	}
	return
}

/*
///尚未测试
func (this *HttpClient) PostValues(surl string, postDict map[string]string) (page string, err error) {
	log4go.Debug("HttpClient.PostValues(%s,%T)", surl, postDict)
	postValues := url.Values{}
	for postKey, PostValue := range postDict {
		postValues.Set(this.Enc(postKey), this.Enc(PostValue))
	}
	postDataStr := postValues.Encode()

	return this.Post(surl, postDataStr)
}
*/
func (this *HttpClient) PostMultipart(u string, w *multipart.Writer, b *bytes.Buffer) (page string, err error) {
	this.clearRedirect()
	//log4go.Debug("HttpClient.PostMultipart(%s,w)", u)
	//自动转化为GB2312编码
	gbkUrl := this.Enc(u)

	req, _ := http.NewRequest("POST", gbkUrl, b)
	for k, v := range this.Header {
		req.Header.Add(k, v[0])
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	if this.Debug {
		headerStr := ""
		for k, v := range req.Header {
			headerStr += fmt.Sprintf("%s:%s\n", k, v)
		}
		log.Printf("Header======\n%sHeader======\n", headerStr)
	}

	//log4go.Finest("PostMultipart Content-Type: %s", w.FormDataContentType())
	resp, err := this.c.Do(req)

	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.PostMultipart ,Response error:%s", err.Error()))
		return
	}
	//log4go.Debug("Response code:%d status:%s", resp.StatusCode, resp.Status)

	defer resp.Body.Close()

	var body string

	if resp.StatusCode == 200 {

		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, e := gzip.NewReader(resp.Body)
			if e != nil {
				err = errors.New(fmt.Sprintf("HttpClient.PostMultipart(%s),Read gzip body error:%s", u, e.Error()))
				return
			}
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					panic(err)
				}

				if n == 0 {
					break
				}
				body += string(buf)
			}
		default:
			bodyByte, e := ioutil.ReadAll(resp.Body)
			if e != nil {
				err = errors.New(fmt.Sprintf("HttpClient.PostMultipart(%s),Read body error:%s", u, e.Error()))
				return
			}
			body = string(bodyByte)
		}
	}

	this.cookies = this.jar.Cookies(req.URL)

	page = this.Dec(body)

	//log4go.Finest("HttpClient.PostMultipart to url :%s returns :", u)
	//log4go.Finest(page)

	if this.Debug {
		log.Printf("================\nHttpClient.PostMultipart(%s)\nStatusCode:%d\nContent:\n%s\n================\n", u, resp.StatusCode, page)
	}
	return

}

func (p *HttpClient) clearRedirect() {
	p.redirect = false
	p.redirectUrl = ""
}
func (p *HttpClient) CheckRedirect() (b bool, surl string) {
	return p.redirect, p.redirectUrl
}
