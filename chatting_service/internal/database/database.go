package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DbConfiguration struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
	SslMode  string
}

func NewDB() (*sql.DB, error) {
	file, err := os.Open("../internal/database/postgres_config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := DbConfiguration{}
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	var defaultConnString = fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s", conf.Host, conf.User, conf.Password, conf.Port, conf.SslMode)
	defaultDB, err := sql.Open("postgres", defaultConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to default database: %w", err)
	}
	defer defaultDB.Close()

	// Проверка и создание базы данных
	var dbExists bool
	err = defaultDB.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", conf.DbName).Scan(&dbExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	// Создание базы данных, если она не существует
	if !dbExists {
		_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", conf.DbName))
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	// Подключение к созданной базе данных
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", conf.Host, conf.User, conf.Password, conf.DbName, conf.Port, conf.SslMode)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	//Проверка подключения
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Создание кэш-таблицы Users, если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
    		id INTEGER PRIMARY KEY, --not serial because primary key must math main table
    		name TEXT NOT NULL,
			username TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache-table users: %w", err)
	}

	// Создание таблицы chats, если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chats (
    		id SERIAL PRIMARY KEY,
    		createdat DATETIME NOT NULL,
    		owner INT REFERENCES users(id) ON DELETE SET NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table chats: %w", err)
	}

	// Создание таблицы chat_to_participants, если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_to_participants (
    		id SERIAL PRIMARY KEY,
    		chat INT REFERENCES chats(id) ON DELETE SET NULL,
			participant INT REFERENCES users(id) ON DELETE SET NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table chat_to_participants: %w", err)
	}

	// Создание таблицы messages, если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
    		id SERIAL PRIMARY KEY,
    		sendtime DATETIME NOT NULL,
    		content TEXT NOT NULL,
    		author INT REFERENCES users(id) ON DELETE SET NULL,
			chat INT REFERENCES chats(id) ON DELETE SET NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table messages: %w", err)
	}

	return db, nil
}
