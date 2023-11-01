package elasticsearch

import "github.com/olivere/elastic/v7"

func NewBoolQuery(column string, keyword string) *elastic.BoolQuery {
	return elastic.NewBoolQuery().Should(
		elastic.NewMatchQuery(column, keyword).
			Fuzziness("AUTO").
			Operator("and").
			FuzzyTranspositions(true),
	)
}
