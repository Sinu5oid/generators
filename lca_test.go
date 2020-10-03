package main

import "testing"

func runTest(t *testing.T, modulus, multiplier, additiveComponent, x0, expectedPeriod int) {
	done := make(chan bool, 1)
	outChan := GenerateRPN(modulus, multiplier, additiveComponent, x0, done)

	t.Logf("x(0) = %d", x0)

	i := 1
	for n := range outChan {
		if n == x0 {
			done <- true
			t.Logf("x(%d) = %d", i, n)
			break
		}
		i += 1
	}

	if i != expectedPeriod {
		t.Errorf("Period got %d, want %d", i, expectedPeriod)
	} else {
		t.Logf("Period got %d, want %d", i, expectedPeriod)
	}
}

func TestLongPeriodRPN(t *testing.T) {
	runTest(t, 6075, 106, 1283, 7, 6075)
}

func TestShortPeriodRPN(t *testing.T) {
	runTest(t, 10, 7,7, 7, 4)
}

//func TestVeryLongPeriodRPN(t *testing.T) {
//	runTest(t, 995300, 199061, 11, 15, 995300)
//}