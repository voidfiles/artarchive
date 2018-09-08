package main

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

// ItemStorage manages storage of items to db
type ItemStorage struct {
	db *sqlx.DB
}

// MustNewItemStorage creates and returns a new ItemStorage
func MustNewItemStorage(db *sqlx.DB) *ItemStorage {
	return &ItemStorage{
		db: db,
	}
}

type storageDocument struct {
	Id         int64          `db:"id"`
	ImageSaved bool           `db:"image_saved"`
	ItemSaved  bool           `db:"item_saved"`
	Key        string         `db:"key"`
	Data       types.JSONText `db:"data"`
}

// MigrateDB will implment schema changes
func (i *ItemStorage) MigrateDB() {
	tx := i.db.MustBegin()
	tx.MustExec("CREATE SEQUENCE item_id_seq;")
	tx.MustExec(`CREATE TABLE items {
      id bigint NOT NULL DEFAULT nextval('item_id_seq');
			image_saved bool
			item_saved bool
      key varchar(500) NULL;
      data jsonb DEFAULT '{}';
  };`)
	tx.MustExec("ALTER SEQUENCE item_id_seq OWNED BY items.id;")
	tx.Commit()
}

// Store will store an item in the database
func (i *ItemStorage) Store(key string, data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("store err: %v", err)
	}

	j := types.JSONText(string(b))

	v, err := j.Value()
	if err != nil {
		return "", fmt.Errorf("store err: %v", err)
	}
	tx := i.db.MustBegin()
	result := tx.MustExec("INSERT INTO items (key, data) VALUES (?, ?) RETURNING id", key, v)
	tx.Commit()

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("store err: %v", err)
	}

	returnID := fmt.Sprintf("%v", id)

	return returnID, nil
}

//FindByKey finds an item in the database
func (i *ItemStorage) FindByKey(key string) ([]byte, error) {

	target := storageDocument{}
	tx := i.db.MustBegin()
	tx.Get(&target, "SELECT * FROM items WHERE key = ?", key)
	tx.Commit()

	if target.Id == 0 {
		return []byte(""), fmt.Errorf("store err: missing")
	}

	data, err := target.Data.MarshalJSON()

	if err != nil {
		return []byte(""), fmt.Errorf("store err: %v", err)
	}
	return data, nil
}

//FindByKey finds an item in the database
func (i *ItemStorage) UpdateFlag(flag, key string) error {

	tx := i.db.MustBegin()
	tx.MustExec("UPDATE items SET (?) VALUES (true) WHERE key = ?", flag, key)
	tx.Commit()

	return nil
}
