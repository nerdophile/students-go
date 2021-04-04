package models

import (
	"fmt"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func ConnectDatabase() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "restfulapi"
	cluster.ProtoVersion = 4
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("cassandra well initialized")
}
