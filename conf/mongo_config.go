package conf

import "go.mongodb.org/mongo-driver/mongo/options"

type MongoCfg struct {
	URI            string
	DB             string
	Collection     string
	ConnectTimeout int
	ReadTimeout    int
	FindOpt        *options.FindOptions
}
