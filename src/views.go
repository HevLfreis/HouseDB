/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-28 13:58:26
 */
package main

import (
	"net/http"
	"sort"

	"github.com/flosch/pongo2"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("access index page", "addr", r.Header.Get("X-Forwarded-For"))

	district := r.URL.Query().Get("district")
	area := r.URL.Query().Get("area")
	comp := r.URL.Query().Get("complex")
	log.Debug("get url params", "district", district, "area", area, "complex", comp)

	var cur_loc string
	if comp != "" {
		cur_loc = comp
	} else if area != "" {
		cur_loc = area
	} else {
		cur_loc = district
	}

	var houses []*House
	if comp != "" {
		houses, _ = qHouseList(district, area, comp, NUM_PER_PAGE_LIMIT, 0)
	}
	areas, _ := qAreasInDistrict(district)
	comps, _ := qComplexsInArea(district, area)
	mean, _ := qMeanRecentPerM2(district, area, comp)
	max, min, _ := qMaxMinHistoryPerM2(district, area, comp)
	hids, _ := qRecentHids(district, area, comp)

	sort.Strings(areas)
	sort.Strings(comps)

	context := map[string]interface{}{
		"districts":    DISTRICTS[1:],
		"cur_district": district,
		"areas":        areas,
		"cur_area":     area,
		"complexs":     comps,
		"cur_complex":  comp,
		"cur_loc":      cur_loc,
		"mean":         mean,
		"max":          max,
		"min":          min,
		"count":        len(hids),
		"houses":       houses,
	}

	indexTpl := pongo2.Must(pongo2.FromFile("templates/index.html"))
	render(w, indexTpl, context)
}

func seriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("access series api", "addr", r.Header.Get("X-Forwarded-For"))

	district := r.URL.Query().Get("district")
	area := r.URL.Query().Get("area")
	comp := r.URL.Query().Get("complex")
	groupby := r.URL.Query().Get("groupby")
	log.Debug("get url params", "district", district, "area", area, "complex", comp, "groupby", groupby)

	series, _ := qSeriesPerM2(district, area, comp, groupby)

	jsonp(w, ERRNO_OK, "query ok", series)
}

func houseHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("access house api", "addr", r.Header.Get("X-Forwarded-For"))

	hid := r.URL.Query().Get("hid")
	log.Debug("get url params", "hid", hid)

	series, _ := qSeriesHouse(hid)

	jsonp(w, ERRNO_OK, "query ok", series)
}
