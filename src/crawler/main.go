package main

import (
	"crawler/engine"
	"crawler/scheduler"
	"crawler/netease/parser"
	"database/sql"
	"log"
	"crawler/storage"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tebeka/selenium"
	"flag"
)

var Url = "http://comment.tie.163.com/E2N3BDNK0001875P.html"

func main()  {
	url := flag.String("url", Url, "网易评论地址")

	flag.Parse()

	Url = *url;

	db, err := sql.Open("mysql", "root:123456@/163_comment")
	if err != nil {
		log.Println("数据库链接失败")
		return ;
	}
	//log.Println("数据库连接成功")
	storage.Db = *db;
	defer storage.Db.Close()

	e := engine.ConcurrentEngine{
		Scheduler : &scheduler.QueuedScheduler{},
		WorkerCount:1 ,
	}

	e.Run(engine.Request{
		Url:Url,
		ParserFunc:func(element selenium.WebElement) engine.ParseResult {
			parseResult :=  parser.ParseComment(element,Url);
			return parseResult
		},
	})

}



