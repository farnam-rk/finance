package budget

import (
	"company/finance/internal"
	"company/finance/internal/user"
	"errors"

	"log"
)

func Add(budget internal.Budget, email string) error {

	u, err := user.LoadUserByEmail(email)
	if err != nil {
		log.Println("Error in Loading User with Email :", email, ", Error :", err.Error())
		return err
	}

	for _, j := range u.Budget {
		if j.Name == budget.Name {
			return errors.New("this budget already exist")
		}
	}

	u.Budget = append(u.Budget, budget)

	log.Println("Updating User with email :", email)
	if err := user.Update(u, u.ID.Hex()); err != nil {
		log.Println("Error in Updating User with email :", u.Email, ", Error :", err.Error())
		return err
	}

	return nil
}

func History(email string) (internal.Budgets, error) {

	u, err := user.LoadUserByEmail(email)
	if err != nil {
		log.Println("Error in Loading User with Email :", email, ", Error :", err.Error())
		return internal.Budgets{}, err
	}

	return u.Budget, nil
}
