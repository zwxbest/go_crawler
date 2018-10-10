package fetcher

import (
	"time"
	"github.com/tebeka/selenium"
	"log"
)


var rateLimiter = time.Tick(100 * time.Millisecond)
func Fetch(url string,webDriver selenium.WebDriver) (selenium.WebElement, error)  {
	<-rateLimiter
	// 导航到目标网站
	err := webDriver.Get(url)
	if err != nil {
		return nil,err
	}
	bodyEle,err := webDriver.FindElement(selenium.ByTagName,"body");
	if err != nil {
		log.Fatal(err)
		return nil,err
	}
	return bodyEle,nil

}

