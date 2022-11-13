package mysql

import (
	"strings"

	"github.com/gba-3/tweat/domain/entity"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type TweatMysql struct {
	db *sqlx.Tx
}

type TweatGormMysql struct {
	db *gorm.DB
}
type ExecQueryer interface {
	sqlx.Queryer
	sqlx.Execer
}

func NewMysql(db *sqlx.Tx) *TweatMysql {
	return &TweatMysql{db}
}

func NewTweatGormMysql(db *gorm.DB) *TweatGormMysql {
	return &TweatGormMysql{db}
}

func (m *TweatMysql) GetAll(userID string) (entity.TweatLikesList, error) {
	query := strings.Join([]string{
		"SELECT tweats.id, tweats.text, tweats.user_id, count(likes.id) likes_count, users.name user_name",
		"FROM tweats",
		"LEFT JOIN likes ON `tweats`.id=`likes`.tweat_id",
		"LEFT JOIN users ON `tweats`.user_id=`users`.id",
		"WHERE tweats.user_id IN",
		"(SELECT follow_user_id FROM follows WHERE follows.user_id=?)",
		"GROUP BY tweats.id",
	}, " ")
	likesRows, err := m.db.Queryx(query, userID)
	if err != nil {
		return nil, err
	}

	tweats := entity.TweatLikesList{}
	for likesRows.Next() {
		var tweat entity.TweatLikes
		if err := likesRows.StructScan(&tweat); err != nil {
			return nil, err
		}
		tweats = append(tweats, tweat)
	}
	return tweats, nil
}

func (m *TweatGormMysql) GetAll(userID string) (entity.Tweats, error) {
	var tweats []entity.Tweat
	// m.db.Model(&tweats).
	// 	Select("tweats.id, tweats.text, tweats.user_id, count(likes.id) likes_count, users.name user_name").
	// 	Joins("LEFT JOIN likes ON `tweats`.id=`likes`.tweat_id").
	// 	Joins("LEFT JOIN users ON `tweats`.user_id=`users`.id").
	// 	Where("tweats.user_id")
	err := m.db.Preload("Likes").Where("user_id = ?", userID).Find(&tweats).Error
	return tweats, err
}
