package builder

import (
	"context"
	"github.com/MintegralTech/juno/helpers"
	"github.com/MintegralTech/juno/index"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoIndexBuilder struct {
	ops        *MongoIndexManagerOps
	innerIndex index.Index
	totalNum   int64
	errorNum   int64
	client     *mongo.Client
	collection *mongo.Collection
	findOpt    *options.FindOptions
	start      int64
	end        int64
}

func NewMongoIndexBuilder(ops *MongoIndexManagerOps) (*MongoIndexBuilder, error) {
	if ops == nil {
		return nil, helpers.MongoCfgError
	}

	mib := &MongoIndexBuilder{
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
		if mib.ops.Logger != nil {
			mib.ops.Logger.Warnf("mongo connect failed: [%s]", err.Error())
		}
		return nil, helpers.ConnectError
	}

	mib.client = client
	mib.findOpt = options.MergeFindOptions(ops.FindOpt)
	d := time.Duration(ops.ReadTimeout) * time.Microsecond
	mib.findOpt.MaxTime = &d

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		if mib.ops.Logger != nil {
			mib.ops.Logger.Warnf("mongo ping failed: [%s]", err.Error())
		}
		return nil, helpers.PingError
	}
	mib.collection = client.Database(ops.DB).Collection(ops.Collection)
	if mib.collection == nil {
		if mib.ops.Logger != nil {
			mib.ops.Logger.Warnf("mongo database[%s] collection[%s] not found", ops.DB, ops.Collection)
		}
		return nil, helpers.CollectionNotFound
	}
	return mib, nil
}

func (mib *MongoIndexBuilder) GetIndex() index.Index {
	return mib.innerIndex
}

func (mib *MongoIndexBuilder) update(ctx context.Context, name string) error {
	mib.start = time.Now().UnixNano()
	err := mib.base(name)
	mib.end = time.Now().UnixNano()
	if err != nil {
		mib.WarnStatus("base load failed: "+err.Error(), mib.end-mib.start)
		return err
	} else {
		mib.InfoStatus("base load success", mib.end-mib.start)
	}
	go func() {
		var (
			base = time.After(time.Duration(mib.ops.BaseInterval) * time.Second)
			inc  = time.After(time.Duration(mib.ops.IncInterval) * time.Second)
		)
		for {
			select {
			case <-ctx.Done():
				mib.start = time.Now().UnixNano()
				mib.end = time.Now().UnixNano()
				mib.InfoStatus("finish: ", mib.end-mib.start)
				return
			case <-base:
				mib.start = time.Now().UnixNano()
				err := mib.base(name)
				mib.end = time.Now().UnixNano()
				base = time.After(time.Duration(mib.ops.BaseInterval) * time.Second)
				if err != nil {
					mib.WarnStatus("base load failed: "+err.Error(), mib.end-mib.start)
				} else {
					mib.InfoStatus("base load success", mib.end-mib.start)
				}
			case <-inc:
				mib.start = time.Now().UnixNano()
				err := mib.inc(ctx)
				mib.end = time.Now().UnixNano()
				inc = time.After(time.Duration(mib.ops.IncInterval)*time.Second + time.Nanosecond)
				if err != nil {
					mib.WarnStatus("inc failed: "+err.Error(), mib.end-mib.start)
				} else {
					mib.InfoStatus("inc success", mib.end-mib.start)
				}
			}
		}
	}()
	return nil
}

func (mib *MongoIndexBuilder) base(name string) (err error) {
	mib.totalNum, mib.errorNum = 0, 0
	if mib.ops.OnBeforeBase != nil {
		mib.ops.BaseQuery = mib.ops.OnBeforeBase(mib.ops.UserData)
	}
	cur, err := mib.collection.Find(nil, mib.ops.BaseQuery, mib.ops.FindOpt)
	if err != nil {
		return err
	}
	defer func() {
		_ = cur.Close(nil)
	}()

	var baseIndex = index.NewIndex(name)
	for cur.Next(nil) {
		if cur.Err() != nil {
			mib.errorNum++
			continue
		}
		r := mib.ops.BaseParser.Parse(cur.Current, mib.ops.UserData)
		if r == nil {
			mib.errorNum++
			continue
		}
		mib.totalNum++
		_ = baseIndex.Add(r.Value)
	}
	mib.innerIndex = baseIndex
	if mib.ops.OnFinishBase != nil {
		mib.ops.OnFinishBase(mib)
	}
	return err
}

func (mib *MongoIndexBuilder) inc(ctx context.Context) (err error) {

	if mib.ops.OnBeforeInc != nil {
		mib.ops.IncQuery = mib.ops.OnBeforeInc(mib.ops.UserData)
	}
	c, cancel := context.WithTimeout(ctx, time.Duration(mib.ops.ReadTimeout)*time.Microsecond)
	defer cancel()

	cur, err := mib.collection.Find(c, mib.ops.IncQuery, mib.ops.FindOpt)
	if err != nil {
		return err
	}
	defer func() {
		_ = cur.Close(c)
	}()

	for cur.Next(nil) {
		if cur.Err() != nil {
			mib.errorNum++
			continue
		}
		r := mib.ops.IncParser.Parse(cur.Current, mib.ops.UserData)
		if r == nil {
			mib.errorNum++
			continue
		}
		if r.DataMod == DataAddOrUpdate {
			mib.innerIndex.Del(r.Value)
			_ = mib.innerIndex.Add(r.Value)
			mib.totalNum++
		} else {
			mib.innerIndex.Del(r.Value)
			mib.totalNum--
		}
	}
	if mib.ops.OnFinishInc != nil {
		mib.ops.OnFinishInc(mib)
	}
	return err
}

func (mib *MongoIndexBuilder) Build(ctx context.Context, name string) error {
	return mib.update(ctx, name)
}

func (mib *MongoIndexBuilder) InfoStatus(s string, t int64) {
	if mib.ops.Logger != nil {
		mib.ops.Logger.Infof("mongo_[%s]:[%s], totalNum:[%d], errorNum:[%d], "+
			"invert index:[%d], storage index:[%d], load time:[%dms]", mib.innerIndex.GetName(), s, mib.totalNum,
			mib.errorNum, mib.innerIndex.GetInvertedIndex().Count(), mib.innerIndex.GetStorageIndex().Count(), t/1e6)
	}
}

func (mib *MongoIndexBuilder) WarnStatus(s string, t int64) {
	if mib.ops.Logger != nil {
		mib.ops.Logger.Warnf("reason: %s, time: %dms", s, t)
	}
}
