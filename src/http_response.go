/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-23 14:45:27
 */
package main

import (
	"encoding/json"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/influxdata/influxdb/client/v2"
)

// var indexTpl = pongo2.Must(pongo2.FromFile("templates/index.html"))

type Response struct {
	Errno   int             `json:"errno"`
	Message string          `json:"message"`
	Results []client.Result `json:"results"`
}

func render(w http.ResponseWriter, tpl *pongo2.Template, ctx map[string]interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	ctx = pongo2.Context(ctx)
	err := tpl.ExecuteWriter(ctx, w)
	return err
}

func jsonp(w http.ResponseWriter, e int, m string, r []client.Result) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{e, m, r})
}
