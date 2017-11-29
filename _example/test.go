package main

import (
	"flag"
	"log"

	"github.com/soopsio/zlog"
	"go.uber.org/config"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func main() {
	cfgfile := flag.String("conf", "cfg.yml", "main log config path")
	flag.Parse()
	p, err := config.NewYAMLProviderFromFiles(*cfgfile)
	if err != nil {
		log.Fatalln(err)
	}

	sw := zlog.NewWriteSyncer(p)
	conf := zap.NewProductionConfig()
	conf.DisableCaller = true
	conf.Encoding = "json"

	logger, _ := conf.Build(zlog.SetOutput(sw, conf))
	for i := 1; i <= 1000; i++ {
		logger.Info("aaaaaaaaaaaa", zap.String("keyaaa", "valueaaa"), zap.Int("key", i))
		logger.Sync()
	}
}
