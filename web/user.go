package web

import (
	model "company/finance/internal"
	"company/finance/internal/user"

	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func addUser(c echo.Context) error {

	form := model.User{}

	if err := c.Bind(&form); err != nil {
		log.Println("Error in Binding Input Form of addUser : ", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Invalid Input"})
	}

	log.Println("Validating Input Form of addUser")
	if err := form.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Service is not Available"})
	}
	log.Println("Validating Input Form of addUser Done")

	er := user.Add(form)
	if er != nil {
		log.Println("Error in Adding User :", er)
		return c.JSON(http.StatusBadRequest, echo.Map{"err": "Service is not Available"})
	}

	return c.JSON(http.StatusOK, echo.Map{"msg": "Your User Added Successfully"})

}

func login(c echo.Context) error {

	form := make(map[string]string)

	log.Println("Binding Input Form ...")
	if err := c.Bind(&form); err != nil {
		log.Println("Error in Binding Input Form in Login :", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Bad Request"})
	}

	email := form["email"]
	email = strings.ToLower(email)
	password := form["password"]

	log.Println(email, "Loging in")

	u, err := user.LoadUserByEmail(email)
	if err != nil {
		log.Println("Error in Loading User with Email :", email, ", Error :", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if u.Password != password {
		log.Println("Password is Wrong for email :", email)

		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Wrong username or password"})

	}

	token, err := user.CreateToken(u.ID.Hex(), email)
	if err != nil {
		log.Println("Error in Creating Token for Username :", email)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal Server Error"})
	}

	u.Token = token

	log.Println("Updating User with email :", email)
	if err := user.Update(u, u.ID.Hex()); err != nil {
		log.Println("Error in Updating User with email :", u.Email, ", Error :", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal Server Error"})
	}

	log.Println(email, "Loged in")

	return c.JSON(http.StatusOK, echo.Map{
		"user": u.RestLogin(),
	})

}
