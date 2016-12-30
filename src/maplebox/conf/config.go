package conf

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/jinzhu/configor.v0"
	"path/filepath"
	"runtime"
	"time"
)

type Config struct {
	AppMode      string
	AppDeveloper string
	AppName      string `default:"maplebox"`
	AppPath      string
	Server       *ServerConfig
	Log          *LogConfig
}

type ServerConfig struct {
	Host           string        `default:""`
	Port           uint          `default:3210`
	StaticPath     string        `default:"static"`
	ViewsPath      string        `default:"views"`
	ReadTimeout    time.Duration `default:10`
	WriteTimeout   time.Duration `default:10`
	MaxHeaderBytes int           `default:1`
}

type LogConfig struct {
	Level              string `default:"debug"`
	Formatter          string `default:"json"`
	WriteFile          bool   `default:false`
	LogDir             string `default:"."`
	AppLogFileFormat   string `default:"app.log"`
	RouteLogFileFormat string `default:"route.log"`
	LocalTime          bool   `default:true`
	MaxAge             int    `default:30`
	MaxBackups         int    `default:50`
	MaxSize            int64  `default:1048576`
}

var (
	AppConfigFile string
	AppConfig     *Config
)

func init() {
	var err error
	AppConfig, err = newConfig()
	if err != nil {
		panic(err)
	}
}

func newConfig() (*Config, error) {
	appPath, err := appPath()
	appMode := *flag.String("app_mode", "dev", "app mode")
	appDeveloper := *flag.String("app_developer", "", "app developer")

	return &Config{
		AppPath:      appPath,
		AppMode:      appMode,
		AppDeveloper: appDeveloper,
	}, err
}

func InitConfig() {
	AppConfigFile = appConfigFilePath(AppConfig.AppPath)

	if err := configor.Load(AppConfig, AppConfigFile); err != nil {
		panic(err)
	}
}

func appPath() (appPath string, err error) {
	if _, file, _, ok := runtime.Caller(1); ok {
		appPath, err = filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	} else {
		err = errors.New("Get App Path ErrorTest")
	}
	return
}

func appConfigFilePath(appPath string) string {
	appConfigName := "app"
	if AppConfig.AppMode != "" {
		appConfigName += "." + AppConfig.AppMode
	}
	if AppConfig.AppDeveloper != "" {
		appConfigName += "." + AppConfig.AppDeveloper
	}
	appConfigName += ".conf"
	appConfigFile := filepath.Join(appPath, "conf", appConfigName)
	return appConfigFile
}

func absPath(appPath string, path string) string {
	if filepath.IsAbs(path) == false {
		return filepath.Join(appPath, path)
	}
	return path
}

func separator(str string) string {
	return string(filepath.Separator) + str
}

func GetStaticPath() string {
	return absPath(AppConfig.AppPath, AppConfig.Server.StaticPath)
}

func GetViewsPath() string {
	return absPath(AppConfig.AppPath, AppConfig.Server.ViewsPath)
}

func GetAddr() string {
	return fmt.Sprintf("%s:%d", AppConfig.Server.Host, AppConfig.Server.Port)
}

func GetReadTimeout() time.Duration {
	return AppConfig.Server.ReadTimeout * time.Second
}

func GetWriteTimeout() time.Duration {
	return AppConfig.Server.WriteTimeout * time.Second
}

func GetMaxHeaderBytes() int {
	return AppConfig.Server.MaxHeaderBytes * (1 << 20)
}
