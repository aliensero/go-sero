package light

import (
	"bytes"
	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/log"
	"github.com/sero-cash/go-sero/rlp"
	"github.com/sero-cash/go-sero/zero/txtool"
)

var current_light *LightNode

func (self *LightNode) CurrentLight() *LightNode {
	return current_light
}

func (self *LightNode) GetOutsByPKr(pkrs []keys.PKr, start, end uint64) (br BlockOutResp, e error) {
	br.CurrentNum = self.getLastNumber()
	blockOuts := []BlockOut{}
	for _, pkr := range pkrs {
		//uPKr := pkr.ToUint512()
		prefix := append(pkrPrefix, pkr[:]...)
		iterator := self.db.NewIteratorWithPrefix(prefix)

		for ok := iterator.Seek(pkrKey(pkr, start)); ok; ok = iterator.Next() {

			key := iterator.Key()
			num := bytesToUint64(key[99:107])
			if num > end {
				break
			}
			var outs []txtool.Out
			if err := rlp.Decode(bytes.NewReader(iterator.Value()), &outs); err != nil {
				log.Error("Light Invalid block RLP", "Num:", num, "err:", err)
				return br, err
			} else {
				blockOut := BlockOut{Num: num, Outs: outs}
				blockOuts = append(blockOuts, blockOut)
			}
		}
	}
	br.BlockOuts = blockOuts
	return br, nil
}

func (self *LightNode) CheckNil(Nils []keys.Uint256) (nilResps []NilValue, e error) {
	if len(Nils) == 0 {
		return
	}
	for _, Nil := range Nils {
		if data, err := self.db.Get(nilKey(Nil)); err != nil {
			continue
		} else {

			nilResp := NilValue{}
			if err:=rlp.DecodeBytes(data,&nilResp);err!=nil{
				continue
			}else{
				nilResp.Nil = Nil
				nilResps = append(nilResps, nilResp)
			}
		}
	}
	return nilResps, nil
}

type BlockOutResp struct {
	CurrentNum uint64
	BlockOuts  []BlockOut
}

type BlockOut struct {
	Num  uint64
	Outs []txtool.Out
}

