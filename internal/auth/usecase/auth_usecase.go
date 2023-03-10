package usecase

import (
	"context"
	"net/http"

	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/auth"
	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/internal/monster"
	httpErr "github.com/iamaul/go-pokedex/pkg/error"
	"github.com/iamaul/go-pokedex/pkg/logger"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthUsecase struct {
	cfg         *config.Config
	authRepo    auth.Repository
	monsterRepo monster.MonsterRepository
	logger      logger.Logger
}

func NewAuthUsecase(cfg *config.Config, authRepo auth.Repository, monsterRepo monster.MonsterRepository, log logger.Logger) auth.Usecase {
	return &AuthUsecase{cfg: cfg, authRepo: authRepo, monsterRepo: monsterRepo, logger: log}
}

func (u *AuthUsecase) UserRegistration(ctx context.Context, user *domain.User) (*domain.UserWithToken, error) {
	_, err := u.authRepo.FindByUsername(ctx, user.Username)
	if err == nil {
		return nil, httpErr.NewRestErrorWithMessage(http.StatusBadRequest, httpErr.ErrUserAlreadyExists, err)
	}

	if err := user.PrepareCreate(); err != nil {
		return nil, httpErr.NewBadRequestError(errors.Wrap(err, "AuthUsecase.UserRegistration.PrepareCreate"))
	}

	createdUser, err := u.authRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Set value of password payload to empty for a security reason
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErr.NewInternalServerError(errors.Wrap(err, "AuthUsecase.UserRegistration.GenerateJWTToken"))
	}

	return &domain.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (u *AuthUsecase) UserAuthentication(ctx context.Context, user *domain.User) (*domain.UserWithToken, error) {
	foundUser, err := u.authRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httpErr.NewUnauthorizedError(errors.Wrap(err, "AuthUsecase.UserAuthentication.ComparePasswords"))
	}

	// Set value of password payload to empty for a security reason
	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, httpErr.NewInternalServerError(errors.Wrap(err, "AuthUsecase.UserAuthentication.GenerateJWTToken"))
	}

	return &domain.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (u *AuthUsecase) UserUpdate(ctx context.Context, user *domain.UserUpdate) (*domain.UserUpdate, error) {
	if err := user.PrepareUpdate(); err != nil {
		return nil, httpErr.NewBadRequestError(errors.Wrap(err, "AuthUsecase.UserUpdate.PrepareUpdate"))
	}

	updatedUser, err := u.authRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *AuthUsecase) UserDeletion(ctx context.Context, userID primitive.ObjectID) error {
	if err := u.authRepo.DeleteUser(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) UserList(ctx context.Context, pq *utils.PaginationQuery) (*domain.UserList, error) {
	return u.authRepo.FetchUsers(ctx, pq)
}

func (u *AuthUsecase) UserCatchMonster(ctx context.Context, userID primitive.ObjectID, monsterID primitive.ObjectID) error {
	monster, err := u.monsterRepo.FindByID(ctx, monsterID)
	if err != nil {
		return err
	}

	if err := u.authRepo.AddMonster(ctx, userID, monster.ID); err != nil {
		return err
	}
	return nil
}

func (u *AuthUsecase) GetByID(ctx context.Context, userID primitive.ObjectID) (*domain.User, error) {
	user, err := u.authRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Set value of password payload to empty for a security reason
	user.SanitizePassword()

	return user, nil
}
