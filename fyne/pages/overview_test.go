package pages

import (
	"testing"
)

func TestBreakBalance(t *testing.T) {
	balance := "155.9999909 DCR"
	b1, b2 := breakBalance(balance)
	if b1 != "155.99" {
		t.Errorf("breakBalance 1 is wrong! got: %v want: %v", b1, "155.99")
	}
	if b2 != "99909 DCR" {
		t.Errorf("breakBalance 2 is wrong! got: %v want: %v", b2, "99909 DCR")
	}
	balanceTwo := "155 DCR"
	b3, b4 := breakBalance(balanceTwo)
	if b1 != "155.99" {
		t.Errorf("breakBalance 3 is wrong! got: %v want: %v", b3, "155.99")
	}
	if b2 != "99909 DCR" {
		t.Errorf("breakBalance 4 is wrong! got: %v want: %v", b4, "99909 DCR")
	}
}
