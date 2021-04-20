// Code generated by pggen. DO NOT EDIT.

package nested

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
	Nested3(ctx context.Context) ([]Qux, error)
	// Nested3Batch enqueues a Nested3 query into batch to be executed
	// later by the batch.
	Nested3Batch(batch *pgx.Batch)
	// Nested3Scan scans the result of an executed Nested3Batch query.
	Nested3Scan(results pgx.BatchResults) ([]Qux, error)
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
	if _, err := p.Prepare(ctx, nested3SQL, nested3SQL); err != nil {
		return fmt.Errorf("prepare query 'Nested3': %w", err)
	}
	return nil
}

// InventoryItem represents the Postgres composite type "inventory_item".
type InventoryItem struct {
	ItemName *string `json:"item_name"`
	Sku      Sku     `json:"sku"`
}

// Qux represents the Postgres composite type "qux".
type Qux struct {
	InvItem InventoryItem `json:"inv_item"`
	Foo     *int          `json:"foo"`
}

// Sku represents the Postgres composite type "sku".
type Sku struct {
	SkuID *string `json:"sku_id"`
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

const nested3SQL = `SELECT ROW (ROW ('item_name', ROW ('sku_id')::sku)::inventory_item, 88)::qux AS qux;`

// Nested3 implements Querier.Nested3.
func (q *DBQuerier) Nested3(ctx context.Context) ([]Qux, error) {
	rows, err := q.conn.Query(ctx, nested3SQL)
	if err != nil {
		return nil, fmt.Errorf("query Nested3: %w", err)
	}
	defer rows.Close()
	items := []Qux{}
	quxRow := newCompositeType(
		"qux",
		[]string{"inv_item", "foo"},
		newCompositeType(
			"inventory_item",
			[]string{"item_name", "sku"},
			&pgtype.Text{},
			newCompositeType(
				"sku",
				[]string{"sku_id"},
				&pgtype.Text{},
			),
		),
		&pgtype.Int8{},
	)
	for rows.Next() {
		var item Qux
		if err := rows.Scan(quxRow); err != nil {
			return nil, fmt.Errorf("scan Nested3 row: %w", err)
		}
		if err := quxRow.AssignTo(&item); err != nil {
			return nil, fmt.Errorf("assign Nested3 row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close Nested3 rows: %w", err)
	}
	return items, err
}

// Nested3Batch implements Querier.Nested3Batch.
func (q *DBQuerier) Nested3Batch(batch *pgx.Batch) {
	batch.Queue(nested3SQL)
}

// Nested3Scan implements Querier.Nested3Scan.
func (q *DBQuerier) Nested3Scan(results pgx.BatchResults) ([]Qux, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query Nested3Batch: %w", err)
	}
	defer rows.Close()
	items := []Qux{}
	quxRow := newCompositeType(
		"qux",
		[]string{"inv_item", "foo"},
		newCompositeType(
			"inventory_item",
			[]string{"item_name", "sku"},
			&pgtype.Text{},
			newCompositeType(
				"sku",
				[]string{"sku_id"},
				&pgtype.Text{},
			),
		),
		&pgtype.Int8{},
	)
	for rows.Next() {
		var item Qux
		if err := rows.Scan(quxRow); err != nil {
			return nil, fmt.Errorf("scan Nested3Batch row: %w", err)
		}
		if err := quxRow.AssignTo(&item); err != nil {
			return nil, fmt.Errorf("assign Nested3 row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close Nested3Batch rows: %w", err)
	}
	return items, err
}
