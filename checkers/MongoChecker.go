package checkers

import (
	"gopkg.in/mgo.v2"
)

type MongoChecker struct {
	HealthChecker HealthChecker
	Session       *mgo.Session
}

func NewMongoChecker(name string, sess *mgo.Session) MongoChecker {
	return MongoChecker{
		HealthChecker: NewHealthChecker(name),
		Session:       sess,
	}
}

func (mc *MongoChecker) Check() (Health, error) {
	err := mc.Session.Ping()
	if err != nil {
		logger.Log("lvl", "ERROR", "msg", err.Error())
		mc.HealthChecker.DownError(err)
		return mc.HealthChecker.Health, err
	}
	mc.HealthChecker.Up()
	return mc.HealthChecker.Health, nil
}

func (mc *MongoChecker) Name() string {
	return mc.HealthChecker.Name
}
