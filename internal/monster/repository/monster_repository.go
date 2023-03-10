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

type MonsterRepo struct {
	db *mongo.Collection
}

func NewMonsterRepo(db *mongo.Database) monster.MonsterRepository {
	return &MonsterRepo{
		db: db.Collection("monsters"),
	}
}

func (r *MonsterRepo) CreateMonster(ctx context.Context, monster *domain.Monster) (*domain.Monster, error) {
	result, err := r.db.InsertOne(ctx, monster)
	if mongodb.IsDuplicate(err) {
		return nil, errors.Wrap(err, httpErr.ErrUserAlreadyExists)
	}

	monster.ID = result.InsertedID.(primitive.ObjectID)

	return monster, err
}

func (r *MonsterRepo) UpdateMonster(ctx context.Context, monster *domain.MonsterUpdate) (*domain.MonsterUpdate, error) {
	updateQuery := bson.M{}

	if monster.Name != "" {
		updateQuery["name"] = monster.Name
	}

	_, err := r.db.UpdateOne(ctx, bson.M{"_id": monster.ID}, bson.M{"$set": updateQuery})
	return monster, err
}

func (r *MonsterRepo) DeleteMonster(ctx context.Context, monsterID primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": monsterID})

	return err
}

func (r *MonsterRepo) FetchMonsters(ctx context.Context, pq *utils.PaginationQuery) (*domain.MonsterList, error) {
	totalCount, err := r.db.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "db.CountDocuments")
	}

	if totalCount == 0 {
		return &domain.MonsterList{
			TotalCount: 0,
			TotalPages: 0,
			Page:       0,
			Size:       0,
			HasMore:    false,
			Monsters:   make([]*domain.Monster, 0),
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

	monsters := make([]*domain.Monster, 0, pq.GetSize())
	for cursor.Next(ctx) {
		var monster domain.Monster
		if err := cursor.Decode(&monster); err != nil {
			return nil, errors.Wrap(err, "cursor.Decode")
		}
		monsters = append(monsters, &monster)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return &domain.MonsterList{
		TotalCount: int(totalCount),
		TotalPages: utils.GetTotalPages(int(totalCount), pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), int(totalCount), pq.GetSize()),
		Monsters:   monsters,
	}, nil
}

func (r *MonsterRepo) FindByID(ctx context.Context, monsterID primitive.ObjectID) (*domain.Monster, error) {
	var monster domain.Monster

	if err := r.db.FindOne(ctx, bson.M{"_id": monsterID}).Decode(&monster); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &domain.Monster{}, errors.Wrap(err, httpErr.ErrNotFound)
		}

		return &domain.Monster{}, err
	}

	return &monster, nil
}

func (r *MonsterRepo) FindByName(ctx context.Context, monsterName string) (*domain.Monster, error) {
	var monster domain.Monster

	err := r.db.FindOne(ctx, bson.M{"name": monsterName}).Decode(&monster)

	return &monster, err
}

func (r *MonsterRepo) AddMonsterType(ctx context.Context, monsterID, monsterTypeID primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": monsterID}, bson.M{"$push": bson.M{"monster_types": monsterTypeID}})

	return err
}
