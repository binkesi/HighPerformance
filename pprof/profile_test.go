package pprof

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/pkg/profile"
)

func randString(n int) string {
	const letterString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterString[rand.Intn(len(letterString))]
	}
	return string(b)
}

func concat(n int) string {
	a := ""
	for i := 0; i < n; i++ {
		a += randString(i)
	}
	return a
}

func builderConcat(n int) string {
	var s strings.Builder
	for i := 0; i < n; i++ {
		s.WriteString(randString(i))
	}
	return s.String()
}

func TestProfile(t *testing.T) {
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	concat(100)
	builderConcat(100)
}
