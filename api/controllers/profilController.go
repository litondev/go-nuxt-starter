package controllers 

import (	
	"os"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"plato.com/plato/helpers"
	"plato.com/plato/models"
	"path/filepath"
	"github.com/disintegration/imaging"
	"gorm.io/gorm"
)

type UpdateProfilRequest struct {
    Name    string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
    Password *string `json:"password"`
	Password_Confirm *string `json:"password_confirm""`
}

func UpdateProfilData(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB);	
	
	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)

	var id uint = uint(claims["id"].(float64))

	var req UpdateProfilRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	
	name := req.Name;
	email := req.Email;
	password := req.Password;
	password_confirm := req.Password_Confirm

	result := map[string]interface{}{}
	
	query := database.Model(&models.User{})
	query.Select("id")
	query.Where("email = ?",email)
	query.Not("id = ?",id)
	query.First(&result)
	if len(result) > 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "Email telah terpakai",
		})
	}

	result = map[string]interface{}{}

	query = database.Model(&models.User{})
	query.Where("id = ?",id)
	query.First(&result)

	if(password_confirm == nil){
		return c.Status(500).JSON(fiber.Map{
			"message": "Password Konfirmasi Tidak Valid",
		})	
	}

	var isValidPassword bool = helpers.CheckPasswordHash(
		*password_confirm,
		result["password"].(string),
	)
	if isValidPassword == false {
		return c.Status(500).JSON(fiber.Map{
			"message": "Password Konfirmasi Tidak Valid",
		})		
	}

	update := &models.User{
		Email : email,
		Name : name,
	}
	
	if(password != nil && *password != ""){
		hash, errHash := helpers.HashPassword(*password)

		if(errHash != nil){
			return c.Status(500).JSON(fiber.Map{
				"message" : "Terjadi Kesalahan",
			})			
		}

		update.Password = hash
	}
		
	query = database.Model(&models.User{})
	query.Where("id = ?",id)
	query.Updates(&update);
	if query.Error != nil {		
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	result = map[string]interface{}{}

	query = database.Model(&models.User{})
	// queryUser.Select("id","name","type")
	query.Where("id = ?",id)
	query.First(&result)

	delete(result,"password")

	return c.Status(200).JSON(result)
}

func UpdateProfilPhoto(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB);	

	file, errGetFile := c.FormFile("photo")
	if errGetFile != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})		
	}

	user := c.Locals("user").(*jwt.Token)
	
	claims := user.Claims.(jwt.MapClaims)

	var id uint = uint(claims["id"].(float64))

    var stringID string = strconv.FormatUint(uint64(id),10)

	ext := filepath.Ext(file.Filename)

	filename := stringID + "-default" + ext;
	
	pathname := filepath.Base("") + "/assets/users/" + filename

	if _,errFileExists := os.Stat(pathname); errFileExists == nil {
		errRemoveFile := os.Remove(pathname)
		if errRemoveFile != nil {
			return c.Status(500).JSON(fiber.Map{
				"message" : "Terjadi Kesalahan",
			})
		}
	}

	if errUploadFile := c.SaveFile(file, pathname); errUploadFile != nil {				
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	openFile, errOpenFile := imaging.Open(pathname)
	if(errOpenFile != nil){
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	resizeFile := imaging.Resize(openFile, 128, 128, imaging.Lanczos)
	if errRessizeFile := imaging.Save(resizeFile, pathname);errRessizeFile != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})	
	}

	query := database.Model(&models.User{})
	query.Select("photo")
	query.Where("id = ?",id)
	query.Update("photo",filename)
	if query.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}
		
	return c.Status(200).JSON(fiber.Map{
		"message" : "success",
	})			
}