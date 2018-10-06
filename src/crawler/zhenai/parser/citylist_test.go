package parser

import (
	"crawler/fetcher"
	"testing"
	"fmt"
)

func TestParseCityList(t *testing.T) {
	contents, e := fetcher.Fetch("http://www.zhenai.com/zhenghun")

	if e!=nil {
		panic(e)
	}
	fmt.Printf("%s\n",contents)
}
