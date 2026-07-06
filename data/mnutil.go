// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/donnie4w/gofer/hashmap"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tim/sys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
)

func connectToMongo(host, port, username, password, authDB string) (*mongo.Client, error) {
	var uri string
	if username != "" && password != "" {
		if authDB != "" {
			uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=%s", username, password, host, port, authDB)
		} else {
			uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
		}
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", host, port)
	}
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type mnClient struct {
	client        *mongo.Client
	dbname        string
	extent        int
	collectionmap *hashmap.Map[string, *mongo.Collection]
}

type mongoManager struct {
	clients []*mnClient
}

func newMongoManager() *mongoManager {
	return &mongoManager{clients: make([]*mnClient, 0)}
}

func (manager *mongoManager) init() error {
	if len(sys.Conf.MongodbExtent) > 0 {
		for _, me := range sys.Conf.MongodbExtent {
			if cli, err := connectToMongo(me.Host, me.Port, me.Username, me.Password, me.AuthDB); err == nil {
				extent := sys.MB
				if me.ExtentMax > 0 {
					extent = me.ExtentMax
				}
				manager.clients = append(manager.clients, &mnClient{extent: extent, client: cli, dbname: me.DbName, collectionmap: hashmap.NewMap[string, *mongo.Collection]()})
				return manager.initTable(cli.Database(me.DbName))
			} else {
				return err
			}
		}
		sort.Slice(manager.clients, func(i, j int) bool { return manager.clients[i].extent < manager.clients[j].extent })
	} else if me := sys.Conf.Mongodb; me != nil {
		if cli, err := connectToMongo(me.Host, me.Port, me.Username, me.Password, me.AuthDB); err == nil {
			extent := sys.MB
			if me.ExtentMax > 0 {
				extent = me.ExtentMax
			}
			manager.clients = append(manager.clients, &mnClient{extent: extent, client: cli, dbname: me.DbName, collectionmap: hashmap.NewMap[string, *mongo.Collection]()})
			return manager.initTable(cli.Database(me.DbName))
		} else {
			return fmt.Errorf("%s", "mongo init error:"+err.Error())
		}
	}
	return errors.New("mongo init error")
}

type Index struct {
	Key    map[string]int `json:"key"`
	Unique bool           `json:"unique,omitempty"`
}

type CollectionConfig struct {
	Collection string  `json:"collection"`
	Indexes    []Index `json:"indexes"`
}

func (manager *mongoManager) initTable(db *mongo.Database) error {
	for _, defJSON := range MongoDB("").CreateSql() {
		var config CollectionConfig
		if err := json.Unmarshal([]byte(defJSON), &config); err != nil {
			return err
		}

		if collections, err := db.ListCollectionNames(context.Background(), bson.M{}); err == nil {
			if contains(collections, config.Collection) {
				continue
			}
		} else {
			return err
		}

		collection := db.Collection(config.Collection)

		existingIndexes, err := collection.Indexes().List(context.Background())
		if err != nil {
			return err
		}
		idxmap := parseIndexs(existingIndexes)
		for _, idx := range config.Indexes {
			if len(idxmap) > 0 {
				exist := false
				for key := range idx.Key {
					for _, m := range idxmap {
						if _, ok := m[key]; ok {
							exist = true
							break
						}
					}
					break
				}
				if exist {
					continue
				}
			}
			keys := bson.D{}
			for key, value := range idx.Key {
				keys = append(keys, bson.E{Key: key, Value: value})
			}
			if len(keys) > 0 {
				_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
					Keys:    keys,
					Options: options.Index().SetUnique(idx.Unique),
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func contains(collections []string, collectionName string) bool {
	for _, c := range collections {
		if c == collectionName {
			return true
		}
	}
	return false
}

func parseIndexs(existingIndexes *mongo.Cursor) []map[string]interface{} {
	defer existingIndexes.Close(context.Background())
	r := make([]map[string]interface{}, 0)
	for existingIndexes.Next(context.Background()) {
		var idx bson.M
		if err := existingIndexes.Decode(&idx); err == nil {
			r = append(r, idx)
		}
	}
	return r
}

func (manager *mongoManager) getMnClient(tid uint64) *mnClient {
	if len(manager.clients) == 0 {
		return nil
	}
	if len(manager.clients) == 1 || tid == 0 {
		return manager.clients[0]
	}
	idx := 0
	if idx = sort.Search(len(manager.clients), func(i int) bool { return manager.clients[i].extent >= int(tid%sys.MB) }); idx >= len(manager.clients) {
		idx = 0
	}
	return manager.clients[idx]
}

func (manager *mongoManager) getCollection(tid uint64, collectionname string) *mongo.Collection {
	mc := manager.getMnClient(tid)
	if r, ok := mc.collectionmap.Get(collectionname); ok {
		return r
	}
	r := mc.client.Database(mc.dbname).Collection(collectionname)
	mc.collectionmap.Put(collectionname, r)
	return r
}

var manager = newMongoManager()

type timMessage struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Mid        int64              `bson:"mid"`
	ChatID     primitive.Binary   `bson:"chatid"`
	FID        int32              `bson:"fid"`
	Stanza     primitive.Binary   `bson:"stanza"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimMessage(tid uint64) *timMessage {
	return &timMessage{collection: manager.getCollection(tid, "timmessage")}
}

func (t *timMessage) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timMessage) GetById(id primitive.ObjectID) (*timMessage, error) {
	var msg timMessage
	err := t.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&msg)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}

func (t *timMessage) Get(where bson.M, field bson.M) (*timMessage, error) {
	var u timMessage
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timMessage) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timMessage) Delete(chatID []byte, mid int64) (*mongo.DeleteResult, error) {
	return t.collection.DeleteMany(context.TODO(), bson.M{"chatid": binary(chatID), "mid": mid})
}

func (t *timMessage) List(where bson.M, field bson.M) ([]*timMessage, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timMessage
	for cursor.Next(ctx) {
		var msg timMessage
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (t *timMessage) ListOptions(where bson.M, options *options.FindOptions) ([]*timMessage, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*timMessage
	for cursor.Next(ctx) {
		var msg timMessage
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timUser struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UUID       int64              `bson:"uuid"`
	Pwd        string             `bson:"pwd"`
	Createtime int64              `bson:"createtime"`
	Ubean      primitive.Binary   `bson:"ubean"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimUser(tid uint64) *timUser {
	return &timUser{collection: manager.getCollection(tid, "timuser")}
}

func (t *timUser) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timUser) Get(where bson.M, field bson.M) (*timUser, error) {
	var u timUser
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timUser) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timUser) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timUser) List(where bson.M, field bson.M) ([]*timUser, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timUser
	for cursor.Next(ctx) {
		var msg timUser
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timGroup struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Gtype      int32              `bson:"gtype"`
	UUID       int64              `bson:"uuid"`
	Createtime int64              `bson:"createtime"`
	Status     int32              `bson:"status"`
	Rbean      primitive.Binary   `bson:"rbean"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimGroup(tid uint64) *timGroup {
	return &timGroup{collection: manager.getCollection(tid, "timgroup")}
}

func (t *timGroup) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timGroup) Get(where bson.M, field bson.M) (*timGroup, error) {
	var u timGroup
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timGroup) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timGroup) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timGroup) List(where bson.M, field bson.M) ([]*timGroup, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timGroup
	for cursor.Next(ctx) {
		var msg timGroup
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timOffline struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UUID       int64              `bson:"uuid"`
	Chatid     primitive.Binary   `bson:"chatid"`
	Stanza     primitive.Binary   `bson:"stanza"`
	Mid        int64              `bson:"mid"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimOffline(tid uint64) *timOffline {
	return &timOffline{collection: manager.getCollection(tid, "timoffline")}
}

func (t *timOffline) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timOffline) Get(where bson.M, field bson.M) (*timOffline, error) {
	var u timOffline
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timOffline) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timOffline) Delete(id primitive.ObjectID) (result *mongo.DeleteResult, err error) {
	return t.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
}

func (t *timOffline) DeleteMany(filter bson.M) (result *mongo.DeleteResult, err error) {
	return t.collection.DeleteMany(context.TODO(), filter)
}

func (t *timOffline) List(where bson.M, field bson.M) ([]*timOffline, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timOffline
	for cursor.Next(ctx) {
		var msg timOffline
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (t *timOffline) ListOptions(where bson.M, options *options.FindOptions) ([]*timOffline, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timOffline
	for cursor.Next(ctx) {
		var msg timOffline
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timRelate struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UUID       primitive.Binary   `bson:"uuid"`
	Status     int32              `bson:"status"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimRelate(tid uint64) *timRelate {
	return &timRelate{collection: manager.getCollection(tid, "timrelate")}
}

func (t *timRelate) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timRelate) Get(where, field bson.M) (*timRelate, error) {
	var u timRelate
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timRelate) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timRelate) UpdateOption(where bson.M, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateMany(context.TODO(), where, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timRelate) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timRelate) List(where, field bson.M) ([]*timRelate, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timRelate
	for cursor.Next(ctx) {
		var msg timRelate
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timRoster struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Unikid     primitive.Binary   `bson:"unikid"`
	UUID       int64              `bson:"uuid"`
	Tuuid      int64              `bson:"tuuid"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimRoster(tid uint64) *timRoster {
	return &timRoster{collection: manager.getCollection(tid, "timroster")}
}

func (t *timRoster) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timRoster) Get(where, field bson.M) (*timRoster, error) {
	var u timRoster
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timRoster) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timRoster) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timRoster) DeleteOption(filter bson.M) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timRoster) List(where, field bson.M) ([]*timRoster, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timRoster
	for cursor.Next(ctx) {
		var msg timRoster
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timMucroster struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Unikid     primitive.Binary   `bson:"unikid"`
	UUID       int64              `bson:"uuid"`
	Tuuid      int64              `bson:"tuuid"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimMucroster(tid uint64) *timMucroster {
	return &timMucroster{collection: manager.getCollection(tid, "timmucroster")}
}

func (t *timMucroster) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timMucroster) Get(where, field bson.M) (*timMucroster, error) {
	var u timMucroster
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timMucroster) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timMucroster) DeleteOption(filter bson.M) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timMucroster) List(where, field bson.M) ([]*timMucroster, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timMucroster
	for cursor.Next(ctx) {
		var msg timMucroster
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timBlock struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Unikid     primitive.Binary   `bson:"unikid"`
	UUID       int64              `bson:"uuid"`
	Tuuid      int64              `bson:"tuuid"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimBlock(tid uint64) *timBlock {
	return &timBlock{collection: manager.getCollection(tid, "timblock")}
}

func (t *timBlock) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timBlock) Get(where, field bson.M) (*timBlock, error) {
	var u timBlock
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timBlock) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timBlock) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timBlock) List(where, field bson.M) ([]*timBlock, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timBlock
	for cursor.Next(ctx) {
		var msg timBlock
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timBlockroom struct {
	collection *mongo.Collection
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Unikid     primitive.Binary   `bson:"unikid"`
	UUID       int64              `bson:"uuid"`
	Tuuid      int64              `bson:"tuuid"`
	Timeseries int64              `bson:"timeseries"`
}

func newMnTimBlockroom(tid uint64) *timBlockroom {
	return &timBlockroom{collection: manager.getCollection(tid, "timblockroom")}
}

func (t *timBlockroom) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timBlockroom) Get(where, field bson.M) (*timBlockroom, error) {
	var u timBlockroom
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timBlockroom) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timBlockroom) DeleteOption(filter bson.M) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timBlockroom) List(where, field bson.M) ([]*timBlockroom, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timBlockroom
	for cursor.Next(ctx) {
		var msg timBlockroom
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type timDomain struct {
	collection    *mongo.Collection
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Adminaccount  string             `bson:"adminaccount"`
	Adminpassword string             `bson:"adminpassword"`
	Timdomain     string             `bson:"timdomain"`
	Createtime    int64              `bson:"createtime"`
	Timeseries    int64              `bson:"timeseries"`
}

func newMnTimDomain(tid uint64) *timDomain {
	return &timDomain{collection: manager.getCollection(tid, "timdomain")}
}

func (t *timDomain) Create() (*mongo.InsertOneResult, error) {
	t.ID = primitive.NewObjectID()
	result, err := t.collection.InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timDomain) Get(where, field bson.M) (*timDomain, error) {
	var u timDomain
	err := t.collection.FindOne(context.TODO(), where, options.FindOne().SetProjection(field)).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (t *timDomain) Update(id primitive.ObjectID, updates bson.M) (*mongo.UpdateResult, error) {
	result, err := t.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timDomain) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := t.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *timDomain) List(where, field bson.M) ([]*timDomain, error) {
	ctx := context.TODO()
	cursor, err := t.collection.Find(ctx, where, options.Find().SetProjection(field))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*timDomain
	for cursor.Next(ctx) {
		var msg timDomain
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		list = append(list, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func checkuseruuidMn(uuids ...uint64) bool {
	for _, uuid := range uuids {
		idbs := goutil.Int64ToBytes(int64(uuid))
		if uuidCache.Contains(idbs) {
			continue
		}
		if tu, _ := newMnTimUser(uuid).Get(bson.M{"uuid": int64(uuid)}, bson.M{"uuid": 1}); tu != nil && tu.UUID != 0 {
			uuidCache.Add(idbs)
			continue
		} else {
			return false
		}
	}
	return true
}

func binary(data []byte) primitive.Binary {
	return primitive.Binary{Subtype: 0x00, Data: data}
}
