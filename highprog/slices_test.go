package highprog

import (
	"fmt"
	"testing"
)

func TrimSpace(s []byte) []byte {
	// 利用了空切片的特性，此处声明的b的len为0但是cap与s相同，避免了append过程中的扩容的内存分配
	b := s[:0]
	for _, v := range s {
		if v != ' ' {
			b = append(b, v)
		}
	}
	return b
}

func TestEmpSlice(t *testing.T) {
	var spac string = "hello world   ,  sungn!"
	fmt.Println(string(TrimSpace([]byte(spac))))
}
