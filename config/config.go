package config

import "fmt"

type ConfigDB struct {
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}

func (c *ConfigDB) GetDataSorceConnectionName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.DB_USER, c.DB_PASSWORD, c.DB_HOST, c.DB_PORT, c.DB_NAME)
}
