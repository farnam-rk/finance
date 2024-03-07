package expense

import (
	"company/finance/internal"
	"company/finance/internal/user"
	"errors"
	"fmt"

	"log"
)

func Add(expense internal.Expense, email string) (string, error) {

	var warning string
	var validBudget int
	var validBudgetExist bool

	var validAccount int
	var validAccountExit bool

	var spend int

	u, err := user.LoadUserByEmail(email)
	if err != nil {
		log.Println("Error in Loading User with Email :", email, ", Error :", err.Error())
		return warning, err
	}

	for i, j := range u.Budget {
		if j.Name == expense.Budget {
			validBudget = i
			validBudgetExist = true
		}
	}
	if !validBudgetExist {
		return warning, errors.New("there is no budget with specified name")
	}

	for _, n := range u.Expense {
		if expense.Budget == n.Budget {
			spend = spend + expense.Amount
		}
	}
	if spend+expense.Amount > u.Budget[validBudget].CapAmount {
		warning = fmt.Sprintf("you reached your cap amount for budget :%s", u.Budget[validBudget].Name)
		log.Println("you reached your cap amount for budget :", u.Budget[validBudget].Name)
	}

	for p, q := range u.Account {
		if q.Name == expense.UsedAccount {
			validAccount = p
			validAccountExit = true
		}
	}

	if !validAccountExit {
		return warning, errors.New("there is no account with specified name")
	}

	if u.Account[validAccount].Balance < expense.Amount {
		log.Println("Account limit reached with name :", u.Account[validAccount].Name)
		return warning, errors.New("account limit reached")
	}

	u.Account[validAccount].Balance = u.Account[validAccount].Balance - expense.Amount
	u.Expense = append(u.Expense, expense)

	log.Println("Updating User with email :", email)
	if err := user.Update(u, u.ID.Hex()); err != nil {
		log.Println("Error in Updating User with email :", u.Email, ", Error :", err.Error())
		return warning, err
	}

	return warning, nil
}

func History(email string, start, end int64) (internal.Expenses, error) {

	expenses := internal.Expenses{}

	u, err := user.LoadUserByEmail(email)
	if err != nil {
		log.Println("Error in Loading User with Email :", email, ", Error :", err.Error())
		return internal.Expenses{}, err
	}

	for _, j := range u.Expense {
		if start < j.Date && j.Date < end {
			expenses = append(expenses, j)
		}
	}

	return expenses, nil
}
