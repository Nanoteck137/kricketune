package hook

import (
	"context"
	"fmt"
	"sync"
)

type HookFunc[T any] func(ctx context.Context, data T) error

type Hook[T any] struct {
	mu       sync.RWMutex
	handlers []HookFunc[T]
	name     string
}

func NewHook[T any](name string) *Hook[T] {
	return &Hook[T]{
		name:     name,
	}
}

func (h *Hook[T]) Register(handler HookFunc[T]) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers = append(h.handlers, handler)
}

func (h *Hook[T]) Call(ctx context.Context, data T) error {
	h.mu.RLock()
	handlers := make([]HookFunc[T], len(h.handlers))
	copy(handlers, h.handlers)
	h.mu.RUnlock()

	for i, handler := range handlers {
		if err := handler(ctx, data); err != nil {
			return fmt.Errorf("hook '%s' handler %d failed: %w", h.name, i, err)
		}
	}
	return nil
}

func (h *Hook[T]) CallAsync(ctx context.Context, data T) []error {
	h.mu.RLock()
	handlers := make([]HookFunc[T], len(h.handlers))
	copy(handlers, h.handlers)
	h.mu.RUnlock()

	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for i, handler := range handlers {
		wg.Add(1)
		go func(idx int, h HookFunc[T]) {
			defer wg.Done()
			if err := h(ctx, data); err != nil {
				errChan <- fmt.Errorf("hook handler %d: %w", idx, err)
			}
		}(i, handler)
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}
	return errors
}

func (h *Hook[T]) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.handlers)
}
