package main

import (
	"fmt"
	"sync"
)

type MyStruct struct {
	data string
}

var pool = sync.Pool{
	New: func() any {
		return &MyStruct{}
	},
}

func main() {
	for i := range 10 {
		item := pool.Get().(*MyStruct)
		item.data = fmt.Sprintf("Data %d", i)
		fmt.Println(item.data)
		pool.Put(item)
	}

}
