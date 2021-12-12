package datastruct

import (
	"fmt"
	"testing"
	"unsafe"
)

// 空结构体不占据内存空间，因此被广泛作为各种场景下的占位符使用。
// 一是节省资源，二是空结构体本身就具备很强的语义，即这里不需要任何值，仅作为占位符
// Go 语言标准库没有提供 Set 的实现，通常使用 map 来代替
// 将 map 作为集合(Set)使用时，可以将值类型定义为空结构体，仅作为占位符使用即可
func TestStructSize(t *testing.T) {
	fmt.Println(unsafe.Sizeof(struct{}{}))
}

type Set map[string]struct{}

func (s Set) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) del(key string) {
	delete(s, key)
}

func TestSet(t *testing.T) {
	s := make(Set)
	s.Add("Tom")
	s.Add("Sam")
	fmt.Println(s.Has("Tom"))
	s.del("Sam")
	fmt.Println(s.Has("Jack"))
	fmt.Println(s.Has("Sam"))
}

// 有时候使用 channel 不需要发送任何的数据，只用来通知子协程(goroutine)执行任务，或只用来控制协程并发度。
// 这种情况下，使用空结构体作为占位符就非常合适了

func worker(ch chan struct{}) {
	<-ch
	fmt.Println("recieve signal from producer")
	close(ch)
}

func TestWorker(t *testing.T) {
	ch := make(chan struct{})
	go worker(ch)
	ch <- struct{}{}
}

// 在部分场景下，结构体只包含方法，不包含任何的字段。例如例子中的 Door，
// 在这种情况下，Door 事实上可以用任何的数据结构替代
// 但无论是 int 还是 bool 都会浪费额外的内存，因此呢，这种情况下，声明为空结构体是最合适的
type Door struct{}

func (d Door) Open() {
	fmt.Println("Open the door")
}

func (d Door) Close() {
	fmt.Println("Close the door")
}

// 逃逸分析：
// 传值会拷贝整个对象，而传指针只会拷贝指针地址，指向的对象是同一个。
// 传指针可以减少值的拷贝，但是会导致内存分配逃逸到堆中，增加垃圾回收(GC)的负担。
// 在对象频繁创建和删除的场景下，传递指针导致的 GC 开销可能会严重影响性能。
// 一般情况下，对于需要修改原对象值，或占用内存比较大的结构体，选择传指针。对于只读的占用内存较小的结构体，直接传值能够获得更好的性能。

// 因此，在声明全局变量时，如果能够确定为常量，尽量使用 const 而非 var，这样很多运算在编译器即可执行。

// 我们可以在源代码中，定义全局常量 debug，值设置为 false，在需要增加调试代码的地方，使用条件语句 if debug 包裹
// 如果是正常编译，常量 debug 始终等于 false，调试语句在编译过程中会被消除，不会影响最终的二进制大小，也不会对运行效率产生任何影响
// 如果我们想编译出 debug 版本的二进制呢？可以将 debug 修改为 true 之后编译。这对于开发者日常调试是非常有帮助的
