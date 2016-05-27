package show

import (
	"time"
)

type Show struct {
	ID        int        `jsonapi:"primary,shows" gorm:"primary_key"`
	CreatedAt time.Time  `jsonapi:"attr,created_at"`
	UpdatedAt time.Time  `jsonapi:"attr,updated_at"`
	DeletedAt *time.Time `jsonapi:"" sql:"index"`
	Title     string     `jsonapi:"attr,title" valid:"ascii,required"`
	Year      int64      `jsonapi:"attr,year" valid:"required"`
	// EpisodesCount uint             `json:"episodes-count"`
	// Episodes      episode.Episodes `json:"episodes"`
}

type Shows []*Show

type ShowResource struct {}

func (Show) TableName() string {
	return "shows";
}