package sqlc

//go:generate sqlc generate -f sqlc.yml

func (q *Queries) WithConn(tx DBTX) *Queries {
	return &Queries{
		db: tx,
	}
}
