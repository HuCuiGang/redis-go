package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

func main()  {
	//配置并获得一个连接对象的指针
	pool := &redis.Pool{
		//最大闲置链接数
		MaxIdle: 20,
		//最大活动链接数，0=无限
		MaxActive: 0,
		//闲置链接的超时时间
		IdleTimeout: time.Second * 100,
		//定义拨号获得链接的函数
		Dial:func()(redis.Conn,error){
			return redis.Dial("tcp","127.0.0.1:6379")
		},
	}

	//延时关闭连接词
	defer pool.Close()

	//10并发链接
	for i := 0; i < 10; i++ {
		go getCounFromPoolAandHappy(pool,i)
	}

	//保持主协成存活
	time.Sleep(3 * time.Second)

}

func getCounFromPoolAandHappy(pool *redis.Pool,i int)  {
	//通过连接池获得链接
	conn := pool.Get()
	//延时关闭连接
	defer conn.Close()
	//使用链接操作数据
	do, err := conn.Do("auth", "root")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(do)
	reply,err := conn.Do("set","conn"+strconv.Itoa(i),i)

	s, _ := redis.String(reply,err)
	fmt.Println(s)
}
