// SPDX-License-Identifier: MIT

// Package fetch 数据的获取
package fetch

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/issue9/errwrap"
)

// ErrInvalidYear 无效的年份
//
// 年份只能介于 [2009, 当前) 的区间之间。
var ErrInvalidYear = errors.New("无效的年份")

const (
	baseURL   = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/"
	startYear = 2009
)

var lastYear = time.Now().Year() - 1

func allYears() []int {
	years := make([]int, 0, lastYear-startYear)
	for year := lastYear; year >= startYear; year-- {
		years = append(years, year)
	}

	return years
}

// Fetch 拉取指定年份的数据
//
// years 为指定的一个或多个年份，如果为空，则表示所有的年份。
// 年份时间为 http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/
// 上存在的时间，从 2009 开始，到当前年份的上一年。
func Fetch(dir string, years ...int) error {
	if len(years) == 0 {
		years = allYears()
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, year := range years {
		if err := fetchYear(dir, year); err != nil {
			return err
		}
	}
	return nil
}

func fetchYear(dir string, year int) error {
	if year < startYear || year > lastYear {
		return ErrInvalidYear
	}

	buf := &errwrap.Buffer{Buffer: bytes.Buffer{}}

	y := strconv.Itoa(year)

	dir = filepath.Join(dir, y)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return collect(dir, buf, baseURL+y)
}
