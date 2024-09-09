package job

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

type ReportError interface {
	Report(ctx context.Context, err error) error
}

type Worker struct {
	name           string
	job            Handler
	limiter        *rate.Limiter
	reStartLimiter *rate.Limiter
	reportError    ReportError
	sleep          time.Duration
}

type Handler func(ctx context.Context) error

type WorkerOption func(worker *Worker)

func WithReport(report ReportError) WorkerOption {
	return func(worker *Worker) {
		worker.reportError = report
	}
}

func WithSleep(duration time.Duration) WorkerOption {
	return func(worker *Worker) {
		worker.sleep = duration
	}
}

func WithLimiter(limiter *rate.Limiter) WorkerOption {
	return func(worker *Worker) {
		worker.limiter = limiter
	}
}

func WithLimiterDuration(duration time.Duration) WorkerOption {
	return func(worker *Worker) {
		worker.limiter = rate.NewLimiter(rate.Every(duration), 1)
	}
}

func NewWorker(name string, job Handler, opts ...WorkerOption) *Worker {
	w := &Worker{
		name: name,
		job:  job,
		// 默认每分钟执行1次任务
		limiter: rate.NewLimiter(rate.Every(time.Minute), 1),
		// 默认每小时最多重启5次任务
		reStartLimiter: rate.NewLimiter(rate.Every(time.Hour), 5),
		// 限流不通过时的休眠时间，默认为 1s
		sleep: time.Second,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}
