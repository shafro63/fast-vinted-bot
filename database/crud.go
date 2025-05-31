package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"fast-vinted-bot/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUser(db string, coll string, data *utils.DiscordUserData) error {
	users := Client.Database(db).Collection(coll)
	m := data.Member
	filter := bson.M{"discord_id": m.User.ID}

	err := users.FindOne(context.TODO(), filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newUser := &utils.User{
				DiscordID: m.User.ID,
				Username:  m.User.Username,
				Channels:  make(map[string]*utils.ChannelInfo),
			}
			newUser.Guild = &utils.GuildInfo{
				GuildID:  data.GuildID,
				JoinedAt: m.JoinedAt,
			}
			_, err := users.InsertOne(context.TODO(), newUser)
			if err != nil {
				return fmt.Errorf("error while creating user %v : %v", newUser.Username, err)
			}
			slog.Info("user created", "user", m.User.Username, "database", db)
			return nil
		} else {
			return fmt.Errorf("error while retrieving user data %v : %v", m.User.Username, err)
		}
	}

	slog.Debug("Create user success", "db", db, "user", m.User.ID)
	return nil
}

func GetUserByID(db string, coll string, userID string) (*utils.User, error) {
	users := Client.Database(db).Collection(coll)
	filter := bson.M{"_id": userID}

	var user utils.User
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func GetUser(db string, coll string, data *utils.DiscordUserData) (*utils.User, error) {
	users := Client.Database(db).Collection(coll)
	u := data.Member.User
	filter := bson.M{"discord_id": u.ID}

	var user utils.User
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func SetChannel(db string, coll string, data *utils.DiscordUserData) error {
	users := Client.Database(db).Collection(coll)
	u := data.Member.User
	filter := bson.M{"discord_id": u.ID}
	fieldpath := "channels." + data.ChannelID

	update := bson.M{
		"$set": bson.M{
			fieldpath: bson.D{
				{Key: "name", Value: data.ChannelName},
				{Key: "links", Value: bson.M{}},
			},
		},
	}

	upsert := options.Update().SetUpsert(true)
	_, err := users.UpdateOne(context.TODO(), filter, update, upsert)
	if err != nil {
		return fmt.Errorf("error while creating channel : %v", err)
	}

	slog.Debug("Set channel success", "db", db, "user", u.ID, "channel", data.ChannelID)
	return nil
}

func DeleteChannel(db string, coll string, data *utils.DiscordUserData) error {
	users := Client.Database(db).Collection(coll)
	u := data.Member.User
	fieldpath := "channels." + data.ChannelID
	filter := bson.M{
		"discord_id": u.ID,
		fieldpath:    bson.M{"$exists": true},
	}

	update := bson.M{
		"$unset": bson.M{
			fieldpath: "",
		},
	}

	_, err := users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error while deleting channel : %v", err)
	}

	slog.Debug("Delete Channel success", "db", db, "user", u.ID, "channel", data.ChannelID)
	return nil
}

func GetChannels(db string, coll string, data *utils.DiscordUserData) (map[string]*utils.ChannelInfo, error) {
	users := Client.Database(db).Collection(coll)
	u := data.Member.User
	filter := bson.M{"discord_id": u.ID}

	var user utils.User
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		} else {
			return nil, fmt.Errorf("failed to get user data : %v", err)
		}
	}
	if user.Channels == nil {
		user.Channels = make(map[string]*utils.ChannelInfo)
	}

	slog.Debug("Get channels success", "db", db, "user", u.ID)
	return user.Channels, nil
}

func SetLink(db string, coll string, data *utils.DiscordUserData) error {
	users := Client.Database(db).Collection(coll)
	u := data.Member.User
	channelpath := "channels." + data.ChannelID
	filter := bson.M{
		"discord_id": u.ID,
		channelpath:  bson.M{"$exists": true},
	}

	fieldpath := "channels." + data.ChannelID + ".links." + data.LinkName
	update := bson.M{
		"$set": bson.M{
			fieldpath: data.Link,
		},
	}

	upsert := options.Update().SetUpsert(true)
	_, err := users.UpdateOne(context.TODO(), filter, update, upsert)
	if err != nil {
		return fmt.Errorf("error while adding link to user channel in database : %v", err)
	}

	slog.Debug("Set link success", "db", db, "user", u.ID, "link", data.LinkName)
	return nil
}

func GetLinks(db string, coll string, data *utils.DiscordUserData) (map[string]string, error) {
	users := Client.Database(db).Collection(coll)
	u := data.Member.User

	fieldpath := "channels." + data.ChannelID
	filter := bson.M{
		"discord_id": u.ID,
		fieldpath:    bson.M{"$exists": true},
	}

	user := &utils.User{}
	err := users.FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("channel not found")
		} else {
			return nil, fmt.Errorf("failed to get user data : %v", err)
		}
	}

	userLinks := user.Channels[data.ChannelID].Links
	if userLinks == nil {
		userLinks = make(map[string]string)
	}

	slog.Debug("Get Links Success", "db", db, "user", u.ID)
	return userLinks, nil
}

func DeleteLink(db string, coll string, data *utils.DiscordUserData) error {
	users := Client.Database(db).Collection(coll)
	u := data.Member.User
	fieldpath := "channels." + data.ChannelID + ".links." + data.LinkName
	filter := bson.M{
		"discord_id": u.ID,
	}

	update := bson.M{
		"$unset": bson.M{
			fieldpath: "",
		},
	}

	_, err := users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("link not found")
		} else {
			return fmt.Errorf("failed to delete user data : %v", err)
		}
	}

	slog.Debug("Delete link success", "db", db, "user", u.ID, "link", data.LinkName)
	return nil
}

func GetAllActiveChannels(db string, coll string) ([]utils.SessionData, error) {
	users := Client.Database(db).Collection(coll)
	filter := bson.M{
		"discord_id": bson.M{"$exists": true},
		"channels":   bson.M{"$exists": true},
	}
	project := bson.M{
		"_id":      0,
		"channels": 1,
	}
	opts := options.Find().SetProjection(project)

	cursor, err := users.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var docs []utils.SessionData
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		return nil, err
	}

	slog.Debug("Get active channels success")
	return docs, nil
}
