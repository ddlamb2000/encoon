// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package configuration

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

func TestLoadConfiguration1(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	if err := LoadConfiguration(fileName); err != nil {
		t.Errorf("Can't load configuration %q: %v.", fileName, err)
	}
}

func TestLoadConfiguration2(t *testing.T) {
	fileName := "../xxx/validConfiguration1.yml"
	if err := LoadConfiguration(fileName); err == nil {
		t.Errorf("Can load configurations from %q: %v!", fileName, err)
	}
}

func TestLoadConfiguration3(t *testing.T) {
	secret := GetJWTSecret("xxx")
	if secret != nil {
		t.Errorf("Invalid Jwt secret: %q.", secret)
	}
}

func TestLoadMainConfiguration4(t *testing.T) {
	fileName := "configuration.go"
	if err := LoadConfiguration(fileName); err == nil {
		t.Errorf("Can load configuration %q.", fileName)
	}
}

func TestLoadConfiguration5(t *testing.T) {
	fileName := "../xxxx.yml"
	if err := LoadConfiguration(fileName); err == nil {
		t.Errorf("Can load configuration %q.", fileName)
	}
}

func TestLoadConfiguration6(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	if err := LoadConfiguration(fileName); err != nil {
		t.Errorf("Can't load configurations from %q: %v.", fileName, err)
	}

	dbName := "test"
	if !IsDatabaseEnabled(dbName) {
		t.Errorf("Database %q isn't enabled.", dbName)
	}

	secret := GetJWTSecret(dbName)
	expect := []byte{
		116, 101, 115, 116, 36, 50, 97, 36, 48, 56, 36,
		100, 99, 110, 50, 50, 118, 82, 70, 73, 90, 109,
		121, 119, 118, 100, 89, 66, 70, 118, 53, 121, 79,
		82, 99, 79, 105, 71, 85, 46, 90, 116, 113, 66,
		57, 83, 49, 100, 84, 99, 120, 115, 112, 86, 108,
		122, 97, 101, 108, 109, 90, 85, 80, 97,
	}
	if !reflect.DeepEqual(secret, expect) {
		t.Errorf("Jwt secret is wrong: %v instead of %v.", secret, expect)
	}

	root, password := GetRootAndPassword(dbName)
	expectPassword := "$2a$08$40D/LcEidSirsqMSQcfc9.DAPTBOpPBelNik5.ppbLwSodxczbNWa"
	if root != "root" || password != expectPassword {
		t.Errorf("Root or password isn't correct for database %q: %q and %q.", dbName, root, password)
	}
}

func TestGetRootAndPassword(t *testing.T) {
	root, password := GetRootAndPassword("xxx")
	if root != "" || password != "" {
		t.Errorf("Root or password isn't correct for database %q: %q and %q.", "xxx", root, password)
	}
}

func TestValidateConfiguration1(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	appConfiguration.AppName = ""
	got := validateConfiguration(&appConfiguration)
	expect := "Missing application name"
	if got == nil || !strings.Contains(got.Error(), expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestValidateConfiguration4(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	appConfiguration.JwtExpiration = 0
	got := validateConfiguration(&appConfiguration)
	expect := "Missing expiration"
	if got == nil || !strings.Contains(got.Error(), expect) {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestValidateConfiguration5(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	got := validateConfiguration(&appConfiguration)
	if got != nil {
		t.Errorf("Got error: %q.", got)
	}
}

func TestInvalidConfiguration(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	if GetConfiguration().AppName != "valid 1" {
		t.Errorf("Configuration 1 doesn't have the expected name: %v.", appConfiguration)
	}
	fileName = "../testData/invalidConfiguration.yml"
	LoadConfiguration(fileName)
	if GetConfiguration().AppName != "valid 1" {
		t.Errorf("Configuration 2 doesn't have the expected name: %v.", appConfiguration)
	}
}

func TestReloadConfiguration(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	if GetConfiguration().AppName != "valid 1" {
		t.Errorf("Configuration 1 doesn't have the expected name: %v.", appConfiguration)
	}
	fileName = "../testData/validConfiguration2.yml"
	LoadConfiguration(fileName)
	if GetConfiguration().AppName != "valid 2" {
		t.Errorf("Configuration 2 doesn't have the expected name: %v.", appConfiguration)
	}
}

func TestGetContextWithTimeOut1(t *testing.T) {
	ctx, cancel := GetContextWithTimeOut(context.Background(), "test")
	defer cancel()
	_, ok := ctx.Deadline()
	if !ok {
		t.Errorf("Context isn't set with a deadline: %v.", ctx)
	}
}

func TestGetContextWithTimeOut2(t *testing.T) {
	ctx, cancel := GetContextWithTimeOut(context.Background(), "xxxx")
	defer cancel()
	_, ok := ctx.Deadline()
	if !ok {
		t.Errorf("Context isn't set with a deadline: %v.", ctx)
	}
}

func TestConfigurationAutoUpdates(t *testing.T) {
	fileName := "/tmp/testConfiguration.yml"
	conf := Configuration{
		AppName:       "testA",
		JwtExpiration: 10,
	}
	out, err := yaml.Marshal(conf)
	if err != nil {
		t.Errorf("Can't marshal yaml: %v.", err)
		return
	}
	err = os.WriteFile(fileName, out, 0666)
	if err != nil {
		t.Errorf("Can't create file %q: %v.", fileName, err)
		return
	}
	WatchConfigurationChanges(fileName)
	time.Sleep(2 * time.Second)
	got := GetConfiguration().AppName
	expect := "testA"
	if got != expect {
		t.Errorf("Application name is %q while it should be: %q.", got, expect)
		return
	}
	conf = Configuration{
		AppName:       "testB",
		JwtExpiration: 10,
	}
	out, err = yaml.Marshal(conf)
	if err != nil {
		t.Errorf("Can't marshal yaml: %v.", err)
		return
	}
	err = os.WriteFile(fileName, out, 0666)
	if err != nil {
		t.Errorf("Can't create file %q: %v.", fileName, err)
		return
	}
	time.Sleep(2 * time.Second)
	got = GetConfiguration().AppName
	expect = "testB"
	if got != expect {
		t.Errorf("Application name is %q while it should be: %q.", got, expect)
		return
	}
}

func TestLog(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	previous := appConfiguration.Trace
	appConfiguration.Log = true
	Log("x", "y", "Test: %v", "test")
	appConfiguration.Log = previous
}

func TestLogError(t *testing.T) {
	LogError("x", "y", "Test: %v", "test")
}

func TestLogAndReturnError(t *testing.T) {
	got := LogAndReturnError("x", "y", "Test: %v", "test")
	expect := "Test: test"
	if got.Error() != expect {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestTrace(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	previous := appConfiguration.Trace
	appConfiguration.Trace = true
	Trace("x", "y", "Test: %v", "test")
	appConfiguration.Trace = previous
}

func TestGetLogPrefix(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	got := getLogPrefix("test", "root")
	expect := "[valid 1] [test] [root] "
	if got != expect {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestGetLogPrefix2(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	got := getLogPrefix("test", "")
	expect := "[valid 1] [test] "
	if got != expect {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestGetLogPrefix3(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	got := getLogPrefix("", "")
	expect := "[valid 1] "
	if got != expect {
		t.Errorf("Got %q instead of %q.", got, expect)
	}
}

func TestTiming(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	previous := appConfiguration.ShowTiming
	appConfiguration.ShowTiming = true
	s := StartTiming()
	StopTiming("test", "root", "test()", s)
	appConfiguration.ShowTiming = previous
}

func TestGetSeedDataFile(t *testing.T) {
	fileName := "../testData/validConfiguration1.yml"
	LoadConfiguration(fileName)
	got := GetSeedDataFile()
	expect := "../seedData.json"
	if got != expect {
		t.Errorf("Got %v instead of %v.", got, expect)
	}
}
