package builder

import (
	"context"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type MongoIndexManager struct {
	ops        *MongoIndexManagerOps
	innerIndex *index.IndexImpl
	totalNum   int64
	errorNum   int64
	client     *mongo.Client
	collection *mongo.Collection
	cursor     *mongo.Cursor
	result     []*ParserResult
	findOpt    *options.FindOptions
}

func NewMongoIndexManager(ops *MongoIndexManagerOps) *MongoIndexManager {
	if ops == nil {
		return nil
	}

	mongoIndexManager := &MongoIndexManager{
		ops:        ops,
		innerIndex: nil,
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Duration(ops.ConnectTimeout)*time.Microsecond)
	opt := options.Client().ApplyURI(ops.URI)
	opt.SetReadPreference(readpref.SecondaryPreferred())

	direct := true
	opt.Direct = &direct
	client, err := mongo.Connect(ctx, opt)

	if err != nil {
		return nil
	}

	mongoIndexManager.client = client
	mongoIndexManager.findOpt = options.MergeFindOptions(ops.FindOpt)
	d := time.Duration(ops.ReadTimeout) * time.Microsecond
	mongoIndexManager.findOpt.MaxTime = &d

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil
	}
	mongoIndexManager.collection = client.Database(ops.DB).Collection(ops.Collection)
	if mongoIndexManager.collection == nil {
		return nil
	}
	return mongoIndexManager
}

func (mim *MongoIndexManager) GetIndex() *index.IndexImpl {
	return mim.innerIndex
}

func (mim *MongoIndexManager) Update(ctx context.Context) error {
	if e := mim.base(ctx); e != nil {
		return e
	}
	for {
		inc := time.After(time.Duration(mim.ops.IncInterval) * time.Second)
		base := time.After(time.Duration(mim.ops.BaseInterval) * time.Second)
		now := time.Now().Unix()
		select {
		case <-ctx.Done():
			return nil
		case <-inc:
			if err := mim.inc(ctx, now); err != nil {
				return err
			}
		case <-base:
			if e := mim.base(ctx); e != nil {
				return e
			}
		}
	}
}

func (mim *MongoIndexManager) base(ctx context.Context) error {
	mim.totalNum = 0
	mim.errorNum = 0
	if mim.ops.OnBeforeBase != nil {
		mim.ops.BaseQuery = mim.ops.OnBeforeBase(mim.ops.UserData)
	}
	context.WithTimeout(ctx, time.Duration(mim.ops.ReadTimeout)*time.Microsecond)
	cur, err := mim.collection.Find(nil, mim.ops.BaseQuery, mim.ops.FindOpt)
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		if cur.Err() != nil {
			mim.errorNum++
			return cur.Err()
		}
		r, err := mim.ops.BaseParser.Parse(cur.Current, true)
		if err != nil {
			mim.errorNum++
			log.Println(err)
			continue
		}
		mim.result = append(mim.result, r)
	}

	var baseIndex = index.NewIndex("base")
	for i := 0; i < len(mim.result); i++ {
		mim.totalNum++
		_ = baseIndex.Add(mim.result[i].Value)
	}
	mim.innerIndex = baseIndex
	return err
}

func (mim *MongoIndexManager) inc(ctx context.Context, now int64) error {
	if mim.ops.OnBeforeInc != nil {
		mim.ops.IncQuery = mim.ops.OnBeforeInc(mim.ops.UserData)
	}
	context.WithTimeout(ctx, time.Duration(mim.ops.ReadTimeout)*time.Microsecond)
	cur, err := mim.collection.Find(nil, bson.M{"updated": bson.M{"$gt": now}}, mim.ops.FindOpt)
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		if cur.Err() != nil {
			mim.errorNum++
			return cur.Err()
		}
		r, err := mim.ops.IncParser.Parse(cur.Current, false)
		if err != nil {
			mim.errorNum++
			log.Println(err)
			continue
		}
		mim.result = append(mim.result, r)
	}
	for i := 0; i < len(mim.result); i++ {
		if mim.result[i].DataMod == 0 {
			_ = mim.innerIndex.Add(mim.result[i].Value)
		} else if mim.result[i].DataMod == 1 {
			mim.innerIndex.Del(mim.result[i].Value)
			_ = mim.innerIndex.Add(mim.result[i].Value)
		} else {
			mim.innerIndex.Del(mim.result[i].Value)
		}
	}
	return nil
}

func (mim *MongoIndexManager) Find() error {

	findOptions := options.Find()
	cur, err := mim.collection.Find(context.TODO(), bson.M{"status": 1}, findOptions)
	if err != nil {
		return helpers.CollectionNotFound
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		if cur.Err() != nil {
			return cur.Err()
		}
		r, err := mim.ops.BaseParser.Parse(cur.Current, true)
		if err != nil {
			mim.errorNum++
			log.Println(err)
			continue
		}
		mim.result = append(mim.result, r)
	}
	if err := cur.Err(); err != nil {
		return helpers.CursorError
	}
	return nil
}
