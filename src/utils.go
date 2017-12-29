/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-22 13:09:01
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

var REG_NUM = regexp.MustCompile("[0-9]+")
var REG_LIANJIA_ID = regexp.MustCompile("sh[0-9]+")

// shortcut
func print(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// http
func httpGet(url string) (string, error) {
	client := &fasthttp.Client{}
	_, body, err := client.Get(nil, url)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

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

// error
func DomNotFound() error {
	return errors.New(ERR_DOM_NOT_FOUND)
}

// format
func Jtoi(j json.Number) int {
	i, _ := strconv.Atoi(string(j))
	return i
}

func formatUnixTime(s string, format string) string {
	i, _ := strconv.ParseInt(s, 10, 64)
	t := time.Unix(i, 0)
	return t.Format(format)
}
