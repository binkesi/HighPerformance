package concurrency

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func dobadthing(done chan bool) {
	time.Sleep(time.Second)
	done <- true
}

func timeout(f func(chan bool)) error {
	// 创建channel done 时，缓冲区设置为 1，即使没有接收方，发送方也不会发生阻塞
	// 使用 select 尝试向信道 done 发送信号，如果发送失败，则说明缺少接收者(receiver)，即超时了，那么直接退出即可。
	/*
			func doGoodthing(done chan bool) {
				time.Sleep(time.Second)
				select {
				case done <- true:
				default:
					return
			}
		}
	*/
	done := make(chan bool, 1)
	go f(done)
	select {
	case <-done:
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond):
		return fmt.Errorf("timeout")
	}
}

func TestTimeout(t *testing.T) {
	fmt.Println(timeout(dobadthing))
}

func test(t *testing.T, f func(chan bool)) {
	t.Helper()
	for i := 0; i < 1000; i++ {
		timeout(f)
	}
	time.Sleep(2 * time.Second)
	fmt.Println(runtime.NumGoroutine())
}

// 最终程序中存在着 1002 个子协程，说明即使是函数执行完成，协程也没有正常退出。
// 那如果在实际的业务中，我们使用了上述的代码，那越来越多的协程会残留在程序中，最终会导致内存耗尽
func TestBadTimeout(t *testing.T) { test(t, dobadthing) }

// 当超时发生时，select 接收到 time.After 的超时信号就返回了，done 没有了接收方(receiver)，
// 而 doBadthing 在执行 1s 后向 done 发送信号，由于没有接收者且无缓存区，发送者(sender)会一直阻塞，导致协程不能退出

func do2phases(phase1, done chan bool) {
	time.Sleep(time.Second)
	select {
	case phase1 <- true:
	default:
		return
	}
	time.Sleep(time.Second)
	done <- true
}

func timeoutFirstPhase() error {
	phase1 := make(chan bool)
	done := make(chan bool)
	go do2phases(phase1, done)
	select {
	case <-phase1:
		<-done
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond):
		return fmt.Errorf("timeout")
	}
}

// 我们将服务端接收请求后的任务拆分为 2 段，一段是执行任务，一段是发送结果。那么就会有两种情况：
// 1. 任务正常执行，向客户端返回执行结果。
// 2. 任务超时执行，向客户端返回超时。
// 这种情况下，就只能够使用 select，而不能能够设置缓冲区的方式了。因为如果给信道 phase1 设置了缓冲区，phase1 <- true 总能执行成功，
// 那么无论是否超时，都会执行到第二阶段，而没有即时返回，这是我们不愿意看到的。对应到上面的业务，就可能发生一种异常情况，向客户端发送了 2 次响应：
// 任务超时执行，向客户端返回超时，一段时间后，向客户端返回执行结果。
// 缓冲区不能够区分是否超时了，但是 select 可以
func Test2phasesTimeout(t *testing.T) {
	for i := 0; i < 1000; i++ {
		timeoutFirstPhase()
	}
	time.Sleep(time.Second * 3)
	fmt.Println(runtime.NumGoroutine())
}

// 因为 goroutine 不能被强制 kill，在超时或其他类似的场景下，为了 goroutine 尽可能正常退出，建议如下：
// 1. 尽量使用非阻塞 I/O（非阻塞 I/O 常用来实现高性能的网络库），阻塞 I/O 很可能导致 goroutine 在某个调用一直等待，而无法正确结束。
// 2. 业务逻辑总是考虑退出机制，避免死循环。
// 3. 任务分段执行，超时后即时退出，避免 goroutine 无用的执行过多，浪费资源。
