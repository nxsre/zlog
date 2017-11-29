package main

import (
	"flag"
	"log"

	"github.com/soopsio/zlog"
	"github.com/soopsio/zlog/zlogbeat/cmd"
	"go.uber.org/config"
	"go.uber.org/zap"
)

var (
	logger  *zap.Logger
	cfgfile = flag.String("logconf", "cfg.yml", "main log config file.")
)

func main() {
	cmd.RootCmd.Flags().AddGoFlag(flag.CommandLine.Lookup("logconf"))
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
