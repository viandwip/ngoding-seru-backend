package repository

import (
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
)

type RepoProductsIF interface {
	CreateProduct(data *models.Product) (*config.Result, error)
	FetchProducts(page, limit int) (*config.Result, error)
	SearchProducts(search string, page, limit int) (*config.Result, error)
	FetchProduct(id, slug string) (*config.Result, error)
	UpdateProduct(data *models.Product) (*config.Result, error)
	RemoveProduct(id string) (*config.Result, error)
}

type RepoProducts struct {
	*sqlx.DB
}

func NewProduct(db *sqlx.DB) *RepoProducts {
	return &RepoProducts{db}
}

// Get All Products
func (r *RepoProducts) FetchProducts(page, limit int) (*config.Result, error) {
	offset := (page - 1) * limit
	var result []models.Product

	q := `SELECT * FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	if err := r.Select(&result, r.Rebind(q), limit, offset); err != nil {
		return nil, err
	}

	// Pagination
	var totalCount int

	totalCountQuery := `SELECT COUNT(*) FROM products`

	if err := r.Get(&totalCount, totalCountQuery); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	var next, prev interface{}
	if page < totalPages {
		next = page + 1
	}
	if page > 1 {
		prev = page - 1
	}

	meta := &config.Metas{
		Total: totalCount,
		Next:  next,
		Prev:  prev,
	}

	return &config.Result{Data: result, Meta: meta}, nil
}

// Search Products
func (r *RepoProducts) SearchProducts(search string, page, limit int) (*config.Result, error) {
	searchStr := fmt.Sprintf("%%%s%%", search)
	offset := (page - 1) * limit
	var result []models.Product

	q := `SELECT * FROM products WHERE name ILIKE $1 ORDER BY name DESC LIMIT $2 OFFSET $3`

	if err := r.Select(&result, r.Rebind(q), searchStr, limit, offset); err != nil {
		return nil, err
	}

	// Pagination
	var totalCount int

	totalCountQuery := `SELECT COUNT(*) FROM products WHERE name ILIKE $1`

	if err := r.Get(&totalCount, totalCountQuery, searchStr); err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	var next, prev interface{}
	if page < totalPages {
		next = page + 1
	}
	if page > 1 {
		prev = page - 1
	}

	meta := &config.Metas{
		Total: totalCount,
		Next:  next,
		Prev:  prev,
	}

	return &config.Result{Data: result, Meta: meta}, nil
}

// Get Product
func (r *RepoProducts) FetchProduct(id, slug string) (*config.Result, error) {
	var result models.Product
	data := id
	column := "id"
	if id == "" {
		data = slug
		column = "slug"
	}
	fmt.Println(data, column)

	q := fmt.Sprintf("SELECT * FROM products WHERE %s = $1", column)

	if err := r.Get(&result, r.Rebind(q), data); err != nil {
		return nil, err
	}

	return &config.Result{Data: result}, nil
}

// Create Product
func (r *RepoProducts) CreateProduct(data *models.Product) (*config.Result, error) {
	q := `INSERT INTO products (
				image, 
				name, 
				description,
				price,
				stock,
				size,
				slug
			)
			VALUES(
				:image,
				:name,
				:description,
				:price,
				:stock,
				:size,
				:slug
			)`

	_, err := r.NamedExec(q, data)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "1 data product created"}, nil
}

// Update Product
func (r *RepoProducts) UpdateProduct(data *models.Product) (*config.Result, error) {
	q := `UPDATE products 
				SET
					image = COALESCE(NULLIF(:image, ''), image), 
					name = COALESCE(NULLIF(:name, ''), name),
					description = COALESCE(NULLIF(:description, ''), description),
					stock = COALESCE(CAST(NULLIF(:stock, '') AS INT), stock),
					price = COALESCE(CAST(NULLIF(:price, '') AS INT), price),
					slug = COALESCE(NULLIF(:slug, ''), slug),
					size = COALESCE(NULLIF(:size, ''), size),
					updated_at = NOW()		
				WHERE
					id = CAST(:id AS UUID)
			`

	_, err := r.NamedExec(q, data)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "1 data product updated"}, nil
}

// Delete Product
func (r *RepoProducts) RemoveProduct(id string) (*config.Result, error) {
	q := `DELETE FROM products 
				WHERE 
					id = $1`

	_, err := r.Exec(q, id)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "1 data product deleted"}, nil
}
