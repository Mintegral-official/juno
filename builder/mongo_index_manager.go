package builder

import (
	"context"
	"fmt"
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
	go func() {
		for {
			inc := time.After(time.Duration(mim.ops.IncInterval) * time.Second)
			base := time.After(time.Duration(mim.ops.BaseInterval) * time.Second)
			select {
			case <-ctx.Done():
				return
			case <-inc:
				fmt.Println(111)
				if err := mim.inc(ctx); err != nil {
					log.Println(err)
					return
				}
			case <-base:
				if err := mim.base(ctx); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}()
	return nil
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
	defer cur.Close(context.TODO())
	var baseResult []*ParserResult
	for cur.Next(context.TODO()) {
		if cur.Err() != nil {
			mim.errorNum++
			continue
		}
		r, err := mim.ops.BaseParser.Parse(cur.Current, true)
		if err != nil {
			mim.errorNum++
			log.Println(err)
			continue
		}
		baseResult = append(baseResult, r)
	}

	var baseIndex = index.NewIndex("base")
	for i := 0; i < len(baseResult); i++ {
		mim.totalNum++
		_ = baseIndex.Add(baseResult[i].Value)
	}
	mim.result = baseResult
	mim.innerIndex = baseIndex
	return err
}

func (mim *MongoIndexManager) inc(ctx context.Context) error {
	if mim.ops.OnBeforeInc != nil {
		mim.ops.IncQuery = mim.ops.OnBeforeInc(mim.ops.UserData)
	}
	context.WithTimeout(ctx, time.Duration(mim.ops.ReadTimeout)*time.Microsecond)
	cur, err := mim.collection.Find(nil, mim.ops.IncQuery, mim.ops.FindOpt)
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())
	var tmpResults []*ParserResult

	for cur.Next(context.TODO()) {
		if cur.Err() != nil {
			continue
		}
		r, err := mim.ops.IncParser.Parse(cur.Current, false)
		if err != nil {
			log.Println(err)
			continue
		}
		tmpResults = append(tmpResults, r)
	}
	for i := 0; i < len(tmpResults); i++ {
		if mim.result[i].DataMod == 0 {
			_ = mim.innerIndex.Add(tmpResults[i].Value)
		} else if mim.result[i].DataMod == 1 {
			mim.innerIndex.Del(tmpResults[i].Value)
			_ = mim.innerIndex.Add(tmpResults[i].Value)
		} else {
			mim.innerIndex.Del(tmpResults[i].Value)
		}
	}
	fmt.Println(mim.innerIndex.GetBitMap().Count())
	return nil
}

func (mim *MongoIndexManager) find(m bson.M) error {

	findOptions := options.Find()
	cur, err := mim.collection.Find(context.TODO(), m, findOptions)
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

func (mim *MongoIndexManager) Build() *index.IndexImpl {
	_ = mim.find(mim.ops.BaseQuery.(bson.M))
	if mim == nil || mim.result == nil || len(mim.result) == 0 {
		return index.NewIndex("empty")
	}
	mim.innerIndex = index.NewIndex("index")
	c := mim.result
	for i := 0; i < len(c); i++ {
		_ = mim.innerIndex.Add(c[i].Value)
	}
	return mim.innerIndex
}
