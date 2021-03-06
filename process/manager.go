package process

import (
	"context"
	"time"

	"github.com/mylxsw/go-toolkit/log"
)

var logger = log.Module("toolkit.process")

// Manager is process manager
type Manager struct {
	programs          map[string]*Program
	restartProcess    chan *Process
	closeTimeout      time.Duration
	processOutputFunc OutputFunc
}

// NewManager create a new process manager
func NewManager(closeTimeout time.Duration, processOutputFunc OutputFunc) *Manager {
	return &Manager{
		programs:          make(map[string]*Program),
		restartProcess:    make(chan *Process),
		closeTimeout:      closeTimeout,
		processOutputFunc: processOutputFunc,
	}
}

// AddProgram add a new program to manager
func (manager *Manager) AddProgram(name string, command string, procNum int, username string) {
	manager.programs[name] = NewProgram(name, command, username, procNum).initProcesses(manager.processOutputFunc)
}

// Watch start watch process
func (manager *Manager) Watch(ctx context.Context) {
	for _, program := range manager.programs {
		for _, proc := range program.processes {
			go manager.startProcess(proc, 0)
		}
	}

	for {
		select {
		case process := <-manager.restartProcess:
			go manager.startProcess(process, process.retryDelayTime())
		case <-ctx.Done():
			logger.Debug("it's time to close all processes...")
			for _, program := range manager.programs {
				for _, proc := range program.processes {
					proc.stop(manager.closeTimeout)
				}
			}
			return
		}
	}
}

func (manager *Manager) startProcess(process *Process, delay time.Duration) {
	if delay > 0 {
		logger.Debugf("process %s will start after %.2fs", process.GetName(), delay.Seconds())
	}

	process.lock.Lock()
	defer process.lock.Unlock()

	process.timer = time.AfterFunc(delay, func() {
		process.removeTimer()

		logger.Debugf("process %s starting...", process.GetName())
		manager.restartProcess <- <-process.start()
	})

}

// Programs return all programs
func (manager *Manager) Programs() map[string]*Program {
	return manager.programs
}
