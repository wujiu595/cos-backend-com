package pgbuilder

import (
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
)

// SQL return sql string and error built from stmt
func SQL(stmt dbr.Builder) (string, error) {
	return dbr.InterpolateForDialect("?", []interface{}{stmt}, dialect.PostgreSQL)
}

// MustSQL return built sql, panic we got err
func MustSQL(stmt dbr.Builder) string {
	s, err := SQL(stmt)
	if err != nil {
		panic(err)
	}

	return s
}
