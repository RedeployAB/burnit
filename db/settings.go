package db

// Settings defines the settings for the application.
type Settings struct {
	Security Security
}

// Security defines the security settings for the application.
type Security struct {
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	EncryptionKey string `bson:"encryptionKey,omitempty"`
}
