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
	INTERVAL_SCHED_H  uint64
	INTERVAL_CRAWL_MS int

	INDEX_DISTRICT_AREAS map[string][]string

	log log15.Logger
)

func main() {

	router := initRouter()
	initLogger()

	INTERVAL_CRAWL_MS = 250
	INTERVAL_SCHED_H = 12
	NUM_CRAWLER = 32
	PORT = "8082"

	// INDEX_DISTRICT_AREAS = make(map[string][]string)
	// districts, _ := sqlGetAllDistricts()
	// for _, d := range districts {
	// 	areas, _ := sqlGetAreasInDistrict(d)
	// 	INDEX_DISTRICT_AREAS[d] = areas
	// }

	go schedCrawl()
	// go crawl()

	log.Info("start http server", "port", PORT)
	// http.ListenAndServe(":"+PORT, nil)
	fasthttp.ListenAndServe(":"+PORT, fasthttpadaptor.NewFastHTTPHandler(router))
}
