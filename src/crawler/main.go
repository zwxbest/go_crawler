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
)

func main()  {

	db, err := sql.Open("mysql", "root:123456@/163_comment")
	if err != nil {
		log.Println("数据库链接失败")
		return ;
	}
	log.Println("数据库连接成功")
	storage.Db = *db;
	defer storage.Db.Close()

	e := engine.ConcurrentEngine{
		Scheduler : &scheduler.QueuedScheduler{},
		WorkerCount:1 ,
	}

	var url = "http://comment.tie.163.com/E1A9013K0001875P.html";
	e.Run(engine.Request{
		Url:url,
		ParserFunc:func(element selenium.WebElement) engine.ParseResult {
			parseResult :=  parser.ParseComment(element,url);
			return parseResult
		},
	})

}



