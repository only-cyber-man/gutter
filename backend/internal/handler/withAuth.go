package handler

import (
	"github.com/gin-gonic/gin"
)

func (c *Client) WithAuth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(401, map[string]interface{}{
				"success": false,
				"message": "missing authorization header",
				"data":    nil,
			})
			ctx.Abort()
			return
		}
		user, err := c.UserService.Verify(token)
		if err != nil {
			ctx.JSON(401, map[string]interface{}{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}
