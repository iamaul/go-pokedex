package auth

import (
	"context"

	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.UserUpdate) (*domain.UserUpdate, error)
	DeleteUser(ctx context.Context, userID primitive.ObjectID) error
	FetchUsers(ctx context.Context, pq *utils.PaginationQuery) (*domain.UserList, error)
	FindByID(ctx context.Context, userID primitive.ObjectID) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	AddMonster(ctx context.Context, userID, monsterID primitive.ObjectID) error
}
