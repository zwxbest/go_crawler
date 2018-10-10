package engine

import "github.com/tebeka/selenium"

type Request struct {
	Url string
	ParserFunc func(element selenium.WebElement) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items    []interface{}
}

