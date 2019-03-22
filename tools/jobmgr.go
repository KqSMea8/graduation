package tools

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
	"time"
)

type Job interface {
	Run() (interface{}, error)
	GetName() string
}

type Result struct {
	Result interface{}
	Err    error
}

// 用于执行并发任务的 JobMgr
type JobMgr struct {
	jobs      []Job
	timeout   time.Duration
	lock      sync.Mutex
	doneCh    chan bool
	errMap    map[string]error
	resultMap map[string]*Result
}

func NewJobMgr(timeout time.Duration) *JobMgr {
	return &JobMgr{
		timeout: timeout,
		jobs:    make([]Job, 0),
	}
}

func (mgr *JobMgr) AddJob(jobs ...Job) {
	mgr.jobs = append(mgr.jobs, jobs...)
}

// 只要有一个 Job 返回 Error 则返回 error
func (mgr *JobMgr) Start(ctx context.Context) error {
	mgr.parallel(ctx)
	mgr.join(ctx)
}

// 每一个 job 对应一个 goroutine 执行
func (mgr *JobMgr) parallel(ctx context.Context) {
	for _, job := range mgr.jobs {
		go func() {
			var err error
			// 不能让当天 job 影响到其他 goroutine
			defer func() {
				if e := recover(); e != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					logrus.Panicf("job: %s panic", job.GetName())
					e = fmt.Errorf("job: %s panic", job.GetName())
				}
				// 通知 join 当天 job 结束
				mgr.doneCh <- true
			}()
			startTime := time.Now().UnixNano()
			defer func() {
				logrus.Debugf("job: %s consume: %dns", job.GetName(), time.Now().UnixNano()-startTime)
			}()
			result, err := job.Run()
			mgr.SetResult(job.GetName(), &Result{
				Result: result,
				Err:    err,
			})
		}()
	}
}

func (mgr *JobMgr) join(ctx context.Context) error {
	timeout := time.After(mgr.timeout)
	for i := 0; i < len(mgr.jobs); {
		// 等待所有 job 退出或者超时
		select {
		case <-mgr.doneCh:
			i++
		case <-timeout:
			logrus.Errorf("jobmgr timeout: %v", mgr.timeout)
			return errors.New("jobmgr timeout")
		}
	}
	return nil
}

func (mgr *JobMgr) SetResult(jobName string, result *Result) {
	// golang map 并非并发安全
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.resultMap[jobName] = result
}

func (mgr *JobMgr) GetResult(jobName string) *Result {
	return mgr.resultMap[jobName]
}
