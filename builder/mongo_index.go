package builder

import (
	"context"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type MongoIndexBuilder struct {
	ops        *MongoIndexManagerOps
	innerIndex *index.IndexImpl
	totalNum   int64
	errorNum   int64
	client     *mongo.Client
	collection *mongo.Collection
	findOpt    *options.FindOptions
}

func NewMongoIndexBuilder(ops *MongoIndexManagerOps) (*MongoIndexBuilder, error) {
	if ops == nil {
		return nil, helpers.MongoCfgError
	}

	mongoIndexBuilder := &MongoIndexBuilder{
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
		return nil, helpers.ConnectError
	}

	mongoIndexBuilder.client = client
	mongoIndexBuilder.findOpt = options.MergeFindOptions(ops.FindOpt)
	d := time.Duration(ops.ReadTimeout) * time.Microsecond
	mongoIndexBuilder.findOpt.MaxTime = &d

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, helpers.PingError
	}
	mongoIndexBuilder.collection = client.Database(ops.DB).Collection(ops.Collection)
	if mongoIndexBuilder.collection == nil {
		return nil, helpers.CollectionNotFound
	}
	return mongoIndexBuilder, nil
}

func (mib *MongoIndexBuilder) GetIndex() *index.IndexImpl {
	return mib.innerIndex
}

func (mib *MongoIndexBuilder) update(ctx context.Context) error {
	if err := mib.base(ctx); err != nil {
		return err
	}
	go func() {
		for {
			inc := time.After(time.Duration(mib.ops.IncInterval) * time.Second)
			base := time.After(time.Duration(mib.ops.BaseInterval) * time.Second)
			select {
			case <-ctx.Done():
				return
			case <-inc:
				if err := mib.inc(ctx); err != nil {
					//log.Println(err) TODO log
				}
			case <-base:
				if err := mib.base(ctx); err != nil {
					// log.Println(err) TODO log
				}
			}
		}
	}()
	return nil
}

func (mib *MongoIndexBuilder) base(ctx context.Context) error {
	mib.totalNum = 0
	mib.errorNum = 0
	if mib.ops.OnBeforeBase != nil {
		mib.ops.BaseQuery = mib.ops.OnBeforeBase(mib.ops.UserData)
	}
	c, cancel := context.WithTimeout(ctx, time.Duration(mib.ops.ReadTimeout)*time.Microsecond)
	defer cancel()
	cur, err := mib.collection.Find(nil, mib.ops.BaseQuery, mib.ops.FindOpt)
	if err != nil {
		return err
	}
	defer cur.Close(c)
	var baseIndex = index.NewIndex("base")
	for cur.Next(c) {
		if cur.Err() != nil {
			mib.errorNum++
			continue
		}
		r, err := mib.ops.BaseParser.Parse(cur.Current)
		if err != nil {
			mib.errorNum++
			//log.Println(err) TODO add log
			continue
		}
		mib.totalNum++
		_ = baseIndex.Add(r.Value)
	}
	mib.innerIndex = baseIndex
	return err
}

func (mib *MongoIndexBuilder) inc(ctx context.Context) error {
	if mib.ops.OnBeforeInc != nil {
		mib.ops.IncQuery = mib.ops.OnBeforeInc(mib.ops.UserData)
	}
	c, cancel := context.WithTimeout(ctx, time.Duration(mib.ops.ReadTimeout)*time.Microsecond)
	defer cancel()
	cur, err := mib.collection.Find(c, mib.ops.IncQuery, mib.ops.FindOpt)

	if err != nil {
		return err
	}
	defer cur.Close(c)

	for cur.Next(c) {
		if cur.Err() != nil {
			continue
		}
		r, err := mib.ops.IncParser.Parse(cur.Current)
		if err != nil {
			log.Println(err)
			continue
		}
		if r.DataMod == 1 {
			mib.innerIndex.Del(r.Value)
			_ = mib.innerIndex.Add(r.Value)
		} else {
			mib.innerIndex.Del(r.Value)
		}
	}
	return nil
}

func (mib *MongoIndexBuilder) Build(ctx context.Context) error {
	return mib.update(ctx)
}
