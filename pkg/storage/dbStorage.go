package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/crypto"
	"github.com/KaziPHone/go-musthave-diploma-tpl/pkg/user"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage interface {
	IsConnected() bool
	GetErr() error
	Insert(query string, args ...interface{}) error
	InsertWithReturning(query string, args ...interface{}) *sql.Row
	Select(query string, args ...interface{}) (sql.Result, error)
	CountRows(query string, args ...interface{}) (int, error)
	GetRows(r context.Context, query string, args ...interface{}) (*sql.Rows, error)
	RegisterUserWithBalance(req user.UserRequest, isHashed bool) (int, error)
	UserIsRegistred(login string) (bool, error)
}

type DataBase struct {
	Err         error
	dataBaseDsn string
	isConnected bool
	db          *sql.DB
	migratePath string
}

func newDatabase(dataBaseDsn string) DBStorage {

	db := &DataBase{
		dataBaseDsn: dataBaseDsn,
	}

	db.initDataBase()

	return db
}

func (d *DataBase) initDataBase() {

	if d.dataBaseDsn == "" {
		return
	}

	var db *sql.DB
	var err error

	db, err = sql.Open("pgx", d.dataBaseDsn)

	if err != nil {
		d.Err = err
		log.Printf("Unable to establish a connection with the database after retries: %v\n", err)
		return
	}

	err = db.Ping()
	if err != nil {
		d.Err = err
		log.Print("Database not connected...")
		return
	}

	d.isConnected = true
	d.db = db
	log.Print("Database connected")

	d.migrateDB()

}

func (d *DataBase) IsConnected() bool {
	return d.isConnected
}

func (d *DataBase) GetErr() error {
	return d.Err
}

func (d *DataBase) execContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error

	result, err = d.db.ExecContext(ctx, query, args...)
	return result, err
}

func (d *DataBase) Insert(query string, args ...interface{}) error {
	_, err := d.execContext(context.Background(), query, args...)
	return err
}

func (d *DataBase) InsertWithReturning(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(context.Background(), query, args...)
}

func (d *DataBase) Select(query string, args ...interface{}) (sql.Result, error) {
	return d.execContext(context.Background(), query, args...)
}

func (d *DataBase) CountRows(query string, args ...interface{}) (int, error) {
	var count int
	err := d.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (d *DataBase) GetRows(r context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(r, query, args...)
}

func (d *DataBase) migrateDB() {

	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		log.Printf("error creating migration driver: %v\n", err)
		return
	}

	d.setMigratePath()
	migrator, err := migrate.NewWithDatabaseInstance(
		d.migratePath,
		"postgres",
		driver,
	)

	if err != nil {
		log.Printf("failed to create migrator instance: %v\n", err)
		return
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Printf("migration failed: %v\n", err)
		return
	}

	log.Print("migration success")
}

// setMigratePath установка пути миграции
func (d *DataBase) setMigratePath() {
	if d.migratePath == "" {
		d.migratePath = "file://./migrations"
	}
}

func (d *DataBase) RegisterUserWithBalance(req user.UserRequest, isHashed bool) (int, error) {

	transaction, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer transaction.Rollback()

	pass := req.Password
	if !isHashed {
		pass = crypto.HashString(req.Password) // хеширование пароля
	}

	var userID int
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`
	err = transaction.QueryRow(query, req.Login, pass).Scan(&userID)
	if err != nil {
		return 0, err
	}

	t := time.Now()
	_, err = transaction.Exec(`INSERT INTO user_balance (user_id, current, withdrawn, updated_at, uploaded_at) 
		VALUES ($1, $2, $3, $4, $5)`, userID, 0, 0, t, t)
	if err != nil {
		return 0, err
	}

	// Фиксируем транзакцию
	if err = transaction.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

func (d *DataBase) UserIsRegistred(login string) (bool, error) {
	query := `SELECT COUNT(id) FROM users WHERE login = $1`
	result, err := d.CountRows(query, login)
	return result > 0, err
}
