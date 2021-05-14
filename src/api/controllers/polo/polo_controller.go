package polo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	polo = "polo"
)

func Marco(c *gin.Context) {
	c.String(http.StatusOK, polo)
}
