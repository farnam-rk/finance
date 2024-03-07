package web

import (
	model "company/finance/internal"
	"company/finance/internal/budget"
	"company/finance/internal/user"

	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func addBudget(c echo.Context) error {

	form := model.Budget{}

	if err := c.Bind(&form); err != nil {
		log.Println("Error in Binding Input Form of addBudget : ", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Invalid Input"})
	}

	_, claims, err := user.CheckToken(c)
	if err != nil {
		log.Println("Error in Checking Token in Add Reward", err.Error())
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	err = budget.Add(form, claims["email"].(string))
	if err != nil {
		log.Println("Error in Adding Budget :", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"err": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"msg": "Your Budget Added Successfully"})

}

func budgetHistory(c echo.Context) error {

	_, claims, err := user.CheckToken(c)
	if err != nil {
		log.Println("Error in Checking Token in Add Reward", err.Error())
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	budgets, er := budget.History(claims["email"].(string))
	if er != nil {
		log.Println("Error in Adding Budget :", er)
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Service is not Available"})
	}

	return c.JSON(http.StatusOK, echo.Map{"budgets": budgets})

}
