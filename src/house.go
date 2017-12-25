/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-22 09:23:45
 */
package main

import (
	"container/list"
	"net/http"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type House struct {
	Hid         string `json:"hid"`
	Url         string `json:"url"`
	District    string `json:"district"`
	Area        string `json:"area"`
	Complex     string `json:"complex"`
	Address     string `json:"address"`
	Title       string `json:"title"`
	BuildYear   int    `json:"build_year"`
	Layout      string `json:"layout"`
	Total       int    `json:"total"`
	PerM2       int    `json:"per_m2`
	Downpayment int    `json:"downpayment"`
	Metro       string `json:"metro"`
	HotTotal    int    `json:"hot_total"`
	Hot7Days    int    `json:"hot_7days"`
}

func NewHouse(hid string, url string, district string,
	area string, comp string, addr string,
	title string, build_year int, layout string,
	total int, per_m2 int, downpayment int,
	metro string, hot_total int, hot_7days int) *House {

	h := House{
		Hid:         hid,
		Url:         url,
		District:    district,
		Area:        area,
		Complex:     comp,
		Address:     addr,
		Title:       title,
		BuildYear:   build_year,
		Layout:      layout,
		Total:       total,
		PerM2:       per_m2,
		Downpayment: downpayment,
		Metro:       metro,
		HotTotal:    hot_total,
		Hot7Days:    hot_7days,
	}
	return &h
}

func houseHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("access index page", "host", r.Host)

	// context := map[string]interface{}{
	// 	"district_areas": INDEX_DISTRICT_AREAS,
	// }
	// render(w, indexTpl, context)
	render(w, indexTpl, nil)
}

func saveHousesToInflux(houses *list.List) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return err
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  DB,
		Precision: DB_PRECISION,
	})
	if err != nil {
		return err
	}

	for houses.Len() > 0 {

		item := houses.Front()
		house := item.Value.(*House)
		houses.Remove(item)

		tags := map[string]string{
			"hid":      house.Hid,
			"district": house.District,
			"area":     house.Area,
			"complex":  house.Complex,
		}
		fields := map[string]interface{}{
			"url":         house.Url,
			"address":     house.Address,
			"title":       house.Title,
			"build_year":  house.BuildYear,
			"layout":      house.Layout,
			"total":       house.Total,
			"per_m2":      house.PerM2,
			"downpayment": house.Downpayment,
			"metro":       house.Metro,
			"hot_total":   house.HotTotal,
			"hot_7days":   house.Hot7Days,
		}

		pt, err := client.NewPoint(MS_HOUSE, tags, fields, time.Now())
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}

	if err := c.Write(bp); err != nil {
		return err
	}

	return nil
}
