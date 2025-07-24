package person

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	domain "github.com/EstebanGitPro/motogo-backend/internal/domain/person"
)

const (
	QueryByEmail = `SELECT id, identity_number, first_name, last_name, second_last_name,
                    email, phone_number, email_verified, phone_number_verified, password, role
                    FROM persons WHERE email = ?`

	queryGetByID = `SELECT id, identity_number, first_name, last_name, second_last_name,
                    email, phone_number, email_verified, phone_number_verified, role
                    FROM persons WHERE id = ?`

	queryGetByIdentityNumber = `SELECT id, identity_number, first_name, last_name, second_last_name,
                    email, phone_number, email_verified, phone_number_verified, role
                    FROM persons WHERE identity_number = ?`

	querySave = `INSERT INTO persons (id, identity_number, first_name, last_name, second_last_name,
                 email, phone_number, email_verified, phone_number_verified, password, role)
                 VALUES (?,?,?,?,?,?,?,?,?,?,?)`

	queryUpdateEmailVerified = `UPDATE persons SET email_verified = TRUE WHERE id = ?`

	queryEmailVerificationTokensExpiration = `DELETE FROM email_verification_tokens WHERE expires_at < NOW()`

	queryTokenVerificationByHash = `SELECT id, user_id, token, expires_at, used, created_at FROM email_verification_tokens WHERE token = ?`

	queryMarkTokenAsUsed = `UPDATE email_verification_tokens SET used = TRUE WHERE id = ?`

	querySaveVerificationToken = `INSERT INTO email_verification_tokens (id, user_id, token, expires_at) VALUES (?, ?, ?, ?)`
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) domain.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *repository) CleanupExpiredTokens() error {
	stmt, err := r.db.Prepare(queryEmailVerificationTokensExpiration)
	if err != nil {
		return domain.ErrCleanupExpiredTokens
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return domain.ErrCleanupExpiredTokens
	}
	return nil
}

func (r *repository) GetByID(id string) (*domain.Person, error) {
	log.Printf("GetByID received ID: '%s'", id)
	stmt, err := r.db.Prepare(queryGetByID)
	if err != nil {
		log.Printf("Error preparing queryGetByID: %v", err)
		return nil, domain.ErrGetUsers
	}
	defer stmt.Close()

	var person Person
	err = stmt.QueryRow(id).Scan(
		&person.ID,
		&person.IdentityNumber,
		&person.FirstName,
		&person.LastName,
		&person.SecondLastName,
		&person.Email,
		&person.PhoneNumber,
		&person.EmailVerified,
		&person.PhoneNumberVerified,
		&person.Role,
	)
	if err != nil {
		log.Printf("Error executing query or scanning row for ID '%s': %v", id, err)
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserCannotFound
		}
		return nil, domain.ErrUserCannotGet
	}

	personDomain := person.ToDomain()
	return &personDomain, nil
}

func (r *repository) GetByIdentityNumber(identityNumber string) (*domain.Person, error) {
	stmt, err := r.db.Prepare(queryGetByIdentityNumber)
	if err != nil {
		return nil, domain.ErrGetUsers
	}
	defer stmt.Close()

	var person Person
	err = stmt.QueryRow(identityNumber).Scan(
		&person.ID,
		&person.IdentityNumber,
		&person.FirstName,
		&person.LastName,
		&person.SecondLastName,
		&person.Email,
		&person.PhoneNumber,
		&person.EmailVerified,
		&person.PhoneNumberVerified,
		&person.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserCannotFound
		}
		return nil, domain.ErrUserCannotGet
	}

	personDomain := person.ToDomain()
	return &personDomain, nil
}

func (r *repository) GetPersonByEmail(email string) (*domain.Person, error) {
	stmt, err := r.db.Prepare(QueryByEmail)
	if err != nil {
		return nil, domain.ErrGettingUserByEmail
	}
	defer stmt.Close()

	var person Person
	err = stmt.QueryRow(email).Scan(
		&person.ID,
		&person.IdentityNumber,
		&person.FirstName,
		&person.LastName,
		&person.SecondLastName,
		&person.Email,
		&person.PhoneNumber,
		&person.EmailVerified,
		&person.PhoneNumberVerified,
		&person.Password,
		&person.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFoundUserByEmail
		}
		return nil, domain.ErrGettingUserByEmail
	}

	personDomain := person.ToDomain()
	return &personDomain, nil
}

func (r *repository) GetVerificationTokenByHash(hashedToken string) error {
	stmt, err := r.db.Prepare(queryTokenVerificationByHash)
	if err != nil {
		return domain.ErrGetVerificationToken
	}
	defer stmt.Close()

	var token domain.EmailVerificationToken

	err = stmt.QueryRow(hashedToken).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.Used,
		&token.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrVerificationTokenNotFound
		}
		return domain.ErrGetVerificationToken
	}

	if time.Now().After(token.ExpiresAt) {
		return domain.ErrTokenExpired
	}
	if token.Used {
		return domain.ErrTokenAlreadyUsed
	}

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	rollbackWithLog := func(err error) error {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Error durante rollback: %v (error original: %v)", rbErr, err)
		} else {
			log.Printf("Rollback exitoso después de error: %v", err)
		}
		return err
	}

	if err := r.MarkTokenAsUsedTx(tx, token.ID); err != nil {
		return rollbackWithLog(err)
	}

	if err := r.MarkEmailAsVerifiedTx(tx, token.UserID); err != nil {
		return rollbackWithLog(err)
	}

	if err := tx.Commit(); err != nil {
		return rollbackWithLog(err)
	}

	return nil
}

func (r *repository) MarkEmailAsVerifiedTx(tx *sql.Tx, userID string) error {
	stmt, err := tx.Prepare(queryUpdateEmailVerified)
	if err != nil {
		return domain.ErrUserCannotUpdateEmailVerified
	}
	defer stmt.Close()

	result, err := stmt.Exec(userID)
	if err != nil {
		return domain.ErrUserCannotUpdateEmailVerified
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrUserCannotFound
	}

	return nil
}

func (r *repository) MarkTokenAsUsedTx(tx *sql.Tx, tokenID string) error {
	stmt, err := tx.Prepare(queryMarkTokenAsUsed)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(tokenID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrTokenNotFound
	}

	return nil
}

func (r *repository) Save(person domain.Person) error {
	personToSave := Person{
		ID:                  person.ID,
		IdentityNumber:      person.IdentityNumber,
		FirstName:           person.FirstName,
		LastName:            person.LastName,
		SecondLastName:      person.SecondLastName,
		Email:               person.Email,
		PhoneNumber:         person.PhoneNumber,
		EmailVerified:       person.EmailVerified,
		PhoneNumberVerified: person.PhoneNumberVerified,
		Password:            person.Password,
		Role:                person.Role,
	}

	stmt, err := r.db.Prepare(querySave)
	if err != nil {
		return domain.ErrUserCannotSave
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		personToSave.ID,
		personToSave.IdentityNumber,
		personToSave.FirstName,
		personToSave.LastName,
		personToSave.SecondLastName,
		personToSave.Email,
		personToSave.PhoneNumber,
		personToSave.EmailVerified,
		personToSave.PhoneNumberVerified,
		personToSave.Password,
		personToSave.Role,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Duplicate"):
			return domain.ErrDuplicateUser
		default:
			return domain.ErrUserCannotSave
		}
	}

	return nil
}

func (r *repository) SaveVerificationToken(token *domain.EmailVerificationToken) error {
	stmt, err := r.db.Prepare(querySaveVerificationToken)
	if err != nil {
		return domain.ErrUserCannotSaveVerificationToken
	}
	defer stmt.Close()

	_, err = stmt.Exec(token.ID, token.UserID, token.Token, token.ExpiresAt)
	if err != nil {
		return domain.ErrUserCannotSaveVerificationToken
	}

	return nil
}
