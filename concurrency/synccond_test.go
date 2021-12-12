package concurrency

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

// 一句话总结：sync.Cond 条件变量用来协调想要访问共享资源的那些 goroutine，当共享资源的状态发生变化的时候，它可以用来通知被互斥锁阻塞的 goroutine。
// sync.Cond 经常用在多个 goroutine 等待，一个 goroutine 通知（事件发生）的场景。如果是一个通知，一个等待，使用互斥锁或 channel 就能搞定了。

var done = false

func read(name string, c *sync.Cond) {
	c.L.Lock()
	for !done {
		c.Wait()
	}
	fmt.Println(name, "start to read")
	c.L.Unlock()
}

func write(name string, c *sync.Cond) {
	log.Println(name, "start to write")
	time.Sleep(time.Second)
	c.L.Lock()
	done = true
	c.L.Unlock()
	fmt.Println("wakes all")
	c.Broadcast()
}

// done 即互斥锁需要保护的条件变量。
// read() 调用 Wait() 等待通知，直到 done 为 true。
// write() 接收数据，接收完成后，将 done 置为 true，调用 Broadcast() 通知所有等待的协程。
// write() 中的暂停了 1s，一方面是模拟耗时，另一方面是确保前面的 3 个 read 协程都执行到 Wait()，处于等待状态。main 函数最后暂停了 3s，确保所有操作执行完毕。
func TestSyncCond(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	go read("reader1", cond)
	go read("reader2", cond)
	go read("reader3", cond)
	write("writer", cond)
	time.Sleep(time.Second * 3)
}
