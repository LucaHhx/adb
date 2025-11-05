package uixml

import (
	"fmt"
	"regexp"
	"strconv"
)

var boundsRe = regexp.MustCompile(`\[(\d+),(\d+)]\[(\d+),(\d+)]`)

// Rect 表示一个矩形
type Rect struct{ X1, Y1, X2, Y2 int }

// ParseBounds 解析类似 "[32,1410][570,1512]"
func ParseBounds(b string) (Rect, error) {
	m := boundsRe.FindStringSubmatch(b)
	if len(m) != 5 {
		return Rect{}, fmt.Errorf("bad bounds format: %q", b)
	}
	x1, _ := strconv.Atoi(m[1])
	y1, _ := strconv.Atoi(m[2])
	x2, _ := strconv.Atoi(m[3])
	y2, _ := strconv.Atoi(m[4])
	return Rect{X1: x1, Y1: y1, X2: x2, Y2: y2}, nil
}
