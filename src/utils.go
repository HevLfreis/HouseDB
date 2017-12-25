/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-22 13:09:01
 */
package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

var REG_NUM = regexp.MustCompile("[0-9]+")
var REG_LIANJIA_ID = regexp.MustCompile("sh[0-9]+")

// http
func FastGet(url string) (string, error) {
	client := &fasthttp.Client{}
	_, body, err := client.Get(nil, url)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// advance

//regex
func trimSpaceAndNewLineAndTab(s string) string {
	return strings.Replace(strings.Replace(strings.Replace(s, " ", "", -1), "\n", "", -1), "\t", "", -1)
}

func extractNum(s string) (int, error) {
	return strconv.Atoi(REG_NUM.FindString(s))
}

func extractLianjiaId(s string) string {
	return REG_LIANJIA_ID.FindString(s)
}

// os
func dirExistedOrCreate(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
		return nil
	} else {
		return err
	}
}
