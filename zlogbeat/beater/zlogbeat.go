package beater

import (
	"fmt"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/tidwall/gjson"

	"github.com/soopsio/zlog/zlogbeat/config"
)

var once sync.Once
var instance *Zlogbeat

func GetZlogbeat() *Zlogbeat {
	once.Do(func() {
		instance = &Zlogbeat{
			wg: &sync.WaitGroup{},
		}
	})
	return instance
}

func (zl *Zlogbeat) SetLogCh(ch chan []byte) error {
	zl.logch = ch
	return nil
}

func (zl *Zlogbeat) GetLogCh() (<-chan []byte, error) {
	return zl.logch, nil
}

type Zlogbeat struct {
	logch  chan []byte
	done   chan struct{}
	wg     *sync.WaitGroup
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	bt := GetZlogbeat()
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
	ticker := time.NewTicker(bt.config.Period)

	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		case val, ok := <-bt.logch:
			if ok {
				// bt.wg.Add(1)
				res := gjson.ParseBytes(val)
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
				// bt.wg.Done()
			}
			counter++
		}
	}
}

func (bt *Zlogbeat) Sync() {
	// TODO: 等到 Publish 成功执行
	bt.wg.Wait()
}

func (bt *Zlogbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
