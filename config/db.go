package config

import (
	"log"

	"github.com/gocroot/helper/atdb"
)

var MongoString string = "mongodb+srv://idbiz:OTBmC4Bcs9AyUdjw@idbiz.bphmr.mongodb.net/"

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "idbizdevelop",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)

func init() {
	if ErrorMongoconn != nil {
		log.Printf("Failed to connect to MongoDB (idbizdevelop): %v\n", ErrorMongoconn)
	} else {
		log.Println("Successfully connected to MongoDB (idbizdevelop)")
	}
}

// Geospacial Database
var MongoStringGeo string = "mongodb+srv://idbiz:OTBmC4Bcs9AyUdjw@idbiz.bphmr.mongodb.net/"

var mongoinfoGeo = atdb.DBInfo{
	DBString: MongoStringGeo,
	DBName:   "geo",
}

var MongoconnGeo, ErrorMongoconnGeo = atdb.MongoConnect(mongoinfoGeo)
