package relational

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type codeRelatedQuery interface {
	InsertCodes(codes []string) error
	GetCodes() (map[string]struct{}, error)
	DeleteCodes() error
}

type codeTable struct {
	conn postgresConn
}

func (c *codeTable) InsertCodes(codes []string) error {
	values := strings.Join(codes, ",")
	query := fmt.Sprintf(insertCodesFormat, values)

	_, err := c.conn.pool.Exec(c.conn.ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (c *codeTable) GetCodes() (map[string]struct{}, error) {
	codes := make(map[string]struct{})

	rows, err := c.conn.pool.Query(c.conn.ctx, getCodesQuery)
	if err != nil {
		if err == pgx.ErrNoRows {
			return codes, nil
		}

		return nil, err
	}

	for rows.Next() {
		var code string
		err = rows.Scan(&code)
		if err != nil {
			return nil, err
		}

		codes[code] = struct{}{}
	}

	return codes, nil
}

func (c *codeTable) DeleteCodes() error {
	_, err := c.conn.pool.Exec(c.conn.ctx, deleteCodesQuery)
	if err != nil {
		return err
	}

	return nil
}

var (
	insertCodesFormat = `
	INSERT INTO hsr_related.code (key) VALUES (%s);`

	getCodesQuery = `
	SELECT * FROM hsr_related.code;`

	deleteCodesQuery = `
	DELETE FROM hsr_related.code;`
)
