package mongo

import (
	"context"
	"fmt"
	"go-auth-user/domain"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(client *mongo.Client, dbName, collectionName string) *UserRepo {
	coll := client.Database(dbName).Collection(collectionName)
	if err := ensureUserEmailIndex(coll); err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
	return &UserRepo{collection: coll}
}

func (r *UserRepo) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, bson.M{
		"id":       user.Id,
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	})
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}
	if err != nil {
		return err
	}

	seq, err := getNextSequenceValue(ctx, r.collection.Database(), "users")
	if err != nil {
		return err
	}

	user.Id = int(seq)
	_, err = r.collection.UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": bson.M{"id": user.Id}})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) Find(id int) *domain.User {
	var user struct {
		Id    int    `bson:"id"`
		Name  string `bson:"name"`
		Email string `bson:"email"`
	}

	err := r.collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil
	}
	return &domain.User{
		Id:    int(user.Id),
		Name:  user.Name,
		Email: user.Email,
	}
}

func (r *UserRepo) FindByEmail(email string) *domain.User {
	var user struct {
		Id       int    `bson:"id"`
		Name     string `bson:"name"`
		Email    string `bson:"email"`
		Password string `bson:"password"`
	}

	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil
	}
	return &domain.User{
		Id:       int(user.Id),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
}

func (r *UserRepo) FindAll(page int, perPage int) []domain.User {
	ctx := context.Background()

	opts := options.Find().
		SetSkip((int64(page) - 1) * int64(perPage)).
		SetLimit(int64(perPage)).
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetProjection(bson.M{
			"password": 0,
		})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	var users []domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil
	}
	return users
}

func (r *UserRepo) Update(id int, name string, email string) (result *mongo.UpdateResult, err error) {
	ctx := context.Background()

	updateData := bson.M{}
	if name != "" {
		updateData["name"] = name
	}
	if email != "" {
		updateData["email"] = email
	}
	result, err = r.collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": updateData})
	if mongo.IsDuplicateKeyError(err) {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserRepo) Delete(id int) error {
	ctx := context.Background()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *UserRepo) CountUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{})
}

func getNextSequenceValue(ctx context.Context, db *mongo.Database, sequenceName string) (int32, error) {
	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updatedDoc bson.M
	err := db.Collection("counters").FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDoc)
	if err != nil {
		return 0, err
	}

	return updatedDoc["seq"].(int32), nil
}

func ensureUserEmailIndex(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, index)
	return err
}
