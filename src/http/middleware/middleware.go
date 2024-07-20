package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/oimoyu/OimoAdmin/src/utils/_type"
	"github.com/oimoyu/OimoAdmin/src/utils/functions"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("admin_token")
		if err != nil {
			restful.UnLoginErr(c, "empty token")
			c.Abort()
			return
		}

		token := cookie.Value
		claims, err := functions.DecodeToken(token, g.SiteSecret)
		if err != nil {
			restful.UnLoginErr(c, "invalid token")
			c.Abort()
			return
		}

		role, _ := claims["role"].(string)
		if role != _type.ADMIN {
			restful.UnLoginErr(c, "invalid role")
			c.Abort()
			return
		}

		if g.AdminLoginData.Token == "" {
			restful.UnLoginErr(c, "no login record")
			c.Abort()
			return
		}

		if g.AdminLoginData.Token != token {
			restful.UnLoginErr(c, "kick out by another login")
			c.Abort()
			return
		}

		if g.AdminLoginData.IP != c.ClientIP() {
			g.OimoAdmin.Logger.Info("IP does not match the last login, last: [%s], now: [%s]", g.AdminLoginData.IP, c.ClientIP())
			restful.UnLoginErr(c, "IP does not match the last login.")
			c.Abort()
			return
		}

		c.Next()
	}
}
