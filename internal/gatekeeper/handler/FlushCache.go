package gateKeeperHandler

import (
	"github.com/gin-gonic/gin"
)

func (h *GateKeeperHandler) FlushAllCacheHandler(c *gin.Context) {

	h.GateKeepingService.FlushAllCache(c.Request.Context())

	h.ResponseWriter.Success(c, gin.H{"message": "Cache cleared successfully"})
}
