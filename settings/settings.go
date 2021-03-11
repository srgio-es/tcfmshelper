package settings

import (
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type Application struct {
	JavaHome string
}

type Server struct {
	RunMode      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type FSC struct {
	FscLocation string
	FscFromUrl  string
}

var AppSettings = &Application{}
var ServerSettings = &Server{}
var FscSettings = &FSC{}

var cfg *ini.File

//Setup sets the configured application parameters
func Setup() {
	var err error

	cfg, err = ini.Load("conf/settings.ini")
	if err != nil {
		log.Fatalf("failed to load settings. Error parsing 'conf/settings.ini': %v", err)
	}

	err = cfg.Section("server").MapTo(ServerSettings)
	if err != nil {
		log.Fatalf("Failed to map server settings from ini file: %v", err)
	}

	err = cfg.Section("fsc").MapTo(FscSettings)
	if err != nil {
		log.Fatalf("Failed to map fsc settings from ini file: %v", err)
	}

	AppSettings.JavaHome = os.Getenv("JAVA_HOME")
	if AppSettings.JavaHome == "" {
		log.Println("JAVA_HOME is not defined")
	}

	if ServerSettings.RunMode == "debug" {
		PrintDebugConfig()
	}
}

//PrintDebugConfig prints the configuration if debug is enabled
func PrintDebugConfig() {
	log.Printf("Server config loaded: \n  -- RunMode: %s\n  -- Port: %d\n  -- ReadTimeout: %v\n  -- WriteTimeout: %v\n", ServerSettings.RunMode, ServerSettings.Port, ServerSettings.ReadTimeout, ServerSettings.WriteTimeout)
	log.Printf("FSC config loaded: \n  -- FscLocation: %s\n  -- FscFromUrl: %s\n", FscSettings.FscLocation, FscSettings.FscFromUrl)
}
