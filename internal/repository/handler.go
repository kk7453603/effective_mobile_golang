package repository

import (
	"context"
	"log"
	"os"

	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type SqlHandler struct {
	DB  *pgxpool.Pool
	dsn string
}

func New() *SqlHandler {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")
	log.Println(dsn)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	return &SqlHandler{DB: pool, dsn: dsn}
}

func (h *SqlHandler) Migrate() error {
	m, err := migrate.New(os.Getenv("DB_MIGRATIONS_PATH"), h.dsn+"?sslmode=disable")
	if err != nil {
		log.Fatal(err) // добавить логгер
		return err
	}
	if err := m.Up(); err != nil {
		log.Fatal(err) // добавить логгер
		return err
	}
	return nil
}

func (h *SqlHandler) AddUser(user_name string, user_pass string) error {
	_, err := h.DB.Exec(context.Background(), "INSERT INTO users(name,password) VALUES ($1,$2)", user_name, user_pass)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (h *SqlHandler) GetUsers() ([]model.User, error) {
	var users []model.User
	rows, err := h.DB.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
