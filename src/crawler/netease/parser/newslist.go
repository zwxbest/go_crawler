package parser

import (
	"github.com/tebeka/selenium"
	"crawler/engine"
	"log"
	"strings"
	"crawler/storage"
)

func ParseNewsLists(bodyEle selenium.WebElement)  engine.ParseResult {

	result := engine.ParseResult{}

	navItems,_ := bodyEle.FindElements(selenium.ByClassName,"nav_name")
	moreNav,_ := bodyEle.FindElement(selenium.ByClassName,"more_list")
	moreNavs,_ := moreNav.FindElements(selenium.ByTagName,"a")

	navItems = append(navItems, moreNavs...)


	for _,navItem := range navItems {
		href,err := navItem.GetAttribute("href")
		if err != nil {
			continue
		}
		if !strings.Contains(href, "http") {
			continue
		}
		name,_ := navItem.Text()
		newCategory := storage.NewsCategory{
			Name:name,
			Url:href,
		}
		storage.Insert(newCategory)
		result.Requests = append(result.Requests,engine.Request{
			Url:href,
			ParserFunc:ParseNewsLists,
		})
	}

	topNewUlEles,_ := bodyEle.FindElements(selenium.ByClassName,"top_news_ul")
	topNewTitleEle,_ := bodyEle.FindElement(selenium.ByClassName,"top_news_title")
	if topNewUlEles!=nil {
		topNewUlEles = append(topNewUlEles, topNewTitleEle)

	}

	for _,topNewEle := range topNewUlEles {
		if topNewEle == nil {
			continue
		}
	    parentEle,_ := topNewEle.FindElement(selenium.ByXPATH,"..")
		className,_ := parentEle.GetAttribute("class")
		if strings.Contains(className," none") {
	   	continue
		}
		hrefEles,_ := topNewEle.FindElements(selenium.ByTagName,"a")
		for _,hreEle := range hrefEles {
			result = getHrefFromEle(hreEle,result)
		}
	}

	newsArticles,_ := bodyEle.FindElements(selenium.ByClassName,"news_title")
	for _,newsArtice := range newsArticles {
		newsHref,_ :=	newsArtice.FindElement(selenium.ByXPATH,"h3/a")
		if newsHref == nil {
			continue
		}
		result = getHrefFromEle(newsHref,result)
	}
	return result
}

func getHrefFromEle(element selenium.WebElement,result engine.ParseResult)  engine.ParseResult{

	log.Println(element.Text())
	href,_ := element.GetAttribute("href")
	result.Items = append(result.Items,"news_href:"+href)
	result.Requests = append(result.Requests, engine.Request{
		Url:href,
		ParserFunc:ParseCommentHref,
	})
	return result
}