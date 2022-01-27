package store

import (
	"github.com/moshrank/spacey-backend/pkg/db"
)

const CONFIG_COLLECTION = "configuration"

type ConfigStoreInterface interface {
	GetConfig(configName string) (map[string]interface{}, error)
}

type ConfigStore struct {
	db db.DatabaseInterface
}

func NewConfigStore(db db.DatabaseInterface) ConfigStoreInterface {
	return &ConfigStore{
		db: db,
	}
}

func (s *ConfigStore) GetConfig(configName string) (map[string]interface{}, error) {
	res := s.db.QueryDocument(CONFIG_COLLECTION, map[string]interface{}{"name": configName})
	var config map[string]interface{}
	err := res.Decode(&config)

	return config, err
}
