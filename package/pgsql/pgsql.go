package pgsql

import "sync"

type Pgsql struct {
	Host                 string
	Port                 int
	Username             string
	Password             string
	Timezone             string
	PreferSimpleProtocol bool
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
