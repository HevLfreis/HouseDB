/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-22 20:32:25
 */
package main

import (
	"path/filepath"

	"github.com/inconshreveable/log15"
)

func initLogger() {
	path := "log"
	file := "housedb.log"

	log = log15.New()
	if err := dirExistedOrCreate(path); err != nil {
		log.Warn("create log path fail")
		return
	}

	handler := log15.MultiHandler(
		// log15.StreamHandler(os.Stderr, log15.TerminalFormat()),
		log15.Must.FileHandler(filepath.Join(path, file), log15.TerminalFormat()))
	log.SetHandler(handler)
	log.Info("start logging", "path", filepath.Join(path, file))
}
