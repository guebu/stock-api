package queries

import "github.com/olivere/elastic"

func (esq *EsQuery) Build() elastic.Query {

	query := elastic.NewBoolQuery()

	equalsQueries := make([]elastic.Query, 0)
	for _, eq := range esq.Equals {
		boolQuery := elastic.NewMatchQuery(eq.Field, eq.Value)
		equalsQueries = append(equalsQueries, boolQuery)
	}
	query.Must(equalsQueries...)

	return query
}