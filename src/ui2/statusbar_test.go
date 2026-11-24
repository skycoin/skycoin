package ui2

import "testing"

func TestLoadBarCSRFNoticeNeedsWalletAPI(t *testing.T) {
	if b := LoadBar(&fakeNode{browser: true}); b.NoCSRF {
		t.Fatal("CSRF notice shown for a node with the wallet API off")
	}
	if b := LoadBar(&fakeNode{}); !b.NoCSRF {
		t.Fatal("CSRF notice missing for a node with the wallet API on and CSRF off")
	}
}
