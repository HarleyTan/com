package com

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"strings"
	"testing"
)

/*
func TestHttpClient(t *testing.T) {

	c := NewHttpClient()
	Convey("HttpClient 测试", t, func() {
		Convey("自定义Header测试", func() {
			//c.SetCharSet("gbk")
			page, err := c.Get("http://www.bing.com")
			//log.Println(page)
			So(err, ShouldEqual, nil)
			So(strings.Contains(page, "Bing"), ShouldEqual, true)
		})
	})

}
*/

func TestQfyf(t *testing.T) {
	c := NewHttpClient()
	Convey("清风扬帆", t, func() {
		page, err := c.Get("http://www.qfyf.net:8080/xxgk/visitor/jcmsWdyvote-index.c")
		//c.SetCharSet("UTF-8")
		log.Println(page)
		So(err, ShouldEqual, nil)
		So(strings.Contains(page, "清风扬帆"), ShouldEqual, true)
	})
}
