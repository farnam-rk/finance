package internal

import (
	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Firstnam string             `json:"firstname" bson:"firstname" valid:"required"`
		Lastname string             `json:"lastname" bson:"lastname" valid:"required"`
		Email    string             `json:"email" bson:"email" valid:"required"`
		Password string             `json:"password" bson:"password" valid:"required"`
		Account  Accounts           `json:"account" bson:"account" valid:""`
		Budget   Budgets            `json:"budget" bson:"budget" valid:""`
		Expense  Expenses           `json:"expense" bson:"expense" valid:""`
		Token    string             `json:"token" bson:"token" valid:""`
	}

	Users []User

	Token struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		jwt.StandardClaims
	}

	Account struct {
		No      string `json:"no" bson:"no"`
		Name    string `json:"name" bson:"name"`
		Balance int    `json:"balance" bson:"balance"`
	}

	Accounts []Account

	Budget struct {
		Name      string `json:"name" bson:"name"`
		CapAmount int    `json:"cap_amount" bson:"cap_amount"`
	}

	Budgets []Budget

	Expense struct {
		Title       string `json:"title" bson:"title"`
		Budget      string `json:"budget" bson:"budget"`
		Date        int64  `json:"date" bson:"date"`
		UsedAccount string `json:"used_account" bson:"used_account"`
		Amount      int    `json:"amount" bson:"amount"`
	}

	Expenses []Expense
)

// Validate AddAPIForm
func (f *User) Validate() error {

	_, err := govalidator.ValidateStruct(f)
	if err != nil {
		return err
	}

	return nil
}

// RestLogin Rest of User for Login
func (u *User) RestLogin() echo.Map {

	return echo.Map{
		"token": u.Token,
	}
}
