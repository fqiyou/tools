package util

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var Log *logrus.Logger


func init() {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
	Log.AddHook(NewFileRowHook())
	Log.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
	//Log.Formatter = &logrus.JSONFormatter{}

	Log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))


	// file logger
	// AddHookLocalFileLogger("", "std.log", time.Hour*24, time.Second*20, "%Y%m%d-%H")

	// syslog logger
	//AddHookSyslog("udp","spark003:514",syslog.LOG_LOCAL4,"")

}
// 添加syslog hoot
func AddHookSyslog(network string, raddr string, priority syslog.Priority, tag string)  {
	hook, err := logrus_syslog.NewSyslogHook(network, raddr, priority, tag)
	if err != nil {
		Log.Error(err)
	}
	if err == nil {
		Log.Hooks.Add(hook)
		Log.Formatter = &logrus.TextFormatter{
			ForceColors: false,
		}

	}
}


// hook实现轮训写文件
func AddHookLocalFileLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration,formatTime string) {

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



// hook实现行号
type FileRowHook struct {
	Field     string
	Skip      int
	levels    []logrus.Level
	Formatter func(file, function string, line int) string
}

func (hook *FileRowHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *FileRowHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = hook.Formatter(findCaller(hook.Skip))
	return nil
}

func NewFileRowHook(levels ...logrus.Level) *FileRowHook {
	hook := FileRowHook{
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
