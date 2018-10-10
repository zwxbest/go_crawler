package engine

import (
	"log"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"fmt"
	"os"
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

func (e *ConcurrentEngine) Run(seeds ...Request) {

	fd,_:= os.OpenFile("a.txt",os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	defer fd.Close()

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
		result := <- out
		for _, item := range result.Items {
			content := []byte(fmt.Sprintf("%v\n", item))
			fd.Write(content)
			
			log.Printf("Got item : %v", item)
		}
		fd.Close()
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
