package main

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/cors"	
	"plato.com/plato/config"
	"plato.com/plato/controllers"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	// LOAD ENV
	errEnv := godotenv.Load()
	
	if errEnv != nil {
		fmt.Println(errEnv)
		os.Exit(1)
	}

	// LOAD DB
	var dsn map[string]string = map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
	}

	db, errDb := config.Database(dsn)

	if errDb != nil {
		fmt.Println(errDb)
		os.Exit(1)
	}

	// DEBUG
	debug := os.Getenv("APP_DEBUG")

	appDebug, errDebug := strconv.ParseBool(debug)
	if errDebug != nil {
		fmt.Println(errDebug)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}			

			logFile, logFileError := os.OpenFile(os.Getenv("APP_LOGGER_LOCATION"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if logFileError != nil {
				ctx.Status(500).JSON(fiber.Map{
					"message": "Terjadi Kesalahan",
				})
			}

			logger := log.New(logFile, "Error : ", log.LstdFlags)
			logger.Println(time.Now().String())
			logger.Println(err.Error())
			fmt.Println(err.Error())

			var message string = "Terjadi Kesalahan"

			if appDebug == true {
				message = err.Error()
			}
				
			err = ctx.Status(code).JSON(fiber.Map{
				"message" : message,
				"code" : code,
			})

			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message" : "Terjadi Kesalahan",
				})
			}

			return ctx.Status(500).JSON(fiber.Map{
				"message" : "Terjadi Kesalahan",
			})
		},
	})

	// CORS
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}));
	
	// GLOBAL MIDDLEWARE
	app.Use(func(c *fiber.Ctx) error {
		// fmt.Println("Setiap Request Akan Masuk Sini")
		c.Locals("DB", db)
		c.Locals("DEBUG", debug)

		return c.Next()
	})

	// ASSET
	app.Static("/assets", "./assets")

	// API 
	api := app.Group("/api")  

	v1 := api.Group("/v1")   

	v1.Get("/status",func(c  *fiber.Ctx) error {
		return c.Status(200).JSON(map[string]string{
			"message" : "active",
		})
	})

	// v1.Post("/register",controllers.Register)
	v1.Post("/login",controllers.Login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),		
	}))

	v1.Post("/refresh-token",controllers.RefreshToken)
	v1.Post("/logout",controllers.Logout)
	v1.Get("/me", controllers.Me)
	v1.Put("/profil/update",controllers.UpdateProfilData)
	v1.Post("/profil/upload",controllers.UpdateProfilPhoto)

	v1.Get("/user",controllers.IndexUser)
	v1.Get("/user/excel",controllers.ExcelUser)
	v1.Get("/user/pdf",controllers.PdfUser)
	v1.Post("/user",controllers.StoreUser)
	v1.Put("/user/:id",controllers.UpdateUser)
	v1.Delete("/user/:id",controllers.DestroyUser)
	v1.Get("/user/:id",controllers.ShowUser)

	app.Listen(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
}