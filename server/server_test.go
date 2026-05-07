package main

import (
	"testing"
)

func TestGiveEvent(t *testing.T) {
	t.Run("exit", func(t *testing.T) {
		got := give_event("exit", "message", "address")
		want := "user exitted!: address"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("send", func(t *testing.T) {
		got := give_event("send", "message", "address")
		want := "user address sent: message"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

}
