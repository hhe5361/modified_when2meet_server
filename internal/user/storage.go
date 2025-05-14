package user

import "database/sql"

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// func scanRoom(rows *sql.Rows) (User, error) {
// 	var r Room
// 	err := rows.Scan(&r.ID, &r.Name, &r.URL, &r.StartTime, &r.EndTime, &r.TimeRegion, &r.IsOnline, &r.CreatedAt, &r.UpdatedAt)
// 	if err != nil {
// 		return Room{}, errors.New("server Error: room scan is failed")
// 	}
// 	return r, nil
// }
