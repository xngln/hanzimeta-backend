package hanzidata

import (
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/xngln/hanzimeta-backend/db"
	"github.com/xngln/hanzimeta-backend/graph/model"
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

func buildQuery(sortby *model.SortBy, first int, after *string) (string, error) {
	// no cursor, just return a query for the first n rows
	if after == nil {
		return fmt.Sprintf(
			"SELECT * FROM characters ORDER BY %s %s, char_id %s LIMIT %d;",
			*sortby.Field,
			*sortby.Order,
			*sortby.Order,
			first,
		), nil
	}

	var query string
	var decodedCursor string
	b, err := b64.StdEncoding.DecodeString(*after)
	if err != nil {
		return query, fmt.Errorf("failed to decode the cursor: %s", decodedCursor)
	}
	decodedCursor = string(b)

	// if string starts with "page ", it is an offset cursor (format: "page <number>")
	// otherwise it is a pointer cursor (the cursor is a json of the actual hanzidata element)
	if strings.HasPrefix(decodedCursor, "page") {
		words := strings.Fields(decodedCursor)
		pageNum, err := strconv.Atoi(words[1])
		if err != nil {
			return query, fmt.Errorf("failed to get page num from offset cursor: %s", decodedCursor)
		}
		pageSize := first - 1
		count, err := GetCount()
		if err != nil {
			return query, err
		}
		numPages := int(math.Ceil(float64(count) / float64(pageSize)))
		if pageNum > numPages {
			return query, fmt.Errorf("page number out of range")
		}
		if err != nil {
			return query, fmt.Errorf("invalid format for cursor: %s", decodedCursor)
		}
		query = buildOffsetQuery(pageNum, sortby, first)
		return query, nil
	}

	cursorID, err := strconv.Atoi(decodedCursor)
	if err != nil {
		return query, fmt.Errorf("failed to convert cursor to id: %s", decodedCursor)
	}
	query, err = buildCursorQuery(cursorID, sortby, first)
	if err != nil {
		return query, fmt.Errorf("failed to build cursor query using cursor: %s", decodedCursor)
	}
	return query, nil
}

func buildCursorQuery(charID int, sortby *model.SortBy, first int) (string, error) {
	var c *int
	q := fmt.Sprintf("SELECT %s FROM characters WHERE char_id=%d;", *sortby.Field, charID)
	err := db.DB.Get(&c, q)
	if err != nil {
		return "", fmt.Errorf("failed to get cursor sort field value, charID: %d", charID)
	}

	comparison := ">"
	if *sortby.Order == model.OrderDesc {
		comparison = "<"
	}

	if c == nil {
		return fmt.Sprintf(
			`SELECT * FROM characters 
			WHERE ( %s is NULL AND char_id %s %d)
			ORDER BY %s %s, char_id %s
			LIMIT %d;`,
			*sortby.Field,
			comparison,
			charID,
			*sortby.Field,
			*sortby.Order,
			*sortby.Order,
			first,
		), nil
	}

	sortVal := strconv.Itoa(*c)
	return fmt.Sprintf(
		`SELECT * FROM characters 
		WHERE %s %s %s
		OR (
			%s = %s AND
			char_id %s %d
		)
		ORDER BY %s %s, char_id %s
		LIMIT %d;`,
		*sortby.Field,
		comparison,
		sortVal,
		*sortby.Field,
		sortVal,
		comparison,
		charID,
		*sortby.Field,
		*sortby.Order,
		*sortby.Order,
		first,
	), nil
}

func buildOffsetQuery(pageNum int, sortby *model.SortBy, first int) string {
	var query string
	// first is one large than page size.
	// we want to query one extra element to check if this is the last page
	pagesize := first - 1
	offset := pagesize * (pageNum - 1)

	query = fmt.Sprintf(
		`SELECT * FROM characters
		ORDER BY %s %s, char_id
		LIMIT %d
		OFFSET %d;`,
		*sortby.Field,
		*sortby.Order,
		first,
		offset,
	)

	return query
}

func GetPage(sortby *model.SortBy, first int, after *string) ([]HanziData, *model.PageInfo, error) {
	pageInfo := &model.PageInfo{}
	hanzi := []HanziData{}

	query, err := buildQuery(sortby, first+1, after)
	if err != nil {
		return nil, nil, err
	}
	db.DB.Select(&hanzi, query)

	if len(hanzi) == 0 {
		return nil, nil, fmt.Errorf("got 0 rows from db, something probably went wrong on the server side")
	}

	// we queried first+1 rows from db. if query returned less than that, then this is the last page
	if len(hanzi) <= first {
		pageInfo.HasNextPage = false
	} else {
		pageInfo.HasNextPage = true
		hanzi = hanzi[:len(hanzi)-1] // remove the first+1th element
	}

	pageInfo.StartCursor = b64.StdEncoding.EncodeToString([]byte(hanzi[0].ID))
	pageInfo.EndCursor = b64.StdEncoding.EncodeToString([]byte(hanzi[len(hanzi)-1].ID))

	return hanzi, pageInfo, nil
}

func GetByChar(character string) []HanziData {
	hanzi := []HanziData{}
	query := fmt.Sprintf(
		`SELECT * FROM characters
		WHERE 
		Simplified = '%s' OR
		Traditional = '%s' OR
		Japanese = '%s';`,
		character,
		character,
		character,
	)
	db.DB.Select(&hanzi, query)
	return hanzi
}

func GetCount() (int, error) {
	var count int
	err := db.DB.Get(&count, "SELECT COUNT(char_id) FROM characters;")
	if err != nil {
		return count, fmt.Errorf("failed to get count from 'characters' db")
	}
	return count, nil
}
