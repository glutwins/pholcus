package mgo

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var ErrIdType = errors.New("invalid id type")

type MgoStorage struct {
	sess   *mgo.Session
	dbname string
}

func NewMgoStorage(dsn string) (*MgoStorage, error) {
	store := &MgoStorage{}
	sess, err := mgo.Dial(dsn)
	if err != nil {
		return nil, err
	}

	if err = sess.Ping(); err != nil {
		return nil, err
	}

	store.sess = sess
	return store, nil
}

func (s *MgoStorage) Close() {
	if s.sess != nil {
		s.sess.Clone()
		s.sess = nil
	}
}

func (s *MgoStorage) Update(tblname string, old, data map[string]interface{}) error {
	c := s.sess.DB(s.dbname).C(tblname)

	if id, ok := old["_id"]; ok {
		old["_id"] = bson.ObjectIdHex(id.(string))
	}

	return c.Update(old, data)
}

func (s *MgoStorage) Remove(tblname string, data map[string]interface{}) error {
	c := s.sess.DB(s.dbname).C(tblname)

	if id, ok := data["_id"]; ok {
		if idStr, ok2 := id.(string); !ok2 {
			return ErrIdType
		} else {
			data["_id"] = bson.ObjectIdHex(idStr)
		}
	}

	return c.Remove(data)
}
