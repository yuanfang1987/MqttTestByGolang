package mytest

import (
	"testing"
)

func Test_pbhaha(t *testing.T) {
	pl := pbMarshal()
	pbUnMarshal(pl)
}
