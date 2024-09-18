package middleware

import (

	"github.com/gin-gonic/gin"

	allErrors "server/src/all_errors"
)

func ErrorManager() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		if len(context.Errors) > 0 {
			err := context.Errors.Last().Err

			if twtErr, ok := err.(allErrors.TwitsnapError); ok {
				context.JSON(twtErr.Status(), twtErr)
				return
			}

		}

		context.Abort()
	}
}

