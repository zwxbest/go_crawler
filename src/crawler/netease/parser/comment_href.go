package parser

import (
	"github.com/tebeka/selenium"
	"crawler/engine"
)

func ParseCommentHref(bodyEle selenium.WebElement) engine.ParseResult {

	result := engine.ParseResult{}
	hrefEle,_ := bodyEle.FindElement(selenium.ByClassName,"post_cnum_tie")
	if hrefEle == nil {
		return result
	}
	href,_ := hrefEle.GetAttribute("href")
	result.Requests = append(result.Requests, engine.Request{
		Url:href,
		ParserFunc: func(element selenium.WebElement) engine.ParseResult {
		  	parseResult := ParseComment(element,href);
		  	return parseResult
		},
	})

	return result

}
