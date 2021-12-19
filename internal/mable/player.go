package mable

import (
	"github.com/gitfyu/mable/entity"
	"github.com/google/uuid"
)

type player struct {
	name string
	uid  uuid.UUID
	id   entity.ID
}

func (p *player) GetEntityID() entity.ID {
	return p.id
}
