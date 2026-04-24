package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/pkg/initializer"
)

func (r *SQLite) UpdateUser(user *domain.User) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET username = ?, email = ?, password = ?, role = ?, active = ?, updated_at = ?
		WHERE id = ?
	`,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.Active,
		user.Updated_at,
		user.ID,
	)

	return err
}

func (r *SQLite) UpdateApplication(application *domain.Application) error {
	_, err := r.db.Exec(`
		UPDATE applications
		SET url = ?, country = ?, created_by = ?, updated_at = ?
		WHERE id = ?
	`,
		application.Url,
		application.Country,
		application.Created_by,
		application.Updated_at,
		application.ID,
	)

	return err
}

func (r *SQLite) ChangeVar(env *domain.Env) error {
	name := strings.TrimSpace(env.Name)
	if name == "" {
		return errors.New("name cannot be empty")
	}

	var currentName string
	var currentValue string
	err := r.db.QueryRow(`SELECT name, value FROM env WHERE id = ?`, env.Id).Scan(&currentName, &currentValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("env variable not found")
		}
		return err
	}

	if initializer.IsMasterKeyVar(currentName) && !initializer.IsMasterKeyVar(name) {
		return errors.New("MASTER_KEY name cannot be changed")
	}

	if initializer.IsMasterKeyVar(name) {
		return initializer.RotateMasterKey(r.db, env.Value)
	}

	if initializer.IsAppPortVar(name) {
		newPort, err := initializer.ParseAppPort(env.Value)
		if err != nil {
			return err
		}

		currentPortValue := ""
		if initializer.IsAppPortVar(currentName) {
			currentPortValue, err = initializer.DecryptValue(currentValue)
			if err != nil {
				return err
			}
		}

		currentPort, err := initializer.ParseAppPort(currentPortValue)
		if err == nil && currentPort == newPort {
			encValue, encErr := initializer.EncryptValue(env.Value)
			if encErr != nil {
				return encErr
			}

			_, err = r.db.Exec(`UPDATE env SET name = ?, value = ? WHERE id = ?`,
				name,
				encValue,
				env.Id,
			)
			return err
		}

		if !initializer.IsPortAvailable(newPort) {
			return errors.New("APP_PORT is unavailable")
		}
	}

	encValue, err := initializer.EncryptValue(env.Value)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`UPDATE env SET name = ?, value = ? WHERE id = ?`,
		name,
		encValue,
		env.Id,
	)

	return err
}

func (r *SQLite) ChangeCors(cors *domain.Cors) error {
	_, err := r.db.Exec(`UPDATE cors SET origin = ? WHERE id = ?`,
		cors.Origin, cors.Id,
	)
	return err
}
