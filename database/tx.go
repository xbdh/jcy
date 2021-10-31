package database



type Account string

type TX struct {
	From Account `json:"from"`
	To Account`json:"to"`
	Value uint `json:"value"`
	Data string `json:"data"`
}

func NewTX(from Account, to Account, value uint, data string) TX {
	return TX{From: from, To: to, Value: value, Data: data}
}

func (t TX) IsReward()bool  {
	return t.Data=="reward"
}

func NewAccount(account string)Account  {
	return Account(account)
}
