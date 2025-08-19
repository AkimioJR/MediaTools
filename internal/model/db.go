package model

import (
	"encoding/json"
	"strings"
)

type DBType uint8

const (
	DBTypeUnknown DBType = iota
	DBTypeSQLite
	DBTypeMySQL
	DBTypePostgres
)

func (dbType DBType) String() string {
	switch dbType {
	case DBTypeSQLite:
		return "SQLite"
	case DBTypeMySQL:
		return "MySQL"
	case DBTypePostgres:
		return "Postgres"
	default: // DBTypeUnknown
		return "UnknownDBType"
	}
}

func (dbType *DBType) ParseString(s string) {
	switch strings.ToLower(s) {
	case "sqlite":
		*dbType = DBTypeSQLite
	case "mysql":
		*dbType = DBTypeMySQL
	case "postgres":
		*dbType = DBTypePostgres
	default:
		*dbType = DBTypeUnknown
	}
}

func (dbType DBType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dbType.String() + `"`), nil
}

func (dbType *DBType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	dbType.ParseString(s)
	return nil
}

func (dbType DBType) MarshalYAML() (any, error) {
	return dbType.String(), nil
}

func (dbType *DBType) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	dbType.ParseString(s)
	return nil
}
