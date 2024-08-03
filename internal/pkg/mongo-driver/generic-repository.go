package mongodriver

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GenericRepository[T any] struct {
	collection *mongo.Collection
}

func NewGenericRepository[T any](db *mongo.Database, collectionName string) *GenericRepository[T] {
	collection := db.Collection(collectionName)
	return &GenericRepository[T]{
		collection: collection,
	}
}

func (r *GenericRepository[T]) Add(ctx context.Context, entity *T) error {
	_, err := r.collection.InsertOne(ctx, entity)
	return err
}

func (r *GenericRepository[T]) AddAll(ctx context.Context, entities []*T) error {
	var documents []interface{}
	for _, entity := range entities {
		documents = append(documents, entity)
	}

	_, err := r.collection.InsertMany(ctx, documents)
	return err
}

func (r *GenericRepository[T]) GetById(ctx context.Context, id interface{}) (*T, error) {
	var result T
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Document not found
		}
		return nil, err // Other error
	}
	return &result, nil
}

// Get retrieves a single document from the collection based on a filter.
func (r *GenericRepository[T]) Get(ctx context.Context, filter interface{}) (*T, error) {
	var result T
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Document not found
		}
		return nil, err // Other error
	}
	return &result, nil
}

// GetAll retrieves all documents from the collection based on a filter and options.
func (r *GenericRepository[T]) GetAll(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*T, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var results []*T
	for cursor.Next(ctx) {
		var result T
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *GenericRepository[T]) UpdateById(ctx context.Context, id interface{}, update interface{}) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *GenericRepository[T]) DeleteById(ctx context.Context, id interface{}) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
