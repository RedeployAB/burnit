package mappers

import (
	"testing"
	"time"

	"github.com/RedeployAB/burnit/secretdb/internal/dto"

	"github.com/RedeployAB/burnit/secretdb/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var id1 = "507f1f77bcf86cd799439011"
var oid1, _ = primitive.ObjectIDFromHex(id1)

func TestSecretToPersistance(t *testing.T) {

	dto1 := &dto.Secret{
		Secret: "secret",
	}
	model1 := Secret{}.ToPersistance(dto1)

	if model1.Secret != dto1.Secret {
		t.Errorf("incorrect value, got: %s, want: %s", model1.Secret, dto1.Secret)
	}

	dto2 := &dto.Secret{
		Secret:     "secret",
		Passphrase: "passphrase",
	}
	model2 := Secret{}.ToPersistance(dto2)

	if model2.Secret != dto2.Secret {
		t.Errorf("incorrect value, got: %s, want: %s", model2.Secret, dto2.Secret)
	}
	if model2.Passphrase != dto2.Passphrase {
		t.Errorf("incorrect value, got: %s, want: %s", model2.Passphrase, dto2.Passphrase)
	}

	dto3 := &dto.Secret{
		Secret:     "secret",
		Passphrase: "passphrase",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now(),
	}
	model3 := Secret{}.ToPersistance(dto3)

	if model3.Secret != dto3.Secret {
		t.Errorf("incorrect value, got: %s, want: %s", model3.Secret, dto3.Secret)
	}
	if model3.Passphrase != dto3.Passphrase {
		t.Errorf("incorrect value, got: %s, want: %s", model3.Passphrase, dto3.Passphrase)
	}
	if model3.CreatedAt != dto3.CreatedAt {
		t.Errorf("incorrect value, got %s, want: %s", model3.CreatedAt.String(), dto3.CreatedAt.String())
	}
	if model3.ExpiresAt != dto3.ExpiresAt {
		t.Errorf("incorrect value, got %s, want: %s", model3.ExpiresAt.String(), dto3.ExpiresAt.String())
	}
}

func TestSecretToDTO(t *testing.T) {
	model := &models.Secret{
		ID:         oid1,
		Secret:     "secret",
		Passphrase: "passphrase",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now(),
	}

	dto := Secret{}.ToDTO(model)

	if dto.ID != model.ID.Hex() {
		t.Errorf("incoorect value, got: %s, want: %s", dto.Secret, model.ID.Hex())
	}
	if dto.Secret != model.Secret {
		t.Errorf("incorrect value, got: %s, want: %s", dto.Secret, model.Secret)
	}
	if dto.Passphrase != model.Passphrase {
		t.Errorf("incorrect value, got: %s, want: %s", dto.Passphrase, model.Passphrase)
	}
	if dto.CreatedAt != model.CreatedAt {
		t.Errorf("incorrect value, got: %v, want: %v", dto.CreatedAt, model.CreatedAt)
	}
	if dto.ExpiresAt != model.ExpiresAt {
		t.Errorf("incorrect value, got:%v, want: %v", dto.ExpiresAt, model.ExpiresAt)
	}
}
