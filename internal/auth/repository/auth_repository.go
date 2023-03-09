package repository

import (
	"context"

	"github.com/iamaul/go-pokedex/internal/auth"
	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/pkg/db/mongodb"
	httpErr "github.com/iamaul/go-pokedex/pkg/error"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepo struct {
	db *mongo.Collection
}

func NewAuthRepo(db *mongo.Database) auth.Repository {
	return &AuthRepo{
		db: db.Collection("users"),
	}
}

func (r *AuthRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	result, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return nil, errors.Wrap(err, httpErr.ErrUserAlreadyExists)
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, err
}

func (r *AuthRepo) UpdateUser(ctx context.Context, user *domain.UserUpdate) (*domain.UserUpdate, error) {
	updateQuery := bson.M{}

	if user.Username != "" {
		updateQuery["username"] = user.Username
	}

	_, err := r.db.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": updateQuery})
	return user, err
}

func (r *AuthRepo) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": userID})

	return err
}

func (r *AuthRepo) FetchUsers(ctx context.Context, pq *utils.PaginationQuery) (*domain.UserList, error) {
	totalCount, err := r.db.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "db.CountDocuments")
	}

	if totalCount == 0 {
		return &domain.UserList{
			TotalCount: 0,
			TotalPages: 0,
			Page:       0,
			Size:       0,
			HasMore:    false,
			Users:      make([]*domain.User, 0),
		}, nil
	}

	limit := int64(pq.GetLimit())
	skip := int64(pq.GetOffset())
	cursor, err := r.db.Find(ctx, bson.D{}, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		return nil, errors.Wrap(err, "db.Find")
	}
	defer cursor.Close(ctx)

	users := make([]*domain.User, 0, pq.GetSize())
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, errors.Wrap(err, "cursor.Decode")
		}
		user.SanitizePassword()
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return &domain.UserList{
		TotalCount: int(totalCount),
		TotalPages: utils.GetTotalPages(int(totalCount), pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), int(totalCount), pq.GetSize()),
		Users:      users,
	}, nil
}

func (r *AuthRepo) FindByID(ctx context.Context, userID primitive.ObjectID) (*domain.User, error) {
	var user domain.User

	if err := r.db.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &domain.User{}, errors.Wrap(err, httpErr.ErrNotFound)
		}

		return &domain.User{}, err
	}

	return &user, nil
}

func (r *AuthRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	return &user, err
}

func (r *AuthRepo) AttachMonster(ctx context.Context, userID, monsterID primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$push": bson.M{"monsters": monsterID}})

	return err
}
