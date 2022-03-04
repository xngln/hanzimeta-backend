package hanzidata

import (
	"database/sql"
	"fmt"

	"github.com/xngln/hanzimeta_backend/db"
	"github.com/xngln/hanzimeta_backend/graph/model"
)

type HanziData struct {
	ID          string `db:"char_id"`
	Simplified  string
	Pinyin      string
	Traditional string
	Japanese    string
	JundaFreq   sql.NullInt16 `db:"junda_freq"`
	GSNum       sql.NullInt16 `db:"gs_num"`
	HSKLvl      sql.NullInt16 `db:"hsk_lvl"`
}

func Get(sortby *model.SortBy) []HanziData {
	hanzi := []HanziData{}
	query := fmt.Sprintf("SELECT * FROM characters ORDER BY %s %s LIMIT 50", sortby.Field, sortby.Order)
	db.DB.Select(&hanzi, query)
	return hanzi
}
