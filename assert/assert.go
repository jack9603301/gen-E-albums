package assert

import (
	"fmt"
	"runtime"
)

func Assert(isTrue bool, msg string) {
	if !isTrue {
		var file string
		var line int
		var ok bool
		_, file, line, ok = runtime.Caller(1)
		if !ok {
			fmt.Println("运行时错误：断言函数堆栈捕获失败")
			panic(msg)
		} else {
			errmsg := fmt.Sprintf("%s:%d: 断言失败: %s\n", file, line, msg)
			panic(errmsg)
		}
	}
}
