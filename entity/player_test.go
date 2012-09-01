package entity

import (
	"testing"
)

func Test_PassHashing(t *testing.T) {
	hash := hashPassword("arsenal")
	lhash := hashPassword("arsenalarsenalarsenal")
	//t.Log(hash)
	if hash != lhash {
		t.Fail()
	}
}
