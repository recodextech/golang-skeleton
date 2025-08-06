package container

import (
	"fmt"
	gocon "github.com/wgarunap/goconf"
	"log"
	"os"
	"sync"
)

type Initable interface {
	Init(Container) error
}
type Container interface {
	Init(modules ...string)
	Bind(typ string, obj interface{})
	Resolve(name string) interface{}
	GetGlobalConfig(typ string) interface{}
}

type container struct {
	bindings      map[string]interface{}
	moduleConfigs map[string]interface{}
	stopSigs      []<-chan interface{} // channel for shutdown signals
	stopped       chan struct{}
	lock          sync.Mutex
	logger        *log.Logger
}

func NewContainer() AppContainer {
	return &container{
		bindings:      map[string]interface{}{},
		moduleConfigs: map[string]interface{}{},
		lock:          sync.Mutex{},
		stopSigs:      []<-chan interface{}{},
		stopped:       make(chan struct{}, 1),
		logger:        log.New(os.Stdout, `di`, log.LstdFlags),
	}
}

func (c *container) Bind(typ string, obj interface{}) {
	c.bindings[typ] = obj
}

func (c *container) Init(modules ...string) {
	for _, name := range modules {
		if in, ok := c.bindings[name].(Initable); ok {
			err := in.Init(c)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (c *container) Resolve(name string) interface{} {
	if con, ok := c.bindings[name]; ok {
		return con
	}
	panic(fmt.Sprintf(`%s no module`, name))
}

func (c *container) GetGlobalConfig(typ string) interface{} {
	if config, ok := c.moduleConfigs[typ]; ok {
		return config
	}
	panic(fmt.Sprintf(`%s no module`, typ))
}

func (c *container) Start(modules ...string) {
	for _, sig := range c.stopSigs {
		go func(ch <-chan interface{}) {
			<-ch
			// initiate graceful shutdown
			c.stopped <- struct{}{}
		}(sig)
	}

	for _, module := range modules {
		c.logger.Println(fmt.Sprintf(`module %s starting...`, module))

		m := c.bindings[module]

		runnable, ok := m.(Runnable)
		if !ok {
			panic(fmt.Sprintf(`container: module [%s] is not runnable, starting failed`, module))
		}
		go func(r Runnable) {
			if err := r.Run(); err != nil {
				panic(err)
			}
		}(runnable)

		c.logger.Println(fmt.Sprintf(`module %s started`, module))
	}

	<-c.stopped
}

// SetModuleGlobalConfig adds static configurations of modules in to the container.
func (c *container) SetModuleGlobalConfig(configs ...ModuleConfig) error {
	cfgs := make([]gocon.Configer, 0)
	for _, value := range configs {
		cfgs = append(cfgs, value.Value.(gocon.Configer))
		c.moduleConfigs[value.Key] = value.Value
	}
	return gocon.Load(cfgs...)
}

// Shutdown gracefully shuts down modules in the order they are provided.
func (c *container) Shutdown(modules ...string) {
	// un register channels

	// stop modules
	for _, module := range modules {
		c.logger.Println(fmt.Sprintf(`module %s stopping...`, module))

		m := c.bindings[module]

		stoppable, ok := m.(Stoppable)
		if !ok {
			panic(fmt.Sprintf(`container: module [%s] is not stoppable, stopping failed`, module))
		}
		if err := stoppable.Stop(); err != nil {
			c.logger.Println(err)
		}

		c.logger.Println(fmt.Sprintf(`module %s stopped`, module))
	}

	c.stopped <- struct{}{}
}
