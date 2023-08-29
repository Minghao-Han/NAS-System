package main

import "time"

func main() {
	time.Sleep(5 * time.Second) //应该保证定时任务都远小于5s
}
