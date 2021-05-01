package auth

import (
	"testing"
)

func TestMemoryTokenStoreSetGetVerify(t *testing.T) {
	token := "abcdefg"
	ts := NewMemoryTokenStore()
	ts.Set(token, "userA")

	// Token that exists.
	_, ok := ts.Get(token)
	if ok != true {
		t.Errorf("incorrect value, got: %v, want: %v", ok, true)
	}

	_, ok = ts.Get("dontexist")
	if ok != false {
		t.Errorf("incorrect value, got: %v, want: %v", ok, false)
	}

	ok = ts.Verify(token)
	if ok != true {
		t.Errorf("incorrect value, got: %v, want: %v", ok, true)
	}

	ok = ts.Verify("dontexist")
	if ok != false {
		t.Errorf("incorrect value, got: %v, want: %v", ok, false)
	}
}
