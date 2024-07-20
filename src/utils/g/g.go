package g

import "github.com/oimoyu/OimoAdmin/src/utils/_type"

var OimoAdmin *_type.OimoAdminStruct

var Config *_type.ConfigStruct

var SiteSecret string

var AdminPathSecret string
var AdminUsername string
var AdminPassword string
var AdminLoginData struct {
	IP    string
	Token string
}
