package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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
	slide, ok := data.(slides.Slide)
	if !ok {
		return fmt.Errorf("store err: data failed .(slide.Slide) casting")
	}
	tx := i.db.MustBegin()
	tx.MustExec("UPDATE items SET data = $1 WHERE key = $2", slide, key) // Good protection kind of  AND data->>'guid_hash' = $2;
	tx.Commit()

	return nil
}

// List will find all the keys in the database
func (i *ItemStorage) List(after int64) ([]slides.Slide, int64, error) {
	slides := make([]slides.Slide, 0)
	target := make([]storageDocument, 0)
	tx := i.db.MustBegin()
	err := tx.Select(&target, "SELECT * FROM items WHERE id >= $1 ORDER BY id LIMIT 100;", after)
	if err != nil {
		return slides, 0, err
	}
	tx.Commit()

	if len(target) == 0 {
		return slides, 0, nil
	}
	next := int64(-1)
	if len(target) >= 100 {
		next = target[len(target)-1].ID
	}
	for _, doc := range target {
		slides = append(slides, doc.Data)
	}
	return slides, next, nil
}

// List will find all the keys in the database
func (i *ItemStorage) FindSites(query string, after int64) ([]slides.Site, int64, error) {
	sites := make([]slides.Site, 0)
	target := make([]storageDocument, 0)
	tx := i.db.MustBegin()
	query = "%" + query + "%"
	err := tx.Select(&target, "SELECT * FROM items WHERE id >= $1 AND data->'site'->>'title' LIKE $2 ORDER BY id LIMIT 100;", after, query)
	if err != nil {
		return sites, 0, err
	}
	tx.Commit()

	if len(target) == 0 {
		return sites, 0, nil
	}
	next := int64(-1)
	if len(target) >= 100 {
		next = target[len(target)-1].ID
	}
	for _, doc := range target {
		sites = append(sites, doc.Data.Site)
	}
	return sites, next, nil
}
