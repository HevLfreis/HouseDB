/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-24 20:26:23
 */
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	DB           = "housedb"
	DB_ADDR      = "http://localhost:8086"
	DB_PRECISION = "s"
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

func qSeriesPerM2(district string, area string, comp string) ([]client.Result, error) {
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
	var where string
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	sql := sprintf(`SELECT MEAN("per_m2") FROM "house" %s GROUP BY time(1d) fill(50000)`, where)

	c, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	// if err != nil {
	// 	return "", err
	// }
	defer c.Close()

	q := client.NewQuery(sql, DB, DB_PRECISION)
	response, err := c.Query(q)
	return response.Results, err
}

func formatUnixTime(s json.Number, format string) string {
	i, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm.Format(format)
}

func Jtoi(j json.Number) int {
	i, err := strconv.Atoi(string(j))
	return i
}

func print(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func main() {
	// qGetAllDistricts()
	// sqlGetAllDistricts2()

	// qLastMeanPerM2("", "", "")
	// qLastMeanPerM2("静安", "", "")
	// qLastMeanPerM2("静安", "莘庄", "")
	// qLastMeanPerM2("", "莘庄", "")
	// qLastMeanPerM2("闵行", "莘庄", "")
	// qHouseList("闵行", "莘庄", "", 30, 0)
	// qHouseList("闵行", "九亭", "", 3, 0)
	// fmt.Println("================")
	// qHouseList("", "九亭", "", 3, 0)
	qSeriesPerM2("闵行", "", "")
	qSeriesPerM2("闵行", "莘庄", "")

}
