package expense

import (
	"company/finance/internal"
	"company/finance/internal/user"
	"errors"
	"fmt"

	"log"
)

func Add(expense internal.Expense, email string) (string, error) {

	u, err := user.LoadUserByEmail(email)
	if err != nil {
		log.Println("Error in Loading User with Email :", email, ", Error :", err.Error())
		return "", err
	}

	warning, validAccount, err := checkStatus(expense, u)
	if err != nil {
		return "", err
	}

	u.Account[validAccount].Balance = u.Account[validAccount].Balance - expense.Amount
	u.Expense = append(u.Expense, expense)

	log.Println("Updating User with email :", email)
	if err := user.Update(u, u.ID.Hex()); err != nil {
		log.Println("Error in Updating User with email :", u.Email, ", Error :", err.Error())
		return "", err
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

func checkStatus(expense internal.Expense, u internal.User) (string, int, error) {

	validBudget, err := budgetCheck(expense, u)
	if err != nil {
		return "", 0, err
	}

	warning, err := expenseCheck(expense, u, validBudget)
	if err != nil {
		return "", 0, err
	}

	validAccount, err := accountCheck(expense, u)
	if err != nil {
		return "", 0, err
	}

	return warning, validAccount, nil

}

func budgetCheck(expense internal.Expense, u internal.User) (int, error) {

	var validBudget int
	var validBudgetExist bool

	for i, j := range u.Budget {
		if j.Name == expense.Budget {
			validBudget = i
			validBudgetExist = true
		}
	}
	if !validBudgetExist {
		return 0, errors.New("there is no budget with specified name")
	}
	return validBudget, nil
}

func expenseCheck(expense internal.Expense, u internal.User, validBudget int) (string, error) {

	var spend int
	var warning string

	for _, n := range u.Expense {
		if expense.Budget == n.Budget {
			spend = spend + expense.Amount
		}
	}
	if spend+expense.Amount > u.Budget[validBudget].CapAmount {
		warning = fmt.Sprintf("you reached your cap amount for budget :%s", u.Budget[validBudget].Name)
		log.Println("you reached your cap amount for budget :", u.Budget[validBudget].Name)
	}
	return warning, nil
}

func accountCheck(expense internal.Expense, u internal.User) (int, error) {

	var validAccount int
	var validAccountExit bool

	for p, q := range u.Account {
		if q.Name == expense.UsedAccount {
			validAccount = p
			validAccountExit = true
		}
	}

	if !validAccountExit {
		return 0, errors.New("there is no account with specified name")
	}

	if u.Account[validAccount].Balance < expense.Amount {
		log.Println("Account limit reached with name :", u.Account[validAccount].Name)
		return 0, errors.New("account limit reached")
	}
	return validAccount, nil
}
