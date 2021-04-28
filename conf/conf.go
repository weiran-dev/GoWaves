package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
)

type Conf struct {
	EnableDictionary *[]string
	BlackListDictionary *[]string
	MultiThreadedMode string
}

func InitConf(path string)(*Conf,error) {
	conf,err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	},path)
	if err != nil {
		return nil, err
	}

	// app config
	app := conf.Section("")
	if !app.HasKey("Multi-threaded-Mode") {
		_,err :=app.NewKey("Multi-threaded-Mode","a")
		if err != nil {
			return nil,err
		}
		err = conf.SaveTo("app.conf")
		if err != nil {
			return nil,err
		}
	}
	key,err := app.GetKey("Multi-threaded-Mode")
	if err != nil {
		return nil,err
	}
	multiThreadedMode := key.Value()

	// dictionary config
	dictionary := conf.Section("dictionary")
	dic,bl,err := makeDictionary(conf,dictionary)
	if err != nil {
		return nil, err
	}
	return &Conf{
		EnableDictionary: dic,
		BlackListDictionary: bl,
		MultiThreadedMode: multiThreadedMode,
	},nil
}

func makeDictionary(conf *ini.File,sec *ini.Section) (*[]string,*[]string,error) {
	var enableDictionary []string
	var blackListDictionary []string
	// 将所有扫描到的字典文件写入配置文件,并默认禁用
	files,err := ioutil.ReadDir("dictionary")
	if err != nil {
		return nil,nil,err
	}
	fileMap := make(map[string]int)
	for i, file := range files {
		fileMap[file.Name()] = i
		// 跳过已存在的配置
		if !sec.HasKey(file.Name()) {
			_,err := sec.NewKey(file.Name(),"y")
			fmt.Println("已启用新添加的字典:"+file.Name())
			if err != nil {
				return nil,nil,err
			}
		}
	}
	// 移除已删除的字典文件的配置
	for _, key := range sec.Keys() {
		if  _,ok := fileMap[key.Name()]; !ok {
			sec.DeleteKey(key.Name())
			fmt.Println("已移除失效的字典:"+key.Name())
		}else {
			if key.Value() == "y" {
				enableDictionary = append(enableDictionary, key.Name())
			}
			if key.Value() == "b" {
				blackListDictionary = append(blackListDictionary,key.Name())
			}
		}
	}
	err = conf.SaveTo("app.conf")
	if err != nil {
		return nil,nil,err
	}
	return &enableDictionary,&blackListDictionary,nil
}