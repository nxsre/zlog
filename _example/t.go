package main

import (
	"fmt"

	"github.com/soopsio/zap"
	"go.uber.org/zap"
	"github.com/uber-go/config"
)

type SyncWriter struct{}

func (sw *SyncWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	fmt.Printf(string(p))
	return n, err
}

type ZapWriteSyncer struct {
	SyncWriter
}

func (zws *ZapWriteSyncer) Sync() error {
	// TODO: 如果是写文件，或者 kafka ，此处应该等待写入完毕
	return nil
}

func main() {
	Foo()
}

func Foo() {
	sw := &ZapWriteSyncer{}
	conf := zap.NewProductionConfig()
	conf.DisableCaller = true
	conf.Encoding = "json"
	z, err := conf.Build(zlog.SetOutput(sw, conf))
	_ = err
	z.Info("aaaaaaaaaaaa", zap.String("keyaaa", "valueaaa"))
	z.Sync()
	// Or, if your already have sugared logger.
	// sugar := z.Sugar()
	// sugar = sugar.Desugar().WithOptions(SetOutput(sw, conf)).Sugar()

	// sugar.Info("aaaaaaa", zap.String("cc", "bbbbbbbbb"))
}
