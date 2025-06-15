package handler

import (
	"net/http"

	"github.com/bignyap/go-admin/utils/converter"
)

func ExtractPaginationDetail(w http.ResponseWriter, r *http.Request) (int, int) {

	pageNumberStr := r.URL.Query().Get("page_number")
	itemsPerPageStr := r.URL.Query().Get("items_per_page")

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
			respondWithError(w, StatusBadRequest, "Invalid items_per_page format")
			return 0, 0
		}
	}

	if pageNumberStr == "" {
		offset = ((defaultPageNumber - 1) * limit)
	} else {
		offset, err = converter.StrToInt(pageNumberStr)
		if err != nil {
			respondWithError(w, StatusBadRequest, "Invalid page_number format")
			return 0, 0
		}
		offset = ((offset - 1) * limit)
	}

	return limit, offset
}
