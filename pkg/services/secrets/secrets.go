package secrets

import (
	"context"
	"fmt"
	"strings"

	"xorm.io/xorm"
)

// Service is an envelope encryption service in charge of encrypting/decrypting secrets.
// It is a replacement for encryption.Service
type Service interface {
	// Encrypt MUST NOT be used within database transactions, it may cause database locks.
	// For those specific use cases where the encryption operation cannot be moved outside
	// the database transaction, look at database-specific methods present at the specific
	// implementation present at manager.SecretsService.
	Encrypt(ctx context.Context, payload []byte, opt EncryptionOptions) ([]byte, error)
	Decrypt(ctx context.Context, payload []byte) ([]byte, error)

	// EncryptJsonData MUST NOT be used within database transactions.
	// Look at Encrypt method comment for further details.
	EncryptJsonData(ctx context.Context, kv map[string]string, opt EncryptionOptions) (map[string][]byte, error)
	DecryptJsonData(ctx context.Context, sjd map[string][]byte) (map[string]string, error)

	GetDecryptedValue(ctx context.Context, sjd map[string][]byte, key, fallback string) string
}

// Store defines methods to interact with secrets storage
type Store interface {
	GetDataKey(ctx context.Context, name string) (*DataKey, error)
	GetAllDataKeys(ctx context.Context) ([]*DataKey, error)
	CreateDataKey(ctx context.Context, dataKey DataKey) error
	CreateDataKeyWithDBSession(ctx context.Context, dataKey DataKey, sess *xorm.Session) error
	DeleteDataKey(ctx context.Context, name string) error
}

// Provider is a key encryption key provider for envelope encryption
type Provider interface {
	Encrypt(ctx context.Context, blob []byte) ([]byte, error)
	Decrypt(ctx context.Context, blob []byte) ([]byte, error)
}

type ProviderID string

func (id ProviderID) Kind() (string, error) {
	idStr := string(id)

	parts := strings.SplitN(idStr, ".", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("malformatted provider identifier %s: expected format <provider>.<keyName>", idStr)
	}

	return parts[0], nil
}

// BackgroundProvider should be implemented for a provider that has a task that needs to be run in the background.
type BackgroundProvider interface {
	Run(ctx context.Context) error
}
