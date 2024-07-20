package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_const "github.com/oimoyu/OimoAdmin/src/utils/const"
	"github.com/oimoyu/OimoAdmin/src/utils/fail2ban"
	"github.com/oimoyu/OimoAdmin/src/utils/functions"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
)

func Login(c *gin.Context) {
	if !fail2ban.IsIPValid(c.ClientIP()) {
		restful.ParamErr(c, "Your IP address is banned")
		return
	}

	var requestData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindBodyWith(&requestData, binding.JSON); err != nil {
		restful.ParamErr(c, fmt.Sprintf("parse input params failed: %v", err))
		return
	}

	if g.AdminUsername != requestData.Username || g.AdminPassword != requestData.Password {
		fail2ban.IncrementIPAttempts(c.ClientIP())
		remainTimes := fail2ban.RemainingAttempts(c.ClientIP())
		restful.ParamErr(c, fmt.Sprintf("Invalid Login Credentials, remain retry times: %d", remainTimes))
		return
	}

	tokenDataMap := map[string]interface{}{
		"role": "ADMIN",
	}
	token, err := functions.GenerateToken(tokenDataMap, g.SiteSecret, _const.TokenDuration)
	if err != nil {
		restful.ParamErr(c, "Failed to generate token")
		return
	}

	returnData := map[string]interface{}{
		//"token": token,
	}

	g.AdminLoginData.IP = c.ClientIP()
	g.AdminLoginData.Token = token

	c.SetCookie("admin_token", token, 0, "/", "", false, true)

	fail2ban.ResetIPAttempts(c.ClientIP())
	g.OimoAdmin.Logger.FileLog(fmt.Sprintf("%s: Login", c.ClientIP()))

	restful.Ok(c, returnData)
}
