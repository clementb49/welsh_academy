// Package config contains the configuration settings for the application.
// It imports packages for the Gin web framework, the Viper configuration library, and the Zap logging library.
package config

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Constants to define the modes of the application.
const (
	waModeDev  = "dev"
	waModeProd = "prod"
)

// WaConfig is the main struct that stores the welsh academy configuration.
// It contains fields for database configuration, mode, logging level, and JWT key.
type WaConfig struct {
	DbCfg    *DbConfig
	Mode     string
	LogLevel string
	JwtKey   string
}

// Global variable to store the welsh academy configuration.
var waConfig *WaConfig

// GetWaConfig function retrieves the WaConfig from Viper configuration.
// It sets up Viper to use the environment prefix "wa" and retrieves values for database configuration, mode, logging level, and JWT key.
// If WaConfig has not been initialized, it initializes it and returns it.
func GetWaConfig() *WaConfig {
	if waConfig == nil {
		viper.SetEnvPrefix("wa")
		viper.AutomaticEnv()
		waConfig = &WaConfig{
			DbCfg:    DbConfigFromViper(),
			Mode:     viper.GetString("mode"),
			LogLevel: viper.GetString("logLevel"),
			JwtKey:   viper.GetString("JWTKEY"),
		}
	}
	return waConfig
}

// ConfigureLogger function returns a Zap logging configuration based on the WaConfig settings.
// It creates a new Zap configuration based on the application mode (dev or prod).
// If the WaConfig has a Log Level set, it sets the level on the Zap configuration.
func (cfg *WaConfig) ConfigureLogger() zap.Config {
	var lcfg zap.Config
	switch waConfig.Mode {
	case waModeDev:
		lcfg = zap.NewDevelopmentConfig()
	case waModeProd:
		lcfg = zap.NewProductionConfig()
	default:
		lcfg = zap.NewProductionConfig()
	}
	if cfg.LogLevel != "" {
		level, err := zap.ParseAtomicLevel(cfg.LogLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error when configuring logger: %s", err)
		}
		lcfg.Level = level
	}
	return lcfg
}

// ConfigureGin function sets the Gin mode based on the WaConfig settings.
// It sets the Gin mode to Debug or Release mode based on the WaConfig mode setting.
func (cfg *WaConfig) ConfigureGin() {
	switch cfg.Mode {
	case waModeDev:
		gin.SetMode(gin.DebugMode)
		break
	case waModeProd:
		gin.SetMode(gin.ReleaseMode)
		break
	default:
		gin.SetMode(gin.ReleaseMode)
		break
	}
}
