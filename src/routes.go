/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-22 20:34:29
 */
package main

import (
	"github.com/gorilla/mux"
)

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", houseHandler).Methods("GET").Name("index")
	return router
}
