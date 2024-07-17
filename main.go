package main

import (
	"fmt"
	"go-fiber-test/database"
	"go-fiber-test/routes"

	m "go-fiber-test/models"

	"github.com/gofiber/fiber/v2" //use freamwork router fiber
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDatabase() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"root",
		"",
		"127.0.0.1",
		"3306",
		"golang_test",
	)
	var err error
	database.DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected!")
	database.DBConn.AutoMigrate(&m.Dogs{}, &m.CompanyData{}, &m.ProfileUser{}) // auto migration create table
}

func main() {

	app := fiber.New() // creater fiber app
	initDatabase()
	routes.InetRoutes(app) // routes api in fuction InetRoutes
	app.Listen(":3000")

	//test
}
