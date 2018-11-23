package engine

import (
	"crawler/fetcher"
	"log"
	"github.com/tebeka/selenium"
)


var Driver selenium.WebDriver;

func Worker(r Request,driver selenium.WebDriver) (ParseResult,error) {
	log.Printf("Fetching %s\n", r.Url)
	body, e := fetcher.Fetch(r.Url,driver)
	if e != nil{
		log.Printf("Fetcher: error fetching url %s: %v",
			r.Url, e)
		return ParseResult{},e
	}

	Driver = driver
	return r.ParserFunc(body), nil
}