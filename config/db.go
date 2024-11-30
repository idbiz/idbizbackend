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

// Geospacial Database
var MongoStringGeo string = "mongodb+srv://idbiz:idbizcroot@geoidbiz.z4wo4.mongodb.net/?retryWrites=true&w=majority&appName=geoidbiz"

var mongoinfoGeo = atdb.DBInfo{
	DBString: MongoStringGeo,
	DBName:   "geo",
}

var MongoconnGeo, ErrorMongoconnGeo = atdb.MongoConnect(mongoinfoGeo)
