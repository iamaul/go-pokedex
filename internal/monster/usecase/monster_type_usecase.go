package usecase

import (
	"context"
	"net/http"

	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/internal/monster"
	httpErr "github.com/iamaul/go-pokedex/pkg/error"
	"github.com/iamaul/go-pokedex/pkg/logger"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MonsterTypeUsecase struct {
	cfg             *config.Config
	monsterTypeRepo monster.MonsterTypeRepository
	logger          logger.Logger
}

func NewMonsterTypeUsecase(cfg *config.Config, monsterTypeRepo monster.MonsterTypeRepository, log logger.Logger) monster.MonsterTypeUsecase {
	return &MonsterTypeUsecase{cfg: cfg, monsterTypeRepo: monsterTypeRepo, logger: log}
}

func (u *MonsterTypeUsecase) MonsterTypeCreate(ctx context.Context, monsterType *domain.MonsterType) (*domain.MonsterType, error) {
	_, err := u.monsterTypeRepo.FindByName(ctx, monsterType.Name)
	if err == nil {
		return nil, httpErr.NewRestErrorWithMessage(http.StatusBadRequest, httpErr.ErrMonsterAlreadyExists, err)
	}

	createdMonsterType, err := u.monsterTypeRepo.CreateMonsterType(ctx, monsterType)
	if err != nil {
		return nil, err
	}

	return createdMonsterType, nil
}

func (u *MonsterTypeUsecase) MonsterTypeUpdate(ctx context.Context, monsterType *domain.MonsterTypeUpdate) (*domain.MonsterTypeUpdate, error) {
	updatedMonsterType, err := u.monsterTypeRepo.UpdateMonsterType(ctx, monsterType)
	if err != nil {
		return nil, err
	}

	return updatedMonsterType, nil
}

func (u *MonsterTypeUsecase) MonsterTypeDeletion(ctx context.Context, monsterTypeID primitive.ObjectID) error {
	if err := u.monsterTypeRepo.DeleteMonsterType(ctx, monsterTypeID); err != nil {
		return err
	}

	return nil
}

func (u *MonsterTypeUsecase) GetMonsterTypeList(ctx context.Context, pq *utils.PaginationQuery) (*domain.MonsterTypeList, error) {
	return u.monsterTypeRepo.FetchMonsterTypes(ctx, pq)
}

func (u *MonsterTypeUsecase) GetByID(ctx context.Context, monsterTypeID primitive.ObjectID) (*domain.MonsterType, error) {
	monsterType, err := u.monsterTypeRepo.FindByID(ctx, monsterTypeID)
	if err != nil {
		return nil, err
	}

	return monsterType, nil
}
