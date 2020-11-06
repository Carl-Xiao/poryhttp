package test

import (
	"math/rand"
	"proxyHTTP/handler"
	"strconv"
	"testing"
	"time"
)

func BenchmarkW1_Next(b *testing.B) {
	w := handler.WeightedRR()
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		w.Add("server"+strconv.Itoa(i), rand.Intn(100))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Next()
	}
}
