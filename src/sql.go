/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-24 20:26:23
 */
package main

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
)

func sqlGetAllDistricts() ([]string, error) {

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

		// Todo: maybe no result
		values := response.Results[0].Series[0].Values

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

		res := make([]string, len(values))
		for i, v := range values {
			res[i] = v[1].(string)
		}
		return res, nil
	} else {
		return nil, err
	}
}
