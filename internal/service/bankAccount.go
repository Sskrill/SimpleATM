package service

import (
	"errors"
	"github.com/Sskrill/SimpleATM/internal/repo"
	"log"
	"time"
)

func NewData() map[int]BankAccount {
	return make(map[int]BankAccount)
}
func NewService(data map[int]BankAccount) *User {
	return &User{data: data}
}

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

type User struct {
	data map[int]BankAccount
}

func (u *User) CreateAccount() {
	newId := len(u.data) + 1
	u.data[newId] = &repo.Account{}
	log.Println("Create account id:", newId, "Time:", time.Now())
}
func (u *User) AddBalance(id int, amount float64) error {

	acc, ok := u.data[id]

	if !ok {
		return errors.New("Not exists account")
	}

	err := acc.Deposit(amount)
	if err != nil {

		return errors.New("cant deposit")
	}

	log.Println("Add balance to id: ", id, " Time:", time.Now())
	return err
}

func (u *User) WithdrawBalance(id int, amount float64) error {

	acc, ok := u.data[id]
	if !ok {

		return errors.New("Not exists account")
	}

	err := acc.Withdraw(amount)
	if err != nil {

		return err
	}

	log.Println("Withdraw balance from id: ", id, " Time:", time.Now())
	return err
}
func (u *User) ShowBalance(id int) (float64, error) {

	acc, ok := u.data[id]
	if !ok {

		return 0, errors.New("Not exists account")
	}

	sum := acc.GetBalance()

	log.Println("Show account balance id: ", id, " Time:", time.Now())
	return sum, nil
}
