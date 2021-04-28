package conf

import (
	"bufio"
	"io"
	"os"
	"statistician/util"
	"strings"
)

func LoadDictionary(files *[]string,bFiles *[]string) (map[string]string,map[string]string,error) {
	var dic = make(map[string]string)
	var ric = make(map[string]string)
	for _, f := range *files {
		err := getDictionaryFromFile("dictionary/" + f,dic)
		if err != nil {
			return nil,nil, err
		}
	}
	var bic = make(map[string]string)
	for _, f := range *bFiles {
		err := getDictionaryFromFile("dictionary/" + f,bic)
		if err != nil {
			return nil,nil, err
		}
	}
	for name := range bic {
		delete(dic, name)
	}
	for d := range dic {
		if hasWildcard(d) {
			delete(dic, d)
			ric[d] = ""
		}
	}
	return dic,ric,nil
}

func getDictionaryFromFile(path string,m map[string]string) error {
	f,err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	br := bufio.NewReader(f)
	for  {
		line,_,err := br.ReadLine()
		if err == io.EOF {
			break
		}
		l := string(line)
		if l == "" {
			continue
		}
		if hasQuote(l) {
			//存在引用关系则进行特殊处理
			err := getQuote(l,m)
			if err != nil {
				return err
			}
		}else {
			m[l] = ""
		}
	}
	return nil
}

// 检查字典是否存在引用
func hasQuote(s string) bool {
	return strings.Contains(s,"${")
}


func hasWildcard (s string) bool {
	return strings.Contains(s,".{")
}

// 获取引用数据
func getQuote(s string,m map[string]string) error {
	if !hasQuote(s) {
		return nil
	}
	file := util.GetBetweenStr(s,"${","}")
	cic := make(map[string]string)
	err := getDictionaryFromFile("dictionary/"+file+".txt",cic)
	if err != nil {
		return err
	}
	err = doQuote(s,file,m,cic)
	if err != nil {
		return err
	}
	return nil
}

// 处理引用
func doQuote(s string,file string,m map[string]string,c map[string]string) error {
	for ss := range c {
		e := strings.Replace(s,"${"+file+"}",ss,1)
		if hasQuote(e) {
			err := getQuote(e,m)
			if err != nil {
				return err
			}
		}else {
			m[e] = ""
		}
	}
	return nil
}