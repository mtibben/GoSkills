package numerics

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func Sqr(x float64) float64 {
	return x * x
}

// Round rounds a float to a specified number of fractional digits
func Round(x float64, digits int) float64 {
	frep := strconv.FormatFloat(x, 'f', digits, 64)
	f, _ := strconv.ParseFloat(frep, 64)
	return f
}

func Float64ToBytes(x float64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, x)
	if err != nil {
		panic("binary.Write failed: " + err.Error())
	}
	return buf.Bytes()
}

func BytesToInt32(b []byte) int32 {
	var x int32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &x)
	if err != nil {
		panic("binary.Read failed: " + err.Error())
	}
	return x
}
