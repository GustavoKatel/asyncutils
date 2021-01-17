package supervisor

import (
	"context"

	"github.com/GustavoKatel/asyncutils/supervisor/interfaces"
)

var _ interfaces.Service = (*noInitNoCleanService)(nil)

type noInitNoCleanService struct {
	sf interfaces.ServiceFunc
}

func (s *noInitNoCleanService) Init(ctx context.Context) error {
	return nil
}

func (s *noInitNoCleanService) Clean(ctx context.Context) error {
	return nil
}

func (s *noInitNoCleanService) Run(ctx context.Context) error {
	return s.sf(ctx)
}

// FuncService create a service from a function
func FuncService(sf interfaces.ServiceFunc) interfaces.Service {
	return &noInitNoCleanService{sf: sf}
}
