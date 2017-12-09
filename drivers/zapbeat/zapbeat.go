package zapbeat

import (
	"flag"
	"io"
	"os"
	"sync"

	"github.com/oleiade/lane"
	"github.com/soopsio/zlog/zlogbeat/beater"
	"github.com/soopsio/zlog/zlogbeat/cmd"
	"go.uber.org/config"
	"go.uber.org/zap/zapcore"
)

// 全局日志 channel
var ch chan []byte
var logqueue *lane.Queue
var zlogbt = beater.NewZlogbeat()
var logwriter io.Writer

func init() {
	logwriter, _ = zlogbt.GetLogWriter()
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
	}()
}

type BeatWriter struct {
	lock *sync.RWMutex
}

func (sw *BeatWriter) Write(p []byte) (n int, err error) {
	n, err = logwriter.Write(p)
	return
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
		return &BeatWriteSyncer{
			BeatWriter: BeatWriter{
				lock: &sync.RWMutex{},
			},
		}, nil
	}
	return nil, err
}
