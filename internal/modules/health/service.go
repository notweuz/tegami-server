package health

import (
	"context"
	"tegami-server/internal/util"
)

type DBChecker interface {
	PingContext(ctx context.Context) error
}

type service struct {
	db DBChecker
}

const (
	up   = "UP"
	down = "DOWN"
)

func NewService(db DBChecker) *service {
	return &service{
		db: db,
	}
}

func (s *service) GetInfo(ctx context.Context) *Health {
	status := up

	// we have only 1 «external platform» rn, so no complex checks yet
	if err := s.db.PingContext(ctx); err != nil {
		status = down
	}

	return NewHealth(
		status,
		util.GetCommit(),
		util.GetUptime().String(),
	)
}
