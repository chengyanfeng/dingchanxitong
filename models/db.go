package models

type P map[string]interface{}

var Mysqlconn = P{

	"username": "root",
	"password": "1234",
	"host":     "localhost",
	"port":     3306,
	"name":     "test",
}
