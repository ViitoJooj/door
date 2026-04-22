package initproject

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"
)

func masterKey() []byte {
	key := os.Getenv("MASTER_KEY")
	return []byte(key)
}

func EnsureMasterKey() string {
	if key := os.Getenv("MASTER_KEY"); len(key) == 32 {
		return key
	}

	data, err := os.ReadFile(".masterkey")
	if err == nil {
		key := strings.TrimSpace(string(data))
		if len(key) == 32 {
			os.Setenv("MASTER_KEY", key)
			return key
		}
	}

	key := randomHex(16)

	err = os.WriteFile(".masterkey", []byte(key), 0600)
	if err != nil {
		panic(err)
	}

	os.Setenv("MASTER_KEY", key)
	return key
}

func encrypt(text string) string {
	block, err := aes.NewCipher(masterKey())
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func decrypt(enc string) string {
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(masterKey())
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	if len(data) < gcm.NonceSize() {
		panic(errors.New("dado inválido"))
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return string(plain)
}

func SaveEnv(db *sql.DB, name string, value string) {
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
