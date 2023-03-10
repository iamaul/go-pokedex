package monster

import (
	"context"

	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MonsterTypeRepository interface {
	CreateMonsterType(ctx context.Context, monsterType *domain.MonsterType) (*domain.MonsterType, error)
	UpdateMonsterType(ctx context.Context, monsterType *domain.MonsterTypeUpdate) (*domain.MonsterTypeUpdate, error)
	DeleteMonsterType(ctx context.Context, monsterTypeID primitive.ObjectID) error
	FetchMonsterTypes(ctx context.Context, pq *utils.PaginationQuery) (*domain.MonsterTypeList, error)
	FindByID(ctx context.Context, monsterTypeID primitive.ObjectID) (*domain.MonsterType, error)
	FindByName(ctx context.Context, monsterTypeName string) (*domain.MonsterType, error)
}

type MonsterRepository interface {
	CreateMonster(ctx context.Context, monster *domain.Monster) (*domain.Monster, error)
	UpdateMonster(ctx context.Context, monster *domain.MonsterUpdate) (*domain.MonsterUpdate, error)
	DeleteMonster(ctx context.Context, monsterID primitive.ObjectID) error
	FetchMonsters(ctx context.Context, pq *utils.PaginationQuery) (*domain.MonsterList, error)
	AddMonsterType(ctx context.Context, monsterID, monsterTypeID primitive.ObjectID) error
	FindByID(ctx context.Context, monsterID primitive.ObjectID) (*domain.Monster, error)
	FindByName(ctx context.Context, monsterName string) (*domain.Monster, error)
}
