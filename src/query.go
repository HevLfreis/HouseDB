/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-24 20:26:23
 */
package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/influxdata/influxdb/client/v2"
)

func qTags(sql string) ([]string, error) {
	var res []string

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	q := client.NewQuery(sql, DB, DB_PRECISION)
	log.Debug("query influx", "sql", sql)

	if resp, err := c.Query(q); err == nil && resp.Error() == nil {
		series := resp.Results[0].Series
		if len(series) == 0 {
			log.Debug("query tags empty", "sql", sql)
			return res, nil
		}

		res := make([]string, len(series[0].Values))
		for i, v := range series[0].Values {
			res[i] = v[1].(string)
		}
		return res, nil
	} else {
		return nil, err
	}
}

func qDistricts() ([]string, error) {
	sql := `SHOW TAG VALUES FROM "house" WITH KEY="district"`
	return qTags(sql)
}

func qAreasInDistrict(d string) ([]string, error) {
	sql := sprintf(
		`SHOW TAG VALUES FROM "house" WITH KEY="area" WHERE "district"='%s'`, d)
	return qTags(sql)
}

func qComplexsInArea(d string, a string) ([]string, error) {
	sql := sprintf(
		`SHOW TAG VALUES FROM "house" WITH KEY="complex" WHERE "district"='%s' AND "area"='%s'`, d, a)
	return qTags(sql)
}

func qHouseList(district string, area string, comp string, each int, page int) ([]*House, error) {
	var conds []string
	if district != "" {
		conds = append(conds, sprintf(`"district"='%s'`, district))
	}
	if area != "" {
		conds = append(conds, sprintf(`"area"='%s'`, area))
	}
	if comp != "" {
		conds = append(conds, sprintf(`"complex"='%s'`, comp))
	}
	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	soffset := "SOFFSET " + strconv.Itoa(page*each)
	sql := sprintf(`SELECT "hid","url","district","area","complex","address","title",`+
		`"build_year","layout","total","per_m2","downpayment","metro","hot_total","hot_7days"`+
		` FROM "house" %s GROUP BY "hid" ORDER BY time DESC LIMIT 1 slimit %d %s`, where, each, soffset)

	var res []*House

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return res, err
	}
	defer c.Close()

	q := client.NewQuery(sql, DB, DB_PRECISION)
	log.Debug("query influx", "sql", sql)

	if resp, err := c.Query(q); err == nil && resp.Error() == nil {
		series := resp.Results[0].Series

		res = make([]*House, len(series))
		for i, s := range series {
			h := NewHouseFromQuery(s.Values[0][1:])
			res[i] = h
		}
	}
	return res, nil
}

func qMeanRecentPerM2(district string, area string, comp string) (json.Number, error) {
	var conds []string
	if district != "" {
		conds = append(conds, sprintf(`"district"='%s'`, district))
	}
	if area != "" {
		conds = append(conds, sprintf(`"area"='%s'`, area))
	}
	if comp != "" {
		conds = append(conds, sprintf(`"complex"='%s'`, comp))
	}
	where := ""
	if len(conds) > 0 {
		where = "AND " + strings.Join(conds, " AND ")
	}

	sql := sprintf(`SELECT MEAN("per_m2") FROM "house" WHERE time > now() - 1d %s`, where)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return "", err
	}
	defer c.Close()

	q := client.NewQuery(sql, DB, DB_PRECISION)
	log.Debug("query influx", "sql", sql)

	if resp, err := c.Query(q); err == nil && resp.Error() == nil {
		series := resp.Results[0].Series
		if len(series) == 0 {
			return "", nil
		}

		n := series[0].Values[0][1].(json.Number)
		return n, nil
	} else {
		return "", err
	}
}

func qSeriesPerM2(district string, area string, comp string, groupby string) ([]client.Result, error) {
	var conds []string
	if district != "" {
		conds = append(conds, sprintf(`"district"='%s'`, district))
	}
	if area != "" {
		conds = append(conds, sprintf(`"area"='%s'`, area))
	}
	if comp != "" {
		conds = append(conds, sprintf(`"complex"='%s'`, comp))
	}

	switch groupby {
	case "month":
		groupby = "time > now() - 360d GROUP BY time(30d)"
	case "year":
		groupby = "GROUP BY time(360d)"
	case "day":
		groupby = "time > now() - 30d GROUP BY time(5d)"
	default:
		groupby = "time > now() - 30d GROUP BY time(5d)"
	}

	var where string
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ") + " AND " + groupby
	} else {
		where = "WHERE " + groupby
	}

	sql := sprintf(`SELECT MEAN("per_m2") FROM "house" %s fill(-1)`, where)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	q := client.NewQuery(sql, DB, DB_PRECISION)
	log.Debug("query influx", "sql", sql)

	if resp, err := c.Query(q); err == nil && resp.Error() == nil {
		return resp.Results, nil
	} else {
		return nil, err
	}
}

func qSeriesHouse(hid string) ([]client.Result, error) {

	sql := sprintf(`SELECT MEAN("per_m2") FROM "house" WHERE "hid"='%s' GROUP BY time(1d) fill(-1)`, hid)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	q := client.NewQuery(sql, DB, DB_PRECISION)
	log.Debug("query influx", "sql", sql)

	if resp, err := c.Query(q); err == nil && resp.Error() == nil {
		return resp.Results, nil
	} else {
		return nil, err
	}
}

func qRecentHids(district string, area string, comp string) ([]string, error) {
	var conds []string
	if district != "" {
		conds = append(conds, sprintf(`"district"='%s'`, district))
	}
	if area != "" {
		conds = append(conds, sprintf(`"area"='%s'`, area))
	}
	if comp != "" {
		conds = append(conds, sprintf(`"complex"='%s'`, comp))
	}
	where := ""
	if len(conds) > 0 {
		where = "AND " + strings.Join(conds, " AND ")
	}

	sql := sprintf(
		`SHOW TAG VALUES FROM "house" WITH KEY="hid" WHERE time > now() - 3d %s`, where)
	return qTags(sql)
}

func qMaxMinHistoryPerM2(district string, area string, comp string) (json.Number, json.Number, error) {
	var conds []string
	if district != "" {
		conds = append(conds, sprintf(`"district"='%s'`, district))
	}
	if area != "" {
		conds = append(conds, sprintf(`"area"='%s'`, area))
	}
	if comp != "" {
		conds = append(conds, sprintf(`"complex"='%s'`, comp))
	}
	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	} else {

		// Todo: avoid influx oom
		return "", "", nil
	}

	sql := sprintf(`SELECT MAX("mean"),MIN("mean") FROM (SELECT MEAN("per_m2") FROM "house" %s GROUP BY time(5d))`, where)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return "", "", err
	}
	defer c.Close()

	q := client.NewQuery(sql, DB, DB_PRECISION)
	log.Debug("query influx", "sql", sql)

	if resp, err := c.Query(q); err == nil && resp.Error() == nil {
		series := resp.Results[0].Series
		if len(series) == 0 {
			return "", "", nil
		}

		max := series[0].Values[0][1].(json.Number)
		min := series[0].Values[0][2].(json.Number)
		return max, min, nil
	} else {
		return "", "", err
	}
}
