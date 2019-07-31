// Simple usage for mongo.
package database

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MDB struct {
	db *mgo.Database
}

// Obtain a session using the Dial function:
// This will establish one or more connections with the cluster of
// servers defined by the url parameter.
func NewMongoDB(url, dbname string) (db *MDB, err error) {
	if dbname == "" {
		return nil, fmt.Errorf("Null dbname")
	}
	// all session methods are concurrency-safe and may be called from multiple goroutines.
	session, e := mgo.Dial(url)
	if e != nil {
		err = e
	} else {
		session.SetMode(mgo.Monotonic, true)
		db = &MDB{
			db: session.DB(dbname),
		}
	}
	return
}

func (this *MDB) getCollection(colName string) *mgo.Collection {
	return this.db.C(colName)
}

func (this *MDB) convertM2Map(m bson.M) map[string]interface{} {
	ma := map[string]interface{}(m)
	for k, v := range ma {
		if m, ok := v.(bson.M); ok {
			ma[k] = this.convertM2Map(m)
		}
	}
	return ma
}

func (this *MDB) UpsertOne(colName string, query, data map[string]interface{}) (info interface{}, err error) {
	col := this.getCollection(colName)
	info, err = col.Upsert(bson.M(query), bson.M(data))
	return
}

func (this *MDB) FindOne(colName string, query map[string]interface{}) (res map[string]interface{}, err error) {
	var doc bson.M
	col := this.getCollection(colName)
	if err = col.Find(bson.M(query)).One(&doc); err == nil {
		res = this.convertM2Map(doc)
	}
	return
}

func (this *MDB) Close() {
	this.db.Session.Close()
}
