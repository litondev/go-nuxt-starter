package helpers

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"math"
	// "fmt"
)

func MakePagination(c *fiber.Ctx,resultCount int64) (int,int,int) {
	page := c.Query("page", "1")
	new_page,_ := strconv.Atoi(page)

	per_page := c.Query("per_page", "10")
	new_per_page,_ := strconv.Atoi(per_page)
	
    count_total_page := float64(resultCount) / float64(new_per_page)
    total_page := int( math.Ceil(count_total_page) )
    limit_start := ( new_page - 1 ) * new_per_page ;

	// fmt.Printf("%T\n",count_total_page)
	// fmt.Printf("%T\n",total_page)
	// fmt.Printf("%T\n",limit_start)
	// fmt.Printf("%T\n",new_per_page)

	return new_per_page,total_page,limit_start
}
