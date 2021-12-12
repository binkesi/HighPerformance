package concurrency

// 接收操作可以有 2 个返回值: v, beforeClosed := <-ch
// beforeClosed 代表 v 是否是信道关闭前发送的。true 代表是信道关闭前发送的，false 代表信道已经关闭。
// 如果一个信道已经关闭，<-ch 将永远不会发生阻塞，但是我们可以通过第二个返回值 beforeClosed 得知信道已经关闭，作出相应的处理。
// 下表为channel的操作和对应状态：
// 操作	    空值(nil)	非空已关闭	 非空未关闭
// 关闭	    panic	    panic	    成功关闭
// 发送数据	永久阻塞	 panic	     阻塞或成功发送
// 接收数据	永久阻塞	 永不阻塞	 阻塞或者成功接收

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// 一个通道被其发送数据协程队列和接收数据协程队列中的所有协程引用着。因此，如果一个通道的这两个队列只要有一个不为空，则此通道肯定不会被垃圾回收。
// 另一方面，如果一个协程处于一个通道的某个协程队列之中，则此协程也肯定不会被垃圾回收，即使此通道仅被此协程所引用。事实上，一个协程只有在退出后才能被垃圾回收

func do(taskCh chan int) {
	for {
		select {
		case t, beforeclosed := <-taskCh:
			if !beforeclosed {
				fmt.Println("channel has been closed")
				return
			}
			time.Sleep(time.Millisecond)
			fmt.Printf("task %d is done\n", t)
		}
	}
}

// 一个常用的使用Go通道的原则是不要在数据接收方或者在有多个发送者的情况下关闭通道。换句话说，我们只应该让一个通道唯一的发送者关闭此通道。

func sendTasks() {
	taskCh := make(chan int, 10)
	go do(taskCh)
	for i := 0; i < 100; i++ {
		taskCh <- i
	}
	close(taskCh)
}

func TestDo(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	sendTasks()
	time.Sleep(time.Second)
	fmt.Println(runtime.NumGoroutine())
}

// 使用 sync.Once 或互斥锁(sync.Mutex)确保 channel 只被关闭一次
type MyChannel struct {
	C    chan bool
	once sync.Once
}

func NewMyChannel() *MyChannel {
	return &MyChannel{C: make(chan bool)}
}

func (mc *MyChannel) SafeClose() {
	mc.once.Do(func() {
		close(mc.C)
	})
}

// 典型的关闭通道场景：
// 情形一：M个接收者和一个发送者，发送者通过关闭用来传输数据的通道来传递发送结束信号。
// 情形二：一个接收者和N个发送者，此唯一接收者通过关闭一个额外的信号通道来通知发送者不要再发送数据了。
// 情形三：M个接收者和N个发送者，它们中的任何协程都可以让一个中间调解协程帮忙发出停止数据传送的信号。
