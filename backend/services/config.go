package services


type configs struct {
	// Add configuration fields here
}
type ConfigService struct {
	fileName string
	config   configs

}

func NewConfigService(fileName string) *ConfigService {
	return &ConfigService{
		fileName: fileName,
		config:   configs{},
	}
}

func (c *ConfigService) GetConfig() configs {
	return c.config
}

func (c *ConfigService) LoadConfig() {
}

func (c *ConfigService) SaveConfig() {
}

func (c *ConfigService) EditConfig() {
}

func (c *ConfigService) ValidateConfig() bool {
	return true
}
