package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func main1(){
	//连接Redis数据库
	conn, err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
		fmt.Println("连接redis服务失败：",err)
		return
	}
	//延迟关闭连接
	defer conn.Close()
	//执行redis命令，获取结果
	reply,err := conn.Do("auth","root")
	if err != nil {
		fmt.Println("获取结果失败：",err)
		return
	}
	fmt.Println(reply)
	reply2 ,err := conn.Do("get","name")

	//结果原始类型是[]byte
	fmt.Printf("type=%T,value=%v\n",reply2,reply2)
	//根据具体的业务类型进行进行数据类型转换
	ret,_:=redis.String(reply2,err)
	//ret,_:= redis.Int(reply,err)
	fmt.Println(ret)

}