package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Human struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

var Conn redis.Conn

func init() {
	var conn redis.Conn
	conn, _ = redis.Dial("tcp", "127.0.0.1:6379")
	conn.Do("auth", "root")
	Conn = conn

}

func main() {
	var cmd string

	for {
		fmt.Println("请输入命令：")
		fmt.Scanln(&cmd)
		switch cmd {
		case "getall":
			//显示所有人员信息
			GetAllPepole()
		case "exit":
			goto GAMEOVER
		default:
			fmt.Println("设么破命令，fuckoff！")
		}
	}
GAMEOVER:

	fmt.Println("GAME OVER")
}

func GetAllPepole() {
	//先尝试拿缓存
	peoplestrs := GetPeopleFromRedis()
	fmt.Println("拿到缓存数据：", peoplestrs, len(peoplestrs))
	//如果没有拿到数据
	if peoplestrs == nil || len(peoplestrs) == 0 {
		GetPeopleFromMysql()
	}

}

func GetPeopleFromMysql() {
	db, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/china")
	defer db.Close()
	if err != nil {
		return
	}
	var people []Human
	err = db.Select(&people, "select name ,age from person")
	if err != nil {
		fmt.Println("查询失败，", err)
		return
	}
	fmt.Println(people)
	//缓存查询结果到Redis
	CachePeople2Redis(people)
	return
}

//从redis拿信息
func GetPeopleFromRedis() (peopleStrs []string) {

	//Conn.Do("auth", "root")
	reply, err := Conn.Do("lrange", "people", "0", "-1")
	if err != nil {
		return
	}

	peopleStrs, err = redis.Strings(reply, err)
	fmt.Println("从缓存拿取结果", peopleStrs, err)

	return
}

func CachePeople2Redis(people []Human) {

	//先清除原有缓存
	Conn.Do("del","people")

	for _, human := range people {
		humanStr := fmt.Sprint(human)
		//Conn.Do("auth", "root")
		_, err := Conn.Do("rpush", "people", humanStr)
		if err != nil {
			fmt.Println("缓存失败，", err)
			return
		}
	}
	_, err := Conn.Do("expire", "people", 60)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("缓存people成功")
}
