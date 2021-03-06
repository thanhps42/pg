package orm

import (
	"errors"
	"strings"

	"github.com/go-pg/pg/types"
)

func HasTable(db DB, model interface{}) (bool, error) {
	q := NewQuery(db, model)
	var count int
	_, err := q.db.QueryOne(Scan(&count), hasTableQuery{
		q: q,
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type hasTableQuery struct {
	q *Query
}

func (q hasTableQuery) Copy() QueryAppender {
	return q
}

func (q hasTableQuery) Query() *Query {
	return q.q
}

func (q hasTableQuery) AppendQuery(b []byte) ([]byte, error) {
	if q.q.stickyErr != nil {
		return nil, q.q.stickyErr
	}
	if q.q.model == nil {
		return nil, errors.New("pg: Model(nil)")
	}

	tableName := string(q.q.appendTableName(nil))
	tableName = strings.TrimPrefix(tableName, `"`)
	tableName = strings.TrimSuffix(tableName, `"`)

	b = append(b, "SELECT count(*) FROM pg_tables WHERE schemaname = 'public' AND tablename = "...)
	b = types.AppendString(b, tableName, 1)

	return b, nil
}
