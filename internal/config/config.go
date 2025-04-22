// config - parse date from file
package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/viper"

	"github.com/Ekvo/yandex-practicum-go-final-project/pkg/common"
)

var (
	ErrConfigEmpty = errors.New("empty")

	// ErrConfigPortNoNumeric - if the expected field does not contain only positive numbers
	ErrConfigPortNoNumeric = errors.New(" no numeric")
)

type Config struct {
	DataBaseDataSourceName string `mapstructure:"TODO_DBFILE"`

	ServerPort string `mapstructure:"TODO_PORT"`

	// data for taskServcie see (/internal/services/usecase/taslcase.go)
	TaskNextDate string `mapstructure:"ALGORITHM_TASK_DATE"`

	// loginServcie see (/internal/services/usecase/logincase.go)
	UserPassword string `mapstructure:"TODO_PASSWORD"`

	// secret key for jwt.Token -> (/internal/lib/jwtsign/jwtsign.go)
	JWTSecretKey string `mapstructure:"TODO_SECRET_KEY"`

	// path for file executeble in application
	PathFilesWeb string `mapstructure:"PATH_DIR_WEB"`

	// options - contain data about the file being analyzed (parse) see (internal/config/options.go)
	options
}

// NewConfig - loads config and check valid of fileds
//
// path - empty -> use ENV variables
// (only image in docker)
// in main.go need -> cfg, err := config.NewConfig("")
//
// path - exist -> parse file
// (only localy)
// in main.go need -> cfg, err := config.NewConfig("./init/.env")
func NewConfig(path string) (*Config, error) {
	cfg := &Config{options: options{pathOfFile: path}}
	if err := cfg.parsePath(); err != nil && !errors.Is(err, ErrOptionsEmptyFile) {
		return nil, err
	}
	if err := cfg.loadConfig(); err != nil {
		return nil, err
	}
	return cfg, cfg.ValidConfig()
}

// loadConfig - read file and fill the fields
func (cfg *Config) loadConfig() error {
	if cfg.pathOfFile != "" {
		data, err := os.ReadFile(cfg.pathOfFile)
		if err != nil {
			return fmt.Errorf("config: read file error - %w", err)
		}
		return cfg.setConfig(bytes.NewBuffer(data))
	}
	return cfg.setConfig(nil)
}

// envNames - contain names all 'ENV'
var envNames = []string{
	"TODO_PORT",
	"TODO_DBFILE",
	"ALGORITHM_TASK_DATE",
	"TODO_PASSWORD",
	"TODO_SECRET_KEY",
	"PATH_DIR_WEB",
}

// setConfig - set extension of parse file from 'options'
func (cfg *Config) setConfig(in io.Reader) error {
	if in == nil {
		viper.AutomaticEnv()
		for _, env := range envNames {
			if err := viper.BindEnv(env); err != nil {
				return fmt.Errorf("config: read config error - %w", err)
			}
		}
	} else {
		viper.SetConfigType(cfg.fileExt)
		if err := viper.ReadConfig(in); err != nil {
			return fmt.Errorf("config: read config error - %w", err)
		}
	}
	return viper.Unmarshal(cfg)
}

// ValidConfig - Checking the config for validity
//
// create common.Message - for create full infor about possible erros -> see (pkg/common/common.go)
func (cfg *Config) ValidConfig() error {
	msgErr := make(common.Message)
	cfg.validDataBase(msgErr)
	cfg.validServe(msgErr)
	cfg.validTask(msgErr)
	//cfg.validPassword(msgErr)
	cfg.validJWT(msgErr)
	cfg.validPathOfFiles(msgErr)
	if len(msgErr) > 0 {
		return fmt.Errorf("config: invalid config - %s", msgErr.String())
	}
	return nil
}

func (cfg *Config) validDataBase(msgErr common.Message) {
	if cfg.DataBaseDataSourceName == "" {
		msgErr["source"] = ErrConfigEmpty
	}
}

func (cfg *Config) validServe(msgErr common.Message) {
	if port, err := strconv.Atoi(cfg.ServerPort); err != nil || port < 1 {
		msgErr["port"] = ErrConfigPortNoNumeric
	}
}

func (cfg *Config) validTask(msgErr common.Message) {
	if cfg.TaskNextDate == "" {
		msgErr["task-nextdate"] = ErrConfigEmpty
	}
}

func (cfg *Config) validPassword(msgErr common.Message) {
	if cfg.UserPassword == "" {
		msgErr["password"] = ErrConfigEmpty
	}
}

func (cfg *Config) validJWT(msgErr common.Message) {
	if cfg.JWTSecretKey == "" {
		msgErr["jwtsign-secret-key"] = ErrConfigEmpty
	}
}

func (cfg *Config) validPathOfFiles(msgErr common.Message) {
	if cfg.PathFilesWeb == "" {
		msgErr["dir-web"] = ErrConfigEmpty
	}
}
