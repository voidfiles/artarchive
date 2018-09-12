package storage

import (
	"fmt"
	"log"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {

	testTable := []struct {
		in1            string
		in2            []byte
		pass1          string
		pass2          interface{}
		expectedID     string
		expectError    bool
		expectErrorMsg error
		insertError    string
	}{
		{
			"123",
			[]byte("\"yo\""),
			"123",
			"yo",
			"1",
			false,
			nil,
			"",
		}, {
			"123",
			[]byte("\"yo\""),
			"123",
			make(chan int),
			"",
			true,
			fmt.Errorf("store err: %v", "json: unsupported type: chan int"),
			"",
		}, {
			"123",
			[]byte("\"yo\""),
			"123",
			make(chan int),
			"",
			true,
			fmt.Errorf("store err: %v", "json: unsupported type: chan int"),
			"",
		}, {
			"123",
			[]byte("\"yo\""),
			"123",
			"yo",
			"",
			true,
			fmt.Errorf("store err: An insert error"),
			"An insert error",
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
		dataOut []byte
		err     error
	}{
		{
			"a-key",
			1,
			[]byte("yo"),
			[]byte("yo"),
			nil,
		}, {
			"a-key",
			1,
			[]byte(""),
			[]byte("{}"),
			nil,
		}, {
			"a-key",
			1,
			[]byte(nil),
			[]byte("{}"),
			nil,
		}, {
			"a-key",
			1,
			[]byte("{\"a"),
			[]byte("{\"a"),
			nil,
		}, {
			"a-key",
			0,
			[]byte("yo"),
			[]byte("yo"),
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

func TestUpdateFalg(t *testing.T) {

	testTable := []struct {
		key  string
		flag string
	}{
		{
			"a-key",
			"image_updated",
		}, {
			"a-key",
			"image_updated",
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
		mock.ExpectExec("UPDATE items SET").
			WithArgs(test.flag, test.key).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectClose()
		store.UpdateFlag(test.flag, test.key)

	}

}
