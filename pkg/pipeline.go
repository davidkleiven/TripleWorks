package pkg

import "fmt"

type Step[T any] struct {
	Name string
	Run  func(*T) error
}

func Pipe[T any](ctx *T, steps ...Step[T]) error {
	for i, step := range steps {
		if err := step.Run(ctx); err != nil {
			return fmt.Errorf("Step %d (%s) failed: %w", i, step.Name, err)
		}
	}
	return nil
}
