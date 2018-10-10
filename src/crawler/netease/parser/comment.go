package parser

import (
	"github.com/tebeka/selenium"
	"crawler/engine"
	"log"
	"strconv"
	"strings"
)
/**
解析评论
 */
func ParseComment(bodyEle selenium.WebElement) engine.ParseResult {



	newComment,_ := bodyEle.FindElement(selenium.ByClassName,"tie-new");
	if newComment != nil {
		return  parserTieNew(newComment)
	}
	newComment,_ = bodyEle.FindElement(selenium.ByID,"mainReplies")
	if newComment !=nil {

	}

}

func parserTieNew(tieNew selenium.WebElement) engine.ParseResult {

	result := engine.ParseResult{}


	hideElements,_ := tieNew.FindElements(selenium.ByClassName,"expand-flr");
	for _,hideElement := range hideElements {
		//隐藏内容全部展开
		hideElement.Click();
	}
	//先取出没有盖楼的
	noFloorComments,_:= tieNew.FindElements(selenium.ByCSSSelector,".tie-bdy:not(.has-flr)")
	for _,noFloorComment := range noFloorComments {
		commentEle,_:=  noFloorComment.FindElement(selenium.ByClassName,"tie-cnt")
		comment,_:= commentEle.Text()
		log.Println(comment)
	}

	//盖楼的评论,所有的一楼
	allBuildings,_ := tieNew.FindElements(selenium.ByClassName,"has-flr")
	var thisBuilding []selenium.WebElement
	var allFloor []string

	//遍历楼层
	for _,firstFloorCommentBody := range allBuildings {
		thisBuilding = nil;
		allFloor = nil
		//深度遍历
		firstFloor,_ := firstFloorCommentBody.FindElement(selenium.ByXPATH,"div[1]")
		thisBuilding = append(thisBuilding, firstFloor)

		curCommentEle,_ := firstFloorCommentBody.FindElement(selenium.ByXPATH,"div[last()]")
		curComment,_ := curCommentEle.Text();
		allFloor = append(allFloor,curComment)

		for ; len(thisBuilding) != 0;{
			curFloor := thisBuilding[0];
			thisBuilding = thisBuilding[1:]

			nextFloor,_ := curFloor.FindElement(selenium.ByXPATH,"div[@class='floor']")
			if nextFloor !=  nil {
				thisBuilding = append(thisBuilding, nextFloor)
			}
			curCommentEle,_ := curFloor.FindElement(selenium.ByXPATH,"div[last()]/div[@class='tie-cnt']")
			if curCommentEle != nil {
				curComment,_ := curCommentEle.Text()
				allFloor = append(allFloor,curComment)
			} else {
				//把flat-flrs都提取出来
				flatFlrs,_ := curFloor.FindElements(selenium.ByClassName,"tie-cnt")

				for i := len(flatFlrs)-1; i >= 0; i-- {
					flatComment,_ := flatFlrs[i].Text();
					allFloor = append(allFloor,flatComment)
				}
			}
		}

		//todo:，每栋楼一个item
		var spaceCount = 1;
		for i := len(allFloor)-1; i >= 0; i-- {

			result.Items = append(result.Items,strconv.Itoa(spaceCount) + strings.Repeat( " ",spaceCount) + allFloor[i])
			result.Requests = nil
			spaceCount++
		}


	}
}
