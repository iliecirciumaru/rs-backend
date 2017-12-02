package middleware

import (
	"github.com/iliecirciumaru/rs-backend/repo"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/iliecirciumaru/rs-backend/structs"
)

func AuthValidation(userRepo repo.UserRepo) gin.HandlerFunc {

	return func(c *gin.Context) {
		if (c.Request.RequestURI == "/user" && c.Request.Method == http.MethodPost) || c.Request.RequestURI == "/login" {
			c.Next()
			return
		}

		authToken := c.Request.Header.Get("X-RS-AUTH-TOKEN")


		user, err := userRepo.ValidateUserByToken(authToken)


		if err == nil && user.ID > 0 {
			c.Set("user", user)
			c.Next()

			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, structs.CustomError{"Token 'X-RS-AUTH-TOKEN' is not valid"})

	}
}