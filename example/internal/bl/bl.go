// Package bl is a business logic layer.
package bl

import (
	"context"

	"github.com/z-riley/platform/example/internal/repository"
)

//mockery:generate: true
type SomethingRepository interface {
	repository.Transactor
	UpdateSomething(ctx context.Context) error
}

var _ SomethingRepository = (*repository.Repository)(nil)

type SomethingService struct {
	Repository SomethingRepository
}

func (s *SomethingService) UpdateSomething(ctx context.Context) error {
	return s.Repository.UpdateSomething(ctx)
}

func (s *SomethingService) UpdateSomethings(ctx context.Context) error {
	tx, err := s.Repository.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.Repository.UpdateSomething(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
