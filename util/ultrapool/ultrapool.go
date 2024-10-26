// Package ultrapool implements a blazing fast worker pool with adaptive
// spawning of new workers and cleanup of idle workers
// It was modeled after valyala/fasthttp's worker pool which is one of the
// best worker pools I've seen in the Go world.

// Copyright 2019-2022 Moritz Fain
// Moritz Fain <moritz@fain.io>
//
// Source available at github.com/maurice2k/ultrapool,
// licensed under the MIT license (see LICENSE file).

package ultrapool

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Task interface{}
type TaskHandlerFunc func(task Task)

type WorkerPool struct {
	handlerFunc        TaskHandlerFunc
	idleWorkerLifetime time.Duration
	numShards          int
	shards             []*poolShard
	mutex              spinLocker
	started            bool
	stopped            bool
	_                  [56]byte
	spawnedWorkers     uint64
}

type workerInstance struct {
	taskChan  chan Task
	shard     *poolShard
	lastUsed  time.Time
	isDeleted bool
	_         [16]byte
}

type poolShard struct {
	wp             *WorkerPool
	workerCache    sync.Pool
	idleWorkerList []*workerInstance
	_              [52]byte
	idleWorker1    *workerInstance
	_              [56]byte
	idleWorker2    *workerInstance
	_              [56]byte
	mutex          spinLocker
	_              [56]byte
	stopped        bool
}

const defaultIdleWorkerLifetime = time.Second
const maxShards = 128

// Creates a new workerInstance pool with the given task handling function
func NewWorkerPool(handlerFunc TaskHandlerFunc) *WorkerPool {
	wp := &WorkerPool{
		handlerFunc:        handlerFunc,
		idleWorkerLifetime: defaultIdleWorkerLifetime,
		numShards:          1,
	}

	wp.SetNumShards(runtime.GOMAXPROCS(0))
	return wp
}

// Sets number of shards (default is GOMAXPROCS shards)
func (wp *WorkerPool) SetNumShards(numShards int) {
	if numShards <= 1 {
		numShards = 1
	}

	if numShards > maxShards {
		numShards = maxShards
	}

	wp.numShards = numShards
}

// Sets the time after which idling workers are shut down (default is 15 seconds)
func (wp *WorkerPool) SetIdleWorkerLifetime(d time.Duration) {
	wp.idleWorkerLifetime = d
}

// Returns the number of currently spawned workers
func (wp *WorkerPool) GetSpawnedWorkers() int {
	return int(atomic.LoadUint64(&wp.spawnedWorkers))
}

// Starts the worker pool
func (wp *WorkerPool) Start() {
	wp.mutex.Lock()
	if !wp.started {
		for i := 0; i < wp.numShards; i++ {
			shard := &poolShard{
				wp: wp,
				workerCache: sync.Pool{
					New: func() interface{} {
						return &workerInstance{
							taskChan: make(chan Task, 0),
						}
					},
				},

				idleWorkerList: make([]*workerInstance, 0, 2048),
			}
			wp.shards = append(wp.shards, shard)
		}

		wp.started = true
	}
	wp.mutex.Unlock()

	go wp.cleanup()
}

// Stops the worker pool.
// All tasks that have been added will be processed before shutdown.
func (wp *WorkerPool) Stop() {
	wp.mutex.Lock()
	if !wp.started {
		wp.mutex.Unlock()
		return
	}

	if !wp.stopped {

		for i := 0; i < wp.numShards; i++ {
			shard := wp.shards[i]
			shard.mutex.Lock()
			shard.stopped = true
			for j := 0; j < len(shard.idleWorkerList); j++ {
				if !shard.idleWorkerList[j].isDeleted {
					shard.idleWorkerList[j].isDeleted = true
					close(shard.idleWorkerList[j].taskChan)
				}
			}
			shard.mutex.Unlock()
		}
	}
	wp.stopped = true
	wp.mutex.Unlock()
}

// Adds a new task
func (wp *WorkerPool) AddTask(task Task) error {
	if !wp.started {
		return errors.New("worker pool must be started first")
	}

	shard := wp.shards[randInt()%wp.numShards]
	shard.getWorker(task)

	return nil
}

// Adds a new task
func (wp *WorkerPool) AddTaskForShard(task Task, shardIdx int) error {
	if !wp.started {
		return errors.New("worker pool must be started first")
	}

	shard := wp.shards[shardIdx%wp.numShards]
	shard.getWorker(task)

	return nil
}

// Returns next free worker or spawns a new worker
func (shard *poolShard) getWorker(task Task) (worker *workerInstance) {
	worker = shard.idleWorker1
	if worker != nil && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&shard.idleWorker1)), unsafe.Pointer(worker), nil) {
		worker.taskChan <- task
		return worker
	}

	worker = shard.idleWorker2
	if worker != nil && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&shard.idleWorker2)), unsafe.Pointer(worker), nil) {
		worker.taskChan <- task
		return worker
	}

	shard.mutex.Lock()
	iws := len(shard.idleWorkerList)
	if iws > 0 {
		worker = shard.idleWorkerList[iws-1]
		shard.idleWorkerList[iws-1] = nil
		shard.idleWorkerList = shard.idleWorkerList[0 : iws-1]
		shard.mutex.Unlock()
		worker.taskChan <- task
		return worker
	}
	shard.mutex.Unlock()

	worker = shard.workerCache.Get().(*workerInstance)
	worker.shard = shard
	go worker.run()

	worker.taskChan <- task
	return worker
}

// Main worker runner
func (worker *workerInstance) run() {
	shard := worker.shard
	wp := shard.wp
	atomic.AddUint64(&wp.spawnedWorkers, +1)

	for task := range worker.taskChan {
		if task == nil {
			break
		}
		wp.handlerFunc(task)
		if !shard.setWorkerIdle(worker) {
			break
		}

	}

	atomic.AddUint64(&wp.spawnedWorkers, ^uint64(0))
	shard.workerCache.Put(worker)
}

// Mark worker as idle
func (shard *poolShard) setWorkerIdle(worker *workerInstance) bool {
	worker.lastUsed = time.Now()

	if shard.idleWorker2 == nil && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&shard.idleWorker2)), nil, unsafe.Pointer(worker)) {
		return true
	}
	if shard.idleWorker1 == nil && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&shard.idleWorker1)), nil, unsafe.Pointer(worker)) {
		return true
	}

	worker.shard.mutex.Lock()
	if !worker.shard.stopped {
		worker.shard.idleWorkerList = append(worker.shard.idleWorkerList, worker)
	}
	worker.shard.mutex.Unlock()
	return !worker.shard.stopped
}

// Worker cleanup
func (wp *WorkerPool) cleanup() {
	var toBeCleaned []*workerInstance
	for {
		time.Sleep(wp.idleWorkerLifetime)
		if wp.stopped {
			return
		}

		now := time.Now()
		for i := 0; i < wp.numShards; i++ {
			shard := wp.shards[i]

			shard.mutex.Lock()
			idleWorkerList := shard.idleWorkerList
			iws := len(idleWorkerList)
			j := 0
			s := 0

			if iws > 400 {
				s = (iws - 1) / 2
				for s > 0 && now.Sub(idleWorkerList[s].lastUsed) < wp.idleWorkerLifetime {
					s = s / 2
				}

				if s == 0 {
					shard.mutex.Unlock()
					continue
				}
			}

			for j = s; j < iws; j++ {
				if now.Sub(idleWorkerList[s].lastUsed) < wp.idleWorkerLifetime {
					break
				}
			}

			if j == 0 {
				shard.mutex.Unlock()
				continue
			}

			toBeCleaned = append(toBeCleaned[:0], idleWorkerList[0:j]...)

			numMoved := copy(idleWorkerList, idleWorkerList[j:])
			for j = numMoved; j < iws; j++ {
				idleWorkerList[j] = nil
			}
			shard.idleWorkerList = idleWorkerList[:numMoved]
			shard.mutex.Unlock()

			for j = 0; j < len(toBeCleaned); j++ {
				if !toBeCleaned[j].shard.stopped {
					toBeCleaned[j].taskChan <- nil
				}
				toBeCleaned[j] = nil
			}
		}
	}
}

// Spin locker
type spinLocker struct {
	lock uint64
}

func (s *spinLocker) Lock() {
	schedulerRuns := 1
	for !atomic.CompareAndSwapUint64(&s.lock, 0, 1) {
		for i := 0; i < schedulerRuns; i++ {
			runtime.Gosched()
		}
		if schedulerRuns < 32 {
			schedulerRuns <<= 1
		}
	}
}

func (s *spinLocker) Unlock() {
	atomic.StoreUint64(&s.lock, 0)
}

// SplitMix64 style random pseudo number generator
type splitMix64 struct {
	state uint64
}

// Initialize SplitMix64
func (sm64 *splitMix64) Init(seed int64) {
	sm64.state = uint64(seed)
}

// Uint64 returns the next SplitMix64 pseudo-random number as a uint64
func (sm64 *splitMix64) Uint64() uint64 {
	sm64.state = sm64.state + uint64(0x9E3779B97F4A7C15)
	z := sm64.state
	z = (z ^ (z >> 30)) * uint64(0xBF58476D1CE4E5B9)
	z = (z ^ (z >> 27)) * uint64(0x94D049BB133111EB)
	return z ^ (z >> 31)

}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64
func (sm64 *splitMix64) Int63() int64 {
	return int64(sm64.Uint64() & (1<<63 - 1))
}

var splitMix64Pool sync.Pool = sync.Pool{
	New: func() interface{} {
		sm64 := &splitMix64{}
		sm64.Init(time.Now().UnixNano())
		return sm64
	},
}

func randInt() (r int) {
	sm64 := splitMix64Pool.Get().(*splitMix64)
	r = int(sm64.Int63())
	splitMix64Pool.Put(sm64)
	return
}
