package logger

import (
	"errors"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestNewLoggerFunc_DevMode(t *testing.T) {
	l, err := newLoggerFunc("dev")
	if err != nil {
		t.Fatalf("expected no error for dev mode, got %v", err)
	}
	sugar := l.Sugar()
	sugar.Info("testing dev mode info")
	sugar.Infof("formatted %s", "message")
	sugar.Errorw("error message", "key", "value")

	if err := sugar.Sync(); err != nil {
		t.Logf("sync returned (ignored): %v", err)
	}
}

func TestNewLoggerFunc_ProductionMode(t *testing.T) {
	l, err := newLoggerFunc("prod")
	if err != nil {
		t.Fatalf("expected no error for production mode, got %v", err)
	}
	sugar := l.Sugar()
	sugar.Info("testing prod mode info")
	sugar.Infof("formatted %s", "message")
	sugar.Errorw("error message", "key", "value")

	if err := sugar.Sync(); err != nil {
		t.Logf("sync returned (ignored): %v", err)
	}
}

func TestNewLoggerFunc_Error(t *testing.T) {
	orig := newLoggerFunc
	newLoggerFunc = func(mode string, opts ...zap.Option) (LoggerInterface, error) {
		return nil, errors.New("init failed")
	}
	defer func() { newLoggerFunc = orig }()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on init error, got none")
		}
	}()
	Init("dev")
}

func TestL_InitializesDefault(t *testing.T) {
	mu.Lock()
	log = nil
	mu.Unlock()

	l := L()
	if l == nil {
		t.Fatal("expected non-nil logger from L()")
	}

	l2 := L()
	if l != l2 {
		t.Error("expected L() to return same logger instance")
	}
}

func TestConcurrency_InitAndL(t *testing.T) {
	mu.Lock()
	log = nil
	mu.Unlock()

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = L()
			Init("dev")
		}()
	}
	wg.Wait()

	if L() == nil {
		t.Error("expected non-nil logger after concurrency")
	}
}
