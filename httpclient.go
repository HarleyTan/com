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
	//"net/url"
)

type HttpClient struct {
	c       *http.Client
	cookies []*http.Cookie
	jar     *cookiejar.Jar
	convGbk bool
	ua      string ///user agent
	enc     mahonia.Encoder
	dec     mahonia.Decoder
}

func NewHttpClient() (this *HttpClient) {

	this = &HttpClient{}

	this.cookies = nil
	//var err error;
	this.jar, _ = cookiejar.New(nil)

	c := &http.Client{Jar: this.jar}
	this.c = c

	this.ua = "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)" //默认IE6

	this.convGbk = true
	this.enc = mahonia.NewEncoder("GBK")
	this.dec = mahonia.NewDecoder("GBK")

	return this
}

func (this *HttpClient) SetConvGbk(b bool) {
	this.convGbk = b
}

func (this *HttpClient) SetUa(ua string) {
	this.ua = ua
}

func (this *HttpClient) Enc(in string) string {
	if this.convGbk {
		return this.enc.ConvertString(in)
	}
	return in
}

func (this *HttpClient) Dec(in string) string {
	if this.convGbk {
		return this.dec.ConvertString(in)
	}
	return in
}

func (this *HttpClient) Get(url string) (page string, err error) {
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

/*
func (this *HttpClient) getUrlRespHtml(strUrl string, postDict map[string]string) (page string, status string) {

	log4go.Debug("getUrlRespHtml, strUrl=%s", strUrl)
	log4go.Debug("postDict=%s", postDict)

	var httpReq *http.Request
	//var newReqErr error
	if nil == postDict {
		log4go.Debug("is GET")
		//httpReq, newReqErr = http.NewRequest("GET", strUrl, nil)
		httpReq, _ = http.NewRequest("GET", strUrl, nil)
		// ...
		//httpReq.Header.Add("If-None-Match", `W/"wyzzy"`)
	} else {
		log4go.Debug("is POST")
		postValues := url.Values{}
		for postKey, PostValue := range postDict {
			postValues.Set(postKey, PostValue)
		}
		log4go.Debug("postValues=%s", postValues)
		postDataStr := postValues.Encode()
		log4go.Debug("postDataStr=%s", postDataStr)
		postDataBytes := []byte(postDataStr)
		log4go.Debug("postDataBytes=%s", postDataBytes)
		postBytesReader := bytes.NewReader(postDataBytes)
		//httpReq, newReqErr = http.NewRequest("POST", strUrl, postBytesReader)
		httpReq, _ = http.NewRequest("POST", strUrl, postBytesReader)
		//httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	httpResp, err := this.c.Do(httpReq)
	// ...

	//httpResp, err := http.Get(strUrl)
	//log4go.Debug("http.Get done")
	if err != nil {
		log4go.Warn("http get strUrl=%s response error=%s\n", strUrl, err.Error())
	}
	log4go.Debug("httpResp.Header=%s", httpResp.Header)
	log4go.Debug("httpResp.Status=%s", httpResp.Status)

	status = httpResp.Status

	defer httpResp.Body.Close()
	// log4go.Debug("defer httpResp.Body.Close done")

	body, errReadAll := ioutil.ReadAll(httpResp.Body)
	//log4go.Debug("ioutil.ReadAll done")
	if errReadAll != nil {
		log4go.Warn("get response for strUrl=%s got error=%s\n", strUrl, errReadAll.Error())
	}
	//log4go.Debug("body=%s\n", body)

	//this.cookies = httpResp.Cookies()
	//gCurCookieJar = this.Jar;
	this.cookies = this.jar.Cookies(httpReq.URL)
	//log4go.Debug("httpResp.Cookies done")

	//respHtml = "just for test log ok or not"
	page = string(body)
	//log4go.Debug("httpResp body []byte to string done")

	return
}
*/

//func dbgPrintCurCookies(cookie []*http.Cookie) {
//	var cookieNum int = len(cookie)
//	log4go.Debug("cookieNum=%d", cookieNum)
//	for i := 0; i < cookieNum; i++ {
//		var curCk *http.Cookie = cookie[i]
//		//log4go.Debug("curCk.Raw=%s", curCk.Raw)
//		log4go.Debug("------ Cookie [%d]------", i)
//		log4go.Debug("Name\t\t=%s", curCk.Name)
//		log4go.Debug("Value\t=%s", curCk.Value)
//		log4go.Debug("Path\t\t=%s", curCk.Path)
//		log4go.Debug("Domain\t=%s", curCk.Domain)
//		log4go.Debug("Expires\t=%s", curCk.Expires)
//		log4go.Debug("RawExpires\t=%s", curCk.RawExpires)
//		log4go.Debug("MaxAge\t=%d", curCk.MaxAge)
//		log4go.Debug("Secure\t=%t", curCk.Secure)
//		log4go.Debug("HttpOnly\t=%t", curCk.HttpOnly)
//		log4go.Debug("Raw\t\t=%s", curCk.Raw)
//		log4go.Debug("Unparsed\t=%s", curCk.Unparsed)
//	}
//}
