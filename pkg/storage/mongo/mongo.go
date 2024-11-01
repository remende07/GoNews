package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	db *mongo.Client
}

type Authors struct {
	ID   int    `bson:"_id"`
	Name string `bson:"name"`
}

type Counter struct {
	ID  int
	Seq int `bson:"seq"`
}

func New(connstr string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(connstr)
	client, err := mongo.Connect(context.Background(), mongoOpts)

	if err != nil {
		return nil, err
	}

	s := Store{
		db: client,
	}

	return &s, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.db.Database("posts").Collection("posts")
	filter := bson.D{}

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var posts []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (s *Store) AddPost(post storage.Post) error {
	collection := s.db.Database("posts").Collection("posts")
	post.CreatedAt = time.Now().Unix()

	authors := s.db.Database("posts").Collection("authors")
	filter := bson.D{{Key: "_id", Value: post.AuthorID}}

	res := authors.FindOne(context.Background(), filter)
	log.Println(res)
	var author Authors
	err := res.Decode(&author)
	if err != nil {
		return err
	}
	post.AuthorName = author.Name

	id, err := s.getNextSeq("posts")
	if err != nil {
		return err
	}
	post.ID = id

	_, err = collection.InsertOne(context.Background(), post)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdatePost(post storage.Post) error {
	collection := s.db.Database("posts").Collection("posts")

	authors := s.db.Database("posts").Collection("authors")
	authorsFilter := bson.D{{Key: "_id", Value: post.AuthorID}}

	res := authors.FindOne(context.Background(), authorsFilter)
	log.Println(res)
	var author Authors
	err := res.Decode(&author)
	if err != nil {
		return err
	}
	post.AuthorName = author.Name

	filter := bson.D{{Key: "_id", Value: post.ID}}
	update := bson.D{{Key: "$set", Value: post}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeletePost(post storage.Post) error {
	collection := s.db.Database("posts").Collection("posts")
	filter := bson.D{{Key: "_id", Value: post.ID}}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) getNextSeq(name string) (int, error) {
	collection := s.db.Database("posts").Collection("counters")
	update := bson.M{
		"$inc": bson.M{
			"seq": 1,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	var counter Counter
	err := collection.FindOneAndUpdate(context.Background(), bson.M{"_id": name}, update, opts).Decode(&counter)
	if err != nil {
		return 0, err
	}
	return counter.Seq, nil
}
