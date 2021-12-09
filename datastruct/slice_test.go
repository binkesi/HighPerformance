package datastruct

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func PrintLenCap(nums []int) {
	fmt.Printf("len: %d, cap: %d, %v\n", len(nums), cap(nums), nums)
}

// 因此，为了避免内存发生拷贝，如果能够知道最终的切片的大小，预先设置 cap 的值能够获得最好的性能。
func TestPrint(t *testing.T) {
	nums := []int{1}
	nums = append(nums, 2)
	PrintLenCap(nums)
	nums = append(nums, 3)
	PrintLenCap(nums)
	nums = append(nums, 4)
	PrintLenCap(nums)

	nums1 := make([]int, 0, 8)
	nums1 = append(nums1, 1, 2, 3, 4, 5)
	PrintLenCap(nums1)
	nums2 := nums1[2:4]
	nums2 = append(nums2, 50, 60)
	PrintLenCap(nums2)
	PrintLenCap(nums1)
	// 删除索引3位置的元素
	// 切片的底层是数组，因此删除意味着后面的元素需要逐个向前移位。每次删除的复杂度为 O(N)，因此切片不合适大量随机删除的场景，这种场景下适合使用链表
	nums1 = append(nums1[:3], nums1[4:]...)
	PrintLenCap(nums1)
}

// 在已有切片的基础上进行切片，不会创建新的底层数组。因为原来的底层数组没有发生变化，内存会一直占用，直到没有变量引用该数组
// 因此很可能出现这么一种情况，原切片由大量的元素构成，但是我们在原切片的基础上切片，虽然只使用了很小一段，但底层数组在内存中仍然占据了大量空间，得不到释放。
// 比较推荐的做法，使用 copy 替代 re-slice
func lastNumsBySlice(origin []int) []int {
	return origin[len(origin)-2:]
}

func lastNumsByCopy(origin []int) []int {
	res := make([]int, 2)
	copy(res, origin[len(origin)-2:])
	return res
}

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	res := make([]int, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, rand.Int())
	}
	return res
}

func printMem(t *testing.T) {
	t.Helper()
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	t.Logf("%.2f MB", float64(rtm.Alloc)/1024./1024.)
}

func testLastChars(t *testing.T, f func([]int) []int) {
	t.Helper()
	ans := make([][]int, 0)
	for k := 0; k < 100; k++ {
		origin := generateWithCap(128 * 1024) // 1M
		ans = append(ans, f(origin))
		runtime.GC()
	}
	printMem(t)
	_ = ans
}

func TestLastCharsBySlice(t *testing.T) { testLastChars(t, lastNumsBySlice) }
func TestLastCharsByCopy(t *testing.T)  { testLastChars(t, lastNumsByCopy) }
