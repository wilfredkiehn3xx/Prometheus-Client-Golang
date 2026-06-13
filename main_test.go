package main

import (
	"testing"
	"time"
)

func TestServerRestartDoesNotPanic(t *testing.T) {
	for i := 0; i < 3; i++ {
		srv := NewServer()
		srv.Addr = "127.0.0.1:0"
		go func() {
			_ = srv.Start()
		}()
		time.Sleep(50 * time.Millisecond)
		_ = srv.Stop()
	}
}
