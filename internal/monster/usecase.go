package monster

import (
	"context"

	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MonsterTypeUsecase interface {
	MonsterTypeCreate(ctx context.Context, monsterType *domain.MonsterType) (*domain.MonsterType, error)
	MonsterTypeUpdate(ctx context.Context, monsterType *domain.MonsterTypeUpdate) (*domain.MonsterTypeUpdate, error)
	MonsterTypeDeletion(ctx context.Context, monsterTypeID primitive.ObjectID) error
	GetMonsterTypeList(ctx context.Context, pq *utils.PaginationQuery) (*domain.MonsterTypeList, error)
	GetByID(ctx context.Context, monsterTypeID primitive.ObjectID) (*domain.MonsterType, error)
}

type MonsterUsecase interface {
	MonsterCreate(ctx context.Context, monster *domain.Monster) (*domain.Monster, error)
	MonsterUpdate(ctx context.Context, monster *domain.MonsterUpdate) (*domain.MonsterUpdate, error)
	MonsterDeletion(ctx context.Context, monsterID primitive.ObjectID) error
	GetMonsterList(ctx context.Context, pq *utils.PaginationQuery) (*domain.MonsterList, error)
	GetByID(ctx context.Context, monsterID primitive.ObjectID) (*domain.Monster, error)
	AttachMonsterType(ctx context.Context, monsterID, monsterTypeID primitive.ObjectID) error
}
