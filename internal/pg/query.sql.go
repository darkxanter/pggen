// Code generated by pggen. DO NOT EDIT.

package pg

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
	FindEnumTypes(ctx context.Context, oIDs []uint32) ([]FindEnumTypesRow, error)
	// FindEnumTypesBatch enqueues a FindEnumTypes query into batch to be executed
	// later by the batch.
	FindEnumTypesBatch(batch *pgx.Batch, oIDs []uint32)
	// FindEnumTypesScan scans the result of an executed FindEnumTypesBatch query.
	FindEnumTypesScan(results pgx.BatchResults) ([]FindEnumTypesRow, error)

	// A composite type represents a row or record, defined implicitly for each
	// table, or explicitly with CREATE TYPE.
	// https://www.postgresql.org/docs/13/rowtypes.html
	FindCompositeTypes(ctx context.Context, oIDs []uint32) ([]FindCompositeTypesRow, error)
	// FindCompositeTypesBatch enqueues a FindCompositeTypes query into batch to be executed
	// later by the batch.
	FindCompositeTypesBatch(batch *pgx.Batch, oIDs []uint32)
	// FindCompositeTypesScan scans the result of an executed FindCompositeTypesBatch query.
	FindCompositeTypesScan(results pgx.BatchResults) ([]FindCompositeTypesRow, error)

	FindOIDByName(ctx context.Context, name string) (pgtype.OID, error)
	// FindOIDByNameBatch enqueues a FindOIDByName query into batch to be executed
	// later by the batch.
	FindOIDByNameBatch(batch *pgx.Batch, name string)
	// FindOIDByNameScan scans the result of an executed FindOIDByNameBatch query.
	FindOIDByNameScan(results pgx.BatchResults) (pgtype.OID, error)

	FindOIDName(ctx context.Context, oID pgtype.OID) (pgtype.Name, error)
	// FindOIDNameBatch enqueues a FindOIDName query into batch to be executed
	// later by the batch.
	FindOIDNameBatch(batch *pgx.Batch, oID pgtype.OID)
	// FindOIDNameScan scans the result of an executed FindOIDNameBatch query.
	FindOIDNameScan(results pgx.BatchResults) (pgtype.Name, error)
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

const findEnumTypesSQL = `WITH enums AS (
  SELECT enumtypid::int8                                   AS enum_type,
         -- pg_enum row identifier.
         -- The OIDs for pg_enum rows follow a special rule: even-numbered OIDs
         -- are guaranteed to be ordered in the same way as the sort ordering of
         -- their enum type. That is, if two even OIDs belong to the same enum
         -- type, the smaller OID must have the smaller enumsortorder value.
         -- Odd-numbered OID values need bear no relationship to the sort order.
         -- This rule allows the enum comparison routines to avoid catalog
         -- lookups in many common cases. The routines that create and alter enum
         -- types attempt to assign even OIDs to enum values whenever possible.
         array_agg(oid::int8 ORDER BY enumsortorder)       AS enum_oids,
         -- The sort position of this enum value within its enum type. Starts as
         -- 1..n but can be fractional or negative.
         array_agg(enumsortorder ORDER BY enumsortorder)   AS enum_orders,
         -- The textual label for this enum value
         array_agg(enumlabel::text ORDER BY enumsortorder) AS enum_labels
  FROM pg_enum
  GROUP BY pg_enum.enumtypid)
SELECT
  typ.oid           AS oid,
  -- typename: Data type name.
  typ.typname::text AS type_name,
  enum.enum_oids    AS child_oids,
  enum.enum_orders  AS orders,
  enum.enum_labels  AS labels,
  -- typtype: b for a base type, c for a composite type (e.g., a table's
  -- row type), d for a domain, e for an enum type, p for a pseudo-type,
  -- or r for a range type.
  typ.typtype       AS type_kind,
  -- typdefault is null if the type has no associated default value. If
  -- typdefaultbin is not null, typdefault must contain a human-readable
  -- version of the default expression represented by typdefaultbin. If
  -- typdefaultbin is null and typdefault is not, then typdefault is the
  -- external representation of the type's default value, which can be fed
  -- to the type's input converter to produce a constant.
  typ.typdefault    AS default_expr
FROM pg_type typ
  JOIN enums enum ON typ.oid = enum.enum_type
WHERE typ.typisdefined
  AND typ.typtype = 'e'
  AND typ.oid = ANY ($1::oid[]);`

type FindEnumTypesRow struct {
	OID         pgtype.OID         `json:"oid"`
	TypeName    pgtype.Text        `json:"type_name"`
	ChildOIDs   pgtype.Int8Array   `json:"child_oids"`
	Orders      pgtype.Float4Array `json:"orders"`
	Labels      pgtype.TextArray   `json:"labels"`
	TypeKind    pgtype.QChar       `json:"type_kind"`
	DefaultExpr pgtype.Text        `json:"default_expr"`
}

// FindEnumTypes implements Querier.FindEnumTypes.
func (q *DBQuerier) FindEnumTypes(ctx context.Context, oIDs []uint32) ([]FindEnumTypesRow, error) {
	rows, err := q.conn.Query(ctx, findEnumTypesSQL, oIDs)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("query FindEnumTypes: %w", err)
	}
	items := []FindEnumTypesRow{}
	for rows.Next() {
		var item FindEnumTypesRow
		if err := rows.Scan(&item.OID, &item.TypeName, &item.ChildOIDs, &item.Orders, &item.Labels, &item.TypeKind, &item.DefaultExpr); err != nil {
			return nil, fmt.Errorf("scan FindEnumTypes row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// FindEnumTypesBatch implements Querier.FindEnumTypesBatch.
func (q *DBQuerier) FindEnumTypesBatch(batch *pgx.Batch, oIDs []uint32) {
	batch.Queue(findEnumTypesSQL, oIDs)
}

// FindEnumTypesScan implements Querier.FindEnumTypesScan.
func (q *DBQuerier) FindEnumTypesScan(results pgx.BatchResults) ([]FindEnumTypesRow, error) {
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	items := []FindEnumTypesRow{}
	for rows.Next() {
		var item FindEnumTypesRow
		if err := rows.Scan(&item.OID, &item.TypeName, &item.ChildOIDs, &item.Orders, &item.Labels, &item.TypeKind, &item.DefaultExpr); err != nil {
			return nil, fmt.Errorf("scan FindEnumTypesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

const findCompositeTypesSQL = `WITH table_cols AS (
  SELECT
    cls.relname                                         AS table_name,
    cls.oid                                             AS table_oid,
    array_agg(attr.attname::text ORDER BY attr.attnum)  AS col_names,
    array_agg(attr.atttypid::int8 ORDER BY attr.attnum) AS col_oids,
    array_agg(attr.attnum::int8 ORDER BY attr.attnum)   AS col_orders,
    array_agg(attr.attnotnull ORDER BY attr.attnum)     AS col_not_nulls,
    array_agg(typ.typname::text ORDER BY attr.attnum)   AS col_type_names
  FROM pg_attribute attr
    JOIN pg_class cls ON attr.attrelid = cls.oid
    JOIN pg_type typ ON typ.oid = attr.atttypid
  WHERE attr.attnum > 0 -- Postgres represents system columns with attnum <= 0
    AND NOT attr.attisdropped
  GROUP BY cls.relname, cls.oid
)
SELECT
  typ.typname::text AS table_type_name,
  typ.oid           AS table_type_oid,
  table_name,
  col_names,
  col_oids,
  col_orders,
  col_not_nulls,
  col_type_names
FROM pg_type typ
  JOIN table_cols cols ON typ.typrelid = cols.table_oid
WHERE typ.oid = ANY ($1::oid[])
  AND typ.typtype = 'c';`

type FindCompositeTypesRow struct {
	TableTypeName pgtype.Text      `json:"table_type_name"`
	TableTypeOID  pgtype.OID       `json:"table_type_oid"`
	TableName     pgtype.Name      `json:"table_name"`
	ColNames      pgtype.TextArray `json:"col_names"`
	ColOIDs       pgtype.Int8Array `json:"col_oids"`
	ColOrders     pgtype.Int8Array `json:"col_orders"`
	ColNotNulls   pgtype.BoolArray `json:"col_not_nulls"`
	ColTypeNames  pgtype.TextArray `json:"col_type_names"`
}

// FindCompositeTypes implements Querier.FindCompositeTypes.
func (q *DBQuerier) FindCompositeTypes(ctx context.Context, oIDs []uint32) ([]FindCompositeTypesRow, error) {
	rows, err := q.conn.Query(ctx, findCompositeTypesSQL, oIDs)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("query FindCompositeTypes: %w", err)
	}
	items := []FindCompositeTypesRow{}
	for rows.Next() {
		var item FindCompositeTypesRow
		if err := rows.Scan(&item.TableTypeName, &item.TableTypeOID, &item.TableName, &item.ColNames, &item.ColOIDs, &item.ColOrders, &item.ColNotNulls, &item.ColTypeNames); err != nil {
			return nil, fmt.Errorf("scan FindCompositeTypes row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// FindCompositeTypesBatch implements Querier.FindCompositeTypesBatch.
func (q *DBQuerier) FindCompositeTypesBatch(batch *pgx.Batch, oIDs []uint32) {
	batch.Queue(findCompositeTypesSQL, oIDs)
}

// FindCompositeTypesScan implements Querier.FindCompositeTypesScan.
func (q *DBQuerier) FindCompositeTypesScan(results pgx.BatchResults) ([]FindCompositeTypesRow, error) {
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	items := []FindCompositeTypesRow{}
	for rows.Next() {
		var item FindCompositeTypesRow
		if err := rows.Scan(&item.TableTypeName, &item.TableTypeOID, &item.TableName, &item.ColNames, &item.ColOIDs, &item.ColOrders, &item.ColNotNulls, &item.ColTypeNames); err != nil {
			return nil, fmt.Errorf("scan FindCompositeTypesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

const findOIDByNameSQL = `SELECT oid
FROM pg_type
WHERE typname::text = $1;`

// FindOIDByName implements Querier.FindOIDByName.
func (q *DBQuerier) FindOIDByName(ctx context.Context, name string) (pgtype.OID, error) {
	row := q.conn.QueryRow(ctx, findOIDByNameSQL, name)
	var item pgtype.OID
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindOIDByName: %w", err)
	}
	return item, nil
}

// FindOIDByNameBatch implements Querier.FindOIDByNameBatch.
func (q *DBQuerier) FindOIDByNameBatch(batch *pgx.Batch, name string) {
	batch.Queue(findOIDByNameSQL, name)
}

// FindOIDByNameScan implements Querier.FindOIDByNameScan.
func (q *DBQuerier) FindOIDByNameScan(results pgx.BatchResults) (pgtype.OID, error) {
	row := results.QueryRow()
	var item pgtype.OID
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindOIDByNameBatch row: %w", err)
	}
	return item, nil
}

const findOIDNameSQL = `SELECT typname as name
FROM pg_type
WHERE oid = $1;`

// FindOIDName implements Querier.FindOIDName.
func (q *DBQuerier) FindOIDName(ctx context.Context, oID pgtype.OID) (pgtype.Name, error) {
	row := q.conn.QueryRow(ctx, findOIDNameSQL, oID)
	var item pgtype.Name
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query FindOIDName: %w", err)
	}
	return item, nil
}

// FindOIDNameBatch implements Querier.FindOIDNameBatch.
func (q *DBQuerier) FindOIDNameBatch(batch *pgx.Batch, oID pgtype.OID) {
	batch.Queue(findOIDNameSQL, oID)
}

// FindOIDNameScan implements Querier.FindOIDNameScan.
func (q *DBQuerier) FindOIDNameScan(results pgx.BatchResults) (pgtype.Name, error) {
	row := results.QueryRow()
	var item pgtype.Name
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan FindOIDNameBatch row: %w", err)
	}
	return item, nil
}
