package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"stepikdb.ru/internal/models"
	"stepikdb.ru/internal/config"
)

func New(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func InsertCourse(db *sql.DB, course models.Course) error {
	query := `INSERT INTO courses (course_id, title, description, url) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, course.CourseID, course.Title, course.Description, course.URL)
	return err
}
