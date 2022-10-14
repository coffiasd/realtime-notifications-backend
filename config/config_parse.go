package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var ParseConfig Conf
func Parse(configFilePath string){
	if _,err := toml.DecodeFile(configFilePath,&ParseConfig);err!=nil{
		fmt.Println(err)
		return 
	}
}