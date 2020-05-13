package dbcache

import (
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/wujiu2020/strip"

	"cos-backend-com/src/common/util"
)

const (
	defaultNotifyName     = "db_update_notify"
	defaultBufferSize     = 1000
	defaultUpdateInterval = 3
)

type TriggerConfig struct {
	RedisPool  *pool.Pool
	DataSource string

	Tables  []*Table
	DBCache *DBCache

	UpdateInterval int
	BufferSize     int

	LockName    string
	ChannelName string
}

type UpdateTrigger struct {
	cfg TriggerConfig
	log strip.Logger

	tables  map[string]*Table
	buffers chan *Data
	mux     sync.Mutex
	rmux    *util.RedisMutex
	lis     *pq.ListenerConn
	running bool
}

func NewUpdateTigger(cfg TriggerConfig) (trig *UpdateTrigger) {
	if cfg.UpdateInterval == 0 {
		cfg.UpdateInterval = defaultUpdateInterval
	}
	if cfg.BufferSize == 0 {
		cfg.BufferSize = defaultBufferSize
	}
	if cfg.LockName == "" {
		cfg.LockName = defaultNotifyName
	}
	if cfg.ChannelName == "" {
		cfg.ChannelName = defaultNotifyName
	}

	tables := make(map[string]*Table, len(cfg.Tables))
	for _, table := range cfg.Tables {
		tables[table.cfg.Table] = table
	}

	trig = &UpdateTrigger{
		cfg:     cfg,
		rmux:    util.NewRedisMutex(cfg.LockName, cfg.RedisPool),
		log:     strip.NewReqLogger(strip.DefaultLogPrinter, "dbtrigger]["+cfg.LockName),
		buffers: make(chan *Data, cfg.BufferSize),
		tables:  tables,
	}
	return
}

func (p *UpdateTrigger) Start() {
	p.mux.Lock()
	if p.running {
		p.mux.Unlock()
		return
	}
	p.running = true
	p.mux.Unlock()

	go p.updateRecords()

	flag := 1
	for {
		p.log.Recover(func() {
			locked := p.loopRun()
			if !locked && (flag&2 == 0) {
				p.log.Info("no lucky, failed get lock, continue trying")
				flag = 1 | 2
			} else if locked {
				flag = 1
			}
		})
		time.Sleep(1e8)
	}
}

func (p *UpdateTrigger) loopRun() (locked bool) {
	unlockedCh, err := p.rmux.LockKeeper()
	if err != nil {
		if err != util.ErrMaxTries {
			p.log.Warn(err)
		}
		return
	}
	defer p.rmux.Unlock()
	locked = true

	p.log.Info("lucky, got locked")
	defer func() {
		p.log.Info("release lock")
	}()

	notifyCh := make(chan *pq.Notification, 1)
	lis, err := pq.NewListenerConn(p.cfg.DataSource, notifyCh)
	if err != nil {
		p.log.Warn(err)
		return
	}
	defer lis.Close()

	ok, err := lis.Listen(p.cfg.ChannelName)
	if err != nil || !ok {
		p.log.Warn(ok, err)
		return
	}

	p.log.Info("start receive notify")

	for {
		select {
		case <-unlockedCh:
			return
		case notify := <-notifyCh:
			if notify == nil {
				return
			}
			p.receiveNotify(notify)
		}
	}
	return
}

func (p *UpdateTrigger) receiveNotify(notify *pq.Notification) {
	parts := strings.Split(notify.Extra, ",")
	if len(parts) < 7 {
		p.log.Error("wrong trigger notify:", notify.Channel, notify.Extra)
		return
	}
	name := parts[5]
	id := parts[6]

	table := p.tables[name]
	if table == nil {
		p.log.Error("trigger notify:", notify.Channel, name, id, "table not defined")
		return
	}

	data := table.Data(id)

	select {
	case p.buffers <- data:
		p.log.Debug("trigger notify:", notify.Channel, name, id, "trigger notify quened")
	case <-time.After(time.Minute):
		p.log.Error("trigger notify timeout:", notify.Channel, name, id)
	}
}

func (p *UpdateTrigger) updateRecords() {
	loc := sync.Mutex{}
	records := make([]*Data, p.cfg.BufferSize)
	n := 0
	go func() {
		for {
			loc.Lock()
			if n >= p.cfg.BufferSize {
				time.Sleep(1e6)
				loc.Unlock()
				continue
			}
			loc.Unlock()

			select {
			case rec := <-p.buffers:
				loc.Lock()
				records[n] = rec
				n += 1
				loc.Unlock()
			}
		}
	}()
	for {
		select {
		case <-time.After(time.Duration(p.cfg.UpdateInterval) * time.Second):
			loc.Lock()
			if n == 0 {
				loc.Unlock()
				continue
			}
			err := p.cfg.DBCache.FetchDatasForce(records[:n]...)
			if err != nil {
				p.log.Warn(err)
				continue
			}
			n = 0
			loc.Unlock()
		}
	}
}
