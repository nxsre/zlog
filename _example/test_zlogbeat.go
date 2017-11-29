package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/soopsio/zlog/zlogbeat/beater"
	"github.com/soopsio/zlog/zlogbeat/cmd"
	"go.uber.org/zap/zapcore"
)

func test_zlogbeat() {
	ch := make(chan []byte)
	zlogbt := beater.GetZlogbeat()
	zlogbt.SetLogCh(ch)
	go func() {
		for i := 0; i <= 10; i++ {
			entry := zapcore.Entry{
				Level:   1,
				Time:    time.Now(),
				Message: "test msg",
			}
			bf, err := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				MessageKey:     "msg",
				LevelKey:       "level",
				NameKey:        "name",
				TimeKey:        "ts",
				CallerKey:      "caller",
				StacktraceKey:  "stacktrace",
				LineEnding:     "\n",
				EncodeTime:     zapcore.EpochTimeEncoder,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}).EncodeEntry(entry, nil)
			if err == nil {
				ch <- []byte(bf.String())
			}

			time.Sleep(time.Millisecond * 500)
		}
		fmt.Println("日志写入完成")
		// time.Sleep(1 * time.Second)
		zlogbt.Stop()
	}()

	// 解析 -c 选项，指定配置文件路径
	c := flag.CommandLine.Lookup("c")
	// c.Value.Set("zlogbeat.yml")

	// 解析 -e 选项
	e := flag.CommandLine.Lookup("e")
	e.Value.Set("true")

	cmd.RootCmd.PersistentFlags().AddGoFlag(c)
	cmd.RootCmd.PersistentFlags().AddGoFlag(e)
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	fmt.Println("cmd.RootCmd 结束")
}
