//An easier way to use http.Client
package com

import (
	"bytes"
	"code.google.com/p/mahonia"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
)

type HttpClient struct {
	c       *http.Client
	cookies []*http.Cookie
	jar     *cookiejar.Jar
	ua      string //user agent

	//编码转换相关处理
	conv    bool //conv between utf-8 and charset
	charset string
	enc     mahonia.Encoder
	dec     mahonia.Decoder

	//链接转向相关处理
	redirect    bool   //是否转向了。每次Get之前置为false
	redirectUrl string //转向后的链接
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

	this.ua = "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)" //默认IE6

	this.conv = false

	return this
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

func (this *HttpClient) Get(url string) (page string, err error) {
	this.clearRedirect()
	//自动转化为GB2312编码
	gbkUrl := this.Enc(url)

	req, err := http.NewRequest("GET", gbkUrl, nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.Get(%s),NewRequest error:%s", url, err.Error()))
		return
	}
	req.Header.Add("User-Agent", this.ua)

	resp, err := this.c.Do(req)

	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.Get(%s),Response error:%s", url, err.Error()))
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.Get(%s),Read body error:%s", url, err.Error()))
		return
	}

	this.cookies = this.jar.Cookies(req.URL)

	page = this.Dec(string(body))
	return
}

func (this *HttpClient) Post(url, postdata string) (page string, err error) {
	this.clearRedirect()
	//自动转化为GB2312编码
	gbkUrl := this.Enc(url)
	gbkPostdata := this.Enc(postdata)

	req, _ := http.NewRequest("POST", gbkUrl, bytes.NewReader([]byte(gbkPostdata)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", this.ua)
	resp, err := this.c.Do(req)

	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.Post(%s,%s),Response error:%s", url, postdata, err.Error()))
		return
	}
	//	log4go.Debug("Response code:%d status:%s", resp.StatusCode, resp.Status)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.Post(%s,%s),Read body error:%s", url, postdata, err.Error()))
		return
	}

	this.cookies = this.jar.Cookies(req.URL)

	page = this.Dec(string(body))

	//log4go.Finest("HttpClient.Post(%s,%s) returns:", url, postdata)
	//log4go.Finest(page)
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
	req.Header.Add("User-Agent", this.ua)
	req.Header.Add("Content-Type", w.FormDataContentType())
	//log4go.Finest("PostMultipart Content-Type: %s", w.FormDataContentType())
	resp, err := this.c.Do(req)

	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.PostMultipart ,Response error:%s", err.Error()))
		return
	}
	//log4go.Debug("Response code:%d status:%s", resp.StatusCode, resp.Status)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("HttpClient.PostMultipart ,Read body error:%s", err.Error()))
		return
	}

	this.cookies = this.jar.Cookies(req.URL)

	page = this.Dec(string(body))

	//log4go.Finest("HttpClient.PostMultipart to url :%s returns :", u)
	//log4go.Finest(page)
	return

}

func (p *HttpClient) clearRedirect() {
	p.redirect = false
	p.redirectUrl = ""
}
func (p *HttpClient) CheckRedirect() (b bool, url string) {
	return p.redirect, p.redirectUrl
}
