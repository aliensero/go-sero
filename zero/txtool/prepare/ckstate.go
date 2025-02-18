package prepare

import (
	"fmt"

	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/zero/txs/assets"
	"github.com/sero-cash/go-sero/zero/utils"
)

type cyState struct {
	balance utils.I256
}

type cyStateMap map[keys.Uint256]*cyState

func (self cyStateMap) add(key *keys.Uint256, value *utils.U256) {
	if _, ok := self[*key]; ok {
		self[*key].balance.AddU(value)
	} else {
		self[*key] = &cyState{
			*value.ToI256().ToRef(),
		}
	}
}

func (self cyStateMap) sub(key *keys.Uint256, value *utils.U256) {
	if _, ok := self[*key]; ok {
		self[*key].balance.SubU(value)
	} else {
		self[*key] = &cyState{}
		self[*key].balance.SubU(value)
	}
}

func newcyStateMap() (ret cyStateMap) {
	ret = make(map[keys.Uint256]*cyState)
	return
}

type CKState struct {
	outPlus bool
	cy      cyStateMap
	tk      map[keys.Uint256]keys.Uint256
}

func (self *CKState) GetList() (tkns []assets.Token, tkts []assets.Ticket) {
	for c, v := range self.cy {
		tkns = append(tkns, assets.Token{c, utils.U256(v.balance)})
	}
	for c, v := range self.tk {
		tkts = append(tkts, assets.Ticket{v, c})
	}
	return
}

func NewCKState(outPlus bool, fee *assets.Token) (ret CKState) {
	ret.outPlus = outPlus
	ret.cy = newcyStateMap()
	if outPlus {
		ret.cy.add(&fee.Currency, &fee.Value)
	} else {
		ret.cy.sub(&fee.Currency, &fee.Value)
	}
	ret.tk = make(map[keys.Uint256]keys.Uint256)
	return
}

func (self *CKState) AddIn(asset *assets.Asset) (added bool, e error) {
	added = false
	if asset.Tkn != nil {
		if asset.Tkn.Currency != keys.Empty_Uint256 {
			if asset.Tkn.Value.ToUint256() != keys.Empty_Uint256 {
				if self.outPlus {
					self.cy.sub(&asset.Tkn.Currency, &asset.Tkn.Value)
				} else {
					self.cy.add(&asset.Tkn.Currency, &asset.Tkn.Value)
				}
				added = true
			}
		}
	}
	if asset.Tkt != nil {
		if asset.Tkt.Category != keys.Empty_Uint256 {
			if asset.Tkt.Value != keys.Empty_Uint256 {
				if _, ok := self.tk[asset.Tkt.Value]; !ok {
					if self.outPlus {
						added = true
						delete(self.tk, asset.Tkt.Value)
					} else {
						e = fmt.Errorf("in tkt duplicate: %v", asset.Tkt.Value)
					}
					return
				} else {
					if !self.outPlus {
						added = true
						self.tk[asset.Tkt.Value] = asset.Tkt.Category
					} else {
						e = fmt.Errorf("in tkt duplicate: %v", asset.Tkt.Value)
					}
					return
				}
			} else {
				return
			}
		} else {
			return
		}
	} else {
		return
	}
}

func (self *CKState) AddOut(asset *assets.Asset) (added bool, e error) {
	added = false
	if asset.Tkn != nil {
		if self.outPlus {
			self.cy.add(&asset.Tkn.Currency, &asset.Tkn.Value)
		} else {
			self.cy.sub(&asset.Tkn.Currency, &asset.Tkn.Value)
		}
		added = true
	}
	if asset.Tkt != nil {
		if _, ok := self.tk[asset.Tkt.Value]; !ok {
			if self.outPlus {
				self.tk[asset.Tkt.Value] = asset.Tkt.Category
				added = true
			} else {
				e = fmt.Errorf("out tkt not in ins: %v", asset.Tkt.Value)
			}
		} else {
			if !self.outPlus {
				delete(self.tk, asset.Tkt.Value)
				added = true
			} else {
				e = fmt.Errorf("out tkt not in ins: %v", asset.Tkt.Value)
			}
			return
		}
	}
	return
}

func (self *CKState) CheckToken() (e error) {
	for currency, state := range self.cy {
		if state.balance.Cmp(&utils.I256_0) != 0 {
			e = fmt.Errorf("currency %v banlance != 0", currency)
			return
		}
	}
	return
}

func (self *CKState) CheckTicket() (e error) {
	if len(self.tk) > 0 {
		e = fmt.Errorf("tikect not use %v", self.tk)
		return
	} else {
		return
	}
}

func (self *CKState) Check() (e error) {
	if e = self.CheckToken(); e != nil {
		return
	}
	if e = self.CheckTicket(); e != nil {
		return
	}
	return
}
