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
	router.HandleFunc("/", indexHandler).Methods("GET").Name("index")
	router.HandleFunc("/series", seriesHandler).Methods("GET").Name("series")
	router.HandleFunc("/series/house", houseHandler).Methods("GET").Name("house")
	return router
}
