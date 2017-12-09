package beater

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/soopsio/zlog/zlogbeat/config"
	"github.com/tidwall/gjson"
)

var once sync.Once
var instance *Zlogbeat

// 使用管道通信，channel 有问题，会出现日志截断和回放的问题
var logReader, logWriter, _ = os.Pipe()

func NewZlogbeat() *Zlogbeat {
	once.Do(func() {
		instance = &Zlogbeat{
			logreader: logReader,
			logwriter: logWriter,
			wg:        &sync.WaitGroup{},
		}
	})
	return instance
}

func (zl *Zlogbeat) GetLogWriter() (io.Writer, error) {
	return zl.logwriter, nil
}

type Zlogbeat struct {
	logreader io.Reader
	logwriter io.Writer
	done      chan struct{}
	wg        *sync.WaitGroup
	config    config.Config
	client    beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	logp.Info("初始化 beat")
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	bt := NewZlogbeat()
	bt.done = make(chan struct{})
	bt.config = config
	return bt, nil
}

func (bt *Zlogbeat) Run(b *beat.Beat) error {
	logp.Info("zlogbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	// ticker := time.NewTicker(bt.config.Period)

	counter := 1
	reader := bufio.NewReader(bt.logreader)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		res := gjson.ParseBytes(line)
		bt.config.Fields["level"] = res.Get("level").String()
		event := beat.Event{
			Meta:      common.MapStr{}, // 添加元素据
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"counter": counter,
				"message": res.Value(),
				"fields":  bt.config.Fields,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
	}
	return nil
}

func (bt *Zlogbeat) Sync() {
	// TODO: 等到 Publish 成功执行
	bt.wg.Wait()
}

func (bt *Zlogbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
