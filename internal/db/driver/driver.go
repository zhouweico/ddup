package driver

import (
	"fmt"

	"ddup-apis/internal/config"

	"gorm.io/gorm"
)

type Driver interface {
	Open(cfg *config.Config) (*gorm.DB, error)
	Name() string
}

type Factory struct {
	drivers map[string]Driver
}

func NewFactory() *Factory {
	return &Factory{
		drivers: make(map[string]Driver),
	}
}

func (f *Factory) Register(driver Driver) {
	f.drivers[driver.Name()] = driver
}

func (f *Factory) Get(name string) (Driver, error) {
	if driver, ok := f.drivers[name]; ok {
		return driver, nil
	}
	return nil, fmt.Errorf("数据库驱动 %s 未注册", name)
}
