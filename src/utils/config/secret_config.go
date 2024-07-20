package config

import (
	"fmt"
	_const "github.com/oimoyu/OimoAdmin/src/utils/const"
	"github.com/oimoyu/OimoAdmin/src/utils/functions"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"os"
	"path/filepath"
)

func LoadAllSecret() {
	wwwRootDir := functions.GetEnvDir()

	siteSecretFilePath := filepath.Join(wwwRootDir, ".site_secret")
	LoadSecret(siteSecretFilePath, &g.SiteSecret, 32)

	adminPathSecretFilePath := filepath.Join(wwwRootDir, ".admin_path_secret")
	LoadSecret(adminPathSecretFilePath, &g.AdminPathSecret, 12)

	adminUsernameFilePath := filepath.Join(wwwRootDir, ".admin_username")
	LoadSecret(adminUsernameFilePath, &g.AdminUsername, 12)

	adminPasswordFilePath := filepath.Join(wwwRootDir, ".admin_password")
	LoadSecret(adminPasswordFilePath, &g.AdminPassword, 12)

	fmt.Println("\n==================================================================")
	fmt.Printf("\033[32m%s Login Info\033[0m\n", _const.OimoAdminString)
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("path: %s\n", fmt.Sprintf("%s/%s/", _const.OimoAdminStaticPrefix, g.AdminPathSecret))
	fmt.Printf("username: %s\n", g.AdminUsername)
	fmt.Printf("password: %s\n", g.AdminPassword)
	fmt.Println("==================================================================\n")
}

func LoadSecret(secretFilePath string, secretStorage *string, secretNum int) {
	if _, err := os.Stat(secretFilePath); os.IsNotExist(err) {
		secret := functions.GenerateRandomString(secretNum)
		err = os.WriteFile(secretFilePath, []byte(secret), 0600)
		if err != nil {
			panic(err)
		}
	}

	secret, err := os.ReadFile(secretFilePath)
	if err != nil {
		panic(err)
	}
	*secretStorage = string(secret)
}
