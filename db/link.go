package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert Insert Data to Collection
func Insert(collection string, data ...interface{}) (interface{}, error) {

	coll := DB.Collection(collection)

	id, err := coll.InsertMany(context.TODO(), data)
	if err != nil {
		return nil, err
	}

	return id, nil
}

// Find Find Data from Collection by Filter and Option
func Find(collection string, filter interface{}, option *options.FindOptions) (*mongo.Cursor, error) {

	coll := DB.Collection(collection)

	cur, err := coll.Find(context.TODO(), filter, option)
	if err != nil {
		return nil, err
	}

	return cur, nil
}

// Update Update Data in Collection by Filter and Option
func Update(collection string, filter, update interface{}, option *options.UpdateOptions) (*mongo.UpdateResult, error) {
	coll := DB.Collection(collection)

	result, err := coll.UpdateMany(context.TODO(), filter, update, option)

	if err != nil {
		return nil, err
	}

	return result, err
}

// Delete Remove Data from Collection by Filter
func Delete(collection string, filter interface{}, deleteOne bool, option *options.DeleteOptions) (*mongo.DeleteResult, error) {

	coll := DB.Collection(collection)
	var err error
	var result *mongo.DeleteResult
	if deleteOne {
		result, err = coll.DeleteOne(context.TODO(), filter, option)
	} else {
		result, err = coll.DeleteMany(context.TODO(), filter, option)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ParseBson(model interface{}) bson.M {
	b, _ := bson.Marshal(model)
	var body bson.M
	bson.Unmarshal(b, &body)
	return bson.M{"$set": body}
}
