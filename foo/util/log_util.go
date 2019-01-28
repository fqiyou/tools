package util

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type ContextHook struct {
}

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (hook ContextHook) Fire(entry *logrus.Entry) error {
	if pc, file, line, ok := runtime.Caller(8); ok {
		funcName := runtime.FuncForPC(pc).Name()
		entry.Data["file"] = path.Base(file)
		entry.Data["func"] = path.Base(funcName)
		entry.Data["line"] = line
	}
	return nil
}


var Log *logrus.Logger


func init() {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
	Log.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
	Log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	//ConfigLocalFilesystemLogger("", "std.log", time.Hour*24, time.Second*20)

	//Log.AddHook(ContextHook{})
	Log.AddHook(NewHook())


}

func getLeaveWriter(leave string,logPath string, logFileName string,maxAge time.Duration,rotationTime time.Duration,formatTime string)  (*rotatelogs.RotateLogs, error){

	baseLogPath := path.Join(logPath, leave + "-" + logFileName)
	timeLogPath := path.Join(logPath,leave + "-"+logFileName+"."+formatTime) //formatTime = "%Y%m%d%H%M"
	leaveWriter, err := rotatelogs.New(timeLogPath,
		rotatelogs.WithLinkName(baseLogPath),// 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),// 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime),// 日志切割时间间隔
		)

	if err != nil {
		Log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	return leaveWriter,err
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration,formatTime string) {

	infoWriter, _ := getLeaveWriter("info",logPath,logFileName,maxAge,rotationTime,formatTime)
	warnWriter, _ := getLeaveWriter("warn",logPath,logFileName,maxAge,rotationTime,formatTime)
	errorWriter, _ := getLeaveWriter("error",logPath,logFileName,maxAge,rotationTime,formatTime)
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: infoWriter,
		logrus.InfoLevel:  infoWriter,
		logrus.WarnLevel:  warnWriter,
		logrus.ErrorLevel: errorWriter,
		logrus.FatalLevel: errorWriter,
	},&logrus.JSONFormatter{})
	Log.AddHook(lfHook)
}



// hook实现行号
type Hook struct {
	Field     string
	Skip      int
	levels    []logrus.Level
	Formatter func(file, function string, line int) string
}

func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = hook.Formatter(findCaller(hook.Skip))
	return nil
}

func NewHook(levels ...logrus.Level) *Hook {
	hook := Hook{
		Field:  "source",
		Skip:   5,
		levels: levels,
		Formatter: func(file, function string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		},
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}

	return &hook
}

func findCaller(skip int) (string, string, int) {
	var (
		pc       uintptr
		file     string
		function string
		line     int
	)
	for i := 0; i < 10; i++ {
		pc, file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return file, function, line
}

func getCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return pc, file, line
}
