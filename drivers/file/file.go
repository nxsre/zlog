package file

import (
	"gopkg.in/natefinch/lumberjack.v2"

	"os"
	"path"

	"go.uber.org/config"
	"go.uber.org/zap/zapcore"
)

type FileWriteSyncer struct {
	*lumberjack.Logger
}

func (zws *FileWriteSyncer) Sync() error {
	return nil
}

func NewWriteSyncer(p config.Value) (zapcore.WriteSyncer, error) {
	conf := &FileConfig{}
	err := p.Populate(conf)
	if err != nil {
		return nil, err
	}
	if conf.Enable {
		return &FileWriteSyncer{
			Logger: &lumberjack.Logger{
				Filename:   path.Join(conf.Path, string(os.PathSeparator), conf.Filename),
				MaxBackups: conf.NumberOfFiles,
				MaxSize:    conf.RotateEveryKB / 1024,
			},
		}, nil
	}
	return nil, err
}
