package main

import (
	"github.com/steveoc64/memdebug"
	"time"
)

func main() {
	memdebug.Profile()

	for {
		t1 := time.Now()
		thing := make(map[string]int, 1024)
		for i:=0; i < 10000; i++ {
			for j:= 0; j < 10000; j++ {
				_ = i*j

			}
		}
		time.Sleep(time.Second)
		memdebug.Print(t1, "thing", thing)
		memdebug.DumpProfile()
	}
}
