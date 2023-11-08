package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestKeyStore(t *testing.T) {
	InitKeyStore("./testdata")
	accts := Accounts()
	if len(accts) != 1 {
		t.Fatal("account files mismatch")
	}
	want := common.HexToAddress("7ef5a6135f1fd6a02593eedc869c6d41d934aef8")
	if accts[0].Address != want {
		t.Fatalf("got(%v), expect(%v)", accts[0].Address, want)
	}
}
