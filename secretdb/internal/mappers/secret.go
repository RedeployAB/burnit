package mappers

import (
	"time"

	"github.com/RedeployAB/burnit/secretdb/internal/dto"
	"github.com/RedeployAB/burnit/secretdb/internal/models"
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
func (m Secret) ToPersistance(s *dto.Secret) *models.Secret {
	if s.TTL == 0 {
		s.TTL = 10080
	}

	secretModel := &models.Secret{
		Secret:    s.Secret,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.TTL)),
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
