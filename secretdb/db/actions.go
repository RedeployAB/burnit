package db

import (
	"errors"
	"time"

	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/RedeployAB/redeploy-secrets/secretdb/internal"
	"github.com/RedeployAB/redeploy-secrets/secretdb/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Find queries the collection for an entry by ID.
func Find(id string, collection *mgo.Collection) (models.Secret, error) {
	if !bson.IsObjectIdHex(id) {
		return models.Secret{}, errors.New("not valid ObjectId")
	}

	s := models.Secret{}
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&s)
	if err != nil {
		return models.Secret{}, err
	}

	secret := internal.Decrypt([]byte(s.Secret), config.Config.Passphrase)
	s.Secret = string(secret)

	return s, nil
}

// Insert handles inserts into the database.
func Insert(s Secret, collection *mgo.Collection) (models.Secret, error) {
	id := bson.NewObjectId()
	created := time.Now()
	expires := created.AddDate(0, 0, 7)

	secret := internal.Encrypt([]byte(s.Secret), config.Config.Passphrase)
	sm := &models.Secret{
		ID:        id,
		Secret:    string(secret),
		CreatedAt: created,
		ExpiresAt: expires,
	}

	if len(s.Passphrase) > 0 {
		sm.Passphrase = internal.Hash(s.Passphrase)
	}

	err := collection.Insert(sm)
	if err != nil {
		return models.Secret{}, err
	}

	return *sm, nil
}

// Delete removes an entry from the collection by ID.
func Delete(id string, collection *mgo.Collection) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not valid ObjectId")
	}

	err := collection.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		return err
	}
	return nil
}

// Secret represents a secret to be inserted into the
// database collection.
type Secret struct {
	Secret     string
	Passphrase string
}
