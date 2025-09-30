package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/entity"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"
	"github.com/daiyanuthsa/grpc-ecom-be/pb/common"
)

type IProductRepository interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProductById(ctx context.Context, id string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, DeletedAt time.Time, DeletedBy string, productId string) error
	ListProducts(ctx context.Context, page int32, limit int32, sort []*common.PaginationSortRequest) ([]*entity.Product, int32, error)
	ListProductsAdmin(ctx context.Context, page int32, limit int32, sort []*common.PaginationSortRequest) ([]*entity.Product, int32, error)
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

func (r *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	_, err := r.db.ExecContext(ctx, "UPDATE \"product\" SET name = $1, description = $2, price = $3, image_file_name = $4, updated_at = $5, updated_by = $6 WHERE id = $7",
		product.Name, product.Description, product.Price, product.ImageFileName, product.UpdatedAt, product.UpdatedBy, product.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, DeletedAt time.Time, DeletedBy string, productId string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE \"product\" SET is_deleted = $1, deleted_at = $2, deleted_by = $3 WHERE id = $4",
		true, DeletedAt, DeletedBy, productId)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) ListProducts(ctx context.Context, page int32, limit int32, sort []*common.PaginationSortRequest) ([]*entity.Product, int32, error) {
	// 1. Hitung OFFSET
	offset := (page - 1) * limit

	// 2. Query untuk Menghitung Total Elemen
	var totalElements int32
	countQuery := `SELECT COUNT(id) FROM "product" WHERE is_deleted = FALSE`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalElements); err != nil {
		log.Printf("Error counting products: %v", err)
		return nil, 0, fmt.Errorf("failed to get total product count: %w", err)
	}

	if totalElements == 0 {
		return nil, 0, nil
	}

	// 3. Bangun klausa ORDER BY secara dinamis dan aman menggunakan utilitas
	allowedSortFields := map[string]bool{
		"id":         true,
		"name":       true,
		"price":      true,
	}
	orderByClause, err := utils.BuildOrderByClause(sort, allowedSortFields, "ORDER BY created_at DESC")
	if err != nil {
		// Jika ada field yang tidak valid, kembalikan error. 
		// Service layer bisa menangani ini sebagai Bad Request.
		return nil, 0, fmt.Errorf("invalid sort parameter: %w", err)
	}

	// 4. Query untuk Mengambil Data Produk dengan LIMIT, OFFSET, dan ORDER BY dinamis
	dataQuery := fmt.Sprintf(`
		SELECT id, name, description, price, image_file_name
		FROM "product" 
		WHERE is_deleted = FALSE
		%s 
		LIMIT $1 OFFSET $2
	`, orderByClause)

	rows, err := r.db.QueryContext(ctx, dataQuery, limit, offset)
	if err != nil {
		log.Printf("Error querying products with pagination: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	// 5. Scan Hasil Query
	var products []*entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.ImageFileName); err != nil {
			log.Printf("Error scanning product row: %v", err)
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during iteration: %w", err)
	}

	// 6. Kembalikan data dan total elemen
	return products, totalElements, nil
}

func (r *productRepository) ListProductsAdmin(ctx context.Context, page int32, limit int32, sort []*common.PaginationSortRequest) ([]*entity.Product, int32, error) {
	// 1. Hitung OFFSET
	offset := (page - 1) * limit

	// 2. Query untuk Menghitung Total Elemen
	var totalElements int32
	countQuery := `SELECT COUNT(id) FROM "product"`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalElements); err != nil {
		log.Printf("Error counting products: %v", err)
		return nil, 0, fmt.Errorf("failed to get total product count: %w", err)
	}

	if totalElements == 0 {
		return nil, 0, nil
	}

	// 3. Bangun klausa ORDER BY secara dinamis dan aman menggunakan utilitas
	allowedSortFields := map[string]bool{
		"id":         true,
		"name":       true,
		"price":      true,
		"created_at": true,
		"updated_at": true,
		"deleted_at": true,
		"created_by": true,
		"updated_by": true,
		"deleted_by": true,
	}
	orderByClause, err := utils.BuildOrderByClause(sort, allowedSortFields, "ORDER BY created_at DESC")
	if err != nil {
		// Jika ada field yang tidak valid, kembalikan error. 
		// Service layer bisa menangani ini sebagai Bad Request.
		return nil, 0, fmt.Errorf("invalid sort parameter: %w", err)
	}

	// 4. Query untuk Mengambil Data Produk dengan LIMIT, OFFSET, dan ORDER BY dinamis
	dataQuery := fmt.Sprintf(`
		SELECT id, name, description, price, image_file_name, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted
		FROM "product" 
		%s 
		LIMIT $1 OFFSET $2
	`, orderByClause)

	rows, err := r.db.QueryContext(ctx, dataQuery, limit, offset)
	if err != nil {
		log.Printf("Error querying products with pagination: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	// 5. Scan Hasil Query
	var products []*entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.ImageFileName, &p.CreatedAt, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy, &p.DeletedAt, &p.DeletedBy, &p.IsDeleted); err != nil {
			log.Printf("Error scanning product row: %v", err)
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during iteration: %w", err)
	}

	// 6. Kembalikan data dan total elemen
	return products, totalElements, nil
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &productRepository{db: db}
}
