package fetcher

import (
	"time"
	"github.com/tebeka/selenium"
	"log"
)


var rateLimiter = time.Tick(100 * time.Millisecond)
func Fetch(url string,webDriver selenium.WebDriver) (selenium.WebElement, error)  {
	<-rateLimiter//
	// 导航到目标网站
	err := webDriver.Get(url)
	if err != nil {
		return nil,err
	}
	title,_ := webDriver.Title()
	log.Println(title)
	bodyEle,_ := webDriver.FindElement(selenium.ByTagName,"body");

	return bodyEle,nil

}

