package fetcher

import (
	"net/http"
	"fmt"
	"bufio"
	"golang.org/x/text/transform"
	"io/ioutil"
	"io"
	"golang.org/x/text/encoding"
	"golang.org/x/net/html/charset"
	"time"
)


var rateLimiter = time.Tick(100 * time.Millisecond)
func Fetch(url string) ([]byte, error)  {
	<-rateLimiter//
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf("wrong status code: %d",resp.StatusCode)
	}
	bodyReader :=bufio.NewReader(resp.Body)
	e :=determinEncoding(resp.Body)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func determinEncoding(r io.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
