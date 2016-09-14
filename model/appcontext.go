package model

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// AppContext that holds the application level context information
// such as DBConnection, configuration, logger ...
type AppContext struct {

	//db connection
	DB *gorm.DB

	//global app values
	Config *viper.Viper

	//logger
}
