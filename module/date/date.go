package date

import (
	"fmt"
	"time"
)

const yyyyMMdd = "yyyyMMdd"

// 格式化时间
func Format(f string) string {

	datestr := ""
	if f == yyyyMMdd {
		y, m, d, _, _, _ := Now()
		datestr = y + m + d
	}
	return datestr
}

func Now() (string, string, string, string, string, string) {
	now := time.Now()
	y := fmt.Sprintf("%d", now.Year())
	m := fmt.Sprintf("%d", now.Month())
	d := fmt.Sprintf("%d", now.Day())

	mstr := m
	dstr := d
	if len(m) == 1 {
		mstr = "0" + m
	}
	if len(d) == 1 {
		dstr = "0" + d
	}

	return y, mstr, dstr, "", "", ""
}
