/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-22 10:11:36
 */
package main

import (
	"container/list"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LianJiaCrawler struct {
	Id        string
	Targets   []string
	waitGroup *sync.WaitGroup

	succ int
	fail int
}

func NewCrawler(id string, targets []string, wg *sync.WaitGroup) *LianJiaCrawler {
	c := LianJiaCrawler{
		Id:        id,
		Targets:   targets,
		waitGroup: wg,
	}
	return &c
}

func (c *LianJiaCrawler) Start() {
	defer c.waitGroup.Done()

	res := list.New()

	for _, url := range c.Targets {
		resp, err := FastGet(url)
		if err != nil {
			// Todo: list loss, need retry
			log.Warn("get target list page fail", "url", url)
			continue
		}

		res_doc := HTMLParse(resp)
		items := res_doc.Find("ul", "class", "js_fang_list").FindAll("li")

		for _, item := range items {
			house := c.build(item)
			if house != nil {
				// log.Debug("get house", "house", house)
				res.PushBack(house)
				c.succ++

				if res.Len() >= 20 {
					if err := saveHousesToInflux(res); err != nil {
						log.Warn("save to influx fail", "err", err)
					} else {
						log.Info("save to influx ok")
					}
				}
			} else {
				c.fail++
			}
			time.Sleep(time.Duration(INTERVAL_CRAWL_MS) * time.Millisecond)
		}
	}

	if err := saveHousesToInflux(res); err != nil {
		log.Warn("save to influx fail", "err", err)
	}

	log.Info("crawler finish", "succ", c.succ, "fail", c.fail)
}

func (c *LianJiaCrawler) build(item Root) *House {
	dm_link := item.Find("a", "class", "text link-hover-green js_triggerGray js_fanglist_title")
	if !found(dm_link) {
		log.Warn(ERR_DOM_NOT_FOUND, "class", "rtext link-hover-green js_triggerGray js_fanglist_title")
		return nil
	}

	house_page := URL_SH_LIANJIA + dm_link.Attrs()["href"]
	log.Info("parse page start", "url", house_page)

	// info in search result item
	hid := extractLianjiaId(house_page)
	url := house_page
	title := dm_link.Attrs()["title"]

	dm_layout := item.Find("span", "class", "info-col row1-text")
	if !found(dm_layout) {
		log.Info(ERR_DOM_NOT_FOUND, "class", "info-col row1-text")
		return nil
	}
	layout := trimSpaceAndNewLineAndTab(dm_layout.Text())

	dm_total := item.Find("span", "class", "total-price strong-num")
	if !found(dm_total) {
		log.Warn(ERR_DOM_NOT_FOUND, "class", "total-price strong-num")
		return nil
	}
	total, _ := strconv.Atoi(dm_total.Text())
	total *= 10000

	dm_per_m2 := item.Find("span", "class", "info-col price-item minor")
	if !found(dm_per_m2) {
		log.Warn(ERR_DOM_NOT_FOUND, "class", "info-col price-item minor")
		return nil
	}
	per_m2, _ := extractNum(dm_per_m2.Text())

	var metro string
	tags := item.Find("div", "class", "property-tag-container").FindAll("span")
	if len(tags) > 0 {
		metro = tags[0].Text()
		if !strings.Contains(metro, "米") {
			metro = ""
		}
	}

	resp, err := FastGet(house_page)
	if err != nil {
		log.Warn("get house page fail", "url", house_page)
		return nil
	}
	house_doc := HTMLParse(resp)

	// info in house page
	dm_build_year := house_doc.Find("ul", "class", "maininfo-main maininfo-item").
		Find("li", "class", "main-item u-tr")
	if !found(dm_build_year) {
		log.Warn(ERR_DOM_NOT_FOUND, "class", "maininfo-main maininfo-item")
		return nil
	}
	build_year, _ := extractNum(trimSpaceAndNewLineAndTab(dm_build_year.FindAll("p")[1].Text()))

	dm_info := house_doc.Find("span", "class", "maininfo-estate-name")
	if !found(dm_info) {
		log.Warn(ERR_DOM_NOT_FOUND, "class", "maininfo-estate-name")
		return nil
	}
	info := dm_info.FindAll("a")
	district := info[1].Text()
	area := info[2].Text()
	comp := info[0].Text()

	dm_addr := house_doc.Find("span", "class", "item-cell maininfo-estate-address")
	if !found(dm_addr) {
		log.Warn(ERR_DOM_NOT_FOUND, "class", "item-cell maininfo-estate-address")
		return nil
	}
	addr := dm_addr.Text()

	dm_downpayment := house_doc.Find("ul", "class", "maininfo-minor maininfo-item")
	if !found(dm_downpayment) {
		log.Warn(ERR_DOM_NOT_FOUND, "class", "maininfo-minor maininfo-item")
		return nil
	}
	downpayment, _ := extractNum(dm_downpayment.FindAll("span")[1].Text())

	hot_total := 0
	hot_7days := 0

	house := NewHouse(
		hid,
		url,
		district,
		area,
		comp,
		addr,
		title,
		build_year,
		layout,
		total,
		per_m2,
		downpayment,
		metro,
		hot_total,
		hot_7days)

	return house

}

func crawlHouseLists() ([]string, error) {
	resp, err := FastGet(URL_SH_LIANJIA + URL_PREFIX_LIANJIA_ERSHOUFANG)
	if err != nil {
		log.Warn("get house list fail", "err", err)
		return nil, err
	}
	doc := HTMLParse(resp)

	dm_total_text := doc.Find("span", "class", "result-count strong-num")
	if !found(dm_total_text) {
		log.Warn(ERR_DOM_NOT_FOUND, "class", "result-count strong-num")
		return nil, DomNotFound()
	}

	total_results, _ := strconv.Atoi(strings.TrimSpace(dm_total_text.Text()))
	total_pages := total_results/30 + 1

	// test trick
	// total_pages = 1
	targets := make([]string, total_pages)

	for i := 0; i < total_pages; i++ {
		targets[i] = URL_SH_LIANJIA +
			URL_PREFIX_LIANJIA_ERSHOUFANG +
			URL_SUFFIX_LIANJIA_SEARCHPAGE +
			strconv.Itoa(i+1)
	}

	return targets, nil
}

func found(root Root) bool {
	if root.Pointer == nil {
		return false
	}
	return true
}

func DomNotFound() error {
	return errors.New(ERR_DOM_NOT_FOUND)
}
