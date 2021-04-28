package conf

import (
	"io/ioutil"
)

type Article struct {
	Title string
	Content string
}

func LoadArticles(dirPath string) (*[]Article,error) {
	var articles []Article
	// 扫描目录下所有文件
	infos,err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		title := info.Name()
		d,err := ioutil.ReadFile(dirPath+"/"+title)
		if err != nil {
			return nil,err
		}
		content := string(d)
		if content == "" {
			continue
		}
		articles =  append(articles,Article{Title: title,Content: content} )
	}
	return &articles,nil
}