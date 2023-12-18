package util

import (
	"math/rand"
	"strconv"
	"time"
)

// RandomNumeric 生成指定位数的随机数
func RandomNumeric(size int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if size <= 0 {
		panic("{ size : " + strconv.Itoa(size) + " } must be more than 0 ")
	}
	value := ""
	for index := 0; index < size; index++ {
		value += strconv.Itoa(r.Intn(10))
	}

	return value
}

// EndOfDay 判断一天的结束
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}
