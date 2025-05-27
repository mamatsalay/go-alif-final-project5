package logger

import (
	"errors"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestNewLoggerFunc_DevMode(t *testing.T) {
	// Test dev mode returns a SugaredLoggerInterface without error
	l, err := newLoggerFunc("dev")
	if err != nil {
		t.Fatalf("expected no error for dev mode, got %v", err)
	}
	sugar := l.Sugar()
	sugar.Info("testing dev mode info")
	sugar.Infof("formatted %s", "message")
	sugar.Errorw("error message", "key", "value")

	// Sync should not error
	if err := sugar.Sync(); err != nil {
		t.Errorf("expected no sync error, got %v", err)
	}
}

func TestNewLoggerFunc_ProductionMode(t *testing.T) {
	// Production (default) mode
	l, err := newLoggerFunc("prod")
	if err != nil {
		t.Fatalf("expected no error for production mode, got %v", err)
	}
	sugar := l.Sugar()
	sugar.Info("testing prod mode info")
	sugar.Infof("formatted %s", "message")
	sugar.Errorw("error message", "key", "value")

	if err := sugar.Sync(); err != nil {
		t.Errorf("expected no sync error, got %v", err)
	}
}

func TestNewLoggerFunc_Error(t *testing.T) {
	// Swap out newLoggerFunc to simulate error
	orig := newLoggerFunc
	newLoggerFunc = func(mode string, opts ...zap.Option) (LoggerInterface, error) {
		return nil, errors.New("init failed")
	}
	defer func() { newLoggerFunc = orig }()

	// Init should panic on error
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on init error, got none")
		}
	}()
	Init("dev")
}

func TestL_InitializesDefault(t *testing.T) {
	// Reset global logger
	mu.Lock()
	log = nil
	mu.Unlock()

	l := L()
	if l == nil {
		t.Fatal("expected non-nil logger from L()")
	}

	// Subsequent calls should return same instance
	l2 := L()
	if l != l2 {
		t.Error("expected L() to return same logger instance")
	}
}

func TestConcurrency_InitAndL(t *testing.T) {
	// Test that concurrent Init and L calls are safe
	mu.Lock()
	log = nil
	mu.Unlock()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = L()
			Init("dev")
		}()
	}
	wg.Wait()

	// After concurrent calls, L() should still return non-nil
	if L() == nil {
		t.Error("expected non-nil logger after concurrency")
	}
}
