package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FriendsMongo struct {
	cache *redis.Client
	mdb   *mongo.Collection
	db    *sqlx.DB
}

func NewFriendsMongo(cache *redis.Client, mdb *mongo.Collection, db *sqlx.DB) *FriendsMongo {
	return &FriendsMongo{
		cache: cache,
		mdb:   mdb,
		db:    db,
	}
}

func (r *FriendsMongo) GetFriends(id int) ([]MongoFriend, error) {
	result, err := r.cache.Get(context.Background(), fmt.Sprintf("%x", id)).Result()
	// cursor, er := r.mdb.Find(context.Background(), bson.D{{"Friends", bson.D{{"$exists", true}}}})

	// if er != nil {
	// 	return nil, er
	// }
	// users := make([]MongoUser, 0)
	// defer cursor.Close(context.Background())
	// for cursor.Next(context.Background()) {
	// 	var user MongoUser
	// 	if err := cursor.Decode(&user); err != nil {
	// 		return nil, err
	// 	}
	// 	users = append(users, user)
	// }
	// fmt.Println(users)
	// fmt.Println(users[0].Id, reflect.TypeOf(users[0].Id))
	if err != nil {
		filter := bson.D{{"id", id}}
		var user MongoUser
		err := r.mdb.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(user.Friends)
		if err != nil {
			return nil, err
		}
		err = r.cache.Set(context.Background(), fmt.Sprintf("%x", id), data, 5*24*time.Hour).Err()
		if err != nil {
			return nil, err
		}
		return user.Friends, nil
	}
	var userFriends []MongoFriend
	err = json.Unmarshal([]byte(result), &userFriends)
	if err != nil {
		return nil, err
	}
	return userFriends, nil
}

func (r *FriendsMongo) AddFriend(userId int, friendId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var userName string
	query := fmt.Sprintf("SELECT name FROM %s WHERE id=$1", usersTable)
	row := tx.QueryRow(query, userId)
	err = row.Scan(&userName)
	if err != nil {
		tx.Rollback()
		return err
	}
	var friendName string
	query = fmt.Sprintf("SELECT name FROM %s WHERE id=$1", usersTable)
	row = tx.QueryRow(query, friendId)
	err = row.Scan(&friendName)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE sender = $1 AND recipient = $2", requestsTable)
	_, err = tx.Exec(query, friendId, userId)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	var user MongoUser
	filter := bson.D{{"id", userId}}
	err = r.mdb.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return err
	}
	userFriends := append(user.Friends, MongoFriend{Name: friendName, Id: friendId})
	update := bson.D{
		{"$set", bson.D{
			{"friends", userFriends},
		}},
	}
	_, err = r.mdb.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	_, err = r.cache.Get(context.Background(), fmt.Sprintf("%x", userId)).Result()
	if err == nil {
		data, err := json.Marshal(userFriends)
		if err != nil {
			return err
		}
		err = r.cache.Set(context.Background(), fmt.Sprintf("%x", userId), data, 5*24*time.Hour).Err()
		if err != nil {
			return err
		}
	}
	var friend MongoUser
	filter = bson.D{{"id", friendId}}
	err = r.mdb.FindOne(context.Background(), filter).Decode(&friend)
	if err != nil {
		return err
	}
	friendFriends := append(friend.Friends, MongoFriend{Name: userName, Id: userId})
	update = bson.D{
		{"$set", bson.D{
			{"friends", friendFriends},
		}},
	}
	_, err = r.mdb.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	_, err = r.cache.Get(context.Background(), fmt.Sprintf("%x", friendId)).Result()
	if err == nil {
		data, err := json.Marshal(friendFriends)
		if err != nil {
			return err
		}
		err = r.cache.Set(context.Background(), fmt.Sprintf("%x", friendId), data, 5*24*time.Hour).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *FriendsMongo) DeleteFriend(userId int, friendId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	var userName string
	query := fmt.Sprintf("SELECT name FROM %s WHERE id=$1", usersTable)
	row := tx.QueryRow(query, userId)
	err = row.Scan(&userName)
	if err != nil {
		tx.Rollback()
		return err
	}
	var friendName string
	query = fmt.Sprintf("SELECT name FROM %s WHERE id=$1", usersTable)
	row = tx.QueryRow(query, friendId)
	err = row.Scan(&friendName)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	var user MongoUser
	filter := bson.D{{"id", userId}}
	err = r.mdb.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return err
	}
	userFriends := make([]MongoFriend, 0)
	for _, value := range user.Friends {
		if value.Id != friendId {
			userFriends = append(userFriends, value)
		}
	}
	update := bson.D{
		{"$set", bson.D{
			{"friends", userFriends},
		}},
	}
	_, err = r.mdb.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	_, err = r.cache.Get(context.Background(), fmt.Sprintf("%x", userId)).Result()
	if err == nil {
		data, err := json.Marshal(userFriends)
		if err != nil {
			return err
		}
		err = r.cache.Set(context.Background(), fmt.Sprintf("%x", userId), data, 5*24*time.Hour).Err()
		if err != nil {
			return err
		}
	}
	var friend MongoUser
	filter = bson.D{{"id", friendId}}
	err = r.mdb.FindOne(context.Background(), filter).Decode(&friend)
	if err != nil {
		return err
	}
	friendFriends := make([]MongoFriend, 0)
	for _, value := range friend.Friends {
		if value.Id != userId {
			friendFriends = append(friendFriends, value)
		}
	}
	update = bson.D{
		{"$set", bson.D{
			{"friends", friendFriends},
		}},
	}
	_, err = r.mdb.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	_, err = r.cache.Get(context.Background(), fmt.Sprintf("%x", friendId)).Result()
	if err == nil {
		data, err := json.Marshal(friendFriends)
		if err != nil {
			return err
		}
		err = r.cache.Set(context.Background(), fmt.Sprintf("%x", friendId), data, 5*24*time.Hour).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
