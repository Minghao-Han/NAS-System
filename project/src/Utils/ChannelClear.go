package Utils

func ClearChannel(ch chan int64) {
	select {
	case <-ch:
		// 清空 channel
		for len(ch) > 0 {
			<-ch
		}
		return
	default:
		return
	}
}
