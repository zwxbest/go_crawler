package engine

import (
	"log"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"fmt"
	"crawler/model"
	"io/ioutil"
)

type ConcurrentEngine struct {
	Scheduler        Scheduler
	WorkerCount      int
}

type Scheduler interface {
	Submit(Request)
	WorkerReady(chan Request)
	Run()
}

var NewsName = "";

func (e *ConcurrentEngine) Run(seeds ...Request) {

	out := make(chan ParseResult)
	e.Scheduler.Run()

	caps,_  := createChrome()

	for i := 0; i < e.WorkerCount; i++ {
		createWorker(out,e.Scheduler,caps)
	}

	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	for {
		//拿出最新的接着执行
		result := <- out
		for _, item := range result.Items {

			//content := []byte(fmt.Sprintf("%v\n", item))
			//log.Println(content);
			//building := model.Building(item);

			content :=""

			items,ok := item.([]model.Building)
			if ok == true {
				for _,building :=  range items{
					for _,flatComment  := range building.FlatComments{
						//println(flatComment)
						content += flatComment+"\n"
					}
					for _,floorComent := range building.FloorComments{
						for _,floor := range floorComent{
							//println(floor)
							content += floor+"\n"
						}
					}
				}
				err := ioutil.WriteFile(NewsName+".txt", []byte(content), 0644)
				if err != nil{
					panic(err)
				}

			}else {
				log.Printf("Got item : %v", items)
			}

		}
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}


func createWorker(out chan ParseResult, s Scheduler,caps selenium.Capabilities) {

	webDriver,_ := createSession(caps)

	in := make(chan Request)
	go func() {
		for {
			s.WorkerReady(in)
			request := <- in
			result,e := Worker(request,webDriver)
			if e != nil {
				continue
			}
			//ParserFunc返回的扔到out里
			out <- result
		}
	}()
}

func createChrome() (selenium.Capabilities,error) {
	opts := []selenium.ServiceOption{}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// 禁止加载图片，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		},
	}
	caps.AddChrome(chromeCaps)
	//启动chromedriver，端口号可自定义
	_, err := selenium.NewChromeDriverService("/opt/google/chronium/chromedriver", 9515, opts...)
	if err != nil {
		log.Printf("Error starting the ChromeDriver server: %v", err)
	}
	return caps,nil
}


func createSession(caps selenium.Capabilities) (selenium.WebDriver,error)  {
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	if err != nil {
		return nil,err
	}
	return webDriver,nil
}
