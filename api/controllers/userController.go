package controllers 

import (
	"github.com/gofiber/fiber/v2"
	"plato.com/plato/models"
	"plato.com/plato/helpers"
	"gorm.io/gorm"
	"github.com/xuri/excelize/v2"
	"strconv"
	"path/filepath"
	"github.com/jung-kurt/gofpdf"
	"fmt"
	"os"
	"time"
)

type UserRequest struct {
    Name    string `json:"nama" binding:"required"`
	Email    string `json:"email" binding:"required"`
    Password *string `json:"password"`
}

func IndexQuery(c *fiber.Ctx) *gorm.DB {
	database := c.Locals("DB").(*gorm.DB)

	var modelData models.User

	search := c.Query("search","")

	query := database.Model(&modelData)
	
	query.Select("id","name","email")

	if search != "" {
		query.Where("name LIKE ?", "%" + search + "%")		
	}
	
	return query
}

func IndexUser(c *fiber.Ctx) error {
	var resultCount int64
	queryResultCount := IndexQuery(c)
	queryResultCount.Count(&resultCount)		

	new_per_page,
	total_page,
	limit_start := helpers.MakePagination(c,resultCount)

	result := []map[string]interface{}{}
	query := IndexQuery(c)
	query.Order("id desc")
	query.Offset(limit_start)		
	query.Limit(new_per_page)		
	query.Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"data" : result,
		"per_page" : new_per_page,
		"total_page" : total_page,
		"total_data" : int(resultCount),
	})
}

func ShowUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id_raw := c.Params("id")

	id, err := strconv.Atoi(id_raw)
    
	if err != nil {
        return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	result := map[string]interface{}{}

	query := database.Model(&models.User{})
	query.Where("id = ?",id)
	query.First(&result)

	if len(result) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
	}

	delete(result,"password")

	return c.Status(200).JSON(result)
}

func StoreUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB);	

   	var req UserRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	name := req.Name;
	email := req.Email;
	password := req.Password;

	if(password == nil){
		return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
	}

	result := map[string]interface{}{}
	
	query := database.Model(&models.User{})
	query.Select("id")
	query.Where("email = ?",email)
	query.First(&result)
	if len(result) > 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "Email telah terpakai",
		})
	}

	hash,errorHash := helpers.HashPassword(*password)
	if errorHash != nil {
		return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
	}

	user := models.User{
		Name : name,
		Email : email,
		Password : hash,
	}

	if err := database.Create(&user).Error;err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
	}

	result = map[string]interface{}{}

	query = database.Model(&models.User{})
	// queryUser.Select("id","name","type")
	query.Where("id = ?",user.Id)
	query.First(&result)

	delete(result,"password")

	return c.Status(200).JSON(result)
}

func UpdateUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id_raw := c.Params("id")

	id, err := strconv.Atoi(id_raw)
    
	if err != nil {
        return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }
	
   	var req UserRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	name := req.Name;
	email := req.Email;
	password := req.Password;

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

	update := &models.User{
		Email : email,
		Name : name,
	}
	
	if(password != nil){
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

func DestroyUser(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id_raw := c.Params("id")

	id, err := strconv.Atoi(id_raw)
    
	if err != nil {
        return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	database.Where("id = ?",id).Delete(&models.User{})

	return c.Status(200).JSON(fiber.Map{
		"message" : "success",
	})
}

func ExcelUser(c *fiber.Ctx) error {
	f := excelize.NewFile()
    sheet := "Sheet1"

	_, err := f.NewSheet(sheet)
    if err != nil {
        return c.Status(200).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	users := []models.User{}
	query := IndexQuery(c)
	query.Order("id desc")
	query.Find(&users)

	f.SetCellValue(sheet,"A1", "Name")
    f.SetCellValue(sheet,"B1", "Email")

    for i, user := range users {
        row := i + 2 
        f.SetCellValue(sheet,fmt.Sprintf("A%d", row), user.Name)
        f.SetCellValue(sheet,fmt.Sprintf("B%d", row), user.Email)
    }

	fileName := filepath.Base("") + "/assets/users.xlsx"

    if err := f.SaveAs(fileName); err != nil {
		return c.Status(200).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	go func() {
		time.Sleep(2 * time.Second)
		
		_, err := os.Stat(fileName)

		if !os.IsNotExist(err) {
			os.Remove(fileName)
		}
	}()

	return c.Status(200).JSON(fiber.Map{
		"message" : "success",
	})

	return c.Download(fileName)
}

func PdfUser(c *fiber.Ctx) error {
	users := []models.User{}
	query := IndexQuery(c)
	query.Order("id desc")
	query.Find(&users)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage() 

	// HEADER
	// pdf.Cell(40, 10, "User Table")
	pdf.SetFont("Arial", "", 10)
	pageWidth, pageHeight := pdf.GetPageSize()
	cellWidth := 100.0
	pdf.SetX((pageWidth - cellWidth) / 2)
	pdf.CellFormat(cellWidth, 5, "Company Name", "", 1, "C", false, 0, "")
	pdf.SetX((pageWidth - cellWidth) / 2)
	pdf.CellFormat(cellWidth, 5, "User Report", "", 1, "C", false, 0, "")
	pdf.SetX((pageWidth - cellWidth) / 2)
	pdf.CellFormat(cellWidth, 5, "Date: 2025-05-30", "", 1, "C", false, 0, "")
    pdf.Ln(8)

    // TABLE HEADER
    pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(20, 5, "No", "1", 0, "L", false, 0, "")
    pdf.CellFormat(20, 5, "Name", "1", 0, "L", false, 0, "")
    pdf.CellFormat(30, 5, "Email", "1", 0, "L", false, 0, "")
    pdf.Ln(-1)

    // TABLE CONTENT
    pdf.SetFont("Arial", "", 8)
    for index, user := range users {
		pdf.CellFormat(20, 5, fmt.Sprintf("%d", index+1), "1", 0, "", false, 0, "")
        pdf.CellFormat(20, 5, user.Name, "1", 0, "", false, 0, "")
        pdf.CellFormat(30, 5, user.Email, "1", 0, "", false, 0, "")
        pdf.Ln(-1)
    }

	// SIGNATURE
	pdf.SetY(pageHeight - 50)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(20, 5, "Signature : ", "", 1, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(20, 5, "____________________", "", 1, "L", false, 0, "")
	pdf.CellFormat(20, 5, "Admin", "", 1, "L", false, 0, "")

	fileName := filepath.Base("") + "/assets/users.pdf"
    
    if err := pdf.OutputFileAndClose(fileName);err != nil {
          return c.Status(500).JSON(fiber.Map{
			"message" : "Terjadi Kesalahan",
		})
    }

	// go func() {
	// 	time.Sleep(2 * time.Second)
		
	// 	_, err := os.Stat(fileName)

	// 	if !os.IsNotExist(err) {
	// 		os.Remove(fileName)
	// 	}
	// }()

	return c.Status(200).JSON(fiber.Map{
		"message" : "success",
	})

	return c.Download(fileName)
}