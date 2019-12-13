package builder

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/juno/index"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type MongoBuilder struct {
	ops        *MongoIndexManagerOps
	innerIndex *index.IndexImpl
	totalNum   int64
	errorNum   int64
	client     *mongo.Client
	collection *mongo.Collection
	cursor     *mongo.Cursor
	findOpt    *options.FindOptions
}

func NewMongoIndexBuilder(ops *MongoIndexManagerOps) *MongoBuilder {
	if ops == nil {
		return nil
	}

	mongoIndexManager := &MongoBuilder{
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

func (mim *MongoBuilder) GetIndex() *index.IndexImpl {
	return mim.innerIndex
}

func (mim *MongoBuilder) update(ctx context.Context) error {
	if err := mim.base(ctx); err != nil {
		return err
	}
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
				}
			case <-base:
				if err := mim.base(ctx); err != nil {
					log.Println(err)
				}
			}
		}
	}()
	return nil
}

func (mim *MongoBuilder) base(ctx context.Context) error {
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
	mim.innerIndex = baseIndex
	return err
}

func (mim *MongoBuilder) inc(ctx context.Context) error {
	if mim.ops.OnBeforeInc != nil {
		mim.ops.IncQuery = mim.ops.OnBeforeInc(mim.ops.UserData)
	}
	c, cancal := context.WithTimeout(ctx, time.Duration(mim.ops.ReadTimeout)*time.Microsecond)
	defer cancal()
	cur, err := mim.collection.Find(c, mim.ops.IncQuery, mim.ops.FindOpt)
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		if cur.Err() != nil {
			continue
		}
		r, err := mim.ops.IncParser.Parse(cur.Current, false)
		if err != nil {
			log.Println(err)
			continue
		}
		if r.DataMod == 0 {

			_ = mim.innerIndex.Add(r.Value)
		} else if r.DataMod == 1 {
			mim.innerIndex.Del(r.Value)
			_ = mim.innerIndex.Add(r.Value)
		} else {
			mim.innerIndex.Del(r.Value)
		}
	}
	fmt.Println(mim.innerIndex.GetBitMap().Count())
	return nil
}

func (mim *MongoBuilder) Build(ctx context.Context) error {
	return mim.update(ctx)
}
