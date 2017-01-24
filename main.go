package main

import (
	"github.com/glutwins/pholcus/exec"
	_ "github.com/henrylee2cn/pholcus_lib" // 此为公开维护的spider规则库
	// _ "github.com/henrylee2cn/pholcus_lib_pte" // 同样你也可以自由添加自己的规则库
)

func main() {
	exec.Run()
}
