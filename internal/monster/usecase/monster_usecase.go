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

type MonsterUsecase struct {
	cfg             *config.Config
	monsterRepo     monster.MonsterRepository
	monsterTypeRepo monster.MonsterTypeRepository
	logger          logger.Logger
}

func NewMonsterUsecase(cfg *config.Config, monsterRepo monster.MonsterRepository, monsterTypeRepo monster.MonsterTypeRepository, log logger.Logger) monster.MonsterUsecase {
	return &MonsterUsecase{cfg: cfg, monsterRepo: monsterRepo, monsterTypeRepo: monsterTypeRepo, logger: log}
}

func (u *MonsterUsecase) MonsterCreate(ctx context.Context, monster *domain.Monster) (*domain.Monster, error) {
	_, err := u.monsterRepo.FindByName(ctx, monster.Name)
	if err == nil {
		return nil, httpErr.NewRestErrorWithMessage(http.StatusBadRequest, httpErr.ErrMonsterAlreadyExists, err)
	}

	createdMonsterType, err := u.monsterRepo.CreateMonster(ctx, monster)
	if err != nil {
		return nil, err
	}

	return createdMonsterType, nil
}

func (u *MonsterUsecase) MonsterUpdate(ctx context.Context, monster *domain.MonsterUpdate) (*domain.MonsterUpdate, error) {
	updatedMonster, err := u.monsterRepo.UpdateMonster(ctx, monster)
	if err != nil {
		return nil, err
	}

	return updatedMonster, nil
}

func (u *MonsterUsecase) MonsterDeletion(ctx context.Context, monsterID primitive.ObjectID) error {
	if err := u.monsterRepo.DeleteMonster(ctx, monsterID); err != nil {
		return err
	}

	return nil
}

func (u *MonsterUsecase) GetMonsterList(ctx context.Context, pq *utils.PaginationQuery) (*domain.MonsterList, error) {
	return u.monsterRepo.FetchMonsters(ctx, pq)
}

func (u *MonsterUsecase) GetByID(ctx context.Context, monsterID primitive.ObjectID) (*domain.Monster, error) {
	monster, err := u.monsterRepo.FindByID(ctx, monsterID)
	if err != nil {
		return nil, err
	}

	return monster, nil
}

func (u *MonsterUsecase) AttachMonsterType(ctx context.Context, monsterID primitive.ObjectID, monsterTypeID primitive.ObjectID) error {
	monsterType, err := u.monsterTypeRepo.FindByID(ctx, monsterTypeID)
	if err != nil {
		return err
	}

	if err := u.monsterRepo.AddMonsterType(ctx, monsterID, monsterType.ID); err != nil {
		return err
	}

	return nil
}
