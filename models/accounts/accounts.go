package accounts

type Accounts struct {
	id      int64
	balance float64
}

func (a *Accounts) GetId() int64 {
	return a.id
}

func (a *Accounts) GetBalance() float64 {
	return a.balance
}

func (a *Accounts) SetId(id int64) {
	a.id = id
}

func (a *Accounts) SetBalance(balance float64) {
	a.balance = balance
}
