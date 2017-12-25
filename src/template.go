/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-23 14:45:27
 */
package main

import (
	"github.com/flosch/pongo2"
	"net/http"
)

var indexTpl = pongo2.Must(pongo2.FromFile("templates/index.html"))

func render(w http.ResponseWriter, tpl *pongo2.Template, context map[string]interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	context = pongo2.Context(context)
	err := tpl.ExecuteWriter(context, w)
	return err
}
