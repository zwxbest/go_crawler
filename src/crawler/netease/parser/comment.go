package parser

import (
	"github.com/tebeka/selenium"
	"crawler/engine"
	"strconv"
	"crawler/model"
	"strings"
	"time"
)

/**
解析评论
 */
func ParseComment(bodyEle selenium.WebElement, url string) engine.ParseResult {
	result := engine.ParseResult{}

	buildings := parserTieNew(bodyEle, url, []model.Building{})
	result.Items = append(result.Items, buildings)
	return result;
}

func parserTieNew(bodyEle selenium.WebElement, url string, buildings []model.Building) []model.Building {

	println("这是第"+strconv.Itoa(len(buildings)) +"页")
	tieNew, _ := bodyEle.FindElement(selenium.ByClassName, "tie-new");
	if tieNew == nil {
		return []model.Building{};
	}
	build := model.Building{}
	hideElements, _ := tieNew.FindElements(selenium.ByClassName, "expand-flr");
	for _, hideElement := range hideElements {
		//隐藏内容全部展开
		hideElement.Click();
	}
	//先取出没有盖楼的
	noFloorComments, _ := tieNew.FindElements(selenium.ByCSSSelector, ".tie-bdy:not(.has-flr)")
	for _, noFloorComment := range noFloorComments {
		commentEle, _ := noFloorComment.FindElement(selenium.ByClassName, "tie-cnt")
		if commentEle == nil {
			break
		}
		comment, _ := commentEle.Text()
		build.FlatComments = append(build.FlatComments, comment)
	}

	//盖楼的评论,所有的一楼
	allBuildings, _ := tieNew.FindElements(selenium.ByClassName, "has-flr")
	var thisBuilding []selenium.WebElement
	var allFloor []string

	//遍历楼层
	for _, firstFloorCommentBody := range allBuildings {
		thisBuilding = nil;
		allFloor = nil
		//深度遍历
		firstFloor, _ := firstFloorCommentBody.FindElement(selenium.ByXPATH, "div[1]")
		thisBuilding = append(thisBuilding, firstFloor)

		curCommentEle, _ := firstFloorCommentBody.FindElement(selenium.ByXPATH, "div[last()]")
		curComment, _ := curCommentEle.Text();
		allFloor = append(allFloor, curComment)

		for ; len(thisBuilding) != 0; {
			curFloor := thisBuilding[0];
			thisBuilding = thisBuilding[1:]

			nextFloor, _ := curFloor.FindElement(selenium.ByXPATH, "div[@class='floor']")
			if nextFloor != nil {
				thisBuilding = append(thisBuilding, nextFloor)
			}
			curCommentEle, _ := curFloor.FindElement(selenium.ByXPATH, "div[last()]/div[@class='tie-cnt']")
			if curCommentEle != nil {
				curComment, _ := curCommentEle.Text()
				allFloor = append(allFloor, curComment)
			} else {
				//把flat-flrs都提取出来
				flatFlrs, _ := curFloor.FindElements(selenium.ByClassName, "tie-cnt")

				for i := len(flatFlrs) - 1; i >= 0; i-- {
					flatComment, _ := flatFlrs[i].Text();
					allFloor = append(allFloor, flatComment)
				}
			}
		}
		build.FloorComments = append(build.FloorComments, reverseSlice(allFloor))
	}
	build.Url = url
	buildings = append(buildings, build)
	//获取下一页
	next, _ := tieNew.FindElement(selenium.ByCSSSelector, ".tie-new .next")
	if next != nil {
		display, _ := next.CSSProperty("color")
		if !strings.Contains(display, "204") {
			next.Click()
			time.Sleep(1)
			//engine.Driver.Refresh()
			nextBody, _ := engine.Driver.FindElement(selenium.ByTagName, "body");
			buildings = parserTieNew(nextBody, url, buildings)
		}
	}
	return buildings
}

func reverseSlice(allFloor []string) []string {
	var reversedFloors []string
	var spaceCount = 1
	for i := len(allFloor) - 1; i >= 0; i-- {
		reversedFloors = append(reversedFloors, strings.Repeat("  ", len(allFloor)-1-i)+strconv.Itoa(spaceCount)+" "+allFloor[i])
		spaceCount++
	}
	return reversedFloors

}
