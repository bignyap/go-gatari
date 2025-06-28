package gatekeeping

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func (s *GateKeepingService) ValidateRequestHeader(c *gin.Context) (*ValidateRequestInput, error) {

	org := c.GetHeader("X-Org-Name")
	path := c.GetHeader(":path")

	if org == "" || path == "" {
		return nil, errors.New("missing required headers")
	}

	return &ValidateRequestInput{
		OrganizationName: org,
		EndpointName:     path,
	}, nil
}

func (s *GateKeepingService) RecordUsageValidator(c *gin.Context) (*RecordUsageInput, error) {

	org := c.GetHeader("X-Org-Name")
	path := c.GetHeader(":path")

	if org == "" || path == "" {
		return nil, errors.New("missing required headers")
	}

	return &RecordUsageInput{
		OrganizationName: org,
		EndpointName:     path,
	}, nil
}
