package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

//Clan identifies uniquely one clan in a given game
type Clan struct {
	ID        int    `db:"id"`
	GameID    string `db:"game_id"`
	ClanID    string `db:"clan_id"`
	Name      string `db:"name"`
	OwnerID   int    `db:"owner_id"`
	Metadata  string `db:"metadata"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}

//PreInsert populates fields before inserting a new clan
func (c *Clan) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now().UnixNano()
	c.UpdatedAt = c.CreatedAt
	return nil
}

//PreUpdate populates fields before updating a clan
func (c *Clan) PreUpdate(s gorp.SqlExecutor) error {
	c.UpdatedAt = time.Now().UnixNano()
	return nil
}

//GetClanByID returns a clan by id
func GetClanByID(id int) (*Clan, error) {
	obj, err := db.Get(Clan{}, id)
	if err != nil {
		return nil, err
	}
	return obj.(*Clan), nil
}
