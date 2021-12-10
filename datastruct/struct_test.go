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
