package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	b64 "encoding/base64"

	"github.com/xngln/hanzimeta-backend/db/hanzidata"
	"github.com/xngln/hanzimeta-backend/graph/generated"
	"github.com/xngln/hanzimeta-backend/graph/model"
)

func (r *queryResolver) HanziConnection(ctx context.Context, first *int, after *string, sortBy *model.SortBy) (*model.HanziConnection, error) {
	var err error
	connection := &model.HanziConnection{}
	connection.TotalCount, err = hanzidata.GetCount()
	if err != nil {
		return nil, err
	}

	dbHanzi, pageInfo, err := hanzidata.GetPage(sortBy, *first, after)
	if err != nil {
		return nil, err
	}
	connection.PageInfo = pageInfo
	for _, hanzi := range dbHanzi {
		var jfreq *int
		if hanzi.JundaFreq.Valid {
			jfreq = new(int)
			*jfreq = int(hanzi.JundaFreq.Int16)
		}
		var gsnum *int
		if hanzi.GSNum.Valid {
			gsnum = new(int)
			*gsnum = int(hanzi.GSNum.Int16)
		}
		var hsk *int
		if hanzi.HSKLvl.Valid {
			hsk = new(int)
			*hsk = int(hanzi.HSKLvl.Int16)
		}
		edge := &model.HanziEdge{
			Cursor: b64.StdEncoding.EncodeToString([]byte(hanzi.ID)),
			Node: &model.HanziData{
				ID:          hanzi.ID,
				Simplified:  hanzi.Simplified,
				Pinyin:      hanzi.Pinyin,
				Traditional: hanzi.Traditional,
				Japanese:    hanzi.Japanese,
				JundaFreq:   jfreq,
				GsNum:       gsnum,
				HskLvl:      hsk,
			},
		}
		connection.Edges = append(connection.Edges, edge)
	}

	// return resultHanzi, nil
	return connection, nil
}

func (r *queryResolver) Hanzi(ctx context.Context, character string) ([]*model.HanziData, error) {
	var hanzi []*model.HanziData
	dbHanzi := hanzidata.GetByChar(character)

	for _, h := range dbHanzi {
		var jfreq *int
		if h.JundaFreq.Valid {
			jfreq = new(int)
			*jfreq = int(h.JundaFreq.Int16)
		}
		var gsnum *int
		if h.GSNum.Valid {
			gsnum = new(int)
			*gsnum = int(h.GSNum.Int16)
		}
		var hsk *int
		if h.HSKLvl.Valid {
			hsk = new(int)
			*hsk = int(h.HSKLvl.Int16)
		}

		hanzi = append(hanzi, &model.HanziData{
			ID:          h.ID,
			Simplified:  h.Simplified,
			Pinyin:      h.Pinyin,
			Traditional: h.Traditional,
			Japanese:    h.Japanese,
			JundaFreq:   jfreq,
			GsNum:       gsnum,
			HskLvl:      hsk,
		})
	}
	return hanzi, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
