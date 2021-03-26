package settings

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	FscLocation  string
	FmsMasterURL []string
	MaxParallel  int
}

type loggerOutput struct {
	StdOut io.Writer
	StdErr io.Writer

	Logger *zap.Logger
}

var AppSettings = &Application{}
var ServerSettings = &Server{}
var FscSettings = &FSC{MaxParallel: 1}
var Log *loggerOutput

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

	if FscSettings.FscLocation == "" {
		fmsHome := os.Getenv("FMS_HOME")
		if fmsHome == "" {
			log.Fatalln("FMS tools location is not specified. It must be specified by the FMS_HOME environment variable or in the conf/settings.ini FscLocation parameter")
		} else {
			FscSettings.FscLocation = fmsHome
		}
	}

	AppSettings.JavaHome = os.Getenv("JAVA_HOME")
	if AppSettings.JavaHome == "" {
		log.Println("JAVA_HOME is not defined")
	}

	Log = initLog()

	if ServerSettings.RunMode == "debug" {
		PrintDebugConfig()
	}
}

func initLog() *loggerOutput {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Error while initializing DEBUG logger")
	}
	output := &loggerOutput{
		StdOut: os.Stdout,
		StdErr: os.Stderr,
		Logger: l,
	}
	if ServerSettings.RunMode == "release" {
		_, err := os.ReadDir("logs")
		if err != nil {
			switch err.(type) {
			default:
				log.Printf("Error accessing log folder: %v", err)
			case *fs.PathError:
				os.Mkdir("logs", 0755)
			}

		}
		logFile, err := os.Create("logs" + string(os.PathSeparator) + fmt.Sprintf("tcfmshelper-%s-out.log", time.Now().Format("060102_030405")))
		logErrFile, err := os.Create("logs" + string(os.PathSeparator) + fmt.Sprintf("tcfmshelper-%s-err.log", time.Now().Format("060102_030405")))
		if err != nil {
			log.Fatalln("Error while creating log file")
		}

		output.StdOut = io.MultiWriter(os.Stdout, logFile)
		output.StdErr = io.MultiWriter(os.Stderr, logErrFile)

		core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(output.StdOut), zap.InfoLevel)
		output.Logger = zap.New(core)
	}

	output.Logger.Debug("Logger initialized")

	return output
}

//PrintDebugConfig prints the configuration if debug is enabled
func PrintDebugConfig() {
	Log.Logger.Sugar().Debugf("Server config loaded: \n  -- RunMode: %s\n  -- Port: %d\n  -- ReadTimeout: %v\n  -- WriteTimeout: %v\n", ServerSettings.RunMode, ServerSettings.Port, ServerSettings.ReadTimeout, ServerSettings.WriteTimeout)
	Log.Logger.Sugar().Debugf("FSC config loaded: \n  -- FscLocation: %s\n  -- FmsMasterURL: %#v\n  -- MaxParallel: %d\n", FscSettings.FscLocation, FscSettings.FmsMasterURL, FscSettings.MaxParallel)
}
