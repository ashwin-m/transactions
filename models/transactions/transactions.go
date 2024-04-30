package transactions

type Transactions struct {
	id                   int64
	sourceAccountId      int64
	destinationAccountId int64
	amount               float64
}

func (t *Transactions) GetId() int64 {
	return t.id
}

func (t *Transactions) GetSourceAccountId() int64 {
	return t.sourceAccountId
}

func (t *Transactions) GetDestinationAccountId() int64 {
	return t.destinationAccountId
}

func (t *Transactions) GetAmount() float64 {
	return t.amount
}

func (t *Transactions) SetId(id int64) {
	t.id = id
}

func (t *Transactions) SetSourceAccountId(sourceAccountId int64) {
	t.sourceAccountId = sourceAccountId
}

func (t *Transactions) SetDestinationAccountId(destinationAccountId int64) {
	t.destinationAccountId = destinationAccountId
}

func (t *Transactions) SetAmount(amount float64) {
	t.amount = amount
}
