package system_configs

import "time"

func GetAllSystemConfigsService() ([]SystemConfig, error) {
	return GetAllSystemConfigs()
}

func GetSystemConfigByKeyService(key string) (*SystemConfig, error) {
	return GetSystemConfigByKey(key)
}

func CreateSystemConfigService(req CreateSystemConfigRequest) (*SystemConfig, error) {
	config := &SystemConfig{
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		Description: req.Description,
		UpdatedBy:   req.UpdatedBy,
	}

	if err := CreateSystemConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

func UpdateSystemConfigService(key string, req UpdateSystemConfigRequest) (*SystemConfig, error) {
	config, err := GetSystemConfigByKey(key)
	if err != nil {
		return nil, err
	}

	if req.ConfigValue != nil {
		config.ConfigValue = *req.ConfigValue
	}
	if req.Description != nil {
		config.Description = req.Description
	}
	if req.UpdatedBy != nil {
		config.UpdatedBy = req.UpdatedBy
	}

	now := time.Now()
	config.UpdatedAt = &now

	if err := UpdateSystemConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

func DeleteSystemConfigService(key string) error {
	return DeleteSystemConfig(key)
}
