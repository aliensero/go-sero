// copyright 2018 The sero.cash Authors
// This file is part of the go-sero library.
//
// The go-sero library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-sero library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-sero library. If not, see <http://www.gnu.org/licenses/>.

package stx

import (
	"encoding/hex"
	"sync/atomic"

	"github.com/sero-cash/go-czero-import/cpt"

	"github.com/sero-cash/go-sero/zero/utils"

	"github.com/sero-cash/go-sero/zero/txs/assets"

	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/crypto"
	"github.com/sero-cash/go-sero/crypto/sha3"
)

type Out_O struct {
	Addr  keys.PKr
	Asset assets.Asset
	Memo  keys.Uint512

	//Cache
	assetCC atomic.Value
}

func (self *Out_O) ToAssetCC() keys.Uint256 {
	if cc, ok := self.assetCC.Load().(keys.Uint256); ok {
		return cc
	}
	asset := self.Asset.ToFlatAsset()
	asset_desc := cpt.AssetDesc{
		Tkn_currency: asset.Tkn.Currency,
		Tkn_value:    asset.Tkn.Value.ToUint256(),
		Tkt_category: asset.Tkt.Category,
		Tkt_value:    asset.Tkt.Value,
	}
	cpt.GenAssetCC(&asset_desc)
	v := asset_desc.Asset_cc
	self.assetCC.Store(v)
	return v
}

func (self *Out_O) Clone() (ret Out_O) {
	utils.DeepCopy(&ret, self)
	return
}
func (this Out_O) ToRef() (ret *Out_O) {
	ret = &this
	return
}

func (self *Out_O) ToHash() (ret keys.Uint256) {
	hash := crypto.Keccak256(
		self.Addr[:],
		self.Asset.ToHash().NewRef()[:],
		self.Memo[:],
	)
	copy(ret[:], hash)
	return ret
}

func (self *Out_O) ToHash_for_gen() (ret keys.Uint256) {
	hash := crypto.Keccak256(
		self.Addr[:],
		self.Asset.ToHash().NewRef()[:],
		self.Memo[:],
	)
	copy(ret[:], hash)
	return ret
}

func (self *Out_O) ToHash_for_sign() (ret keys.Uint256) {
	hash := crypto.Keccak256(
		self.Addr[:],
		self.Asset.ToHash().NewRef()[:],
		self.Memo[:],
	)
	copy(ret[:], hash)
	return ret
}

type In_O struct {
	Root keys.Uint256
	Sign keys.Uint512
}

func (ino In_O) MarshalText() ([]byte, error) {
	input := [64]byte{}
	copy(input[:32], ino.Root[:])
	copy(input[32:], ino.Sign[:])
	result := make([]byte, len(input)*2+2)
	copy(result, `0x`)
	hex.Encode(result[2:], input[:])
	return result, nil
}

func (self *In_O) ToHash() (ret keys.Uint256) {
	copy(ret[:], self.Root[:])
	copy(ret[:], self.Sign[:])
	return ret
}

func (self *In_O) ToHash_for_gen() (ret keys.Uint256) {
	copy(ret[:], self.Root[:])
	return ret
}

func (self *In_O) ToHash_for_sign() (ret keys.Uint256) {
	copy(ret[:], self.Root[:])
	return ret
}

type In_S struct {
	Root keys.Uint256
	Nil  keys.Uint256
	Sign keys.Uint512
}

func (self *In_S) ToHash() (ret keys.Uint256) {
	hash := crypto.Keccak256(
		self.Root[:],
		self.Nil[:],
		self.Sign[:],
	)
	copy(ret[:], hash)
	return ret
}

func (self *In_S) ToHash_for_gen() (ret keys.Uint256) {
	hash := crypto.Keccak256(
		self.Root[:],
	)
	copy(ret[:], hash)
	return ret
}

func (self *In_S) ToHash_for_sign() (ret keys.Uint256) {
	hash := crypto.Keccak256(
		self.Root[:],
	)
	copy(ret[:], hash)
	return ret
}

type Desc_O struct {
	Ins  []In_S
	Outs []Out_O
}

func (self *Desc_O) ToHash() (ret keys.Uint256) {
	d := sha3.NewKeccak256()
	for _, in := range self.Ins {
		d.Write(in.ToHash().NewRef()[:])
	}
	for _, out := range self.Outs {
		d.Write(out.ToHash().NewRef()[:])
	}
	copy(ret[:], d.Sum(nil))
	return
}

func (self *Desc_O) ToHash_for_gen() (ret keys.Uint256) {
	d := sha3.NewKeccak256()
	for _, in := range self.Ins {
		d.Write(in.ToHash_for_gen().NewRef()[:])
	}
	for _, out := range self.Outs {
		d.Write(out.ToHash_for_gen().NewRef()[:])
	}
	copy(ret[:], d.Sum(nil))
	return
}

func (self *Desc_O) ToHash_for_sign() (ret keys.Uint256) {
	d := sha3.NewKeccak256()
	for _, in := range self.Ins {
		d.Write(in.ToHash_for_sign().NewRef()[:])
	}
	for _, out := range self.Outs {
		d.Write(out.ToHash_for_sign().NewRef()[:])
	}
	copy(ret[:], d.Sum(nil))
	return
}
