package datastruct

import (
	"fmt"
	"testing"
	"unsafe"
)

type Args struct {
	a int
	b int
}

type Flag struct {
	c int32
	d int16
}

// Flag 由一个 int32 和 一个 int16 的字段构成，成员变量占据的字节数为 4+2 = 6，
// 但是 unsafe.Sizeof 返回的结果为 8 字节，多出来的 2 字节是内存对齐的结果。
// Alignof 方法，可以返回一个类型的对齐值，也可以叫做对齐系数或者对齐倍数
func TestMemAlloc(t *testing.T) {
	fmt.Println(unsafe.Sizeof(Args{}), unsafe.Alignof(Args{}))
	fmt.Println(unsafe.Sizeof(Flag{}), unsafe.Alignof(Flag{}))
}

// CPU 始终以字长访问内存，如果不进行内存对齐，很可能增加 CPU 访问内存的次数
// 合理的内存对齐可以提高内存读写的性能，并且便于实现变量操作的原子性。

type demo1 struct {
	a int8
	b int16
	c int32
}

type demo2 struct {
	a int8
	c int32
	b int16
}

// 每个字段按照自身的对齐倍数来确定在内存中的偏移量，字段排列顺序不同，上一个字段因偏移而浪费的大小也不同。
// 在对内存特别敏感的结构体的设计上，我们可以通过调整字段的顺序，减少内存的占用。
func TestStruct(t *testing.T) {
	fmt.Println(unsafe.Sizeof(demo1{})) // 8
	fmt.Println(unsafe.Sizeof(demo2{})) // 12
}

// 空 struct{} 大小为 0，作为其他 struct 的字段时，一般不需要内存对齐
// 当 struct{} 作为结构体最后一个字段时，需要内存对齐
type demo3 struct {
	c int32
	a struct{}
}

type demo4 struct {
	a struct{}
	c int32
}

func TestEmptyStruct(t *testing.T) {
	fmt.Println(unsafe.Sizeof(demo3{})) // 8
	fmt.Println(unsafe.Sizeof(demo4{})) // 4
}
