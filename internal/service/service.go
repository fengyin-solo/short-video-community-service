package service

import (
	"shortvideo/internal/config"
	"shortvideo/internal/store"
	"shortvideo/pkg/logger"
)

// Service 业务服务层。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 构造业务服务。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
