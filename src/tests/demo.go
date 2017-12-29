/*
* Author: hevlhayt@foxmail.com
* Date:   2017-12-25 10:29:01
 */
package main

import (
	"fmt"
	"strconv"
	// "strings"
	"time"
)

func hello() []string {
	return nil

}

func main() {
	fmt.Println(strconv.FormatInt(time.Now().Unix(), 10))

	// a := []string{"1", "3", nil}
	// fmt.Println(strings.Join(a, ":"))
	hello()

	fmt.Println(fmt.Sprintf(`%s`, "hello"))

	// c := cron.New()
	// c.AddFunc("0 29 21 * * *", func() { fmt.Println("Every hour on the half hour") })
	// c.Start()
	// fmt.Println(strconv.Atoi("1.222"))
	// a := ""
	// if a {
	// 	fmt.Println("he")
	// }

	i, err := strconv.ParseInt("1514246400", 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	fmt.Println(tm.Format("2006-01"))

}
