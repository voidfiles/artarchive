package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/artarchive/slides"
)

func TestStore(t *testing.T) {

	testTable := []struct {
		in1            string
		in2            []byte
		pass1          string
		pass2          slides.Slide
		expectedID     string
		expectError    bool
		expectErrorMsg error
		insertError    string
	}{
		{
			"123",
			[]byte(`{"site":{},"page":{},"edited":"0001-01-01T00:00:00Z","guid_hash":"123"}`),
			"123",
			slides.Slide{GUIDHash: "123"},
			"1",
			false,
			nil,
			"",
		}, {
			"123",
			[]byte(`{"site":{},"page":{},"edited":"0001-01-01T00:00:00Z","guid_hash":"yo"}`),
			"123",
			slides.Slide{GUIDHash: "yo"},
			"",
			true,
			fmt.Errorf("store err: An insert error"),
			"An insert error",
		},
	}
	for _, test := range testTable {
		dump, _ := json.Marshal(test.pass2)
		log.Printf("%v", string(dump))
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			log.Fatal(err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		store := MustNewItemStorage(sqlxDB)

		result := sqlmock.NewResult(1, 1)
		if test.insertError != "" {
			result = sqlmock.NewErrorResult(fmt.Errorf(test.insertError))
		}
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO items").
			WithArgs(test.in1, test.in2).
			WillReturnResult(result)
		val, err := store.Store(test.pass1, test.pass2)
		if test.expectError {
			assert.Equal(t, test.expectErrorMsg, err)
			mock.ExpectRollback()
		} else {
			assert.Equal(t, test.expectedID, val)
			mock.ExpectClose()
		}

	}

}

func TestFindByKey(t *testing.T) {

	testTable := []struct {
		key     string
		numRows int
		dataIn  []byte
		dataOut slides.Slide
		err     error
	}{
		{
			"a-key",
			1,
			[]byte(`{"guid_hash":"yo"}`),
			slides.Slide{GUIDHash: "yo"},
			nil,
		}, {
			"a-key",
			0,
			[]byte(`{"guid_hash":"yo"}`),
			slides.Slide{GUIDHash: "yo"},
			fmt.Errorf("store err: missing"),
		},
	}
	for _, test := range testTable {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			log.Fatal(err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		store := MustNewItemStorage(sqlxDB)

		rows := sqlmock.NewRows([]string{"id", "key", "data"})
		for _ = range make([]string, test.numRows) {
			rows.AddRow(1, test.key, test.dataIn)
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT").
			WithArgs(test.key).
			WillReturnRows(rows)
		mock.ExpectClose()
		data, err := store.FindByKey(test.key)

		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.dataOut, data)
		}

	}

}

func TestUpdateByKey(t *testing.T) {
	testTable := []struct {
		key        string
		dataIn     map[string]string
		dataInsert []byte
		err        error
	}{
		{
			"a-key",
			map[string]string{"key": "123"},
			[]byte(`{"key":"123"}`),
			nil,
		},
	}

	for _, test := range testTable {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			log.Fatal(err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		store := MustNewItemStorage(sqlxDB)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE items").
			WithArgs(test.dataInsert, test.key).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectClose()
		err = store.UpdateByKey(test.key, test.dataIn)

		if err != nil {
			assert.Equal(t, test.err, err)
		}

	}
}
