package storage

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/voidfiles/artarchive/slides"
)

var (
	ErrMissingSlide = fmt.Errorf("store err: missing")
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
	ID   int64        `db:"id"`
	Key  string       `db:"key"`
	Data slides.Slide `db:"data"`
}

// MigrateDB will implment schema changes
func (i *ItemStorage) MigrateDB() {
	tx := i.db.MustBegin()
	tx.MustExec("CREATE SEQUENCE IF NOT EXISTS item_id_seq;")
	tx.MustExec(`CREATE TABLE IF NOT EXISTS  items (
      id bigint NOT NULL DEFAULT nextval('item_id_seq'),
      key varchar(500) NULL,
      data jsonb DEFAULT '{}',
			CONSTRAINT key_v1 UNIQUE(key)
  );`)
	tx.MustExec("ALTER SEQUENCE item_id_seq OWNED BY items.id;")
	tx.Commit()
}

// Store will store an item in the database
func (i *ItemStorage) Store(key string, data slides.Slide) (string, error) {

	tx := i.db.MustBegin()
	result := tx.MustExec("INSERT INTO items (key, data) VALUES ($1, $2) RETURNING id;", key, data)
	tx.Commit()

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("store err: %v", err)
	}

	returnID := fmt.Sprintf("%v", id)

	return returnID, nil
}

//FindByKey finds an item in the database
func (i *ItemStorage) FindByKey(key string) (slides.Slide, error) {

	target := storageDocument{}
	tx := i.db.MustBegin()
	tx.Get(&target, "SELECT * FROM items WHERE key = $1;", key)
	tx.Commit()

	if target.ID == 0 {
		return slides.Slide{}, ErrMissingSlide
	}

	return target.Data, nil
}

// UpdateByKey will store an item in the database
func (i *ItemStorage) UpdateByKey(key string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("store err: %v", err)
	}

	j := types.JSONText(string(b))
	log.Printf("%v", j)
	v, err := j.Value()
	if err != nil {
		return fmt.Errorf("store err: %v", err)
	}

	tx := i.db.MustBegin()
	tx.MustExec("UPDATE items SET data = $1 WHERE key = $2", v, key) // Good protection kind of  AND data->>'guid_hash' = $2;
	tx.Commit()

	return nil
}
