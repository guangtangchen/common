package logs

import (
	"bytes"
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	depth         = 3
	logId         = "LogId"
	line          = "line"
	filename      = "filename"
	keyEnv        = "env"
	envValueLocal = "local"
)

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	realLogId := "?"
	if value, ok := entry.Data[logId]; ok {
		st, ok2 := value.(string)
		if ok2 {
			realLogId = st
		}
	}

	fileName := "?"
	if value, ok := entry.Data[filename]; ok {
		st, ok2 := value.(string)
		if ok2 {
			fileName = st
		}
	}

	num := -1
	if value, ok := entry.Data[line]; ok {
		st, ok2 := value.(int)
		if ok2 {
			num = st
		}
	}

	var newLog string
	if entry.HasCaller() {
		newLog = fmt.Sprintf("%s | %s | %s | %s:%d | %s\n", realLogId, entry.Level, timestamp, fileName, num, entry.Message)
	} else {
		newLog = fmt.Sprintf("%s | %s | %s\n", entry.Level, timestamp, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func getCtxLogId(ctx context.Context) string {
	id := ctx.Value(logId)
	strId, ok := id.(string)
	if ok {
		return strId
	}
	return ""
}

func Init(appName string) {
	// 采用自定义格式输出日志
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&MyFormatter{})

	// 输出到文件，且定时切割
	filePrefix := "./log/" + appName + "."
	fileOutput, err := rotatelogs.New(
		filePrefix+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(filePrefix),
		rotatelogs.WithMaxAge(365*24*time.Hour),   // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 文件切割时间间隔，1h就是每个小时一个文件
	)
	if err != nil {
		panic(err)
	}
	// 本地开发环境，文件+控制台输出，需要配置环境变量
	if os.Getenv(keyEnv) == envValueLocal {
		//mw := io.MultiWriter(os.Stdout, fileOutput)
		logrus.SetOutput(os.Stdout)
	} else {
		logrus.SetOutput(fileOutput)
	}
}

func getFieldValues(ctx context.Context) (string, string, int) {
	var ok bool
	id := getCtxLogId(ctx)
	name := "?"
	num := -1
	_, name, num, ok = runtime.Caller(depth)
	if ok {
		name = filepath.Base(name)
	}
	return id, name, num
}

func genCommonEntry(ctx context.Context, format string, args ...interface{}) (string, *logrus.Entry) {
	msg := format
	if len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}

	id, name, num := getFieldValues(ctx)

	return msg, logrus.WithFields(logrus.Fields{
		logId:    id,
		filename: name,
		line:     num,
	})
}

func CtxDug(ctx context.Context, format string, args ...interface{}) {
	msg, entry := genCommonEntry(ctx, format, args...)
	entry.Debug(msg)
}

func CtxInfo(ctx context.Context, format string, args ...interface{}) {
	msg, entry := genCommonEntry(ctx, format, args...)
	entry.Info(msg)
}

func CtxWarn(ctx context.Context, format string, args ...interface{}) {
	msg, entry := genCommonEntry(ctx, format, args...)
	entry.Warn(msg)
}

func CtxError(ctx context.Context, format string, args ...interface{}) {
	msg, entry := genCommonEntry(ctx, format, args...)
	entry.Error(msg)
}
