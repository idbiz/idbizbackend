package config

import (
	"github.com/gocroot/helper/atdb"
)

var MongoString string = "mongodb+srv://idbiz:OTBmC4Bcs9AyUdjw@idbiz.bphmr.mongodb.net/"

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "idbizdevelop",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
