package store

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestHDFS(t *testing.T) {
	data, err := ioutil.ReadFile("/home/xxx/test.json")
	if err != nil {
		fmt.Println(fmt.Errorf("read file failed:%v", err))
		return
	}

	err = createFromData("xxx", "test.json", data)
	if err != nil {
		fmt.Println(fmt.Errorf("create hdfs file failed:%v", err))
		return
	}
}