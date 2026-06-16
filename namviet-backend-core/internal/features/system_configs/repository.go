package system_configs

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllSystemConfigs() ([]SystemConfig, error) {
	var results []SystemConfig
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetSystemConfigByKey(key string) (*SystemConfig, error) {
	var result SystemConfig
	db := supabase.DB
	err := db.First(&result, "config_key = ?", key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("system config not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateSystemConfig(data *SystemConfig) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateSystemConfig(data *SystemConfig) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteSystemConfig(key string) error {
	db := supabase.DB
	return db.Where("config_key = ?", key).Delete(&SystemConfig{}).Error
}
