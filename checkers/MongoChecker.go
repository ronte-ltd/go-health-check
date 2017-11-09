package checkers

import (
	"gopkg.in/mgo.v2"
)

// MongoChecker struct provide health checker and Session for mongodb
type MongoChecker struct {
	HealthChecker HealthChecker
	Session       *mgo.Session
}

//NewMongoChecker return new instance MongoChecker
func NewMongoChecker(name string, sess *mgo.Session) *MongoChecker {
	return &MongoChecker{
		HealthChecker: NewHealthChecker(name),
		Session:       sess,
	}
}

//Check connect to mongodb and sent ping, after return health or error
func (mc *MongoChecker) Check() (Health, error) {
	err := mc.Session.Ping()
	if err != nil {
		//logger.Log("lvl", "ERROR", "msg", err.Error())
		mc.HealthChecker.DownError(err)
		return mc.HealthChecker.Health, err
	}
	mc.HealthChecker.Up()
	return mc.HealthChecker.Health, nil
}

//Name return name current health check
func (mc *MongoChecker) Name() string {
	return mc.HealthChecker.Health.Name
}
