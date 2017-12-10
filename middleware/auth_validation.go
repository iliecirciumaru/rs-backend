package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/iliecirciumaru/rs-backend/repo"
	"github.com/iliecirciumaru/rs-backend/structs"
	"net/http"
)

const AUTH_HEADER = "X-RS-AUTH-TOKEN"

func AuthValidation(userRepo repo.UserRepo) gin.HandlerFunc {

	return func(c *gin.Context) {
		if (c.Request.RequestURI == "/user" && c.Request.Method == http.MethodPost) ||
			c.Request.RequestURI == "/login" || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		authToken := c.Request.Header.Get(AUTH_HEADER)
		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, structs.CustomError{"Token 'X-RS-AUTH-TOKEN' is not valid"})
		}
		user, err := userRepo.ValidateUserByToken(authToken)

		if err == nil && user.ID > 0 {
			c.Set("user", user)
			c.Next()

			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, structs.CustomError{"Token 'X-RS-AUTH-TOKEN' is not valid"})

	}
}
