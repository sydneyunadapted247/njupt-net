package main

import "testing"

func TestParseFormPairs_OK(t *testing.T) {
	form, err := parseFormPairs([]string{"a=1", "b=2", " spaced = value with space "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if form["a"] != "1" {
		t.Fatalf("expected a=1, got %q", form["a"])
	}
	if form["b"] != "2" {
		t.Fatalf("expected b=2, got %q", form["b"])
	}
	if form["spaced"] != " value with space" {
		t.Fatalf("expected spaced form value preserved, got %q", form["spaced"])
	}
}

func TestParseFormPairs_InvalidNoEqual(t *testing.T) {
	_, err := parseFormPairs([]string{"invalid"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseFormPairs_EmptyKey(t *testing.T) {
	_, err := parseFormPairs([]string{"=value"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
