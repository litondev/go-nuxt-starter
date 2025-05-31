package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/joho/godotenv"
	"plato.com/plato/config"
	"plato.com/plato/models"
	"plato.com/plato/helpers"
	"math/rand"
	"sync"
)

func main() {
	errEnv := godotenv.Load()

	if errEnv != nil {
		fmt.Println(errEnv)
		os.Exit(1)
	}

	var dsn map[string]string = map[string]string{
		"DB_HOST" : os.Getenv("DB_HOST"),
		"DB_PORT" : os.Getenv("DB_PORT"),
		"DB_NAME" : os.Getenv("DB_NAME"),
		"DB_USER" : os.Getenv("DB_USER"),
		"DB_PASSWORD" : os.Getenv("DB_PASSWORD"),
	}

	database, errDb := config.Database(dsn)

	if errDb == nil {
		// MIGRATE AND SEED USER
			tableUser := database.Migrator().HasTable(&models.User{})
			if tableUser == true {
				database.Migrator().DropTable(&models.User{})
			}

			database.AutoMigrate(&models.User{})
		
			// DEFAULT
			// for i := 0; i < 500; i++ {
			// 	hash, _ := helpers.HashPassword("password")
			// 	var id string = "user" + strconv.Itoa(i)

			// 	user := models.User{
			// 		Name:     id,
			// 		Email:    id + "@gmail.com",
			// 		Password: hash,
			// 	}

			// 	result := database.Create(&user)

			// 	if result.Error != nil {
			// 		fmt.Print(result.Error)
			// 		os.Exit(1)
			// 	}
			// }

			// WAITGROUP
			var wg sync.WaitGroup
			sem := make(chan struct{}, 25) // limit to 10 concurrent goroutines

			for i := 0; i < 500; i++ {
				wg.Add(1)
				sem <- struct{}{}  // acquire token

				go func(i int) {
					defer wg.Done()
					defer func() { <-sem }() // release token

					hash, _ := helpers.HashPassword("password")
					id := "user" + strconv.Itoa(rand.Intn(10000))

					user := models.User{
						Name:     id,
						Email:    id + "@gmail.com",
						Password: hash,
					}

					database.Create(&user)

					// result := database.Create(&user)
					// if result.Error != nil {
						// fmt.Println(result.Error)
					// }
				}(i)
			}

			wg.Wait()

			fmt.Println("Success Migrate And Seed User")
		// MIGRATE AND SEED USER 	
	}
}
