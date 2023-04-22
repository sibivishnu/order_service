package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sibivishnu/order_service/models"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/orders_db")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %s", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %s", err))
	}
}

// AddOrder adds a new order to the database.
func AddOrder(order *models.Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO orders (id, status, total, currency_unit) VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(query, order.ID, order.Status, order.Total, order.CurrencyUnit)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		query := `INSERT INTO items (id, order_id, description, price, quantity) VALUES (?, ?, ?, ?, ?)`
		_, err = tx.Exec(query, item.ID, order.ID, item.Description, item.Price, item.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// UpdateOrder updates an existing order in the database.
func UpdateOrder(orderID string, order *models.Order) error {
	query := `UPDATE orders SET status = ?, total = ?, currency_unit = ? WHERE id = ?`
	_, err := db.Exec(query, order.Status, order.Total, order.CurrencyUnit, orderID)
	if err != nil {
		return err
	}

	// Update the items for the order
	// This example assumes that the items have been replaced completely
	_, err = db.Exec("DELETE FROM items WHERE order_id = 		?", orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		query := `INSERT INTO items (id, order_id, description, price, quantity) VALUES (?, ?, ?, ?, ?)`
		_, err = db.Exec(query, item.ID, order.ID, item.Description, item.Price, item.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOrders retrieves orders from the database based on the provided filter.
func GetOrders(filter Filter) ([]models.Order, error) {
	var orders []models.Order

	query := `SELECT id, status, total, currency_unit FROM orders`
	if filter.Status != "" {
		query += " WHERE status = ?"
	}

	rows, err := db.Query(query, filter.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.Status, &order.Total, &order.CurrencyUnit)
		if err != nil {
			return nil, err
		}

		items, err := getItemsForOrder(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

// getItemsForOrder retrieves items for a specific order from the database.
func getItemsForOrder(orderID string) ([]models.Item, error) {
	var items []models.Item

	query := `SELECT id, description, price, quantity FROM items WHERE order_id = ?`
	rows, err := db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.ID, &item.Description, &item.Price, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Filter represents a filter for orders.
type Filter struct {
	Status string
}
