package main

import (
	"crawler/engine"
	"crawler/netease/parser"
	"crawler/scheduler"
)

func main()  {
	e := engine.ConcurrentEngine{
		Scheduler : &scheduler.QueuedScheduler{},
		WorkerCount:10 ,
	}
	e.Run(engine.Request{
		Url:"https://news.163.com",
		ParserFunc:parser.ParseNewsLists,
	})

}



