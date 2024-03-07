package user

import (
	"company/finance/config"
	"company/finance/db"
	model "company/finance/internal"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Add(user model.User) error {

	_, err := db.Insert("user", user)
	if err != nil {
		log.Println("user insertion failed")
		return err
	}

	return nil
}

// CreateToken Create JWT Token
func CreateToken(id, email string) (string, error) {

	claims := model.Token{
		ID:    id,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(48 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.ConfigValue.JWTSecret))

	return t, err
}

// CheckToken Check Token is Valid or not
func CheckToken(c echo.Context) (string, jwt.MapClaims, error) {

	var token string
	token = c.Request().Header.Get("token")

	if token == "" {
		token = c.QueryParam("token")
		if token == "" {
			log.Println("Token is Nil")
			return "", nil, errors.New("Token is Nil")
		}
	}

	claims := jwt.MapClaims{}

	// Check Token by User Secret Key
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.ConfigValue.JWTSecret), nil
	})

	if t.Valid == false || err != nil { // Token is Invalid -> Check Token by Admin Secret Key

		return "", nil, err

	}

	// User
	return token, claims, nil
}

// LoadUserByEmail Load User by Email
func LoadUserByEmail(email string) (model.User, error) {

	filter := bson.M{"email": email}

	u, err := db.Find("user", filter, nil)
	if err != nil {
		return model.User{}, errors.New("User not Found")
	}
	if u == nil {
		return model.User{}, errors.New("User not Found")
	}

	users, err := convertToStruct(u)
	if users == nil {
		return model.User{}, errors.New("User not Found")
	}

	return users[0], nil

}

// convertToStruct Convert Mongo Curser to User Struct
func convertToStruct(cur *mongo.Cursor) (model.Users, error) {

	var users model.Users

	for cur.Next(nil) {

		var user model.User

		err := cur.Decode(&user)
		if err != nil {
			return model.Users{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Update Update User
func Update(data model.User, id string) error {

	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID is not Correct")
	}
	log.Println("Updating User with email :", data.Email, "in Update Func")

	filter := bson.M{"_id": ObjectID}
	newData := db.ParseBson(data)
	_, err = db.Update("user", filter, newData, nil)
	if err != nil {
		return err
	}

	return nil
}
