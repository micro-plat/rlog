package services

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
)

type Logging struct {
	bufferChan  chan []byte
	cacheBuffer [][]byte
	conf        *Conf
	logger      *logger.Logger
	lock        sync.Mutex
	client      *Client
	once        sync.Once
	w           sync.WaitGroup
	done        bool
	closeChan   chan struct{}
}

//NewLogging 创建日志组件
func NewLogging(client *Client, c *Conf, l *logger.Logger) (r *Logging, err error) {
	r = &Logging{
		client:      client,
		conf:        c,
		logger:      l,
		bufferChan:  make(chan []byte, 100000),
		cacheBuffer: make([][]byte, 0, 1000),
		closeChan:   make(chan struct{}),
	}
	go r.addToBuffer()
	go r.loopWrite()
	return r, nil
}

//Save 保存日志
func (l *Logging) Save(data []byte) error {
	if l.done {
		return errs.NewError(901, "服务已关闭，日志未处理")
	}
	var buff bytes.Buffer
	if err := json.Compact(&buff, data); err != nil {
		return err
	}
	l.bufferChan <- buff.Bytes()
	return nil
}

func (l *Logging) addToBuffer() {
LOOP:
	for {
		select {
		case buff, ok := <-l.bufferChan:
			if !ok {
				break LOOP
			}
			if len(buff) <= 2 {
				continue
			}
			nbuff := make([][]byte, 0, 1)
			if buff[0] == '[' {
				sections := bytes.Split(buff[1:len(buff)-1], []byte("},"))
				for _, s := range sections {
					if string(s[len(s)-1]) != "}" {
						s = append(s, []byte("}")...)
					}

					nbuff = append(nbuff, s)
				}
			} else {
				nbuff = append(nbuff, buff)
			}
			l.lock.Lock()
			l.cacheBuffer = append(l.cacheBuffer, nbuff...)
			l.lock.Unlock()
		}
	}
}

func (l *Logging) loopWrite() {
Loop:
	for {
		select {
		case <-l.closeChan:
			break Loop
		case <-time.After(time.Second * time.Duration(l.conf.Cron)): //定时写入数据
			if l.done {
				break Loop
			}
			l.writeNow()
		}
	}

	//最后写入所有数据
	l.writeNow()
	l.w.Done()
}

func (l *Logging) writeNow() {
	l.lock.Lock()
	buff := l.cacheBuffer[0:]
	l.cacheBuffer = l.cacheBuffer[:0]
	l.lock.Unlock()
	if len(buff) > 0 {
		go l.Write(buff)
	}
}

func (l *Logging) Write(p [][]byte) (n int, err error) {
	l.w.Add(1)
	defer l.w.Done()
	l.logger.Info(" --> logging request")
	start := time.Now()
	n, err = l.client.BenchAddData(p, l.conf.WriteTimeout)
	if err != nil {
		l.logger.Errorf("-> logging response %d条 %v %v", len(p), time.Since(start), err)
		return 0, err
	}
	l.logger.Infof(" --> logging response %d条 %v %v", len(p), n, time.Since(start))
	return len(p) - 1, nil
}

//Close 关闭当前日志组件
func (l *Logging) Close() error {
	l.done = true
	l.once.Do(func() {
		l.w.Add(1)
		close(l.bufferChan)
		close(l.closeChan)
	})
	l.w.Wait()
	return nil
}
