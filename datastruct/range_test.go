package datastruct

import (
	"fmt"
	"testing"
)

// 变量 words 在循环开始前，仅会计算一次，如果在循环中修改切片的长度不会改变本次循环的次数
// 如果删除元素的话会影响迭代输出结果
func TestRangeSlice(t *testing.T) {
	slice := []string{"go", "language", "coding", "practice"}
	fmt.Println(&slice[2])
	for i, v := range slice {
		if i == 0 {
			// 新创建的切片只是对原有切片的引用，所以查看元素地址可以发现，其实原来的slice切片元素变成了{"go", "language", "practice", "practice"}
			slice = append(slice[:len(slice)-2], slice[len(slice)-1:]...)
		}
		fmt.Printf("%d %v\n", i, v)
	}
	fmt.Println(slice)
	fmt.Println(&slice[2])
}

// 和切片不同的是，迭代过程中，删除还未迭代到的键值对，则该键值对不会被迭代。
// 在迭代过程中，如果创建新的键值对，那么新增键值对，可能被迭代，也可能不会被迭代。
func TestRangeMap(t *testing.T) {
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	for k, v := range m {
		delete(m, 2)
		m[4] = "four"
		fmt.Printf("%d: %s\n", k, v)
	}
}

// 发送给信道(channel) 的值可以使用 for 循环迭代，直到信道被关闭。
// 如果是 nil 信道，循环将永远阻塞
func TestChannel(t *testing.T) {
	ch := make(chan string)
	go func() {
		ch <- "Go"
		ch <- "语言"
		ch <- "高性能"
		ch <- "编程"
		close(ch)
	}()
	for n := range ch {
		fmt.Println(n)
	}
}

func BenchmarkForIntSlice(b *testing.B) {
	nums := generateWithCap(1024 * 1024)
	for i := 0; i < b.N; i++ {
		length := len(nums)
		var tmp int
		for k := 0; k < length; k++ {
			tmp = nums[k]
		}
		_ = tmp
	}
}

func BenchmarkRangeIntSlice(b *testing.B) {
	nums := generateWithCap(1024 * 1024)
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, num := range nums {
			tmp = num
		}
		_ = tmp
	}
}

// 与 for 不同的是，range 对每个迭代值都创建了一个拷贝。因此如果每次迭代的值内存占用很小的情况下，for 和 range 的性能几乎没有差异，
// 但是如果每个迭代值内存占用很大，例如上面的例子中，每个结构体需要占据 4KB 的内存，这种情况下差距就非常明显了
type Item struct {
	id  int
	val [4096]byte
}

func BenchmarkForStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		length := len(items)
		var tmp int
		for k := 0; k < length; k++ {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangeIndexStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		var tmp int
		for k := range items {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangeStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, item := range items {
			tmp = item.id
		}
		_ = tmp
	}
}

func TestValueCopyRange(t *testing.T) {
	persons := []struct{ no int }{{no: 1}, {no: 2}, {no: 3}}
	for _, v := range persons {
		v.no += 100
	}
	fmt.Println("use range to modify value:", persons)
	for i := 0; i < len(persons); i++ {
		persons[i].no += 100
	}
	fmt.Println("use for to modify value:", persons)
}

// 切片元素从结构体 Item 替换为指针 *Item 后，for 和 range 的性能几乎是一样的。而且使用指针还有另一个好处，可以直接修改指针对应的结构体的值。
func generateItems(n int) []*Item {
	items := make([]*Item, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, &Item{id: i})
	}
	return items
}

func BenchmarkForPointer(b *testing.B) {
	items := generateItems(1024)
	for i := 0; i < b.N; i++ {
		length := len(items)
		var tmp int
		for k := 0; k < length; k++ {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangePointer(b *testing.B) {
	items := generateItems(1024)
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, item := range items {
			tmp = item.id
		}
		_ = tmp
	}
}

// range 在迭代过程中返回的是迭代值的拷贝，如果每次迭代的元素的内存占用很低，那么 for 和 range 的性能几乎是一样，例如 []int。
// 但是如果迭代的元素内存占用较高，例如一个包含很多属性的 struct 结构体，那么 for 的性能将显著地高于 range，有时候甚至会有上千倍的性能差异。
