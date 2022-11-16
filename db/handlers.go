package db

import (
	"log"
	"strconv"

	m "example.com/models"
	u "example.com/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

// ErrorResponse is error response
type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

// ValidateStruct validate struct
func ValidateStruct(dog m.Dog) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(dog)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

// AddDog do insert one dog
func AddDog(c *fiber.Ctx) error {
	dog := new(m.Dog)

	if err := c.BodyParser(dog); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateStruct(*dog)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	Database.Create(&dog)
	return c.Status(201).JSON(dog)
}

// UpdateDog do update one dog with all fields
func UpdateDog(c *fiber.Ctx) error {
	dog := new(m.Dog)
	if err := c.BodyParser(dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	log.Printf("update: %+v", dog)

	// Transaction return nil or error
	err := Database.Transaction(func(tx *gorm.DB) error {
		id := c.Params("id")
		// 모든 필드에 대해 업데이트 (ID 제외)
		if err := tx.Model(&dog).Select("*").Where("ID = ?", id).Omit("ID").Updates(dog).Error; err != nil {
			return err
		}
		return nil // commit
	})

	if err != nil {
		log.Fatalln(err)
		return c.Status(503).SendString(err.Error())
	}
	return c.Status(200).JSON(dog)
}

// UpdateDogPartial do update one dog with partial fields
func UpdateDogPartial(c *fiber.Ctx) error {
	dog := new(m.Dog)
	if err := c.BodyParser(dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	tableName, fieldNames := GetTableJSONTags(Database, dog)
	if fieldNames == nil {
		return c.Status(503).SendString("Any JSON tag is not defined")
	}

	// 업데이트 대상 json tag 추출
	var tags []string = u.ExtractFields(c.Body())
	// json tag 를 field name 로 변환
	var fields []string = u.ReplaceSliceByMap(tags, u.MapS(fieldNames).Reverse())
	log.Printf("%s: tags %v => fields %+v", tableName, tags, fields)

	err := Database.Transaction(func(tx *gorm.DB) error {
		id := c.Params("id")
		// 업데이트 대상 필드(fields)들만 업데이트
		if err := tx.Model(&dog).Select(fields).Where("ID = ?", id).Updates(dog).Error; err != nil {
			return err
		}
		return nil // commit
	})

	if err != nil {
		log.Fatalln(err)
		return c.Status(503).SendString(err.Error())
	}
	return c.Status(200).JSON(dog)
}

// RemoveDog do delete one dog
func RemoveDog(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(406).SendString(err.Error())
	}

	var dog m.Dog = m.Dog{ID: id} // for BeforeDelete
	result := Database.Model(&dog).Delete(&dog, id)
	if result.Error != nil {
		return c.Status(403).SendString(result.Error.Error()) // 403 Forbidden
	}

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.SendStatus(200)
}

// getParams get All query params
// **NOTE: c.AllParams() is not working
// https://stackoverflow.com/a/73736090/6811653
func getParams(c *fiber.Ctx) map[string]string {
	params := make(map[string]string)
	var err error
	c.Context().QueryArgs().VisitAll(func(key, val []byte) {
		if err != nil {
			return
		}
		k := utils.UnsafeString(key)
		v := utils.UnsafeString(val)
		params[k] = v
	})
	return params
}

// GetDogs do select all dogs
func GetDogs(c *fiber.Ctx) error {
	var dogs []m.Dog

	params := getParams(c)
	for k, v := range params {
		if k == "age" {
			Database.Where("age >= ?", v).Find(&dogs)
			return c.Status(200).JSON(dogs)
		}
	}

	Database.Find(&dogs) // All dogs
	return c.Status(200).JSON(dogs)
}

// GetDog do select one dog
func GetDog(c *fiber.Ctx) error {
	id := c.Params("id")
	var dog m.Dog

	result := Database.Find(&dog, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(&dog)
}

/*
// DogParams is QueryParams of FindDogs
type DogParams struct {
	Name  string `query:"name"`
	Breed string `query:"breed"`
	Age   int    `query:"age"`
}
*/
