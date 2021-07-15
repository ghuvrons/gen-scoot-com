package decoder

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
)

func readUint(buf io.Reader, rk reflect.Kind, len int) uint64 {
	var maxlen int = 1
	switch rk {
	case reflect.Uint8, reflect.Int8:
		maxlen = 1
	case reflect.Uint16, reflect.Int16:
		maxlen = 2
	case reflect.Uint32, reflect.Int32, reflect.Float32:
		maxlen = 4
	case reflect.Uint64, reflect.Uint, reflect.Int64, reflect.Int, reflect.Float64:
		maxlen = 8
	}
	if len == 0 {
		len = maxlen
	}
	b := make([]byte, len)
	newb := make([]byte, 8)
	binary.Read(buf, binary.BigEndian, &b)
	for i := 0; i < len; i++ {
		newb[i] = b[i]
	}
	return binary.LittleEndian.Uint64(newb)
}

func Decode(buf io.Reader, v interface{}) {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr {
		return
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < rv.NumField(); i++ {
		var len int = 0
		var factor float64 = 1
		rvField := rv.Field(i)
		typeField := rv.Type().Field(i)
		tag := typeField.Tag
		if tagLen := tag.Get("len"); tagLen != "" {
			tmplen, _ := strconv.ParseInt(tagLen, 10, 64)
			len = int(tmplen)
		}
		if tagF := tag.Get("factor"); tagF != "" {
			factor, _ = strconv.ParseFloat(tagF, 64)
		}

		if rvField.IsValid() {
			if rvField.CanSet() {
				switch rk := rvField.Kind(); rk {

				case reflect.Ptr:
					rvField.Set(reflect.New(rvField.Type().Elem()))
					Decode(buf, rvField.Interface())

				case reflect.String:
					x := make([]byte, len)
					binary.Read(buf, binary.LittleEndian, &x)
					rvField.SetString(string(x))

				case reflect.Bool:
					var x bool
					binary.Read(buf, binary.LittleEndian, &x)
					rvField.SetBool(x)

				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
					x := readUint(buf, rk, len)
					if !rvField.OverflowUint(x) {
						rvField.SetUint(x)
					}

				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
					x := int64(readUint(buf, rk, len))
					if !rvField.OverflowInt(x) {
						rvField.SetInt(x)
					}

				case reflect.Float32:
					var x float32
					var x64 float64
					n := uint32(readUint(buf, rk, len))
					if factor != 1 {
						x64 = float64(n) * factor
						fmt.Println("floating", n, "x", factor, "=", x64)

					} else {
						x = math.Float32frombits(n)
						x64 = float64(x)
					}
					if !rvField.OverflowFloat(x64) {
						rvField.SetFloat(x64)
					} else {
						fmt.Println("  ", "overflow")
					}

				case reflect.Float64:
					var x float64
					n := readUint(buf, rk, len)
					if factor == 1 {
						x = float64(n) * factor

					} else {
						x = math.Float64frombits(n)
						x = float64(x)
					}
					if !rvField.OverflowFloat(x) {
						rvField.SetFloat(x)
					} else {
						fmt.Println("  ", "overflow")
					}

				default:
				}
			}
		}
	}
}
