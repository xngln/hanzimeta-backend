package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/xngln/hanzimeta_backend/db/hanzidata"
	"github.com/xngln/hanzimeta_backend/graph/generated"
	"github.com/xngln/hanzimeta_backend/graph/model"
)

func (r *queryResolver) HanziConnection(ctx context.Context, first *int, after *string, sortBy *model.SortBy) (*model.HanziConnection, error) {
	var resultHanzi []*model.HanziData

	dbHanzi := hanzidata.Get(sortBy)
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
		resultHanzi = append(resultHanzi, &model.HanziData{
			ID:          hanzi.ID,
			Simplified:  hanzi.Simplified,
			Pinyin:      hanzi.Pinyin,
			Traditional: hanzi.Traditional,
			Japanese:    hanzi.Japanese,
			JundaFreq:   jfreq,
			GsNum:       gsnum,
			HskLvl:      hsk,
		})
	}

	// return resultHanzi, nil
	return nil, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
