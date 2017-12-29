/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-22 20:27:06
 */
package main

import (
	"github.com/robfig/cron"
)

func cronCrawl() {
	c := cron.New()
	c.AddFunc(CRON_CRAWL, crawl)
	c.Start()
}
