package hanzidata

import (
	"database/sql"

	"github.com/xngln/hanzimeta_backend/db"
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

func Get() []HanziData {
	hanzi := []HanziData{}
	db.DB.Select(&hanzi, "SELECT * FROM characters ORDER BY junda_freq ASC LIMIT 10")
	return hanzi
}
