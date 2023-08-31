package Utils

import "time"

func MakeTimeout(timeoutChan chan bool, t int, unit time.Duration) { //配合select-case使用
	sleepTime := time.Duration(t) * unit
	time.Sleep(sleepTime)
	timeoutChan <- true
}
