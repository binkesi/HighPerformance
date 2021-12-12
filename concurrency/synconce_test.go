package concurrency

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 在多数情况下，sync.Once 被用于控制变量的初始化，这个变量的读写满足如下三个条件：
// 当且仅当第一次访问某个变量时，进行初始化（写）；
// 变量初始化过程中，所有读都被阻塞，直到初始化完成；
// 变量仅初始化一次，初始化完成后驻留在内存里。
// sync.Once 仅提供了一个方法 Do，参数 f 是对象初始化函数。

type Config struct {
	Server string
	Port   int64
}

var (
	once   sync.Once
	config *Config
)

func ReadConfig() *Config {
	once.Do(func() {
		var err error
		config = &Config{Server: os.Getenv("TT_SERVER_URL")}
		config.Port, err = strconv.ParseInt(os.Getenv("TT_PORT"), 10, 0)
		if err != nil {
			config.Port = 8080
		}
		fmt.Println("init config")
	})
	return config
}

func TestReadConfig(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			_ = ReadConfig()
		}()
	}
	time.Sleep(time.Second)
}

// sync.Once实现源码, 代码位于 $(dirname $(which go))/../src/sync/once.go
// 首先：保证变量仅被初始化一次，需要有个标志来判断变量是否已初始化过，若没有则需要初始化。
// 第二：线程安全，支持并发，无疑需要互斥锁来实现。
type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

// 为什么把done放在第一个字段
// 热路径(hot path)是程序非常频繁执行的一系列指令，sync.Once 绝大部分场景都会访问 o.done，在热路径上是比较好理解的，如果 hot path 编译后的机器码指令更少，更直接，必然是能够提升性能的。
// 为什么放在第一个字段就能够减少指令呢？因为结构体第一个字段的地址和结构体的指针是相同的，如果是第一个字段，直接对结构体的指针解引用即可。
// 如果是其他的字段，除了结构体指针外，还需要计算与第一个值的偏移(calculate offset)。在机器码中，偏移量是随指令传递的附加值，CPU 需要做一次偏移值与指针的加法运算，才能获取要访问的值的地址。
// 因为，访问第一个字段的机器代码更紧凑，速度更快
