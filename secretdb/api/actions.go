package api

import (
	"errors"
	"os"
	"time"

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

	secret := internal.Decrypt([]byte(s.Secret), os.Getenv("PASSPHRASE"))
	s.Secret = string(secret)

	return s, nil
}

// Insert handles inserts into the database.
func Insert(sb SecretBody, collection *mgo.Collection) (models.Secret, error) {
	id := bson.NewObjectId()
	created := time.Now()
	expires := created.AddDate(0, 0, 7)

	secret := internal.Encrypt([]byte(sb.Secret), os.Getenv("PASSPHRASE"))
	s := &models.Secret{
		ID:        id,
		Secret:    string(secret),
		CreatedAt: created,
		ExpiresAt: expires,
	}

	if len(sb.Passphrase) > 0 {
		s.Passphrase = internal.Hash(sb.Passphrase)
	}

	err := collection.Insert(s)
	if err != nil {
		return models.Secret{}, err
	}

	return *s, nil
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
