package config

import (
	"errors"
	"fmt"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Database(env map[string]string) (*gorm.DB, error) {
	dsn := os.Getenv("DB_USER") + ":" + 
		os.Getenv("DB_PASSWORD") + "@tcp(" + 
		os.Getenv("DB_HOST") + ":" + 
		os.Getenv("DB_PORT")+")/" + 
		os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, errDb := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if errDb != nil {
		fmt.Println(errDb)

		return nil, errors.New("Can't Connect To Database")
	}
	
	return db, nil
}