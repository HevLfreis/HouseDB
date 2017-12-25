/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-24 20:26:23
 */
package main

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	DB           = "housedb_test"
	DB_ADDR      = "http://localhost:8086"
	DB_PRECISION = "s"
)

func sqlGetAllDistricts() ([]string, error) {

	sql := `SHOW TAG VALUES FROM "house" WITH KEY="district"`

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	q := client.NewQuery(sql, "housedb_test", "ns")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		values := response.Results[0].Series[0].Values

		fmt.Println(values)

		res := make([]string, len(values))
		for i, v := range values {
			res[i] = v[1].(string)
		}
		return res, nil
	} else {
		return nil, err
	}
}

func sqlGetAllDistricts2() ([]string, error) {

	sql := `SHOW TAG VALUES FROM "house" WITH KEY="district"`

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	q := client.NewQuery(sql, DB, DB_PRECISION)
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		values := response.Results[0].Series[0].Values

		fmt.Println(values)

		res := make([]string, len(values))
		for i, v := range values {
			res[i] = v[1].(string)
		}
		return res, nil
	} else {
		return nil, err
	}
}

func sqlGetAreasInDistrict(d string) ([]string, error) {

	sql := fmt.Sprintf(
		`SHOW TAG VALUES FROM "house" WITH KEY="area" WHERE "district"='%s'`, d)

	fmt.Println(sql)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: DB_ADDR,
	})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	q := client.NewQuery(sql, DB, "ns")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		values := response.Results[0].Series[0].Values
		fmt.Println("res", response.Results)

		res := make([]string, len(values))
		for i, v := range values {
			res[i] = v[1].(string)
		}
		return res, nil
	} else {
		return nil, err
	}
}

func main() {
	sqlGetAllDistricts()
	sqlGetAllDistricts2()
	// INDEX_DISTRICT_AREAS := make(map[string][]string)

	// districts, _ := sqlGetAllDistricts()
	// for _, d := range districts {
	// 	areas, err := sqlGetAreasInDistrict(d)
	// 	fmt.Println(err)
	// 	fmt.Println(areas)
	// 	INDEX_DISTRICT_AREAS[d] = areas
	// }

	// fmt.Println(keys(INDEX_DISTRICT_AREAS))
}
