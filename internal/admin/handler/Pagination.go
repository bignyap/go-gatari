package adminHandler

import (
	"fmt"

	converter "github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

func ExtractPaginationDetail(c *gin.Context) (int, int, error) {

	pageNumberStr := c.Query("page_number")
	itemsPerPageStr := c.Query("items_per_page")

	defaultPageNumber := 1
	defaultItemsPerPage := 25

	var limit int
	var offset int
	var err error

	if itemsPerPageStr == "" {
		limit = defaultItemsPerPage
	} else {
		limit, err = converter.StrToInt(itemsPerPageStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid items_per_page format")
		}
	}

	if pageNumberStr == "" {
		offset = ((defaultPageNumber - 1) * limit)
	} else {
		offset, err = converter.StrToInt(pageNumberStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid page_number format")
		}
		offset = ((offset - 1) * limit)
	}

	return limit, offset, nil
}
