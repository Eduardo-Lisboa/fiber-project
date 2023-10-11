package controllers

import (
	"codigo-fluente/database"
	"codigo-fluente/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
}

func Register(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	if data["password"] != data["confirm_password"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	if len(strings.TrimSpace(user.Email)) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messagem": "email é obrigatorio",
		})
	}

	if result := database.DB.Where("email = ?", data["email"]).First(&user); result.RowsAffected > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messagem": "email invalido",
		})
	}
	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})

	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Incorret Password",
		})
	}

	claims := jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"jwt": token,
	})

}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unautenticated",
		})
	}
	claims := token.Claims.(*Claims)
	id := claims.Issuer
	var user models.User
	database.DB.Where("id = ?", id).First(&user)
	return c.JSON(user)

}

func Getusers(c *fiber.Ctx) error {

	var users []models.User

	database.DB.Find(&users)
	if len(users) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "nenhum usuario encontrado",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(users)

}
