/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-21 22:09:55
 */
package main

import (
	"github.com/inconshreveable/log15"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	PORT              string
	NUM_CRAWLER       int
	CRON_CRAWL        string
	INTERVAL_CRAWL_MS int

	DISTRICTS = []string{"全区域", "浦东", "闵行",
		"宝山", "徐汇", "普陀", "杨浦", "长宁", "松江",
		"嘉定", "黄浦", "静安", "闸北", "虹口", "青浦",
		"奉贤", "金山", "崇明", "上海周边"}

	log log15.Logger
)

func main() {
	router := initRouter()
	initLogger()

	INTERVAL_CRAWL_MS = 250
	CRON_CRAWL = "0 30 20 * * *" // every day at 20:30
	NUM_CRAWLER = 32
	PORT = "8082"

	cronCrawl()

	log.Info("start http server", "port", PORT)
	fasthttp.ListenAndServe(":"+PORT, fasthttpadaptor.NewFastHTTPHandler(router))
}
