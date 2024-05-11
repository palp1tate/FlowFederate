package util

import "strconv"

func ParsePageAndPageSize(page, pageSize string) (p, pSize int64) {
	p, _ = strconv.ParseInt(page, 10, 64)
	pSize, _ = strconv.ParseInt(pageSize, 10, 64)
	if p <= 0 {
		p = 1
	}
	switch {
	case pSize > 100:
		pSize = 100
	case pSize <= 0:
		pSize = 10
	}
	return
}
