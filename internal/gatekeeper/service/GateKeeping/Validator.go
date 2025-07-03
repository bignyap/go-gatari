package gatekeeping

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (s *GateKeepingService) ValidateRequestHeader(c *gin.Context) (*ValidateRequestInput, error) {

	var input ValidateRequestInput
	if err := c.ShouldBind(&input); err != nil {
		return nil, fmt.Errorf("input validation failed %s", err)
	}

	// org := c.GetHeader("X-Org-Name")
	// method := c.Request.Method
	// path := c.Request.URL.Path

	if input.OrganizationName == "" || input.Path == "" {
		return nil, fmt.Errorf("missing required headers")
	}

	return &input, nil
}

func (s *GateKeepingService) RecordUsageValidator(c *gin.Context) (*RecordUsageInput, error) {

	var input RecordUsageInput
	if err := c.ShouldBind(&input); err != nil {
		return nil, fmt.Errorf("input validation failed %s", err)
	}

	// org := c.GetHeader("X-Org-Name")
	// method := c.Request.Method
	// path := c.Request.URL.Path

	if input.OrganizationName == "" || input.Path == "" {
		return nil, fmt.Errorf("missing required headers")
	}

	return &input, nil
}
