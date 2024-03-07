package web

import (
	model "company/finance/internal"
	"company/finance/internal/expense"
	"company/finance/internal/user"
	"time"

	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func addExpense(c echo.Context) error {

	form := model.Expense{}

	if err := c.Bind(&form); err != nil {
		log.Println("Error in Binding Input Form of addExpense : ", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Invalid Input"})
	}

	_, claims, err := user.CheckToken(c)
	if err != nil {
		log.Println("Error in Checking Token in Add Reward", err.Error())
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	form.Date = time.Now().Unix()

	warning, err := expense.Add(form, claims["email"].(string))
	if err != nil {
		log.Println("Error in Adding Expense :", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"err": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"msg": "Your Expense Added Successfully", "warning": warning})

}

func expenseHistory(c echo.Context) error {

	form := make(map[string]int64)

	log.Println("Binding Input Form ...")
	if err := c.Bind(&form); err != nil {
		log.Println("Error in Binding Input Form in expense history :", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Bad Request"})
	}

	start := form["start"]
	end := form["end"]

	_, claims, err := user.CheckToken(c)
	if err != nil {
		log.Println("Error in Checking Token in Add Reward", err.Error())
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	expenses, er := expense.History(claims["email"].(string), start, end)
	if er != nil {
		log.Println("Error in Adding Expense :", er)
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Service is not Available"})
	}

	return c.JSON(http.StatusOK, echo.Map{"expenses": expenses})

}
