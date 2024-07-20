package main

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/oimoyu/OimoAdmin"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type MyModel interface {
}

var MyModels = []MyModel{
	&Order{},
	&User{},
}

type User struct {
	ID       int    `gorm:"primaryKey;not null" json:"id"`
	Username string `gorm:"not null" json:"username"`
}

type Order struct {
	ID         int   `gorm:"primaryKey;not null" json:"id"`
	Status     int   `gorm:"default:0;not null" json:"status"`
	CreateTime int64 `gorm:"index;autoCreateTime;not null" json:"create_time"`
	EndTime    int64 `gorm:"index;not null" json:"end_time"`
}

func GetExecuteDir() string {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	rootDir := filepath.Dir(executable)
	return rootDir
}

var DB *gorm.DB

func initDB() {
	// connect
	dsn := filepath.Join(GetExecuteDir(), "db.db")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db

	// migrate
	for _, model := range MyModels {
		var _ MyModel = model // validate interface
		if err := DB.AutoMigrate(model); err != nil {
			panic(err)
		}
	}
}

func main() {
	initDB()
	r := gin.Default()
	OimoAdmin.Init(r, DB)

	r.Run("0.0.0.0:8098")
}
