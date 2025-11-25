package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gyyyy/cornpuzzle-go/cornpuzzle"
)

func main() {
	var (
		x int
		y int
		b string
		v bool
	)
	flag.IntVar(&x, "x", 14, "玉米拼图长度")
	flag.IntVar(&y, "y", 7, "玉米拼图高度")
	flag.StringVar(&b, "b", "", "拼图块，使用'||'分隔各个拼图块")
	flag.BoolVar(&v, "v", false, "输出解题过程")
	flag.Parse()
	if x <= 0 || y <= 0 {
		log.Fatalln("无效参数 [x, y]")
	}
	if b = strings.TrimSpace(b); b == "" {
		log.Fatalln("无效参数 [b]")
	}
	cornpuzzle.Verbose = v
	pzl, err := cornpuzzle.Create(x, y, strings.Split(b, "||"))
	if err != nil {
		log.Fatalln(err)
	}
	var (
		startTime = time.Now()
		ok        = cornpuzzle.Resolve(pzl)
		duration  = time.Since(startTime)
	)
	if ok {
		fmt.Printf("谜题答案:\n%s\n", pzl.Corn)
	} else {
		fmt.Println("解谜失败")
	}
	fmt.Printf("求解耗时: %v\n", duration)
}
