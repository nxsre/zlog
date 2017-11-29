package zapbeat

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/soopsio/zlog/zlogbeat/beater"
	"github.com/soopsio/zlog/zlogbeat/cmd"
	"go.uber.org/config"
	"go.uber.org/zap/zapcore"
)

// 全局日志 channel
var ch = make(chan []byte)
var zlogbt = beater.GetZlogbeat()

func init() {

	zlogbt.SetLogCh(ch)

	go func() {
		// 启动 beater
		// 添加 -c 选项，指定配置文件路径
		c := flag.CommandLine.Lookup("c")
		// c.Value.Set("zlogbeat.yml")

		// 添加 -e 选项
		e := flag.CommandLine.Lookup("e")
		// e.Value.Set("true")

		cmd.RootCmd.PersistentFlags().AddGoFlag(c)
		cmd.RootCmd.PersistentFlags().AddGoFlag(e)
		if err := cmd.RootCmd.Execute(); err != nil {
			os.Exit(1)
		}
		fmt.Println("cmd.RootCmd 结束")
	}()
}

type BeatWriter struct{}

func (sw *BeatWriter) Write(p []byte) (n int, err error) {
	timer := time.NewTimer(time.Second * 5)
	n = len(p)
	select {
	case ch <- p:
	case <-timer.C:
		timer.Reset(time.Second * 5)
	}
	return n, err
}

type BeatWriteSyncer struct {
	BeatWriter
}

func (zws *BeatWriteSyncer) Sync() error {
	zlogbt.Sync()
	return nil
}

func NewWriteSyncer(p config.Value) (zapcore.WriteSyncer, error) {
	conf := &BeatConfig{}
	err := p.Populate(conf)
	if err != nil {
		return nil, err
	}
	if conf.Enable {
		return &BeatWriteSyncer{}, nil
	}
	return nil, err
}
