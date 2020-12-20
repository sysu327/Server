package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/sysu327/Server/dal/db"

	"github.com/sysu327/Server/dal/model"
)

func main() {
	articles := make([]model.Article, 3)
	users := make([]model.User, 3)
	titles := []string{"开发者点评GitHub 暗黑模式：太暗了", "从“卡脖子”到“主导”，国产数据库 40 年的演变！", "JavaScript 25 岁了！"}
	author := []string{"CSDN资讯", "CSDN资讯", "启舰"}
	times := []string{"2020-12-15 14:33:17", "2020-12-18 11:11:28", "2020-12-17 08:29:22"}
	for i := 0; i < 3; i++ {
		f, err := os.Open(strconv.FormatInt(int64(i), 10) + ".html")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		content, _ := ioutil.ReadAll(f)
		a1 := model.Article{int64(i + 1), titles[i], author[i], nil, times[i], string(content), nil}
		articles = append(articles, a1)
		u := model.User{author[i], "123"}
		users = append(users, u)
	}
	db.PutUsers(users)
	db.PutArticles(articles)
}
