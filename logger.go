// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package zlog

import (
	"log"

	"github.com/soopsio/zlog/drivers/file"
	"github.com/soopsio/zlog/drivers/zapbeat"

	"go.uber.org/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetOutput replaces existing Core with new, that writes to passed WriteSyncer.
func SetOutput(ws zapcore.WriteSyncer, conf zap.Config) zap.Option {
	var enc zapcore.Encoder
	// Copy paste from zap.Config.buildEncoder.
	switch conf.Encoding {
	case "json":
		enc = zapcore.NewJSONEncoder(conf.EncoderConfig)
	case "console":
		enc = zapcore.NewConsoleEncoder(conf.EncoderConfig)
	default:
		panic("unknown encoding")
	}
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(enc, ws, conf.Level)
	})
}

type ZapWriteSyncerInterface interface {
	Sync() error
}

type ZapWriteSyncer struct {
	writeSyncer []zapcore.WriteSyncer
}

// Write 调用底层驱动 io.Writer 实现
func (zws *ZapWriteSyncer) Write(p []byte) (n int, err error) {
	for _, v := range zws.writeSyncer {
		v.Write(p)
	}
	return len(p), nil
}

// Sync 调用底层驱动 Sync 方法
func (zws *ZapWriteSyncer) Sync() error {
	// TODO: 如果是写文件，或者 kafka ，此处应该等待写入完毕
	for _, v := range zws.writeSyncer {
		v.Sync()
	}
	return nil
}

func NewWriteSyncer(provider config.Provider) zapcore.WriteSyncer {
	zapWriteSyncer := &ZapWriteSyncer{}
	var syncWriters []zapcore.WriteSyncer
	if output_zapbeat := provider.Get("zap.zapbeat"); output_zapbeat.HasValue() {
		syncWriter, err := zapbeat.NewWriteSyncer(output_zapbeat)
		if err != nil {
			log.Println("zapbeat 初始化失败:", err)
		} else {
			if syncWriter != nil {
				log.Println(syncWriter)
				syncWriters = append(syncWriters, syncWriter)
			}
		}
	}

	if output_file := provider.Get("zap.file"); output_file.HasValue() {
		syncWriter, err := file.NewWriteSyncer(output_file)
		if err != nil {
			log.Println("file 初始化失败:", err)
		} else {
			if syncWriter != nil {
				syncWriters = append(syncWriters, syncWriter)
			}
		}
	}
	zapWriteSyncer.writeSyncer = syncWriters
	return zapWriteSyncer
}
