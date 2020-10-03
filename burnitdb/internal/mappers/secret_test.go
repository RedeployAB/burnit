package mappers

import (
	"testing"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/db"
	"github.com/RedeployAB/burnit/burnitdb/internal/dto"
)

var id1 = "507f1f77bcf86cd799439011"

func TestSecretToPersistance(t *testing.T) {

	dto1 := &dto.Secret{
		Secret: "secret",
	}
	secret1 := Secret{}.ToPersistance(dto1)

	if secret1.Secret != dto1.Secret {
		t.Errorf("incorrect value, got: %s, want: %s", secret1.Secret, dto1.Secret)
	}

	dto2 := &dto.Secret{
		Secret:     "secret",
		Passphrase: "passphrase",
	}
	secret2 := Secret{}.ToPersistance(dto2)

	if secret2.Secret != dto2.Secret {
		t.Errorf("incorrect value, got: %s, want: %s", secret2.Secret, dto2.Secret)
	}
	if secret2.Passphrase != dto2.Passphrase {
		t.Errorf("incorrect value, got: %s, want: %s", secret2.Passphrase, dto2.Passphrase)
	}

	dto3 := &dto.Secret{
		Secret:     "secret",
		Passphrase: "passphrase",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now(),
	}
	secret3 := Secret{}.ToPersistance(dto3)

	if secret3.Secret != dto3.Secret {
		t.Errorf("incorrect value, got: %s, want: %s", secret3.Secret, dto3.Secret)
	}
	if secret3.Passphrase != dto3.Passphrase {
		t.Errorf("incorrect value, got: %s, want: %s", secret3.Passphrase, dto3.Passphrase)
	}
	if secret3.CreatedAt != dto3.CreatedAt {
		t.Errorf("incorrect value, got %s, want: %s", secret3.CreatedAt.String(), dto3.CreatedAt.String())
	}
	if secret3.ExpiresAt != dto3.ExpiresAt {
		t.Errorf("incorrect value, got %s, want: %s", secret3.ExpiresAt.String(), dto3.ExpiresAt.String())
	}
}

func TestSecretToDTO(t *testing.T) {
	secret := &db.Secret{
		ID:         id1,
		Secret:     "secret",
		Passphrase: "passphrase",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now(),
	}

	dto := Secret{}.ToDTO(secret)

	if dto.ID != secret.ID {
		t.Errorf("incoorect value, got: %s, want: %s", dto.Secret, secret.ID)
	}
	if dto.Secret != secret.Secret {
		t.Errorf("incorrect value, got: %s, want: %s", dto.Secret, secret.Secret)
	}
	if dto.Passphrase != secret.Passphrase {
		t.Errorf("incorrect value, got: %s, want: %s", dto.Passphrase, secret.Passphrase)
	}
	if dto.CreatedAt != secret.CreatedAt {
		t.Errorf("incorrect value, got: %v, want: %v", dto.CreatedAt, secret.CreatedAt)
	}
	if dto.ExpiresAt != secret.ExpiresAt {
		t.Errorf("incorrect value, got:%v, want: %v", dto.ExpiresAt, secret.ExpiresAt)
	}
}
