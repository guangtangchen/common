package util

import (
	"context"
	"github.com/guangtangchen/common/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"os"
	"strconv"
	"time"
)

const (
	logId = "LogId"
)

func JsonMarshal(ctx context.Context, v interface{}) string {
	var ob = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, err := ob.Marshal(v)
	if err != nil {
		logs.CtxError(ctx, "JsonMarshal err = %v", err)
		return err.Error()
	}
	return string(bytes)
}

// JsonUnmarshal 反序列化，to为空结构体的指针
func JsonUnmarshal(ctx context.Context, from string, to interface{}) error {
	var ob = jsoniter.ConfigCompatibleWithStandardLibrary
	err := ob.Unmarshal([]byte(from), to)
	if err != nil {
		logs.CtxInfo(ctx, "JsonUnmarshal err = %v", err)
		return err
	}
	return nil
}

func TimeFormatString(target time.Time) string {
	return target.Format("2006-01-02 15:04:05")
}

func TimeFormatSerialString(target time.Time) string {
	return target.Format("20060102150405")
}

func NewCtxWithLogId(ctx context.Context) context.Context {
	return context.WithValue(ctx, logId, GenLogId())
}

func GetCtxLogId(ctx context.Context) string {
	id := ctx.Value(logId)
	strId, ok := id.(string)
	if ok {
		return strId
	}
	return ""
}

func GenId() string {
	return xid.New().String()
}

func GenLogId() string {
	id := xid.New()
	now := strconv.FormatInt(time.Now().UnixMicro(), 10)
	return TimeFormatSerialString(time.Now()) + now[len(now)-7:] + id.String()
}

func IsLocalEnv() bool {
	return os.Getenv("env") == "local"
}
