package concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRoutineNum(t *testing.T) {
	ch := make(chan struct{}, 3)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
			time.Sleep(time.Second)
			<-ch
		}(i)
	}
	wg.Wait()
}

// 虚拟内存是一项非常常见的技术了，即在内存不足时，将磁盘映射为内存使用，比如 linux 下的交换分区(swap space)。
/* 在 linux 上创建并使用交换分区是一件非常简单的事情：
1 sudo fallocate -l 20G /mnt/.swapfile # 创建 20G 空文件
2 sudo mkswap /mnt/.swapfile    # 转换为交换分区文件
3 sudo chmod 600 /mnt/.swapfile # 修改权限为 600
4 sudo swapon /mnt/.swapfile    # 激活交换分区
5 free -m # 查看当前内存使用情况(包括交换分区)

关闭交换分区也非常简单：
1 sudo swapoff /mnt/.swapfile
2 rm -rf /mnt/.swapfile

磁盘的 I/O 读写性能和内存条相差是非常大的，例如 DDR3 的内存条读写速率很容易达到 20GB/s，但是 SSD 固态硬盘的读写性能通常只能达到 0.5GB/s，相差 40倍之多。
因此，使用虚拟内存技术将硬盘映射为内存使用，显然会对性能产生一定的影响。如果应用程序只是在较短的时间内需要较大的内存，那么虚拟内存能够有效避免 out of memory 的问题。
如果应用程序长期高频度读写大量内存，那么虚拟内存对性能的影响就比较明显了
*/
