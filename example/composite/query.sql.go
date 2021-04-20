// Code generated by pggen. DO NOT EDIT.

package composite

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	SearchScreenshots(ctx context.Context, params SearchScreenshotsParams) ([]SearchScreenshotsRow, error)
	// SearchScreenshotsBatch enqueues a SearchScreenshots query into batch to be executed
	// later by the batch.
	SearchScreenshotsBatch(batch *pgx.Batch, params SearchScreenshotsParams)
	// SearchScreenshotsScan scans the result of an executed SearchScreenshotsBatch query.
	SearchScreenshotsScan(results pgx.BatchResults) ([]SearchScreenshotsRow, error)

	SearchScreenshotsOneCol(ctx context.Context, params SearchScreenshotsOneColParams) ([][]Blocks, error)
	// SearchScreenshotsOneColBatch enqueues a SearchScreenshotsOneCol query into batch to be executed
	// later by the batch.
	SearchScreenshotsOneColBatch(batch *pgx.Batch, params SearchScreenshotsOneColParams)
	// SearchScreenshotsOneColScan scans the result of an executed SearchScreenshotsOneColBatch query.
	SearchScreenshotsOneColScan(results pgx.BatchResults) ([][]Blocks, error)

	InsertScreenshotBlocks(ctx context.Context, screenshotID int, body string) (InsertScreenshotBlocksRow, error)
	// InsertScreenshotBlocksBatch enqueues a InsertScreenshotBlocks query into batch to be executed
	// later by the batch.
	InsertScreenshotBlocksBatch(batch *pgx.Batch, screenshotID int, body string)
	// InsertScreenshotBlocksScan scans the result of an executed InsertScreenshotBlocksBatch query.
	InsertScreenshotBlocksScan(results pgx.BatchResults) (InsertScreenshotBlocksRow, error)
}

type DBQuerier struct {
	conn genericConn
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{
		conn: conn,
	}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// preparer is any Postgres connection transport that provides a way to prepare
// a statement, most commonly *pgx.Conn.
type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

// PrepareAllQueries executes a PREPARE statement for all pggen generated SQL
// queries in querier files. Typical usage is as the AfterConnect callback
// for pgxpool.Config
//
// pgx will use the prepared statement if available. Calling PrepareAllQueries
// is an optional optimization to avoid a network round-trip the first time pgx
// runs a query if pgx statement caching is enabled.
func PrepareAllQueries(ctx context.Context, p preparer) error {
	if _, err := p.Prepare(ctx, searchScreenshotsSQL, searchScreenshotsSQL); err != nil {
		return fmt.Errorf("prepare query 'SearchScreenshots': %w", err)
	}
	if _, err := p.Prepare(ctx, searchScreenshotsOneColSQL, searchScreenshotsOneColSQL); err != nil {
		return fmt.Errorf("prepare query 'SearchScreenshotsOneCol': %w", err)
	}
	if _, err := p.Prepare(ctx, insertScreenshotBlocksSQL, insertScreenshotBlocksSQL); err != nil {
		return fmt.Errorf("prepare query 'InsertScreenshotBlocks': %w", err)
	}
	return nil
}

// newBlocksArrayDecoder creates a new decoder for the Postgres '_blocks' array type.
func newBlocksArrayDecoder() pgtype.ValueTranscoder {
	return pgtype.NewArrayType("_blocks", ignoredOID, newBlocksDecoder)
}

// Blocks represents the Postgres composite type "blocks".
type Blocks struct {
	ID           int    `json:"id"`
	ScreenshotID int    `json:"screenshot_id"`
	Body         string `json:"body"`
}

// newBlocksDecoder creates a new decoder for the Postgres 'blocks' composite type.
func newBlocksDecoder() pgtype.ValueTranscoder {
	return newCompositeType(
		"blocks",
		[]string{"id", "screenshot_id", "body"},
		&pgtype.Int4{},
		&pgtype.Int8{},
		&pgtype.Text{},
	)
}

// ignoredOID means we don't know or care about the OID for a type. This is okay
// because pgx only uses the OID to encode values and lookup a decoder. We only
// use ignoredOID for decoding and we always specify a concrete decoder for scan
// methods.
const ignoredOID = 0

func newCompositeType(name string, fieldNames []string, vals ...pgtype.ValueTranscoder) *pgtype.CompositeType {
	fields := make([]pgtype.CompositeTypeField, len(fieldNames))
	for i, name := range fieldNames {
		fields[i] = pgtype.CompositeTypeField{Name: name, OID: ignoredOID}
	}
	// Okay to ignore error because it's only thrown when the number of field
	// names does not equal the number of ValueTranscoders.
	rowType, _ := pgtype.NewCompositeTypeValues(name, fields, vals)
	return rowType
}

const searchScreenshotsSQL = `SELECT
  ss.id,
  array_agg(bl) AS blocks
FROM screenshots ss
  JOIN blocks bl ON bl.screenshot_id = ss.id
WHERE bl.body LIKE $1 || '%'
GROUP BY ss.id
ORDER BY ss.id
LIMIT $2 OFFSET $3;`

type SearchScreenshotsParams struct {
	Body   string
	Limit  int
	Offset int
}

type SearchScreenshotsRow struct {
	ID     int      `json:"id"`
	Blocks []Blocks `json:"blocks"`
}

// SearchScreenshots implements Querier.SearchScreenshots.
func (q *DBQuerier) SearchScreenshots(ctx context.Context, params SearchScreenshotsParams) ([]SearchScreenshotsRow, error) {
	rows, err := q.conn.Query(ctx, searchScreenshotsSQL, params.Body, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshots: %w", err)
	}
	defer rows.Close()
	items := []SearchScreenshotsRow{}
	blocksRow := newCompositeType(
		"blocks",
		[]string{"id", "screenshot_id", "body"},
		&pgtype.Int4{},
		&pgtype.Int8{},
		&pgtype.Text{},
	)
	blocksArray := pgtype.NewArrayType("_blocks", ignoredOID, func() pgtype.ValueTranscoder {
		return blocksRow.NewTypeValue().(*pgtype.CompositeType)
	})
	for rows.Next() {
		var item SearchScreenshotsRow
		if err := rows.Scan(&item.ID, blocksArray); err != nil {
			return nil, fmt.Errorf("scan SearchScreenshots row: %w", err)
		}
		if err := blocksArray.AssignTo(&item.Blocks); err != nil {
			return nil, fmt.Errorf("assign SearchScreenshots row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close SearchScreenshots rows: %w", err)
	}
	return items, err
}

// SearchScreenshotsBatch implements Querier.SearchScreenshotsBatch.
func (q *DBQuerier) SearchScreenshotsBatch(batch *pgx.Batch, params SearchScreenshotsParams) {
	batch.Queue(searchScreenshotsSQL, params.Body, params.Limit, params.Offset)
}

// SearchScreenshotsScan implements Querier.SearchScreenshotsScan.
func (q *DBQuerier) SearchScreenshotsScan(results pgx.BatchResults) ([]SearchScreenshotsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshotsBatch: %w", err)
	}
	defer rows.Close()
	items := []SearchScreenshotsRow{}
	blocksRow := newCompositeType(
		"blocks",
		[]string{"id", "screenshot_id", "body"},
		&pgtype.Int4{},
		&pgtype.Int8{},
		&pgtype.Text{},
	)
	blocksArray := pgtype.NewArrayType("_blocks", ignoredOID, func() pgtype.ValueTranscoder {
		return blocksRow.NewTypeValue().(*pgtype.CompositeType)
	})
	for rows.Next() {
		var item SearchScreenshotsRow
		if err := rows.Scan(&item.ID, blocksArray); err != nil {
			return nil, fmt.Errorf("scan SearchScreenshotsBatch row: %w", err)
		}
		if err := blocksArray.AssignTo(&item.Blocks); err != nil {
			return nil, fmt.Errorf("assign SearchScreenshots row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close SearchScreenshotsBatch rows: %w", err)
	}
	return items, err
}

const searchScreenshotsOneColSQL = `SELECT
  array_agg(bl) AS blocks
FROM screenshots ss
  JOIN blocks bl ON bl.screenshot_id = ss.id
WHERE bl.body LIKE $1 || '%'
GROUP BY ss.id
ORDER BY ss.id
LIMIT $2 OFFSET $3;`

type SearchScreenshotsOneColParams struct {
	Body   string
	Limit  int
	Offset int
}

// SearchScreenshotsOneCol implements Querier.SearchScreenshotsOneCol.
func (q *DBQuerier) SearchScreenshotsOneCol(ctx context.Context, params SearchScreenshotsOneColParams) ([][]Blocks, error) {
	rows, err := q.conn.Query(ctx, searchScreenshotsOneColSQL, params.Body, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshotsOneCol: %w", err)
	}
	defer rows.Close()
	items := [][]Blocks{}
	blocksRow := newCompositeType(
		"blocks",
		[]string{"id", "screenshot_id", "body"},
		&pgtype.Int4{},
		&pgtype.Int8{},
		&pgtype.Text{},
	)
	blocksArray := pgtype.NewArrayType("_blocks", ignoredOID, func() pgtype.ValueTranscoder {
		return blocksRow.NewTypeValue().(*pgtype.CompositeType)
	})
	for rows.Next() {
		var item []Blocks
		if err := rows.Scan(blocksArray); err != nil {
			return nil, fmt.Errorf("scan SearchScreenshotsOneCol row: %w", err)
		}
		if err := blocksArray.AssignTo(&item); err != nil {
			return nil, fmt.Errorf("assign SearchScreenshotsOneCol row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close SearchScreenshotsOneCol rows: %w", err)
	}
	return items, err
}

// SearchScreenshotsOneColBatch implements Querier.SearchScreenshotsOneColBatch.
func (q *DBQuerier) SearchScreenshotsOneColBatch(batch *pgx.Batch, params SearchScreenshotsOneColParams) {
	batch.Queue(searchScreenshotsOneColSQL, params.Body, params.Limit, params.Offset)
}

// SearchScreenshotsOneColScan implements Querier.SearchScreenshotsOneColScan.
func (q *DBQuerier) SearchScreenshotsOneColScan(results pgx.BatchResults) ([][]Blocks, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshotsOneColBatch: %w", err)
	}
	defer rows.Close()
	items := [][]Blocks{}
	blocksRow := newCompositeType(
		"blocks",
		[]string{"id", "screenshot_id", "body"},
		&pgtype.Int4{},
		&pgtype.Int8{},
		&pgtype.Text{},
	)
	blocksArray := pgtype.NewArrayType("_blocks", ignoredOID, func() pgtype.ValueTranscoder {
		return blocksRow.NewTypeValue().(*pgtype.CompositeType)
	})
	for rows.Next() {
		var item []Blocks
		if err := rows.Scan(blocksArray); err != nil {
			return nil, fmt.Errorf("scan SearchScreenshotsOneColBatch row: %w", err)
		}
		if err := blocksArray.AssignTo(&item); err != nil {
			return nil, fmt.Errorf("assign SearchScreenshotsOneCol row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close SearchScreenshotsOneColBatch rows: %w", err)
	}
	return items, err
}

const insertScreenshotBlocksSQL = `WITH screens AS (
  INSERT INTO screenshots (id) VALUES ($1)
    ON CONFLICT DO NOTHING
)
INSERT
INTO blocks (screenshot_id, body)
VALUES ($1, $2)
RETURNING id, screenshot_id, body;`

type InsertScreenshotBlocksRow struct {
	ID           int    `json:"id"`
	ScreenshotID int    `json:"screenshot_id"`
	Body         string `json:"body"`
}

// InsertScreenshotBlocks implements Querier.InsertScreenshotBlocks.
func (q *DBQuerier) InsertScreenshotBlocks(ctx context.Context, screenshotID int, body string) (InsertScreenshotBlocksRow, error) {
	row := q.conn.QueryRow(ctx, insertScreenshotBlocksSQL, screenshotID, body)
	var item InsertScreenshotBlocksRow
	if err := row.Scan(&item.ID, &item.ScreenshotID, &item.Body); err != nil {
		return item, fmt.Errorf("query InsertScreenshotBlocks: %w", err)
	}
	return item, nil
}

// InsertScreenshotBlocksBatch implements Querier.InsertScreenshotBlocksBatch.
func (q *DBQuerier) InsertScreenshotBlocksBatch(batch *pgx.Batch, screenshotID int, body string) {
	batch.Queue(insertScreenshotBlocksSQL, screenshotID, body)
}

// InsertScreenshotBlocksScan implements Querier.InsertScreenshotBlocksScan.
func (q *DBQuerier) InsertScreenshotBlocksScan(results pgx.BatchResults) (InsertScreenshotBlocksRow, error) {
	row := results.QueryRow()
	var item InsertScreenshotBlocksRow
	if err := row.Scan(&item.ID, &item.ScreenshotID, &item.Body); err != nil {
		return item, fmt.Errorf("scan InsertScreenshotBlocksBatch row: %w", err)
	}
	return item, nil
}
