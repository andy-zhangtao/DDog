// 维护Metadata的channel
package bridge

var metaChan = make(chan int)

func GetMetaChan() chan int {
	return metaChan
}
