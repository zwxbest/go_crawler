package parser

import (
	"github.com/tebeka/selenium"
	"crawler/engine"
)

func ParseCommentHref(bodyEle selenium.WebElement) engine.ParseResult {

	result := engine.ParseResult{}
	hrefEle,_ := bodyEle.FindElement(selenium.ByClassName,"post_cnum_tie")
	href,_ := hrefEle.GetAttribute("href")
	result.Items = append(result.Items,"news_href:"+href)
	result.Requests = append(result.Requests, engine.Request{
		Url:href,
		ParserFunc:ParseComment,
	})

	return result

}
