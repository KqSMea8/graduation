package tools

import (
	"github.com/g10guang/graduation/constdef"
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	r, err := os.Open("/Users/g10guang/Public/png.png")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	jw, err := os.OpenFile("/Users/g10guang/Public/output.jpeg", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer jw.Close()
	pw, err := os.OpenFile("/Users/g10guang/Public/output.png", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer pw.Close()
	if err := ImageCompress(r, jw, pw, constdef.Png); err != nil {
		panic(err)
	}
}

