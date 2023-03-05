package repository

import (
	"context"

	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/internal/monster"
	"github.com/iamaul/go-pokedex/pkg/db/mongodb"
	httpErr "github.com/iamaul/go-pokedex/pkg/error"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MonsterTypeRepo struct {
	db *mongo.Collection
}

func NewMonsterTypeRepo(db *mongo.Database) monster.MonsterTypeRepository {
	return &MonsterTypeRepo{
		db: db.Collection("monster_types"),
	}
}

func (r *MonsterTypeRepo) CreateMonsterType(ctx context.Context, monsterType *domain.MonsterType) (*domain.MonsterType, error) {
	_, err := r.db.InsertOne(ctx, monsterType)
	if mongodb.IsDuplicate(err) {
		return nil, errors.Wrap(err, httpErr.ErrEmailAlreadyExists)
	}

	return monsterType, err
}

func (r *MonsterTypeRepo) UpdateMonsterType(ctx context.Context, monsterType *domain.MonsterTypeUpdate) (*domain.MonsterTypeUpdate, error) {
	updateQuery := bson.M{}

	if monsterType.Name != "" {
		updateQuery["name"] = monsterType.Name
	}

	_, err := r.db.UpdateOne(ctx, bson.M{"_id": monsterType.ID}, bson.M{"$set": updateQuery})
	return monsterType, err
}

func (r *MonsterTypeRepo) DeleteMonsterType(ctx context.Context, monsterTypeID primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": monsterTypeID})

	return err
}

func (r *MonsterTypeRepo) FetchMonsterTypes(ctx context.Context, pq *utils.PaginationQuery) (*domain.MonsterTypeList, error) {
	totalCount, err := r.db.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "db.CountDocuments")
	}

	if totalCount == 0 {
		return &domain.MonsterTypeList{
			TotalCount:   0,
			TotalPages:   0,
			Page:         0,
			Size:         0,
			HasMore:      false,
			MonsterTypes: make([]*domain.MonsterType, 0),
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

	monsterTypes := make([]*domain.MonsterType, 0, pq.GetSize())
	for cursor.Next(ctx) {
		var monsterType domain.MonsterType
		if err := cursor.Decode(&monsterType); err != nil {
			return nil, errors.Wrap(err, "cursor.Decode")
		}
		monsterTypes = append(monsterTypes, &monsterType)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return &domain.MonsterTypeList{
		TotalCount:   int(totalCount),
		TotalPages:   utils.GetTotalPages(int(totalCount), pq.GetSize()),
		Page:         pq.GetPage(),
		Size:         pq.GetSize(),
		HasMore:      utils.GetHasMore(pq.GetPage(), int(totalCount), pq.GetSize()),
		MonsterTypes: monsterTypes,
	}, nil
}

func (r *MonsterTypeRepo) FindByID(ctx context.Context, monsterTypeID primitive.ObjectID) (*domain.MonsterType, error) {
	var monsterType domain.MonsterType

	if err := r.db.FindOne(ctx, bson.M{"_id": monsterTypeID}).Decode(&monsterType); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &domain.MonsterType{}, errors.Wrap(err, httpErr.ErrNotFound)
		}

		return &domain.MonsterType{}, err
	}

	return &monsterType, nil
}

func (r *MonsterTypeRepo) FindByName(ctx context.Context, monsterTypeName string) (*domain.MonsterType, error) {
	var monsterType domain.MonsterType

	err := r.db.FindOne(ctx, bson.M{"name": monsterTypeName}).Decode(&monsterType)
	if err != nil {
		return &domain.MonsterType{}, err
	}

	return &monsterType, nil
}
