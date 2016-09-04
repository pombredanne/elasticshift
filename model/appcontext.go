package model

import "github.com/spf13/viper"

// AppContext that holds the application level context information
// such as DBConnection, configuration ...
type AppContext struct {
	//db connection
	//global app values
	Config *viper.Viper
}
