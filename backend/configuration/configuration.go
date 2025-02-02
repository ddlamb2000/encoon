// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package configuration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"d.lambert.fr/encoon/utils"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	AppName       string                   `yaml:"appName"`
	Log           bool                     `yaml:"log"`
	Trace         bool                     `yaml:"trace"`
	ShowTiming    bool                     `yaml:"showTiming"`
	GridCacheSize int                      `yaml:"gridCacheSize"`
	SeedDataFile  string                   `yaml:"seedDataFile"`
	Databases     []*DatabaseConfiguration `yaml:"database"`
	Kafka         KafkaConfiguration       `yaml:"kafka"`
	JwtExpiration int                      `yaml:"jwtExpiration"`
	AI            AIConfiguration          `yaml:"ai"`
}

type AIConfiguration struct {
	Model             string  `yaml:"model"`
	ApiKeyFile        string  `yaml:"apiKeyFile"`
	Temperature       float32 `yaml:"temperature"`
	TopP              float32 `yaml:"topP"`
	TopK              int32   `yaml:"topK"`
	MaxOutputTokens   int32   `yaml:"maxOutputTokens"`
	SystemInstruction string  `yaml:"systemInstruction"`
}

type KafkaConfiguration struct {
	Brokers     string `yaml:"brokers"`
	GroupID     string `yaml:"groupID"`
	TopicPrefix string `yaml:"topicPrefix"`
}

type DatabaseConfiguration struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	Name             string `yaml:"name"`
	Role             string `yaml:"role"`
	RolePassword     string `yaml:"rolePassword"`
	JwtSecret        string `yaml:"jwtsecret"`
	Root             string `yaml:"root"`
	RootPassword     string `yaml:"rootPassword"`
	TestSleepTime    int    `yaml:"testSleepTime"`
	TimeOutThreshold int    `yaml:"timeOutThreshold"`
}

var (
	appConfigurationMutex sync.RWMutex
	appConfiguration      Configuration
	configurationFileName string
	configurationHash     string
)

func GetConfiguration() Configuration {
	appConfigurationMutex.RLock()
	defer appConfigurationMutex.RUnlock()
	return appConfiguration
}

func LoadConfiguration(fileName string) error {
	appConfigurationMutex.Lock()
	defer appConfigurationMutex.Unlock()
	configurationFileName = fileName
	return loadConfigurationFromFile()
}

func loadConfigurationFromFile() error {
	Log("", "", "Loading configuration from %v.", configurationFileName)
	f, err := os.ReadFile(configurationFileName)
	if err != nil {
		return LogAndReturnError("", "", "Error loading configuration from file %q: %v.", configurationFileName, err)
	}
	hash, _ := utils.CalculateFileHash(configurationFileName)
	newConfiguration := new(Configuration)
	if err = yaml.Unmarshal(f, &newConfiguration); err != nil {
		return LogAndReturnError("", "", "Error parsing configuration from file %q: %v.", configurationFileName, err)
	}
	if err = validateConfiguration(newConfiguration); err != nil {
		return err
	}
	appConfiguration = *newConfiguration
	configurationHash = hash
	Log("", "", "Configuration loaded from file %q.", configurationFileName)
	return nil
}

func validateConfiguration(conf *Configuration) error {
	if conf.AppName == "" {
		return LogAndReturnError("", "", "Missing application name (appName) from configuration file %v.", configurationFileName)
	}
	if conf.JwtExpiration == 0 {
		return LogAndReturnError("", "", "Missing expiration (jwtExpiration) from configuration file %v.", configurationFileName)
	}
	Log("", "", "Configuration from %v is valid.", configurationFileName)
	return nil
}

func IsDatabaseEnabled(dbName string) bool {
	return dbName != "" && GetDatabaseConfiguration(dbName) != nil
}

func GetDatabaseConfiguration(dbName string) *DatabaseConfiguration {
	for _, dbConfig := range appConfiguration.Databases {
		if dbConfig.Name == dbName {
			return dbConfig
		}
	}
	return nil
}

func GetJWTSecret(dbName string) []byte {
	if !IsDatabaseEnabled(dbName) {
		return nil
	}
	return []byte(dbName + GetDatabaseConfiguration(dbName).JwtSecret)
}

func GetRootAndPassword(dbName string) (string, string) {
	if !IsDatabaseEnabled(dbName) {
		return "", ""
	}
	dbConfiguration := GetDatabaseConfiguration(dbName)
	return dbConfiguration.Root, dbConfiguration.RootPassword
}

func GetContextWithTimeOut(ct context.Context, dbName string) (context.Context, context.CancelFunc) {
	threshold := 0
	dbConfiguration := GetDatabaseConfiguration(dbName)
	if dbConfiguration != nil {
		threshold = dbConfiguration.TimeOutThreshold
	}
	if threshold < 10 {
		threshold = 10
	}
	ctx, ctxFunc := context.WithTimeout(ct, time.Duration(threshold)*time.Millisecond)
	return ctx, ctxFunc
}

func WatchConfigurationChanges(fileName string) {
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			newHash, _ := utils.CalculateFileHash(fileName)
			if newHash != configurationHash {
				configurationHash = newHash
				LoadConfiguration(fileName)
			}
		}
	}()
}

func Log(dbName, userName, format string, a ...any) {
	if appConfiguration.Log {
		fmt.Printf(getLogPrefix(dbName, userName)+format+"\n", a...)
	}
}

func LogError(dbName, userName, format string, a ...any) {
	fmt.Printf(getLogPrefix(dbName, userName)+"[ERROR] "+format+"\n", a...)
}

func LogAndReturnError(dbName, userName, format string, a ...any) error {
	m := fmt.Sprintf(format, a...)
	LogError(dbName, userName, m)
	return errors.New(m)
}

func Trace(dbName, userName, format string, a ...any) {
	if appConfiguration.Trace {
		fmt.Printf(getLogPrefix(dbName, userName)+"[TRACE] "+format+"\n", a...)
	}
}

func getLogPrefix(dbName, userName string) string {
	if userName != "" {
		return "[" + appConfiguration.AppName + "] [" + dbName + "] [" + userName + "] "
	} else if dbName != "" {
		return "[" + appConfiguration.AppName + "] [" + dbName + "] "
	} else if appConfiguration.AppName != "" {
		return "[" + appConfiguration.AppName + "] "
	}
	return ""
}

func StartTiming() time.Time {
	return time.Now()
}

func StopTiming(dbName, userName, funcName string, start time.Time) {
	if appConfiguration.ShowTiming {
		duration := time.Since(start)
		fmt.Printf(getLogPrefix(dbName, userName)+"[TIMING] "+funcName+" - %v\n", duration)
	}
}

func GetSeedDataFile() string {
	return appConfiguration.SeedDataFile
}
