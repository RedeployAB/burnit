package db

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"time"

	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"gopkg.in/mgo.v2"
)

// Secrets represents the DB session to Secrets with methods etc.
var session *mgo.Session

// Connect is used to connect to database with options
// specified in the passed in ConnectionOptions argument.
func Connect(opts config.Database) {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{opts.Address},
		Timeout:  60 * time.Second,
		Database: opts.Database,
		Username: opts.Username,
		Password: opts.Password,
	}

	if opts.SSL == true {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		}
	}
	// Modify these later on.
	dialInfo.Direct = true
	dialInfo.FailFast = true

	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("connection to database failed, error: %v\n", err)
	}
	session.SetSafe(&mgo.Safe{})
}

// GetSession gets an active session if exists,
// creates a new if not.
func GetSession() (*mgo.Session, error) {
	// Make modifications in the futures for setup of connection info.
	// for now program will be terminated if no sessions are found.
	if session == nil {
		return nil, errors.New("db: no existing session")
	}
	return session.Clone(), nil
}
