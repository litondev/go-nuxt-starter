package controllers

import (
	"time"
	"gorm.io/gorm"
	"plato.com/plato/models"
	"plato.com/plato/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func Login(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	var req LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	email := req.Email
	password := req.Password

	result := map[string]interface{}{}

	query := database.Model(&models.User{})
	query.Select("id","password","email")
	query.Where("email = ?", email)
	query.First(&result)
	if len(result) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"message" : "Email tidak ditemukan",
		})
	}
	
	var isValidPassword bool = helpers.CheckPasswordHash(
		password,
		result["password"].(string),
	)

	if isValidPassword == false {
		return c.Status(500).JSON(fiber.Map{
			"message" : "Password Salah",
		})
	}

	claims := jwt.MapClaims{
		"id": result["id"],
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	access_token, errSigned := token.SignedString([]byte("secret"))

	if errSigned != nil {
		return c.Status(200).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"token": access_token,
	})
}

func Me(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB);	

	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	
	id := claims["id"].(float64)

	result := map[string]interface{}{}

	query := database.Model(&models.User{})
	// queryUser.Select("id","name","type")
	query.Where("id = ?",id)
	query.First(&result)

	delete(result,"password")
	
	return c.Status(200).JSON(result)
}

func Logout(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message" : "success",
	})
}

func RefreshToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	
	claimHeaders := user.Claims.(jwt.MapClaims)

	id := claimHeaders["id"].(float64)

	claims := jwt.MapClaims{
		"id": id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	access_token, errSigned := token.SignedString([]byte("secret"))

	if errSigned != nil {
		return c.Status(401).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"token" : access_token,
	})
}