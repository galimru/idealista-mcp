package main

import "testing"

func TestNewServerDoesNotRequireConfiguredRuntime(t *testing.T) {
	s := newServer()
	if s == nil {
		t.Fatal("expected server to be created")
	}
}
