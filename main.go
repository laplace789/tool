package main

import (
	"fmt"
	"sync"
	"time"
	"tool/rdb"
)

func main() {
	redisCfg := rdb.InitRedisConf("./conf")
	client := rdb.NewRedisClient(redisCfg)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		lock := rdb.NewRedisLock("7414", 10000000)
		wg.Add(1)
		go getKey(client, lock, i, &wg)
	}
	wg.Wait()
}

func getKey(cli rdb.RedisClient, key *rdb.RedisLock, index int, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%dth goroutinue[%d] try get lock\n", i, index)
		resp, err := cli.GetLock(key)
		if err != nil {
			fmt.Printf("%dth goroutinue[%d] try get lock fail err = %v\n", i, index, err)
			continue
		} else if resp == false {
			fmt.Printf("%dth goroutinue[%d] try get lock fail resp = %v\n", i, index, resp)
			time.Sleep(1 * time.Second)
			continue
		} else {
			fmt.Printf("%dth goroutinue[%d] try get lock sucess resp = %v\n", i, index, resp)
			fmt.Printf("%dth goroutinue[%v] try release lock \n", i, index)
			respRel, errRel := cli.ReleaseLock(key)
			if errRel != nil {
				fmt.Printf("%dth goroutinue[%d] try Release lock err = %v\n", i, index, errRel)
			}

			if resp {
				fmt.Printf("%dth goroutinue[%d] release lock  sucess resp = %v\n", i, index, respRel)
				fmt.Printf("i'm here mother fucker1 %v\n", index)
				wg.Done()
				return
			}
		}
	}
	fmt.Printf("i'm here mother fucke2r %v\n", index)
	wg.Done()
}
