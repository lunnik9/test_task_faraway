package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbName   = "statistics"
	collName = "user_stat"
)

// as long as we have no idea on optimal difficulty of an algorithm,
// keeping track on statistic of users solving challenge is extremely important
// for further optimisations. I chose mongo because it has no strict requirements for document structure,
// which is perfect for our situation, when parameters we keep track on can be changed without altering table,
// disturbing indexes and so on. for further usage, I'd prefer transition to some more appropriate db like clickhouse

type UserStat struct {
	ID                  primitive.ObjectID `bson:"_id"`
	OS                  string             `bson:"os"`
	Arch                string             `bson:"arch"`
	NumCPU              int                `bson:"num_cpu"`
	UserID              int64              `bson:"user_id"`
	PersonalDataAllowed bool               `bson:"personal_data_allowed"`
	TimeSpent           int64              `bson:"time_spent"`
	Challenge           string             `bson:"challenge"`
	Answer              int64              `bson:"answer"`
}

type StatModule interface {
	SaveUserStat(ctx context.Context, stat UserStat) error
}

type StatModuleStruct struct {
	client *mongo.Client
}

func NewStatModuleStruct(client *mongo.Client) *StatModuleStruct {
	return &StatModuleStruct{client: client}
}

func (s *StatModuleStruct) SaveUserStat(ctx context.Context, stat UserStat) error {
	stat.ID = primitive.NewObjectID()

	_, err := s.client.Database(dbName).Collection(collName).InsertOne(ctx, stat)
	if err != nil {
		return err
	}

	return nil
}
