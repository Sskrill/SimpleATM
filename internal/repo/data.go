package repo

import (
	"errors"
)

type Account struct {
	balance float64
}

func (a *Account) Deposit(amount float64) error {

	a.balance = a.balance + amount
	return nil
}
func (a *Account) Withdraw(amount float64) error {
	if a.balance < amount {
		return errors.New("Not enough")
	}
	a.balance -= amount
	return nil
}
func (a *Account) GetBalance() float64 {
	return a.balance
}
