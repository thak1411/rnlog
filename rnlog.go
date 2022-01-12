package rnlog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 로거 사용 규칙
// 상황에 맞게 로거 레벨을 골라 사용합니다. 로거 레벨은 아래와 같습니다.
// 로거는 Printf를 기반으로 동작하며 문장의 마지막에 \n를 넣어주지 않아도 됩니다.
// Init 함수를 통해 로거를 사용했을 경우 프로그램이 종료해야하는 시점에 로거를 필수로 닫아주어야 합니다.

// 로거 레벨
// [Lv0] Log
// [Lv1] Debug
// [Lv2] Info
// [Lv3] Warn
// [Lv4] Error
// [Lv5] Fatal

type rnLogger struct {
	out      *log.Logger
	err      *log.Logger
	logFile  *os.File
	logLevel int
}

var logger rnLogger

// Init을 호출하지 않았을 때에도 사용할 수 있도록, 콘솔 로그만 남겨주는 로거로 초기화합니다.
func init() {
	WriterOut := io.Writer(os.Stdout)
	WriterErr := io.Writer(os.Stderr)
	logger.out = log.New(WriterOut, "", 0)
	logger.err = log.New(WriterErr, "", 0)
	logger.logLevel = 0
	logger.logFile = nil
}

// 기본 로거를 설정하고 초기화합니다.
// 이 함수를 호출하지 않으면 콘솔 로그만 출력합니다.
func Init(logFilePath string, logFileName string) error {
	if logFilePath == "" {
		logFilePath = "./"
	}
	if logFileName == "" {
		return errors.New("[rnlog error]: logFileName is empty")
	}
	var err error
	err = os.MkdirAll(logFilePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("[rnlog error]: MkdirAll error: %s", err.Error())
	}
	logger.logFile, err = os.OpenFile(filepath.Join(logFilePath, logFileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("[rnlog error]: OpenFile error: %s", err.Error())
	}
	multiWriterOut := io.MultiWriter(os.Stdout, logger.logFile)
	multiWriterErr := io.MultiWriter(os.Stderr, logger.logFile)
	logger.logLevel = 0

	logger.out = log.New(multiWriterOut, "", 0)
	logger.err = log.New(multiWriterErr, "", 0)
	return nil
}

// 로거의 레벨을 설정합니다.
// 설정한 레벨보다 같거나 높은 레벨의 로그만 출력합니다. 기본은 0입니다.
func SetLoggerLevel(level int) {
	logger.logLevel = level
}

// 로그의 맨 앞에 붙을 prefix를 만듭니다.
func makePrefix(level int) string {
	if level < 0 || level > 5 {
		return ""
	}
	levelTag := []string{
		"",
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"FATAL",
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	pc, _, line, ok := runtime.Caller(2)
	if ok {
		spt := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		funcName := spt[len(spt)-1]
		return fmt.Sprintf("[%5s][%s][%s:%d] ", levelTag[level], t, funcName, line)
	}
	return fmt.Sprintf("[%5s][%s][UnkownFunction] ", levelTag[level], t)
}

// [Lv0] prefix가 아무것도 붙지 않는 로거입니다. 구분선 등을 출력할 때 사용합니다.
func Log(format string, v ...interface{}) {
	if logger.logLevel <= 0 {
		logger.out.Printf(format, v...)
	}
}

// [Lv1] 개발시 필요한 정보를 제공하기 위해 사용합니다.
func Debug(format string, v ...interface{}) {
	if logger.logLevel <= 1 {
		logger.out.Printf(makePrefix(1)+format, v...)
	}
}

// [Lv2] 내가 원하는대로 동작하고있는지 정보를 주기위해 사용합니다.
func Info(format string, v ...interface{}) {
	if logger.logLevel <= 2 {
		logger.out.Printf(makePrefix(2)+format, v...)
	}
}

// [Lv3] 예상치 못한 일이 발생하거나, 미래에 발생할 수 있는 일에 대한 경고를 주기위해 사용합니다.
// **이 문제로 프로그램이 죽으면 안됩니다.**
func Warn(format string, v ...interface{}) {
	if logger.logLevel <= 3 {
		logger.out.Printf(makePrefix(3)+format, v...)
	}
}

// [Lv4] 오류로 인해 프로그램이 일부 기능을 수행하지 못할 때 사용합니다.
func Error(format string, v ...interface{}) {
	if logger.logLevel <= 4 {
		logger.err.Printf(makePrefix(4)+format, v...)
	}
}

// [Lv5] 프로그램이 종료될 때 사용합니다.
func Fatal(format string, v ...interface{}) {
	if logger.logLevel <= 5 {
		logger.err.Printf(makePrefix(5)+format, v...)
	}
}

// 로거 파일을 정상적으로 닫고 종료합니다.
func Close() {
	if logger.logFile != nil {
		logger.out.Printf(makePrefix(1) + "Logger closing...")
		logger.logFile.Close()
	}
}
