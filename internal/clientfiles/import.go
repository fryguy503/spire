package clientfiles

import (
	"database/sql"
	"fmt"
	"github.com/Akkadius/spire/internal/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

type Importer struct {
	logger *logrus.Logger
}

func NewImporter(logger *logrus.Logger) *Importer {
	return &Importer{logger: logger}
}

func (i *Importer) getDatabase(g *gorm.DB) *sql.DB {
	// get database instance
	db, err := g.DB()
	if err != nil {
		i.logger.Fatal(err)
	}

	return db
}

const SPELL_COLUMN_LENGTH = 237

type ImportResult struct {
	Table        string
	ImportedRows int
	DroppedRows  int64
}

func (i *Importer) ImportSpells(db *gorm.DB, fileContents string) (ImportResult, error) {

	// purge
	rowsDeleted := db.Exec("DELETE FROM spells_new").RowsAffected

	// build column name list
	dbColumns := database.GetTableColumnsOrdered(db, "spells_new")
	dbColumnsStr := "`" + strings.Join(dbColumns, "`, `") + "`"

	// build parameter bindings string
	params := []string{}
	for _, _ = range dbColumns {
		params = append(params, "?")
	}
	paramsStr := strings.Join(params, ",")

	// vars
	placeholders := []string{}
	var values []interface{}
	chunk := 0

	i.logger.Debugf("[import] Database columns [%v]\n", len(dbColumns))

	// loop through lines
	processedRows := 0
	for _, s := range strings.Split(fileContents, "\n") {

		// process columns
		columns := strings.Split(s, "^")
		rowHasValues := false
		for _, column := range columns {
			if len(columns) == SPELL_COLUMN_LENGTH {
				values = append(values, column)
				rowHasValues = true
			}
		}

		// if row has expected values, append
		if rowHasValues {
			placeholders = append(placeholders, fmt.Sprintf("(%s)", paramsStr))
		}

		// flush chunk
		if chunk >= 250 {
			err := i.insertBulk(db, "spells_new", dbColumnsStr, placeholders, values)
			if err != nil {
				return ImportResult{}, err
			}
			placeholders = nil
			values = nil
			chunk = 0
		}

		chunk++
		processedRows++
	}

	i.logger.Debugf("[import] Processed rows [%v]\n", processedRows)

	err := i.insertBulk(db, "spells_new", dbColumnsStr, placeholders, values)
	if err != nil {
		return ImportResult{}, err
	}

	return ImportResult{
		Table:        "spells_new",
		ImportedRows: processedRows,
		DroppedRows:  rowsDeleted,
	}, nil
}

const DB_STR_COLUMN_LENGTH = 4

func (i *Importer) ImportDbStr(db *gorm.DB, fileContents string) (ImportResult, error) {

	// purge
	rowsDeleted := db.Exec("DELETE FROM db_str").RowsAffected

	placeholders := []string{}
	var values []interface{}
	chunk := 0

	// loop through lines
	processedRows := 0
	for _, s := range strings.Split(fileContents, "\n") {

		// process columns
		columns := strings.Split(s, "^")
		rowHasValues := false
		for count, column := range columns {

			// Only 3 database columns, 4 file columns
			if len(columns) == DB_STR_COLUMN_LENGTH && count < 3 {
				values = append(values, column)
				rowHasValues = true
			}
		}

		// if row has expected values, append
		if rowHasValues {
			placeholders = append(placeholders, "(?, ?, ?)")
		}

		// flush chunk
		if chunk >= 10000 {
			err := i.insertBulk(db, "db_str", "id, type, value", placeholders, values)
			if err != nil {
				return ImportResult{}, err
			}
			placeholders = nil
			values = nil
			chunk = 0
		}

		chunk++
		processedRows++
	}

	err := i.insertBulk(db, "db_str", "id, type, value", placeholders, values)
	if err != nil {
		return ImportResult{}, err
	}

	return ImportResult{
		Table:        "db_str",
		ImportedRows: processedRows,
		DroppedRows:  rowsDeleted,
	}, nil
}

func (i *Importer) insertBulk(db *gorm.DB, table string, columns string, placeholders []string, values []interface{}) error {
	_, err := i.getDatabase(db).Exec(
		fmt.Sprintf("INSERT INTO %s (%s) VALUES %v", table, columns, strings.Join(placeholders, ", ")),
		values...,
	)

	//pp.Println(fmt.Sprintf("INSERT INTO db_str (id, type, value) VALUES %v", strings.Join(placeholders, ", ")))
	//pp.Println(values)
	//fmt.Println(err)

	if err != nil {
		return err
	}

	return nil
}
