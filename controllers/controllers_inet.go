// function
package controllers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"go-fiber-test/database"
	m "go-fiber-test/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func HelloTest(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
func HelloTestV2(c *fiber.Ctx) error {
	return c.SendString("Hello, World! V2")
}

/* -------------------------Body Parser-------------------------*/
func BodyParserTest(c *fiber.Ctx) error {
	p := new(m.Person)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	log.Println(p.Name) // john
	log.Println(p.Pass) // doe

	str := p.Name + p.Pass
	return c.JSON(str)
}

/* -------------------------Params-------------------------*/
func ParamsTest(c *fiber.Ctx) error {

	str := "hello ==> " + c.Params("name")
	return c.JSON(str)
}

/* -------------------------Query-------------------------*/
func QueryTest(c *fiber.Ctx) error {
	a := c.Query("search")
	str := "my search is  " + a
	return c.JSON(str)
}

/* -------------------------Validation-------------------------*/
func ValidTest(c *fiber.Ctx) error {
	//Connect to database

	// รับค่าแบบ BodyParser
	user := new(m.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	validate := validator.New()
	errors := validate.Struct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}
	return c.JSON(user)
}

/* -------------------------5.1 Factorial-------------------------*/
func Factorial(c *fiber.Ctx) error {

	number := c.Params("number")   //parameter from URL
	num, _ := strconv.Atoi(number) //convert string to int
	fact := 1
	for i := 1; i <= num; i++ {
		fact *= i
	}
	result := fmt.Sprintf("%d! = %d", num, fact)
	// fmt.Println(result)
	return c.SendString(result)
}

/* -------------------------5.2 ascii-------------------------*/
func TaxID(c *fiber.Ctx) error {

	valueTax := c.Query("tax_id") //param from URL

	var asciiValue []string

	for _, v := range valueTax {
		asciiValue = append(asciiValue, strconv.Itoa(int(v)))
	}

	result := fmt.Sprintf("tax_id = %s -> %s", valueTax, asciiValue)
	return c.SendString(result)
}

/* -------------------------6 Validate Register-------------------------*/
func Register(c *fiber.Ctx) error {

	regis := new(m.Register)
	if err := c.BodyParser(&regis); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	validate := validator.New()
	validate.RegisterValidation("usernameValidate", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(fl.Field().String())
	})
	validate.RegisterValidation("websiteValidate", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-z0-9_]+$`).MatchString(fl.Field().String())
	})

	errors := validate.Struct(regis)
	if errors != nil {
		fieldErrors := make(map[string]string)
		for _, err := range errors.(validator.ValidationErrors) {
			if err.Field() == "Email" && err.Tag() == "email" {
				fieldErrors[strings.ToLower(err.Field())] = "Invalid email"
			} else if err.Field() == "Username" && err.Tag() == "usernameValidate" {
				fieldErrors[strings.ToLower(err.Field())] = "ใช้อักษรภาษาอังกฤษ (a-z), (A-Z), ตัวเลข (0-9) และเครื่องหมาย (_), (-) เท่านั้น เช่น Example_01"
			} else if err.Field() == "Password" && (err.Tag() == "min" || err.Tag() == "max") {
				fieldErrors[strings.ToLower(err.Field())] = "ความยาว 6-20 อักษร"
			} else if err.Field() == "Website" && (err.Tag() == "min" || err.Tag() == "max") {
				fieldErrors[strings.ToLower(err.Field())] = "ความยาว 2-30 อักษร"
			} else if err.Field() == "Website" && err.Tag() == "websiteValidate" {
				fieldErrors[strings.ToLower(err.Field())] = "ใช้อักษรภาษาอังกฤษตัวเล็ก (a-z), ตัวเลข (0-9) และเครื่องหมายอักขระพิเศษ ยกเว้นเครื่องหมายขีด (-) ห้ามเว้นวรรคและห้ามใช้ภาษาไทย"
			} else {
				fieldErrors[strings.ToLower(err.Field())] = err.Field() + " is required"
			}
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation errors occurred",
			"errors":  fieldErrors,
		})
	}
	return c.JSON(regis)
}

/* -------------------------CRUD-------------------------*/
func GetDogs(c *fiber.Ctx) error {
	db := database.DBConn //ตัวแปรที่ได้รับจากการ connect
	var dogs []m.Dogs

	db.Find(&dogs) //delelete = null
	return c.Status(200).JSON(dogs)
}

func GetDog(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var dog []m.Dogs

	result := db.Find(&dog, "dog_id = ?", search)

	// returns found records count, equals `len(users)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(&dog)
}

func AddDog(c *fiber.Ctx) error {
	//twst3
	db := database.DBConn
	var dog m.Dogs

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&dog)
	return c.Status(201).JSON(dog)
}

func UpdateDog(c *fiber.Ctx) error {
	db := database.DBConn
	var dog m.Dogs
	id := c.Params("id")

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&dog)
	return c.Status(200).JSON(dog)
}

func RemoveDog(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var dog m.Dogs

	result := db.Delete(&dog, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func GetDogsJson(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Find(&dogs) //10ตัว

	var dataResults []m.DogsRes
	for _, v := range dogs { //1 inet 112 //2 inet1 113
		typeStr := ""
		if v.DogID == 111 {
			typeStr = "red"
		} else if v.DogID == 113 {
			typeStr = "green"
		} else if v.DogID == 999 {
			typeStr = "pink"
		} else {
			typeStr = "no color"
		}

		d := m.DogsRes{
			Name:  v.Name,  //inet
			DogID: v.DogID, //112
			Type:  typeStr, //no color
		}
		dataResults = append(dataResults, d)
		// sumAmount += v.Amount
	}

	r := m.ResultData{
		Data:  dataResults,
		Name:  "golang-test",
		Count: len(dogs), //หาผลรวม,
	}
	return c.Status(200).JSON(r)
}

/* -------------------------7.0.1 CRUD Company-------------------------*/
func AddCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var CompanyData m.CompanyData

	if err := c.BodyParser(&CompanyData); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&CompanyData)
	return c.Status(201).JSON(CompanyData)
}

func GetCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var CompanyData []m.CompanyData

	db.Find(&CompanyData)
	return c.Status(200).JSON(CompanyData)
}

func GetCompanyByID(c *fiber.Ctx) error {
	db := database.DBConn
	var CompanyData m.CompanyData
	id := c.Params("id")

	db.Find(&CompanyData, id)
	return c.Status(200).JSON(CompanyData)
}

func UpdateCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var CompanyData m.CompanyData
	id := c.Params("id")

	if err := c.BodyParser(&CompanyData); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&CompanyData) //update
	return c.Status(200).JSON(CompanyData)
}

func RemoveCompany(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var CompanyData m.CompanyData

	result := db.Delete(&CompanyData, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

/* -------------------------7.0.2 Show Deleted-------------------------*/
func GetDeletedDogs(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Unscoped().Where("deleted_at is not null").Find(&dogs)
	return c.Status(200).JSON(dogs)
}

/* -------------------------7.1 Range Search-------------------------*/
func GetRangeDogs(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Unscoped().Where("dog_id > ? && dog_id < ?", 500, 1000).Find(&dogs)
	return c.Status(200).JSON(dogs)
}

/* -------------------------7.2 Get New Dog Json-------------------------*/
func GetNewDogJson(c *fiber.Ctx) error {
	db := database.DBConn

	var dogs []m.Dogs
	db.Find(&dogs)

	var dataResults []m.DogsRes //slice
	countRed := 0
	countGreen := 0
	countPink := 0
	countNoColor := 0

	for _, v := range dogs {
		typeStr := ""
		if v.DogID >= 10 && v.DogID <= 50 {
			typeStr = "red"
			countRed++
		} else if v.DogID >= 100 && v.DogID <= 150 {
			typeStr = "green"
			countGreen++
		} else if v.DogID >= 200 && v.DogID <= 250 {
			typeStr = "pink"
			countPink++
		} else {
			typeStr = "no color"
			countNoColor++
		}

		d := m.DogsRes{
			Name:  v.Name,  //inet
			DogID: v.DogID, //112
			Type:  typeStr, //no color
		}
		dataResults = append(dataResults, d)

	}
	r := m.ResultData{
		Name:       "golang-test",
		Count:      len(dogs), //หาผลรวมใน slice dogs
		Data:       dataResults,
		SumRed:     countRed,
		SumGreen:   countGreen,
		SumPink:    countPink,
		SumNoColor: countNoColor,
	}
	return c.Status(200).JSON(r)
}

/* -------------------------Profile User-------------------------*/

func AddProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var Profile m.ProfileUser

	if err := c.BodyParser(&Profile); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&Profile)
	return c.Status(201).JSON(Profile)

}

func GetProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var Profile []m.ProfileUser

	db.Find(&Profile)
	return c.Status(200).JSON(Profile)
}

func GetProfileByID(c *fiber.Ctx) error {
	db := database.DBConn
	var Profile m.ProfileUser
	id := c.Params("id")

	db.Find(&Profile, id)
	return c.Status(200).JSON(Profile)
}

func UpdateProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var Profile m.ProfileUser
	id := c.Params("id")

	if err := c.BodyParser(&Profile); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&Profile)
	return c.Status(200).JSON(Profile)

}

func RemoveProfile(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var Profile m.ProfileUser

	result := db.Delete(&Profile, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func GetAgeProfile(c *fiber.Ctx) error {
	db := database.DBConn

	var Profile []m.ProfileUser
	db.Find(&Profile)

	var dataResults []m.AgeProfile
	countGenz := 0
	countGenY := 0
	countGenX := 0
	countBabyBoomer := 0
	countGI := 0

	for _, v := range Profile {
		typeStr := ""
		if v.Age < 24 {
			typeStr = "GenZ"
			countGenz++
		} else if v.Age >= 24 && v.Age <= 41 {
			typeStr = "GenY"
			countGenY++
		} else if v.Age >= 42 && v.Age <= 56 {
			typeStr = "GenX"
			countGenX++
		} else if v.Age >= 57 && v.Age <= 75 {
			typeStr = "Baby Boomer"
			countBabyBoomer++
		} else {
			typeStr = "G.I. Generation"
			countGI++
		}

		d := m.AgeProfile{
			Employee: v.Employee_id,
			Name:     v.Firstname,
			Lastname: v.Lastname,
			Age:      v.Age,
			Gen:      typeStr,
		}
		dataResults = append(dataResults, d)
	}
	r := m.ResultAgeProfile{
		Data:    dataResults,
		Name:    "Age Profile",
		Count:   len(Profile),
		SumGenz: countGenz,
		SumGenY: countGenY,
		SumGenX: countGenX,
		SumBB:   countBabyBoomer,
		SumGI:   countGI,
	}
	return c.Status(200).JSON(r)
}

func GetSearchProfile(c *fiber.Ctx) error {
	db := database.DBConn

	search := strings.TrimSpace(c.Query("search"))
	var Profile []m.ProfileUser

	result := db.Find(&Profile, "firstname = ? OR lastname = ? OR employee_id = ?", search, search, search)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(&Profile)
}
