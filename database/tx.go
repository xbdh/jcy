package database



type Account string

type Tx struct {
	From Account `json:"from"`
	To Account`json:"to"`
	Value uint `json:"value"`
	Data string `json:"data"`
}

func NewTx(from Account, to Account, value uint, data string) Tx {
	return Tx{From: from, To: to, Value: value, Data: data}
}

func (t Tx) IsReward()bool  {
	return t.Data=="reward"
}

func NewAccount(account string)Account  {
	return Account(account)
}
