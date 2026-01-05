package pgsql

import "sync"

type Pgsql struct {
	Host                 string `validate:"required,ip"`
	Port                 int    `validate:"required,gt=0,lte=65535"`
	Username             string `validate:"required"`
	Password             string `validate:"required"`
	Timezone             string `validate:"required"`
	PreferSimpleProtocol bool   `validate:"required"`
}

type Database struct {
	db sync.Map
}

var pg Database

// NewPgsql TODO: implement pgsql connection and methods
func NewPgsql(databases map[string]Pgsql) {
	for dbname, v := range databases {
		pg.db.Store(dbname, v)
	}
}
