/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-22 20:27:06
 */
package main

import (
	"sync"

	"github.com/jasonlvhit/gocron"
)

func schedCrawl() {
	gocron.Every(INTERVAL_SCHED_H).Hours().Do(crawl)
	// gocron.Every(2).Minutes().Do(crawl)

	<-gocron.Start()
}

func crawl() {
	log.Info("crawl start")

	targets, err := crawlHouseLists()
	if err != nil {
		log.Warn("get house list fail", "err", err)
		return
	}
	log.Info("get house list ok", "target", len(targets))

	var wg sync.WaitGroup
	wg.Add(NUM_CRAWLER)

	t := len(targets) / NUM_CRAWLER
	for i := 0; i < NUM_CRAWLER-1; i++ {
		c := NewCrawler("", targets[i*t:(i+1)*t], &wg)
		go c.Start()
	}
	c := NewCrawler("", targets[(NUM_CRAWLER-1)*t:], &wg)
	go c.Start()

	wg.Wait()

	log.Info("crawl end")
}
