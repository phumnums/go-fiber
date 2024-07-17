// method
package routes

import (
	c "go-fiber-test/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func InetRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")
	profile := v1.Group("/profile")

	profile.Get("/", c.GetProfile)
	profile.Get("/:id", c.GetProfileByID)

	// middleware
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			// john and admin can access api
			"john":    "doe",
			"admin":   "123456",
			"gofiber": "21022566",
			"testgo":  "23012023",
		},
	}))

	// Group api

	v2 := api.Group("/v2")
	v3 := api.Group("/v3")

	v1.Get("/", c.HelloTest)              // Hello World
	v1.Post("/", c.BodyParserTest)        //  BodyParser
	v1.Get("/user/:name", c.ParamsTest)   // Params
	v1.Post("/inet", c.QueryTest)         // Query
	v1.Post("/valid", c.ValidTest)        // Validation เช็คข้อมูลผิดพลาด
	v1.Post("/fact/:number", c.Factorial) // 5.1 Factorial
	v1.Post("/register", c.Register)      //6 register

	v2.Get("/", c.HelloTestV2)

	v3.Get("/:name", c.TaxID) //QueryParam 5.2 ascii

	//CRUD dogs
	dog := v1.Group("/dog")

	dog.Get("", c.GetDogs)
	dog.Get("/filter", c.GetDog)
	dog.Get("/json", c.GetDogsJson)
	dog.Post("/", c.AddDog)
	dog.Put("/:id", c.UpdateDog)
	dog.Delete("/:id", c.RemoveDog)
	dog.Get("/deleted", c.GetDeletedDogs) //7.0.2
	dog.Get("/range", c.GetRangeDogs)     //7.1
	dog.Get("/newjson", c.GetNewDogJson)

	//CRUD company 7.0.1
	company := v1.Group("/company")

	company.Post("/", c.AddCompany)
	company.Get("/", c.GetCompany)
	company.Get("/:id", c.GetCompanyByID)
	company.Put("/:id", c.UpdateCompany)
	company.Delete("/:id", c.RemoveCompany)

	//CRUD profile User Project

	profile.Post("/", c.AddProfile)
	profile.Put("/:id", c.UpdateProfile)
	profile.Delete("/:id", c.RemoveProfile)

}
