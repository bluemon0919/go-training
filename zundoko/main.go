package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	var zzzzd string
	for {
		if strings.Contains(zzzzd, "ズンズンズンズンドコ") {
			break
		}
		zd := zundoko()
		fmt.Println(zd)
		zzzzd += zd
		time.Sleep(time.Millisecond * 300)
	}
	fmt.Println("キヨシ！")
}

// zundoko randomly returns "ズン" or "ドコ"
func zundoko() string {
	res := "ズン"
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100)%2 == 0 {
		res = "ドコ"
	}
	return res
}
