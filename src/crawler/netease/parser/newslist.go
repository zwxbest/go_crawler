package parser

import (
	"github.com/tebeka/selenium"
	"crawler/engine"
)

func ParseNewsLists(bodyEle selenium.WebElement)  engine.ParseResult {

	result := engine.ParseResult{}

	topNewEles,_ := bodyEle.FindElements(selenium.ByClassName,"top_news")
	for _,topNewEle := range topNewEles {
		hrefEles,_ := topNewEle.FindElements(selenium.ByTagName,"a")
		for _,hreEle := range hrefEles {
			href,_ := hreEle.GetAttribute("href")
			result.Items = append(result.Items,"news_href:"+href)
			result.Requests = append(result.Requests, engine.Request{
				Url:href,
				ParserFunc:ParseCommentHref,
			})
		}
	}
	return result
}