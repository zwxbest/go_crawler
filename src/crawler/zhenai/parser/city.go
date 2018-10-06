package parser

import (
	"crawler/engine"
	"regexp"
)

const cityRe = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`
var (
	profileRe = regexp.MustCompile(cityRe)
)

func ParseCity(contents []byte) engine.ParseResult {

	matches := profileRe.FindAllSubmatch(contents,-1)
	result :=engine.ParseResult{}

	for _, m := range matches{
		name:=string(m[2])
		result.Items = append(result.Items,"User:"+string(name))
		result.Requests = append(result.Requests,
			engine.Request{
				Url:string(m[1]),
				ParserFunc: func(bytes []byte) engine.ParseResult {
				 return ParseProfile(bytes,name)
				},
			})
	}

	return result

}