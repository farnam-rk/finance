package account

import (
	model "company/finance/internal"
	"company/finance/internal/user"
	"errors"

	"log"
)

func Add(account model.Account, email string) error {

	u, err := user.LoadUserByEmail(email)
	if err != nil {
		log.Println("Error in Loading User with Email :", email, ", Error :", err.Error())
		return err
	}

	for _, j := range u.Account {
		if j.Name == account.Name {
			return errors.New("this account already exist")
		}
	}

	u.Account = append(u.Account, account)

	log.Println("Updating User with email :", email)
	if err := user.Update(u, u.ID.Hex()); err != nil {
		log.Println("Error in Updating User with email :", u.Email, ", Error :", err.Error())
		return err
	}

	return nil
}

func History(email string) (model.Accounts, error) {

	u, err := user.LoadUserByEmail(email)
	if err != nil {
		log.Println("Error in Loading User with Email :", email, ", Error :", err.Error())
		return model.Accounts{}, err
	}

	return u.Account, nil
}
