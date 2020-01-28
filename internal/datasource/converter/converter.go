package converter

import (
	"strconv"
	"time"
)

func FromFloat64(v interface{}) string {
	return strconv.FormatFloat(*(v.(*float64)), 'f', -1, 64)
}

func FromDate(v interface{}) string {
	return v.(*time.Time).Format("2006-01-02")
}

func FromDateTime(v interface{}) string {
	return v.(*time.Time).Format(time.RFC3339)
}

func FromString(v interface{}) string {
	return *(v.(*string))
}

func FromBoolean(v interface{}) string {
	return strconv.FormatBool(*v.(*bool))
}

func FromUInt(v interface{}) string {
	return strconv.Itoa(int(*(v.(*uint))))
}

func FromInt(v interface{}) string {
	return strconv.Itoa(int(*(v.(*int))))
}
