package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/namviet/backend-core/internal/models"
)

type WarehouseRepository struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) FindAll(ctx context.Context) ([]models.Warehouse, error) {
	query := `SELECT id, company_id, key, code, name, type, address, manager, phone, latitude, longitude, status, outlet_type, created_at FROM warehouses ORDER BY id DESC`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var w models.Warehouse
		err := rows.Scan(
			&w.ID, &w.CompanyID, &w.Key, &w.Code, &w.Name, &w.Type,
			&w.Address, &w.Manager, &w.Phone, &w.Latitude, &w.Longitude,
			&w.Status, &w.OutletType, &w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, w)
	}

	return warehouses, nil
}

func (r *WarehouseRepository) FindByID(ctx context.Context, id int64) (*models.Warehouse, error) {
	query := `SELECT id, company_id, key, code, name, type, address, manager, phone, latitude, longitude, status, outlet_type, created_at FROM warehouses WHERE id = $1`
	
	var w models.Warehouse
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.CompanyID, &w.Key, &w.Code, &w.Name, &w.Type,
		&w.Address, &w.Manager, &w.Phone, &w.Latitude, &w.Longitude,
		&w.Status, &w.OutletType, &w.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("không tìm thấy chi nhánh")
		}
		return nil, err
	}
	return &w, nil
}

func (r *WarehouseRepository) Create(ctx context.Context, w *models.Warehouse) error {
	query := `
		INSERT INTO warehouses (company_id, key, code, name, type, address, manager, phone, latitude, longitude, status, outlet_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, 
		w.CompanyID, w.Key, w.Code, w.Name, w.Type, w.Address, w.Manager, w.Phone, w.Latitude, w.Longitude, w.Status, w.OutletType,
	).Scan(&w.ID, &w.CreatedAt)

	return err
}

func (r *WarehouseRepository) Update(ctx context.Context, w *models.Warehouse) error {
	query := `
		UPDATE warehouses
		SET company_id = $1, key = $2, code = $3, name = $4, type = $5, address = $6, manager = $7, phone = $8, latitude = $9, longitude = $10, status = $11, outlet_type = $12
		WHERE id = $13
	`
	res, err := r.db.ExecContext(ctx, query,
		w.CompanyID, w.Key, w.Code, w.Name, w.Type, w.Address, w.Manager, w.Phone, w.Latitude, w.Longitude, w.Status, w.OutletType,
		w.ID,
	)
	if err != nil {
		return err
	}
	
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("không tìm thấy chi nhánh")
	}

	return nil
}

func (r *WarehouseRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM warehouses WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("không tìm thấy chi nhánh")
	}
	return nil
}
