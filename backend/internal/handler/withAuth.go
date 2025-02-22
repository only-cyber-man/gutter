package handler

import (
	"github.com/gin-gonic/gin"
)

func (c *Client) WithAuth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(401, map[string]string{
				"error": "missing authorization header",
			})
			ctx.Abort()
			return
		}
		user, err := c.UserService.Verify(token)
		if err != nil {
			ctx.JSON(401, map[string]string{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}
