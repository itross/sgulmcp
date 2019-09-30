package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/itross/sgulmcp/internal/discoverer"
	"golang.org/x/sync/semaphore"
)

// ErrRedSemaphore is returned in case of an execution locked by a red semaphore.
var ErrRedSemaphore = errors.New("semaphore is red")

// OnceExecutor manages a semaphore and a wait group
// to execute only one goroutine at time.
type OnceExecutor struct {
	sem *semaphore.Weighted
}

// Try try and acquire the semaphore and than executes the func().
// If the semaphore is red, returns with a ErrRedSemaphore.
func (oe *OnceExecutor) Try(fn func() error) error {
	if !oe.sem.TryAcquire(1) {
		return ErrRedSemaphore
	}
	defer oe.sem.Release(1)

	return fn()
}

// TryWithErrCB will try and execute the func().
// In case of errors from fn() (not ErrRedSemaphore) the callback func will be executed.
func (oe *OnceExecutor) TryWithErrCB(fn func() error, cb func(err error)) {
	if err := oe.Try(fn); err != nil && err != ErrRedSemaphore {
		cb(err)
	}
}

// MCP .
type MCP struct {
	d        *discoverer.Discoverer
	hm       *sync.Mutex
	executor *OnceExecutor
}

func (mcp *MCP) masterLoop() {
	rate := 1
	ui := time.Second * time.Duration(rate)
	masterTicker := time.NewTicker(ui).C

	dr := rate //* 5
	di := time.Second * time.Duration(dr)
	discoveryTicker := time.NewTicker(di).C

	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sigTerm:
			return
		case <-masterTicker:
			go mcp.healthCheck()
		case <-discoveryTicker:
			log.Println("(*) try and execute Discovery")
			go mcp.executor.TryWithErrCB(mcp.discover, mcp.onDiscoveryError)
		}
	}
}

func (mcp *MCP) discover() error {
	log.Printf("(**) discovering services")
	return mcp.d.Discover()
}

func (mcp *MCP) onDiscoveryError(err error) {
	log.Printf("(**) Discvery error: %s", err)
}

func (mcp *MCP) healthCheck() {
	mcp.hm.Lock()
	defer mcp.hm.Unlock()
	log.Printf("health check")
}

func main() {
	log.Println("initializing Master Control Program")
	mcp := &MCP{
		d:        discoverer.New(),
		hm:       &sync.Mutex{},
		executor: &OnceExecutor{sem: semaphore.NewWeighted(1)},
	}
	mcp.masterLoop()
}
