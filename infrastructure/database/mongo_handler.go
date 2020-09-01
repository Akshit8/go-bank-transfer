package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoHandler struct {
	db      *mongo.Database
	session mongo.Session
}

func NewMongoHandler(c *config) (*mongoHandler, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.ctxTimeout)
	defer cancel()

	uri := fmt.Sprintf(
		"%s://@%s",
		c.host,
		c.host,
	)

	clientOpts := options.Client().ApplyURI(uri).SetDirect(true)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}

	return &mongoHandler{
		db:      client.Database(c.database),
		session: session,
	}, nil
}

func (mgo mongoHandler) Store(ctx context.Context, collection string, data interface{}) error {
	if _, err := mgo.db.Collection(collection).InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}

func (mgo mongoHandler) Update(ctx context.Context, collection string, query interface{}, update interface{}) error {
	if _, err := mgo.db.Collection(collection).UpdateOne(ctx, query, update); err != nil {
		return err
	}

	return nil
}

func (mgo mongoHandler) FindAll(ctx context.Context, collection string, query interface{}, result interface{}) error {
	cur, err := mgo.db.Collection(collection).Find(ctx, query)
	if err != nil {
		return err
	}

	defer cur.Close(ctx)
	if err = cur.All(ctx, result); err != nil {
		return err
	}

	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}

func (mgo mongoHandler) FindOne(
	ctx context.Context,
	collection string,
	query interface{},
	projection interface{},
	result interface{},
) error {
	var err = mgo.db.Collection(collection).
		FindOne(
			ctx,
			query,
			options.FindOne().SetProjection(projection),
		).Decode(result)
	if err != nil {
		return err
	}

	return nil
}
