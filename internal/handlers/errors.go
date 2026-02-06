package handlers

import "github.com/gin-gonic/gin"

const (
	ErrStatusNotFound   = "NOT_FOUND"
	ErrStatusInternal   = "INTERNAL"
	ErrStatusBadRequest = "BAD_REQUEST"
)

// ErrorResponse ответ ошибка
type ErrorResponse struct {
	Err Error
}

type Error struct {
	Code    string
	Message string
}

func respondWithError(c *gin.Context, code int, errStatus string, err error) {
	c.JSON(
		code,
		gin.H{
			"err": gin.H{
				"code":    errStatus,
				"message": err.Error(),
			},
		},
	)
}
