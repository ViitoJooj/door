package initializer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const masterKeyVarName = "MASTER_KEY"

func normalizeMasterKey(key string) string {
	return strings.TrimSpace(key)
}

func validMasterKey(key string) bool {
	return len(key) == 32
}

func IsMasterKeyVar(name string) bool {
	return strings.EqualFold(strings.TrimSpace(name), masterKeyVarName)
}

func loadMasterKeyFromDB(db *sql.DB) (string, error) {
	var key string
	err := db.QueryRow(`SELECT value FROM env WHERE name = ?`, masterKeyVarName).Scan(&key)
	if err != nil {
		return "", err
	}

	key = normalizeMasterKey(key)
	if !validMasterKey(key) {
		return "", fmt.Errorf("invalid MASTER_KEY found in database")
	}
	return key, nil
}

func SaveMasterKey(db *sql.DB, key string) error {
	key = normalizeMasterKey(key)
	if !validMasterKey(key) {
		return fmt.Errorf("MASTER_KEY must have exactly 32 characters")
	}

	_, err := db.Exec(`
		INSERT INTO env (name, value)
		VALUES (?, ?)
		ON CONFLICT(name) DO UPDATE SET value = excluded.value;
	`, masterKeyVarName, key)
	if err != nil {
		return err
	}

	os.Setenv(masterKeyVarName, key)
	return nil
}

func EnsureMasterKey(db *sql.DB) string {
	if key := normalizeMasterKey(os.Getenv(masterKeyVarName)); validMasterKey(key) {
		if err := SaveMasterKey(db, key); err != nil {
			panic(err)
		}
		return key
	}

	if key, err := loadMasterKeyFromDB(db); err == nil {
		os.Setenv(masterKeyVarName, key)
		return key
	}

	data, err := os.ReadFile(".masterkey")
	if err == nil {
		key := normalizeMasterKey(string(data))
		if validMasterKey(key) {
			if err := SaveMasterKey(db, key); err != nil {
				panic(err)
			}
			return key
		}
	}

	key := randomHex(16)
	if err := SaveMasterKey(db, key); err != nil {
		panic(err)
	}
	return key
}

func encryptWithKey(text string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func decryptWithKey(enc string, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(data) < gcm.NonceSize() {
		return "", errors.New("dado inválido")
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}

func EncryptValue(text string) (string, error) {
	key := normalizeMasterKey(os.Getenv(masterKeyVarName))
	if !validMasterKey(key) {
		return "", fmt.Errorf("MASTER_KEY is not configured")
	}
	return encryptWithKey(text, key)
}

func DecryptValue(enc string) (string, error) {
	key := normalizeMasterKey(os.Getenv(masterKeyVarName))
	if !validMasterKey(key) {
		return "", fmt.Errorf("MASTER_KEY is not configured")
	}
	return decryptWithKey(enc, key)
}

func encrypt(text string) string {
	enc, err := EncryptValue(text)
	if err != nil {
		panic(err)
	}
	return enc
}

func decrypt(enc string) string {
	plain, err := DecryptValue(enc)
	if err != nil {
		panic(err)
	}
	return plain
}

func RotateMasterKey(db *sql.DB, newKey string) error {
	newKey = normalizeMasterKey(newKey)
	if !validMasterKey(newKey) {
		return fmt.Errorf("MASTER_KEY must have exactly 32 characters")
	}

	currentKey := EnsureMasterKey(db)
	if currentKey == newKey {
		return SaveMasterKey(db, newKey)
	}

	type envRecord struct {
		name  string
		value string
	}

	var records []envRecord

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	rows, err := tx.Query(`SELECT name, value FROM env WHERE name <> ?`, masterKeyVarName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, enc string
		if err = rows.Scan(&name, &enc); err != nil {
			return err
		}

		plain, decErr := decryptWithKey(enc, currentKey)
		if decErr != nil {
			return decErr
		}

		records = append(records, envRecord{name: name, value: plain})
	}
	if err = rows.Err(); err != nil {
		return err
	}

	if _, err = tx.Exec(`
		INSERT INTO env (name, value)
		VALUES (?, ?)
		ON CONFLICT(name) DO UPDATE SET value = excluded.value;
	`, masterKeyVarName, newKey); err != nil {
		return err
	}

	for _, record := range records {
		enc, encErr := encryptWithKey(record.value, newKey)
		if encErr != nil {
			return encErr
		}

		if _, err = tx.Exec(`UPDATE env SET value = ? WHERE name = ?`, enc, record.name); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	os.Setenv(masterKeyVarName, newKey)
	return nil
}

func SaveEnv(db *sql.DB, name string, value string) {
	if IsMasterKeyVar(name) {
		if err := SaveMasterKey(db, value); err != nil {
			panic(err)
		}
		return
	}

	enc := encrypt(value)

	_, err := db.Exec(`
		UPDATE env SET value = ? WHERE name = ?;
	`, enc, name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		INSERT INTO env (name, value)
		SELECT ?, ?
		WHERE NOT EXISTS (SELECT 1 FROM env WHERE name = ?);
	`, name, enc, name)

	if err != nil {
		panic(err)
	}
}

func LoadEnv(db *sql.DB, name string) string {
	if IsMasterKeyVar(name) {
		key, err := loadMasterKeyFromDB(db)
		if err != nil {
			return ""
		}
		return key
	}

	var enc string
	err := db.QueryRow(`SELECT value FROM env WHERE name = ?`, name).Scan(&enc)
	if err != nil {
		return ""
	}
	return decrypt(enc)
}

func InjectDefaultCors(db *sql.DB) {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO cors (origin) VALUES
		('http://localhost:4200'),
		('http://localhost:3000');
	`)
	if err != nil {
		panic(err)
	}
}

func LoadCors(db *sql.DB) []string {
	rows, err := db.Query(`SELECT origin FROM cors`)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	var origins []string

	for rows.Next() {
		var origin string
		if err := rows.Scan(&origin); err != nil {
			continue
		}

		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}
