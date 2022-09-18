package mysql

import (
	"strings"

	"github.com/gba-3/tweat/domain/entity"
	"github.com/jmoiron/sqlx"
)

type UserMysql struct {
	db *sqlx.DB
}

func NewUserMysql(db *sqlx.DB) *UserMysql {
	return &UserMysql{db}
}

func (um *UserMysql) FindByEmail(email string) (*entity.User, error) {
	query := "SELECT id, email, password FROM users WHERE email=?"

	user := entity.User{}
	if err := um.db.Get(&user, query, email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserMysql) CreateUser(name string, email string, password string) error {
	query := strings.Join([]string{
		"INSERT INTO `users` (`name`, `email`, `password`) VALUES",
		"(?, ?, ?)",
	}, "")

	tx, err := um.db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(query, name, email, password); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (um *UserMysql) GetFollowUsers(userID int) (entity.Follows, error) {
	query := strings.Join([]string{
		"SELECT * FROM follows",
		"WHERE user_id=?",
	}, " ")
	rows, err := um.db.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	follows := entity.Follows{}
	for rows.Next() {
		var follow entity.Follow
		if err := rows.StructScan(&follow); err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}
	return follows, nil
}
