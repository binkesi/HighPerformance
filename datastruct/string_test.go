package datastruct

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

const letterString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	var b = make([]byte, n)
	for i := range b {
		b[i] = letterString[rand.Intn(len(letterString))]
	}
	return string(b)
}

func plusConcat(n int, str string) string {
	var s string = ""
	for i := 0; i < n; i++ {
		s += str
	}
	return s
}

func sprintfConcat(n int, str string) string {
	var s string = ""
	for i := 0; i < n; i++ {
		s = fmt.Sprintf("%s%s", s, str)
	}
	return s
}

func builderConcat(n int, str string) string {
	var s strings.Builder
	s.Grow(n * len(str))
	for i := 0; i < n; i++ {
		s.WriteString(str)
	}
	return s.String()
}

func bufferConcat(n int, str string) string {
	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		buf.WriteString(str)
	}
	return buf.String()
}

func byteConcat(n int, str string) string {
	s := make([]byte, 0)
	for i := 0; i < n; i++ {
		s = append(s, str...)
	}
	return string(s)
}

func preByteConcat(n int, str string) string {
	buf := make([]byte, 0, n*len(str))
	for i := 0; i < n; i++ {
		buf = append(buf, str...)
	}
	return string(buf)
}

func benchmark(b *testing.B, f func(int, string) string) {
	str := randString(10)
	for i := 0; i < b.N; i++ {
		f(10000, str)
	}
}

func BenchmarkPlusConcat(b *testing.B)    { benchmark(b, plusConcat) }
func BenchmarkSprintfConcat(b *testing.B) { benchmark(b, sprintfConcat) }
func BenchmarkBuilderConcat(b *testing.B) { benchmark(b, builderConcat) }
func BenchmarkBufferConcat(b *testing.B)  { benchmark(b, bufferConcat) }
func BenchmarkByteConcat(b *testing.B)    { benchmark(b, byteConcat) }
func BenchmarkPreByteConcat(b *testing.B) { benchmark(b, preByteConcat) }
