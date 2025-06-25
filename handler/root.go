package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) RootHandler(c *gin.Context) {
	fmt.Fprintf(c.Writer, "Welcome !!")
}
