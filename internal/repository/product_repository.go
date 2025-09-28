package repository

import (
	"context"
	"database/sql"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/entity"
)

type IProductRepository interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProductById(ctx context.Context, id string) (*entity.Product, error)
}



type productRepository struct {
	db *sql.DB
}

func (r *productRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO \"product\" (id, name, description, price, image_file_name, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		product.Id, product.Name, product.Description, product.Price, product.ImageFileName, product.CreatedAt, product.CreatedBy, product.UpdatedAt, product.UpdatedBy, product.DeletedAt, product.DeletedBy, product.IsDeleted)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) GetProductById(ctx context.Context, id string) (*entity.Product, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, description, price, image_file_name FROM \"product\" WHERE id = $1 AND is_deleted = FALSE", id)
	var product entity.Product
	if err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.ImageFileName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil

}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &productRepository{db: db}
}
