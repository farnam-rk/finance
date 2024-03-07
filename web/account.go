package web

import (
	"company/finance/internal"
	"company/finance/internal/account"
	"company/finance/internal/user"

	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func addAccount(c echo.Context) error {

	form := internal.Account{}

	if err := c.Bind(&form); err != nil {
		log.Println("Error in Binding Input Form of addAccount : ", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Invalid Input"})
	}

	_, claims, err := user.CheckToken(c)
	if err != nil {
		log.Println("Error in Checking Token in Add Reward", err.Error())
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	err = account.Add(form, claims["email"].(string))
	if err != nil {
		log.Println("Error in Adding Account :", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"err": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"msg": "Your Account Added Successfully"})

}

func accountHistory(c echo.Context) error {

	_, claims, err := user.CheckToken(c)
	if err != nil {
		log.Println("Error in Checking Token in Add Reward", err.Error())
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	accounts, err := account.History(claims["email"].(string))
	if err != nil {
		log.Println("Error in Adding Budget :", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Service is not Available"})
	}

	return c.JSON(http.StatusOK, echo.Map{"accounts": accounts})

}
