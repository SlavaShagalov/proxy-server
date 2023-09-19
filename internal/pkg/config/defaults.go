package config

import (
	"github.com/spf13/viper"
)

func SetDefaultMongoConfig() {
	viper.SetDefault(MongoHost, "storage")
	//viper.SetDefault(MongoHost, "localhost")
	viper.SetDefault(MongoPort, 27017)
	viper.SetDefault(MongoDB, "mitmdb")
}
