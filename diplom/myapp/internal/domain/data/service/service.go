package service

import (
	"myapp/internal/domain/data/source"
	"myapp/internal/entity"
)

type Service interface {
	Collect() (entity.ResultSetT, error)
	ShowData() (entity.ResultSetT, error)
}
type service struct {
	source source.Source
}

func New(src source.Source) Service {
	return &service{
		source: src,
	}
}
