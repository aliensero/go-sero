package txstate

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/sero-cash/go-sero/zero/consensus"

	"github.com/sero-cash/go-czero-import/cpt"
	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/crypto"
	"github.com/sero-cash/go-sero/serodb"
)

type TreeState struct {
	db *consensus.FakeTri
}

func (self *TreeState) TryGet(key []byte) ([]byte, error) {
	return nil, nil
}
func (self *TreeState) TryUpdate(key, value []byte) error {
	return nil
}

func (self *TreeState) SetState(key *keys.Uint256, value *keys.Uint256) {
	self.db.TryUpdate(key[:], value[:])
}
func (self *TreeState) GetState(key *keys.Uint256) (ret keys.Uint256) {
	r, e := self.db.TryGet(key[:])
	if e == nil {
		copy(ret[:], r)
	}
	return
}
func (self *TreeState) GlobalGetter() serodb.Getter {
	return nil
}

func TestOutTree(t *testing.T) {
	// Create an empty state database
	cpt.ZeroInit("", 0)

	ft := consensus.NewFakeTri()
	outState := NewMerkleTree(&TreeState{db: &ft})

	for i := 1; i <= 100; i++ {
		value := crypto.Keccak256Hash(big.NewInt(int64(i)).Bytes()).HashToUint256()
		outState.AppendLeaf(*value)

		/*for i := 1; i <= 15; i++ {
			key := indexPathKey(uint64(i), uint64(0))
			value := outState.db.GetState(&key)
			fmt.Println(i, ":", common.Bytes2Hex(value[:]))
		}*/

		/*if i == 3 {
			current := crypto.Keccak256Hash(big.NewInt(int64(1)).Bytes()).HashToUint256()
			index, getPaths, anchor := outState.GetPaths(*current)
			ret := CalcRoot(current, index, &getPaths)
			if anchor != ret {
				fmt.Println(i, 1)
				t.FailNow()
			}
		}*/

		for j := 1; j <= i; j++ {
			current := crypto.Keccak256Hash(big.NewInt(int64(j)).Bytes()).HashToUint256()
			index, getPaths, anchor := outState.GetPaths(*current)
			ret := CalcRoot(current, index, &getPaths)
			if anchor != ret {
				fmt.Println(i, j)
				t.FailNow()
			}
		}
	}
}
