package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oktaviandwip/musalabel-backend/config"
	models "github.com/oktaviandwip/musalabel-backend/internal/models"
)

type RepoOrdersIF interface {
	FetchOrders(id string) (*config.Result, error)
	CreateOrder(data *models.Order) (*config.Result, error)
	UpdateOrder(data *models.Order) (*config.Result, error)
	RemoveOrder(data *models.Order) (*config.Result, error)
	CreatePurchase(data *models.Purchase, id string) (*config.Result, error)
	FetchPurchases(email, status string) (*config.Result, error)
	FetchPurchasesCount(id, status string) (int, error)
	UpdatePurchaseStatus(status, puchase_id string) (*config.Result, error)
	FetchIncome(interval string) (*config.Result, error)
	FetchQuantity(interval string) (*config.Result, error)
}

type RepoOrders struct {
	*sqlx.DB
}

func NewOrder(db *sqlx.DB) *RepoOrders {
	return &RepoOrders{db}
}

// Get All Cart Orders
func (r *RepoOrders) FetchOrders(id string) (*config.Result, error) {
	var result []models.Order

	q := `SELECT product_id, user_id, quantity, o.size, status, name, slug, description, price, image, stock 
				FROM orders o
				JOIN products p ON o.product_id = p.id 
				WHERE o.user_id = CAST($1 AS UUID) AND status = 'cart'
				ORDER BY o.created_at DESC
			`

	if err := r.Select(&result, r.Rebind(q), id); err != nil {
		return nil, err
	}

	for i, value := range result {
		result[i].Image = strings.Split(value.Image, ",")[0]
	}

	return &config.Result{Data: result}, nil
}

// Create Order
func (r *RepoOrders) CreateOrder(data *models.Order) (*config.Result, error) {
	q := `INSERT INTO orders (
					product_id,
					user_id,
					quantity,
					size,
					status
				)
				VALUES(
					$1,
					$2,
					$3,
					$4,
					$5
				)
				RETURNING id
			`

	args := []interface{}{
		data.Product_id,
		data.User_id,
		data.Quantity,
		data.Size,
		data.Status,
	}

	type orderId struct {
		Id string
	}

	var Id orderId
	err := r.QueryRowx(q, args...).StructScan(&Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &config.Result{Data: Id}, nil
}

// Update Order
func (r *RepoOrders) UpdateOrder(data *models.Order) (*config.Result, error) {
	q := `UPDATE orders 
				SET
					quantity = COALESCE(CAST(NULLIF(:quantity, '') AS INT), quantity),
					updated_at = NOW()		
				WHERE
					product_id = CAST(:product_id AS UUID) AND user_id = CAST(:user_id AS UUID) AND size = :size
			`

	_, err := r.NamedExec(q, data)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "1 data order updated"}, nil
}

// Delete Order
func (r *RepoOrders) RemoveOrder(data *models.Order) (*config.Result, error) {
	q := `DELETE FROM orders 
				WHERE 
					product_id = CAST(:product_id AS UUID) AND user_id = CAST(:user_id AS UUID) AND size = :size AND status = 'cart'
			`

	_, err := r.NamedExec(q, data)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "1 data order deleted"}, nil
}

// Creates or Updates Purchase
func (r *RepoOrders) CreatePurchase(data *models.Purchase, id string) (*config.Result, error) {
	var err error
	var tx *sql.Tx

	// Start a transaction
	tx, err = r.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback() // Rollback in case of error
		} else {
			err = tx.Commit() // Commit if no errors
		}
	}()

	if id == "" {
		// Insert a new record
		q := `
			INSERT INTO orders (
				product_id,
				user_id,
				quantity,
				size,
				status,
				purchase_id,
				invoice_url,
				delivery_option,
				delivery_address,
				total_price
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			);
		`

		_, err = tx.Exec(q, data.Product_id, data.User_id, data.Quantity, data.Size, data.Status, data.Purchase_id, data.Invoice_url, data.Delivery_option, data.Delivery_Address, data.Total_price)
		if err != nil {
			return nil, err
		}

		updateStock := `
			UPDATE products
			SET stock = stock - $1
			WHERE id = $2
		`

		_, err = tx.Exec(updateStock, data.Quantity, data.Product_id)
		if err != nil {
			return nil, err
		}

		return &config.Result{Message: "1 data order created"}, nil
	}

	// Update an existing record
	type result struct {
		Quantity   int
		Product_id string
	}

	res := result{}

	d := `SELECT quantity, product_id FROM orders WHERE id = $1`
	getErr := r.Get(&res, d, id)
	if getErr != nil {
		return nil, getErr
	}

	q := `
		UPDATE orders 
		SET
			status = COALESCE(NULLIF($1, ''), status),
			purchase_id = COALESCE(NULLIF($2, ''), purchase_id),
			invoice_url = COALESCE(NULLIF($3, ''), invoice_url),
			delivery_option = COALESCE(NULLIF($4, ''), delivery_option),
			delivery_address = COALESCE(NULLIF($5, ''), delivery_address),
			total_price = COALESCE(CAST(NULLIF($6, '') AS INT), total_price),
			updated_at = NOW()		
		WHERE
			id = CAST($7 AS UUID)
	`

	_, err = tx.Exec(q, data.Status, data.Purchase_id, data.Invoice_url, data.Delivery_option, data.Delivery_Address, data.Total_price, id)
	if err != nil {
		return nil, err
	}

	updateStock := `
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2
	`

	_, err = tx.Exec(updateStock, res.Quantity, res.Product_id)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "1 data order updated"}, nil
}

// Get Purchases
func (r *RepoOrders) FetchPurchases(email, status string) (*config.Result, error) {
	var result []models.Order
	queryString := fmt.Sprintf(`WHERE status = '%s' AND email = '%s'`, status, email)
	if status == "Semua" {
		queryString = fmt.Sprintf(`WHERE email = '%s'`, email)
	}
	if email == "admin@gmail.com" {
		queryString = fmt.Sprintf(`WHERE status = '%s'`, status)
	}
	if status == "Semua" && email == "admin@gmail.com" {
		queryString = ``
	}
	if status == "Dibatalkan" {
		queryString = fmt.Sprintf(`WHERE email = '%s' AND (status = 'Dibatalkan' OR status = 'EXPIRED')`, email)
	}
	if status == "Dibatalkan" && email == "admin@gmail.com" {
		queryString = `WHERE status = 'Dibatalkan' OR status = 'EXPIRED'`
	}

	q := fmt.Sprintf(`SELECT product_id, user_id, quantity, o.size, status, purchase_id, invoice_url, delivery_option, delivery_address, total_price, name, slug, description, price, p.image, stock, email
				FROM orders o
				JOIN products p ON o.product_id = p.id
				JOIN users u ON o.user_id = u.id
				%s
				ORDER BY o.created_at DESC, o.updated_at DESC
			`, queryString)

	if err := r.Select(&result, r.Rebind(q)); err != nil {
		return nil, err
	}

	for i, value := range result {
		result[i].Image = strings.Split(value.Image, ",")[0]
	}

	return &config.Result{Data: result}, nil
}

// Get Purchases Count
func (r *RepoOrders) FetchPurchasesCount(email, status string) (int, error) {
	var result int
	queryString := fmt.Sprintf(`WHERE status = '%s' AND email = '%s'`, status, email)
	if email == "admin@gmail.com" {
		queryString = fmt.Sprintf(`WHERE status = '%s'`, status)
	}
	if status == "Dibatalkan" {
		queryString = fmt.Sprintf(`WHERE email = '%s' AND (status = 'Dibatalkan' OR status = 'EXPIRED')`, email)
	}
	if status == "Dibatalkan" && email == "admin@gmail.com" {
		queryString = `WHERE status = 'Dibatalkan' OR status = 'EXPIRED'`
	}

	q := fmt.Sprintf(`SELECT COUNT (DISTINCT purchase_id) 
				FROM orders o
				JOIN products p ON o.product_id = p.id
				JOIN users u ON o.user_id = u.id
				%s`, queryString)

	if err := r.Get(&result, q); err != nil {
		return 0, err
	}

	return result, nil
}

// Update Purchase Status
func (r *RepoOrders) UpdatePurchaseStatus(status, puchase_id string) (*config.Result, error) {
	qString := ""

	if status == "Sedang Dikemas" {
		qString = `, purchased_at = NOW()`
	}

	q := fmt.Sprintf(`UPDATE orders 
				SET
					status = COALESCE(NULLIF($1, ''), status),
					updated_at = NOW()
					%s
				WHERE
					purchase_id = $2
			`, qString)

	_, err := r.Exec(q, status, puchase_id)
	if err != nil {
		return nil, err
	}

	return &config.Result{Message: "1 data order updated"}, nil
}

// Fetch Income
func (r *RepoOrders) FetchIncome(interval string) (*config.Result, error) {
	var qString string

	switch interval {
	case "daily":
		qString = `
			WITH last_7_days AS (
			  SELECT
			    CURRENT_DATE - INTERVAL '1 day' * (n - 1) AS day
			  FROM generate_series(1, 7) AS n
			),
			daily_totals AS (
			  SELECT
			    DATE_TRUNC('day', purchased_at) AS day,
			    SUM(total_price) AS total_price
			  FROM (
			    SELECT DISTINCT ON (purchase_id)
			      purchase_id, total_price, purchased_at
			    FROM public.orders
			    WHERE status IN ('Dikirim', 'Sedang Dikemas', 'Selesai')
			  ) AS distinct_orders
			  GROUP BY DATE_TRUNC('day', purchased_at)
			)
			SELECT
			  LEFT(TO_CHAR(last_7_days.day, 'Day'), 3) AS day_name,
			  COALESCE(daily_totals.total_price, 0) AS total_price
			FROM last_7_days
			LEFT JOIN daily_totals
			ON last_7_days.day = daily_totals.day
			ORDER BY last_7_days.day ASC;
		`
	case "weekly":
		qString = `
			WITH last_12_weeks AS (
			  SELECT
			    DATE_TRUNC('week', CURRENT_DATE) - INTERVAL '1 week' * (n - 1) AS week
			  FROM generate_series(0, 11) AS n
			),
			weekly_totals AS (
			  SELECT
			    (DATE_TRUNC('week', purchased_at) + INTERVAL '1 week') AS week,
			    SUM(total_price) AS total_price
			  FROM (
			    SELECT DISTINCT ON (purchase_id)
			      purchase_id, total_price, purchased_at
			    FROM public.orders
			    WHERE status IN ('Dikirim', 'Sedang Dikemas', 'Selesai')
			  ) AS distinct_orders
			  GROUP BY DATE_TRUNC('week', purchased_at) + INTERVAL '1 week'
			)
			SELECT
			  TO_CHAR(last_12_weeks.week, 'FMDD Mon') AS week,
			  COALESCE(weekly_totals.total_price, 0) AS total_price
			FROM last_12_weeks
			LEFT JOIN weekly_totals
			ON last_12_weeks.week = weekly_totals.week
			ORDER BY last_12_weeks.week ASC;
		`
	case "monthly":
		qString = `
			WITH last_12_months AS (
			  SELECT 
			    DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '1 month' * (n - 1) AS month
			  FROM 
			    generate_series(1, 12) AS n
			),
			monthly_totals AS (
			  SELECT
			    DATE_TRUNC('month', purchased_at) AS month,
			    SUM(total_price) AS total_price
			  FROM (
			    SELECT DISTINCT ON (purchase_id)
			      purchase_id, total_price, purchased_at
			    FROM public.orders
			    WHERE status IN ('Dikirim', 'Sedang Dikemas', 'Selesai')
			  ) AS distinct_orders
			  GROUP BY DATE_TRUNC('month', purchased_at)
			)
			SELECT
			  TO_CHAR(last_12_months.month, 'Mon') AS month,
			  COALESCE(monthly_totals.total_price, 0) AS total_price
			FROM last_12_months
			LEFT JOIN monthly_totals
			ON last_12_months.month = monthly_totals.month
			ORDER BY last_12_months.month ASC;
		`
	default:
		return nil, fmt.Errorf("invalid interval: %s", interval)
	}

	rows, err := r.Query(qString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var period string
		var totalPrice float64
		result := make(map[string]interface{})

		if err := rows.Scan(&period, &totalPrice); err != nil {
			return nil, err
		}
		result["Periode"] = period
		result["Total Pendapatan"] = totalPrice

		results = append(results, result)
	}

	return &config.Result{Data: results}, nil
}

// Fetch Quantity
func (r *RepoOrders) FetchQuantity(interval string) (*config.Result, error) {
	qString := fmt.Sprintf(`
		SELECT name, SUM(quantity)
		FROM orders o 
		JOIN products p ON o.product_id = p.id
		WHERE purchased_at BETWEEN CURRENT_DATE - INTERVAL '1 %s' AND CURRENT_DATE
		GROUP BY name
		ORDER BY SUM(quantity)
	`, interval)

	rows, err := r.Query(qString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var product string
		var totalSales float64
		result := make(map[string]interface{})

		if err := rows.Scan(&product, &totalSales); err != nil {
			return nil, err
		}
		result["Produk"] = product
		result["Total Penjualan"] = totalSales

		results = append(results, result)
	}

	return &config.Result{Data: results}, nil
}
