package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mintegral-official/juno/conf"
	"github.com/Mintegral-official/juno/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Mongo struct {
	cfg        *conf.MongoCfg
	client     *mongo.Client
	collection *mongo.Collection
	cursor     *mongo.Cursor
	findOpt    *options.FindOptions
	results    []*CampaignInfo
}

func NewMongo(mongoCfg *conf.MongoCfg) (*Mongo, error) {
	if mongoCfg == nil {
		return nil, helpers.MongoCfgError
	}
	m := &Mongo{
		cfg: mongoCfg,
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Duration(mongoCfg.ConnectTimeout)*time.Microsecond)
	opt := options.Client().ApplyURI(mongoCfg.URI)
	opt.SetReadPreference(readpref.SecondaryPreferred())

	direct := true
	opt.Direct = &direct
	client, err := mongo.Connect(ctx, opt)

	if err != nil {
		return nil, err
	}

	m.client = client
	m.findOpt = options.MergeFindOptions(mongoCfg.FindOpt)
	d := time.Duration(mongoCfg.ReadTimeout) * time.Microsecond
	m.findOpt.MaxTime = &d

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	m.collection = client.Database(mongoCfg.DB).Collection(mongoCfg.Collection)
	if m.collection == nil {
		return nil, errors.New(fmt.Sprintf("[%s.%s] Not found", mongoCfg.DB, mongoCfg.Collection))
	}
	return m, nil
}

func (m *Mongo) Find() ([]*CampaignInfo, error) {

	findOptions := options.Find()
	cur, err := m.collection.Find(context.TODO(), bson.M{"status": 1}, findOptions)
	if err != nil {
		return nil, helpers.CollectionNotFound
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var ele CampaignInfo
		if err := cur.Decode(&ele); err != nil {
			log.Println(err)
			continue
		}
		m.results = append(m.results, &ele)
	}
	if err := cur.Err(); err != nil {
		return nil, helpers.CursorError
	}
	return m.results, nil
}
