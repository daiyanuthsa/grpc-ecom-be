package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/entity"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	InsertUser(ctx context.Context, user *entity.User) error
	UpdateUserPassword(ctx context.Context, userID string, newHashedPassword string, updateBy string) error
}

type authRepository struct {
	db *sql.DB
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	// Implement your logic to get user by email from the database
	row := r.db.QueryRowContext(ctx, "SELECT id, email, full_name, role_code, password, created_at FROM \"user\" WHERE email = $1 AND is_deleted = FALSE", email)
	var user entity.User
	if err := row.Scan(&user.Id, &user.Email, &user.FullName, &user.RoleCode, &user.Password, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) InsertUser(ctx context.Context, user *entity.User) error {
	// Implement your logic to save user to the database
	_, err := r.db.ExecContext(ctx, "INSERT INTO \"user\" (id, email, full_name, password, role_code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		user.Id, user.Email, user.FullName, user.Password, user.RoleCode, user.CreatedAt, user.CreatedBy, user.UpdatedAt, user.UpdatedBy, user.DeletedAt, user.DeletedBy, user.IsDeleted)
	if err != nil {
		return err
	}
	return nil
}
func (r *authRepository) UpdateUserPassword(ctx context.Context, userID string, newHashedPassword string, updateBy string) error {

	result, err := r.db.ExecContext(ctx,
		"UPDATE \"user\" SET password = $1, updated_at = $2 , updated_by=$3 WHERE id = $4", 
		newHashedPassword, 
		time.Now(), 
		updateBy,
		userID)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	// Jika tidak ada baris yang terpengaruh, berarti email tidak ditemukan
	if rowsAffected == 0 {
		log.Println(err)
		return err
	}

	return nil
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &authRepository{db: db}
}
