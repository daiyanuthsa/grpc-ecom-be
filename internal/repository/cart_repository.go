package repository

import (
	"context"
	"database/sql"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/entity"
)

type ICartRepository interface {
	// FindByUserIDAndProductID retrieves a cart item by user ID and product ID.
	FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*entity.CartItem, error)
	// Update updates an existing cart item.
	Update(ctx context.Context, item *entity.CartItem) error
	// Insert inserts a new cart item.
	Insert(ctx context.Context, item *entity.CartItem) error
	// FindByUserID retrieves all cart items for a given user ID.
	FindByUserID(ctx context.Context, userID string) ([]*entity.CartItem, error)
	// FindByID retrieves a single cart item by its ID.
	FindByID(ctx context.Context, cartID string) (*entity.CartItem, error)
	// Delete deletes a cart item by its ID.
	Delete(ctx context.Context, cartID string) error
}


// CartRepository implements ICartRepository for SQL database operations.
type CartRepository struct {
	db *sql.DB
}

// NewCartRepository creates a new instance of CartRepository.
func NewCartRepository(db *sql.DB) ICartRepository {
	return &CartRepository{
		db: db,
	}
}

// FindByUserIDAndProductID retrieves a cart item by user ID and product ID.
func (r *CartRepository) FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*entity.CartItem, error) {
	var item entity.CartItem
	query := `SELECT id, user_id, product_id, quantity, created_at, created_by, updated_at, updated_by
			  FROM public.user_cart WHERE user_id = $1 AND product_id = $2`
	row := r.db.QueryRowContext(ctx, query, userID, productID)
	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
		&item.CreatedBy,
		&item.UpdatedAt,
		&item.UpdatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &item, nil
}

// Update updates an existing cart item in the database.
func (r *CartRepository) Update(ctx context.Context, item *entity.CartItem) error {
	query := `UPDATE public.user_cart
			  SET quantity = $1, updated_at = $2, updated_by = $3
			  WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, item.Quantity, item.UpdatedAt, item.UpdatedBy, item.ID)
	return err
}

// Insert inserts a new cart item into the database.
func (r *CartRepository) Insert(ctx context.Context, item *entity.CartItem) error {
	query := `INSERT INTO public.user_cart (id, user_id, product_id, quantity, created_at, created_by, updated_at, updated_by)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.UserID,
		item.ProductID,
		item.Quantity,
		item.CreatedAt,
		item.CreatedBy,
		item.UpdatedAt,
		item.UpdatedBy,
	)
	return err
}

// FindByUserID retrieves all cart items for a given user ID.
func (r *CartRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.CartItem, error) {
	var items []*entity.CartItem
	query := `SELECT id, user_id, product_id, quantity, created_at, created_by, updated_at, updated_by
			  FROM public.user_cart WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.CartItem
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.CreatedBy,
			&item.UpdatedAt,
			&item.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// FindByID retrieves a single cart item by its ID.
func (r *CartRepository) FindByID(ctx context.Context, cartID string) (*entity.CartItem, error) {
	var item entity.CartItem
	query := `SELECT id, user_id, product_id, quantity, created_at, created_by, updated_at, updated_by
			  FROM public.user_cart WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, cartID)
	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
		&item.CreatedBy,
		&item.UpdatedAt,
		&item.UpdatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &item, nil
}

// Delete deletes a cart item from the database by its ID.
func (r *CartRepository) Delete(ctx context.Context, cartID string) error {
	query := `DELETE FROM public.user_cart WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, cartID)
	return err
}
