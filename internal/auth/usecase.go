package auth

import (
	"context"

	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	UserRegistration(ctx context.Context, user *domain.User) (*domain.UserWithToken, error)
	UserAuthentication(ctx context.Context, user *domain.User) (*domain.UserWithToken, error)
	UserUpdate(ctx context.Context, user *domain.UserUpdate) (*domain.UserUpdate, error)
	UserDeletion(ctx context.Context, userID primitive.ObjectID) error
	UserList(ctx context.Context, pq *utils.PaginationQuery) (*domain.UserList, error)
	UserCatchMonster(ctx context.Context, userID primitive.ObjectID, monsterID primitive.ObjectID) error
	GetUserByID(ctx context.Context, userID primitive.ObjectID) (*domain.User, error)
}
