package mappers

import (
	"time"

	"github.com/RedeployAB/burnit/burnitdb/internal/dto"
	"github.com/RedeployAB/burnit/burnitdb/internal/models"
)

// Mapper is an interface that covers the mapping
// between DTO and Model, and from Model to DTO.
type Mapper interface {
	ToPersistance(interface{}) interface{}
	ToDTO(interface{}) interface{}
}

// Secret implements Mapper interface
// and provides methods ToPersistance and ToDTO.
type Secret struct{}

// ToPersistance transforms a Secret (DTO) to Secret (Model).
// It is the responsibility of the objects implementing
// Repository to handle ID setting and creating.
func (m Secret) ToPersistance(s *dto.Secret) *models.Secret {
	// Fallback for setting CreatedAt and ExpiresAt
	// if those are not set in incoming DTO object.
	var createdAt, expiresAt time.Time
	if s.CreatedAt.IsZero() {
		createdAt = time.Now()
	} else {
		createdAt = s.CreatedAt
	}

	if s.ExpiresAt.IsZero() {
		expiresAt = time.Now()
	} else {
		expiresAt = s.ExpiresAt
	}

	secretModel := &models.Secret{
		Secret:    s.Secret,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	if len(s.Passphrase) > 0 {
		secretModel.Passphrase = s.Passphrase
	}

	return secretModel
}

// ToDTO transforms a Secret (Model) to Secret (DTO).
func (m Secret) ToDTO(s *models.Secret) *dto.Secret {
	return &dto.Secret{
		ID:         s.ID.Hex(),
		Secret:     s.Secret,
		Passphrase: s.Passphrase,
		CreatedAt:  s.CreatedAt,
		ExpiresAt:  s.ExpiresAt,
	}
}
