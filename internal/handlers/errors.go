package handlers

import "github.com/gin-gonic/gin"

const (
	ErrStatusNotFound   = "NOT_FOUND"
	ErrStatusInternal   = "INTERNAL"
	ErrStatusBadRequest = "BAD_REQUEST"
)

func respondWithError(c *gin.Context, code int, errStatus string, err error) {
	c.JSON(
		code,
		gin.H{
			"error": gin.H{
				"code":    errStatus,
				"message": err.Error(),
			},
		},
	)
}
