package concurrency

// 一句话总结：保存和复用临时对象，减少内存分配，降低 GC 压力。
// json 的反序列化在文本解析和网络通信过程中非常常见，当程序并发度非常高的情况下，短时间内需要创建大量的临时对象。
// 而这些对象是都是分配在堆上的，会给 GC 造成很大压力，严重影响程序的性能。
import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
)

// Go 语言从 1.3 版本开始提供了对象重用的机制，即 sync.Pool。sync.Pool 是可伸缩的，同时也是并发安全的，其大小仅受限于内存的大小。
//sync.Pool 用于存储那些被分配了但是没有被使用，而未来可能会使用的值。这样就可以不用再次经过内存分配，可直接复用已有对象，减轻 GC 的压力，从而提升系统的性能。
// sync.Pool 的大小是可伸缩的，高负载时会动态扩容，存放在池中的对象如果不活跃了会被自动清理。

type Student struct {
	Name   string
	Age    int32
	Remark [1024]byte
}

var buf, _ = json.Marshal(Student{Name: "sungn", Age: 24})

// 因为 Student 结构体内存占用较小，内存分配几乎不耗时间。而标准库 json 反序列化时利用了反射，效率是比较低的，占据了大部分时间，因此两种方式最终的执行时间几乎没什么变化。
// 但是内存占用差了一个数量级，使用了 sync.Pool 后，内存占用仅为未使用的 234/5096 = 1/22，对 GC 的影响就很大了
var studentPool = sync.Pool{
	New: func() interface{} {
		return new(Student)
	},
}

func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stu := &Student{}
		json.Unmarshal(buf, stu)
	}
}

func BenchmarkUnmarshalWithpool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stu := studentPool.Get().(*Student)
		json.Unmarshal(buf, stu)
		studentPool.Put(stu)
	}
}

// 在Go语言中五个引用类型变量,其他都是值类型: slice, map, channel, interface, func()
// 由于结构体是值类型,在方法传递时希望传递结构体地址,可以使用时结构体指针完成
// 可以结合new(T)函数创建结构体指针
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

var data = make([]byte, 10000)

func BenchmarkBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		buf.Write(data)
	}
}

func BenchmarkBufferWithpool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Write(data)
		buf.Reset()
		bufferPool.Put(buf)
	}
}

// 参考fmt.Printf的源码
// fmt.Printf 的调用是非常频繁的，利用 sync.Pool 复用 pp 对象能够极大地提升性能，减少内存占用，同时降低 GC 压力
