package main

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"io/ioutil"
	"regexp"
	"runtime"
	"statistician/conf"
	"statistician/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

var dictionary map[string]string
var rDictionary map[string]string
var articles []conf.Article
var config conf.Conf

func main() {
	start()
}

func start() {
	// 初始化数据，如果成功则开始处理数据
	if _init() {
		echo()
	}else {
		// 初始化出现错误，一般控制台同时会有错误日志输出
		fmt.Println("初始化错误，请检查配置文件")
	}
	pause()
	start()
}

func _init() bool {
	// 初始化变量
	dictionary = make(map[string]string)
	articles = []conf.Article{}

	// 检查文件
	ok,err := util.CheckOrCreateDir("dictionary")
	util.HandleError(err)
	if !ok {
		return false
	}
	ok,err = util.CheckOrCreateDir("echo")
	util.HandleError(err)
	if !ok {
		return false
	}
	ok,err = util.CheckOrCreateDir("articles")
	util.HandleError(err)
	if !ok {
		return false
	}
	ok,err = util.CheckOrCreateFile("app.conf")
	util.HandleError(err)
	if !ok {
		return false
	}
	// 加载配置
	fmt.Println("loading Config...")
	c,err := conf.InitConf("app.conf")
	util.HandleError(err)
	config = *c
	// 加载文件
	fmt.Println("loading File...")
	// 字典
	d,r,err := conf.LoadDictionary(c.EnableDictionary,c.BlackListDictionary)
	util.HandleError(err)
	dictionary = d
	rDictionary = r
	if len(dictionary ) == 0 && len(rDictionary) == 0 {
		fmt.Println("错误：字典为空，请添加字典文件")
		return false
	}
	// 文章
	a,err := conf.LoadArticles("articles")
	util.HandleError(err)
	articles = *a
	if len(articles) == 0 {
		fmt.Println("错误：文章为空，请添加文章文件")
		return false
	}
	return true
}

func echo() {
	start := time.Now().Unix()
	for _, art := range articles {
		 wg := &(sync.WaitGroup{})
		var task = 1
		//获取全文文字量
		l := strings.Count(art.Content,"") - 1

		// 根据配置文件和文字量开启多线程
		var multiThreadedMode = false
		if config.MultiThreadedMode == "a" {
			if  l > 1500000 {
				multiThreadedMode = true
				fmt.Println("多线程模式开启")
			}else {
				fmt.Println("多线程模式关闭")
			}
		}else if config.MultiThreadedMode == "y" {
			multiThreadedMode = true
			fmt.Println("多线程模式开启")
		}else if config.MultiThreadedMode == "n" {
			multiThreadedMode = false
			fmt.Println("多线程模式关闭")
		}

		// 统计数据，列map
		var args = cmap.New()
		for s := range dictionary {
			runtime.GOMAXPROCS(task)
			wg.Add(1)
			if multiThreadedMode {
				task ++
			}
			s := s
			go func() {
				count := strings.Count(art.Content,s)
				// 忽略掉没有出现的词
				if count != 0 {
					args.Set(s,count)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		// 统计通配数据
		rMap := make(map[string]int)
		for rd := range rDictionary {
			r,err := regexp.Compile(rd)
			util.HandleError(err)
			s := r.FindAllString(art.Content,-1)
			rMap[rd] = len(s)
		}

		m := make(map[string]int)
		for _, s := range args.Keys() {
			i,ok := args.Get(s)
			if ok {
				m[s] = i.(int)
			}
		}

		// 初始化数据
		var echo string
		var lastCount = 0
		var index = 0
		echo += "文件名：" + art.Title +"\n"
		echo += "文字数量：" + strconv.Itoa(l) +"\n"

		// 通配数据
		echo += "----通配数据------\n"
		data := util.SortMapByValue(rMap)
		for _, datum := range data {
			count := datum.Value
			if lastCount != count {
				index ++
			}
			lastCount = count
			// 计算百分比
			f := float64(count)/float64(l) * 100
			// 记录输出
			echo += strconv.Itoa(index) + "：" + datum.Key + "(" + strconv.Itoa(datum.Value) + ") " + fmt.Sprintf("%.2f", f) + "%\n"
		}

		echo += "----字典数据------\n"
		index = 0
		data = util.SortMapByValue(m)
		for _, datum := range data {
			count := datum.Value
			if lastCount != count {
				index ++
			}
			lastCount = count
			// 计算百分比
			f := float64(count)/float64(l) * 100
			// 记录输出
			echo += strconv.Itoa(index) + "：" + datum.Key + "(" + strconv.Itoa(datum.Value) + ") " + fmt.Sprintf("%.2f", f) + "%\n"
		}

		err := ioutil.WriteFile("echo/"+art.Title+"-"+strconv.Itoa(int(time.Now().Unix()))+".txt",[]byte(echo),0666)
		if err != nil {
			fmt.Println("写出："+art.Title+"数据失败")
		}else {
			fmt.Println("Echo："+art.Title)
		}

	}
	end := time.Now().Unix()
	fmt.Println("end At " + strconv.Itoa(int(start))+" use："+strconv.Itoa(int(end-start)) + "s")
}

func pause() {
	fmt.Println("-按 enter 重载，按 ctrl+c 退出-")
	var key string
	_,_ = fmt.Scanln(&key)
	if key == "\n" {
		return
	}
}
