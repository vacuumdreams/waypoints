package db

import (
	"context"
	"fmt"
	api "github.com/vacuumdreams/waypoints/api"
	"strconv"
	"strings"
)

func (db Database) GetList(userId string) ([]api.SavedWaypoint, error) {
	results := make([]api.SavedWaypoint, 0, 10)
	query := `
	SELECT waypoints.id, name, longitude, latitude, order_id
	FROM waypoints
	LEFT JOIN waypoints_order
	ON waypoints.id = waypoints_order.waypoint_id
	WHERE waypoints.user_id = $1
	`

	rows, err := db.Conn.Query(context.Background(), query, userId)

	if err != nil {
		return results, err
	}
	for rows.Next() {
		var item api.SavedWaypoint
		err := rows.Scan(&item.Id, &item.Name, &item.Longitude, &item.Latitude, &item.Order)

		if err != nil {
			return results, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (db Database) AddItem(userId string, item api.Waypoint) (*api.SavedWaypoint, error) {
	result := &api.SavedWaypoint{}

	query := `
  WITH op1 AS (
    INSERT INTO waypoints (user_id, name, longitude, latitude, created_at)
    VALUES ($1, $2, $3, $4, NOW())
    RETURNING id, name, longitude, latitude
  ), op2 AS (
    SELECT COALESCE((SELECT MAX(order_id) FROM waypoints_order WHERE user_id = $1), 0) + 1 AS order
  ), op3 AS (
    INSERT INTO waypoints_order (waypoint_id, user_id, order_id)
    VALUES ((SELECT op1.id FROM op1), $1, (SELECT op2.order FROM op2))
    RETURNING order_id AS order
  )
  SELECT * FROM op1, op3;
  `

	err := db.Conn.QueryRow(context.Background(), query, userId, item.Name, item.Longitude, item.Latitude).Scan(&result.Id, &result.Name, &result.Longitude, &result.Latitude, &result.Order)

	return result, err
}

func (db Database) DeleteItem(userId string, itemId uint64) error {
	query := `
  DELETE FROM waypoints
  WHERE id IN (
    SELECT id FROM waypoints WHERE user_id = $1 AND id = $2
  )
  `
	_, err := db.Conn.Exec(context.Background(), query, userId, itemId)
	return err
}

func (db Database) UpdateOrders(userId string, orders []api.WaypointOrder) ([]api.WaypointOrder, error) {
	query := ""
	values := []interface{}{}
	results := make([]api.WaypointOrder, 0, 10)

	values = append(values, userId)

	// build query and values dynamically for batch upsert
	for i, row := range orders {
		j := (i+2)*2 - 2
		query += fmt.Sprintf("($%s, $1, $%s),", strconv.Itoa(j), strconv.Itoa(j+1))
		values = append(values, row.Id, row.Order)
	}

	query = strings.TrimSuffix(query, ",")

	query = fmt.Sprintf(`
  INSERT INTO waypoints_order (waypoint_id, user_id, order_id)
  VALUES %s
  ON CONFLICT (waypoint_id)
  DO UPDATE SET order_id = EXCLUDED.order_id
  WHERE waypoints_order.user_id = $1
	RETURNING waypoint_id, order_id
  `, query)

	rows, err := db.Conn.Query(context.Background(), query, values...)

	if err != nil {
		return results, err
	}

	for rows.Next() {
		var item api.WaypointOrder
		err := rows.Scan(&item.Id, &item.Order)
		if err != nil {
			return results, err
		}
		results = append(results, item)
	}

	return results, err
}
