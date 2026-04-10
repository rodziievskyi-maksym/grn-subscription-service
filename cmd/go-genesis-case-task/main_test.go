package main

import "testing"

func TestAlwaysFails(t *testing.T) {
	expected := "success"
	actual := "failure"

	if actual != expected {
		t.Errorf("Test failed: expected %s, but got %s", expected, actual)
	}
}

func TestAddition(t *testing.T) {
	// Навмисно неправильна математика
	sum := 2 + 2
	if sum != 5 {
		t.Errorf("Math is broken: 2+2 should be 4, but we want it to fail!")
	}
}
