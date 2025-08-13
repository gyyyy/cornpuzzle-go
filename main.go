package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gyyyy/cornpuzzle-go/cornpuzzle"
)

func main() {
	var (
		x int    // 14
		y int    // 7
		b string // "  #|###||  #|###|| ##|###|| ##|###||## |###||  #|# #|###||  #|# #|###||# #|# #|###||  #|  #|###||#  |## |###|| # |###| ##|| ##|###|## ||  #|###|## ||#  |#  |## |###|| # | # | # |###|| ##| # | # |## ||  #|  #|###| ##"
		v bool
	)
	flag.IntVar(&x, "x", 14, "length of corn")
	flag.IntVar(&y, "y", 7, "height of corn")
	flag.StringVar(&b, "b", "", "blocks")
	flag.BoolVar(&v, "v", false, "verbose")
	flag.Parse()
	if x <= 0 || y <= 0 {
		log.Fatalln("invalid arg [x, y]")
	}
	if b = strings.TrimSpace(b); b == "" {
		log.Fatalln("invalid arg [b]")
	}
	cornpuzzle.Verbose = v
	pzl, err := cornpuzzle.Create(x, y, strings.Split(b, "||"))
	if err != nil {
		log.Fatalln(err)
	}
	if cornpuzzle.Resolve(pzl) {
		fmt.Printf("解谜成功：\n%s\n", pzl.Corn)
	} else {
		fmt.Println("解谜失败")
	}
}
