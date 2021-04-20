// Code generated by pggen. DO NOT EDIT.

package enums

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"unsafe"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error)
	// FindAllDevicesBatch enqueues a FindAllDevices query into batch to be executed
	// later by the batch.
	FindAllDevicesBatch(batch *pgx.Batch)
	// FindAllDevicesScan scans the result of an executed FindAllDevicesBatch query.
	FindAllDevicesScan(results pgx.BatchResults) ([]FindAllDevicesRow, error)

	InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error)
	// InsertDeviceBatch enqueues a InsertDevice query into batch to be executed
	// later by the batch.
	InsertDeviceBatch(batch *pgx.Batch, mac pgtype.Macaddr, typePg DeviceType)
	// InsertDeviceScan scans the result of an executed InsertDeviceBatch query.
	InsertDeviceScan(results pgx.BatchResults) (pgconn.CommandTag, error)

	// Select an array of all device_type enum values.
	FindOneDeviceArray(ctx context.Context) ([]DeviceType, error)
	// FindOneDeviceArrayBatch enqueues a FindOneDeviceArray query into batch to be executed
	// later by the batch.
	FindOneDeviceArrayBatch(batch *pgx.Batch)
	// FindOneDeviceArrayScan scans the result of an executed FindOneDeviceArrayBatch query.
	FindOneDeviceArrayScan(results pgx.BatchResults) ([]DeviceType, error)

	// Select many rows of device_type enum values.
	FindManyDeviceArray(ctx context.Context) ([][]DeviceType, error)
	// FindManyDeviceArrayBatch enqueues a FindManyDeviceArray query into batch to be executed
	// later by the batch.
	FindManyDeviceArrayBatch(batch *pgx.Batch)
	// FindManyDeviceArrayScan scans the result of an executed FindManyDeviceArrayBatch query.
	FindManyDeviceArrayScan(results pgx.BatchResults) ([][]DeviceType, error)

	// Select many rows of device_type enum values with multiple output columns.
	FindManyDeviceArrayWithNum(ctx context.Context) ([]FindManyDeviceArrayWithNumRow, error)
	// FindManyDeviceArrayWithNumBatch enqueues a FindManyDeviceArrayWithNum query into batch to be executed
	// later by the batch.
	FindManyDeviceArrayWithNumBatch(batch *pgx.Batch)
	// FindManyDeviceArrayWithNumScan scans the result of an executed FindManyDeviceArrayWithNumBatch query.
	FindManyDeviceArrayWithNumScan(results pgx.BatchResults) ([]FindManyDeviceArrayWithNumRow, error)

	// Regression test for https://github.com/jschaf/pggen/issues/23.
	EnumInsideComposite(ctx context.Context) (Device, error)
	// EnumInsideCompositeBatch enqueues a EnumInsideComposite query into batch to be executed
	// later by the batch.
	EnumInsideCompositeBatch(batch *pgx.Batch)
	// EnumInsideCompositeScan scans the result of an executed EnumInsideCompositeBatch query.
	EnumInsideCompositeScan(results pgx.BatchResults) (Device, error)
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
	if _, err := p.Prepare(ctx, findAllDevicesSQL, findAllDevicesSQL); err != nil {
		return fmt.Errorf("prepare query 'FindAllDevices': %w", err)
	}
	if _, err := p.Prepare(ctx, insertDeviceSQL, insertDeviceSQL); err != nil {
		return fmt.Errorf("prepare query 'InsertDevice': %w", err)
	}
	if _, err := p.Prepare(ctx, findOneDeviceArraySQL, findOneDeviceArraySQL); err != nil {
		return fmt.Errorf("prepare query 'FindOneDeviceArray': %w", err)
	}
	if _, err := p.Prepare(ctx, findManyDeviceArraySQL, findManyDeviceArraySQL); err != nil {
		return fmt.Errorf("prepare query 'FindManyDeviceArray': %w", err)
	}
	if _, err := p.Prepare(ctx, findManyDeviceArrayWithNumSQL, findManyDeviceArrayWithNumSQL); err != nil {
		return fmt.Errorf("prepare query 'FindManyDeviceArrayWithNum': %w", err)
	}
	if _, err := p.Prepare(ctx, enumInsideCompositeSQL, enumInsideCompositeSQL); err != nil {
		return fmt.Errorf("prepare query 'EnumInsideComposite': %w", err)
	}
	return nil
}

// Device represents the Postgres composite type "device".
type Device struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// ignoredOID means we don't know or care about the OID for a type. This is okay
// because pgx only uses the OID to encode values and lookup a decoder. We only
// use ignoredOID for decoding and we always specify a concrete decoder for scan
// methods.
const ignoredOID = 0

// DeviceType represents the Postgres enum "device_type".
type DeviceType string

const (
	DeviceTypeUndefined DeviceType = "undefined"
	DeviceTypePhone     DeviceType = "phone"
	DeviceTypeLaptop    DeviceType = "laptop"
	DeviceTypeIpad      DeviceType = "ipad"
	DeviceTypeDesktop   DeviceType = "desktop"
	DeviceTypeIot       DeviceType = "iot"
)

func (d DeviceType) String() string { return string(d) }

var enumDecoderDeviceType = pgtype.NewEnumType("device_type", []string{
	string(DeviceTypeUndefined),
	string(DeviceTypePhone),
	string(DeviceTypeLaptop),
	string(DeviceTypeIpad),
	string(DeviceTypeDesktop),
	string(DeviceTypeIot),
})

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

const findAllDevicesSQL = `SELECT mac, type
FROM device;`

type FindAllDevicesRow struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// FindAllDevices implements Querier.FindAllDevices.
func (q *DBQuerier) FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error) {
	rows, err := q.conn.Query(ctx, findAllDevicesSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindAllDevices: %w", err)
	}
	defer rows.Close()
	items := []FindAllDevicesRow{}
	for rows.Next() {
		var item FindAllDevicesRow
		if err := rows.Scan(&item.Mac, &item.Type); err != nil {
			return nil, fmt.Errorf("scan FindAllDevices row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindAllDevices rows: %w", err)
	}
	return items, err
}

// FindAllDevicesBatch implements Querier.FindAllDevicesBatch.
func (q *DBQuerier) FindAllDevicesBatch(batch *pgx.Batch) {
	batch.Queue(findAllDevicesSQL)
}

// FindAllDevicesScan implements Querier.FindAllDevicesScan.
func (q *DBQuerier) FindAllDevicesScan(results pgx.BatchResults) ([]FindAllDevicesRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindAllDevicesBatch: %w", err)
	}
	defer rows.Close()
	items := []FindAllDevicesRow{}
	for rows.Next() {
		var item FindAllDevicesRow
		if err := rows.Scan(&item.Mac, &item.Type); err != nil {
			return nil, fmt.Errorf("scan FindAllDevicesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindAllDevicesBatch rows: %w", err)
	}
	return items, err
}

const insertDeviceSQL = `INSERT INTO device (mac, type)
VALUES ($1, $2);`

// InsertDevice implements Querier.InsertDevice.
func (q *DBQuerier) InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error) {
	cmdTag, err := q.conn.Exec(ctx, insertDeviceSQL, mac, typePg)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertDevice: %w", err)
	}
	return cmdTag, err
}

// InsertDeviceBatch implements Querier.InsertDeviceBatch.
func (q *DBQuerier) InsertDeviceBatch(batch *pgx.Batch, mac pgtype.Macaddr, typePg DeviceType) {
	batch.Queue(insertDeviceSQL, mac, typePg)
}

// InsertDeviceScan implements Querier.InsertDeviceScan.
func (q *DBQuerier) InsertDeviceScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertDeviceBatch: %w", err)
	}
	return cmdTag, err
}

const findOneDeviceArraySQL = `SELECT enum_range(NULL::device_type) AS device_types;`

// FindOneDeviceArray implements Querier.FindOneDeviceArray.
func (q *DBQuerier) FindOneDeviceArray(ctx context.Context) ([]DeviceType, error) {
	row := q.conn.QueryRow(ctx, findOneDeviceArraySQL)
	item := []DeviceType{}
	deviceTypesArray := &pgtype.EnumArray{}
	if err := row.Scan(deviceTypesArray); err != nil {
		return item, fmt.Errorf("query FindOneDeviceArray: %w", err)
	}
	_ = deviceTypesArray.AssignTo((*[]string)(unsafe.Pointer(&item))) // safe cast; enum array is []string
	return item, nil
}

// FindOneDeviceArrayBatch implements Querier.FindOneDeviceArrayBatch.
func (q *DBQuerier) FindOneDeviceArrayBatch(batch *pgx.Batch) {
	batch.Queue(findOneDeviceArraySQL)
}

// FindOneDeviceArrayScan implements Querier.FindOneDeviceArrayScan.
func (q *DBQuerier) FindOneDeviceArrayScan(results pgx.BatchResults) ([]DeviceType, error) {
	row := results.QueryRow()
	item := []DeviceType{}
	deviceTypesArray := &pgtype.EnumArray{}
	if err := row.Scan(deviceTypesArray); err != nil {
		return item, fmt.Errorf("scan FindOneDeviceArrayBatch row: %w", err)
	}
	_ = deviceTypesArray.AssignTo((*[]string)(unsafe.Pointer(&item))) // safe cast; enum array is []string
	return item, nil
}

const findManyDeviceArraySQL = `SELECT enum_range('ipad'::device_type, 'iot'::device_type) AS device_types
UNION ALL
SELECT enum_range(NULL::device_type) AS device_types;`

// FindManyDeviceArray implements Querier.FindManyDeviceArray.
func (q *DBQuerier) FindManyDeviceArray(ctx context.Context) ([][]DeviceType, error) {
	rows, err := q.conn.Query(ctx, findManyDeviceArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArray: %w", err)
	}
	defer rows.Close()
	items := [][]DeviceType{}
	deviceTypesArray := &pgtype.EnumArray{}
	for rows.Next() {
		var item []DeviceType
		if err := rows.Scan(deviceTypesArray); err != nil {
			return nil, fmt.Errorf("scan FindManyDeviceArray row: %w", err)
		}
		_ = deviceTypesArray.AssignTo((*[]string)(unsafe.Pointer(&item))) // safe cast; enum array is []string
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindManyDeviceArray rows: %w", err)
	}
	return items, err
}

// FindManyDeviceArrayBatch implements Querier.FindManyDeviceArrayBatch.
func (q *DBQuerier) FindManyDeviceArrayBatch(batch *pgx.Batch) {
	batch.Queue(findManyDeviceArraySQL)
}

// FindManyDeviceArrayScan implements Querier.FindManyDeviceArrayScan.
func (q *DBQuerier) FindManyDeviceArrayScan(results pgx.BatchResults) ([][]DeviceType, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArrayBatch: %w", err)
	}
	defer rows.Close()
	items := [][]DeviceType{}
	deviceTypesArray := &pgtype.EnumArray{}
	for rows.Next() {
		var item []DeviceType
		if err := rows.Scan(deviceTypesArray); err != nil {
			return nil, fmt.Errorf("scan FindManyDeviceArrayBatch row: %w", err)
		}
		_ = deviceTypesArray.AssignTo((*[]string)(unsafe.Pointer(&item))) // safe cast; enum array is []string
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindManyDeviceArrayBatch rows: %w", err)
	}
	return items, err
}

const findManyDeviceArrayWithNumSQL = `SELECT 1 AS num, enum_range('ipad'::device_type, 'iot'::device_type) AS device_types
UNION ALL
SELECT 2 as num, enum_range(NULL::device_type) AS device_types;`

type FindManyDeviceArrayWithNumRow struct {
	Num         *int32       `json:"num"`
	DeviceTypes []DeviceType `json:"device_types"`
}

// FindManyDeviceArrayWithNum implements Querier.FindManyDeviceArrayWithNum.
func (q *DBQuerier) FindManyDeviceArrayWithNum(ctx context.Context) ([]FindManyDeviceArrayWithNumRow, error) {
	rows, err := q.conn.Query(ctx, findManyDeviceArrayWithNumSQL)
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArrayWithNum: %w", err)
	}
	defer rows.Close()
	items := []FindManyDeviceArrayWithNumRow{}
	deviceTypesArray := &pgtype.EnumArray{}
	for rows.Next() {
		var item FindManyDeviceArrayWithNumRow
		if err := rows.Scan(&item.Num, deviceTypesArray); err != nil {
			return nil, fmt.Errorf("scan FindManyDeviceArrayWithNum row: %w", err)
		}
		_ = deviceTypesArray.AssignTo((*[]string)(unsafe.Pointer(&item.DeviceTypes))) // safe cast; enum array is []string
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindManyDeviceArrayWithNum rows: %w", err)
	}
	return items, err
}

// FindManyDeviceArrayWithNumBatch implements Querier.FindManyDeviceArrayWithNumBatch.
func (q *DBQuerier) FindManyDeviceArrayWithNumBatch(batch *pgx.Batch) {
	batch.Queue(findManyDeviceArrayWithNumSQL)
}

// FindManyDeviceArrayWithNumScan implements Querier.FindManyDeviceArrayWithNumScan.
func (q *DBQuerier) FindManyDeviceArrayWithNumScan(results pgx.BatchResults) ([]FindManyDeviceArrayWithNumRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindManyDeviceArrayWithNumBatch: %w", err)
	}
	defer rows.Close()
	items := []FindManyDeviceArrayWithNumRow{}
	deviceTypesArray := &pgtype.EnumArray{}
	for rows.Next() {
		var item FindManyDeviceArrayWithNumRow
		if err := rows.Scan(&item.Num, deviceTypesArray); err != nil {
			return nil, fmt.Errorf("scan FindManyDeviceArrayWithNumBatch row: %w", err)
		}
		_ = deviceTypesArray.AssignTo((*[]string)(unsafe.Pointer(&item.DeviceTypes))) // safe cast; enum array is []string
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindManyDeviceArrayWithNumBatch rows: %w", err)
	}
	return items, err
}

const enumInsideCompositeSQL = `SELECT ROW('08:00:2b:01:02:03'::macaddr, 'phone'::device_type) ::device;`

// EnumInsideComposite implements Querier.EnumInsideComposite.
func (q *DBQuerier) EnumInsideComposite(ctx context.Context) (Device, error) {
	row := q.conn.QueryRow(ctx, enumInsideCompositeSQL)
	var item Device
	rowRow := newCompositeType(
		"device",
		[]string{"mac", "type"},
		&pgtype.Macaddr{},
		enumDecoderDeviceType,
	)
	if err := row.Scan(rowRow); err != nil {
		return item, fmt.Errorf("query EnumInsideComposite: %w", err)
	}
	if err := rowRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign EnumInsideComposite row: %w", err)
	}
	return item, nil
}

// EnumInsideCompositeBatch implements Querier.EnumInsideCompositeBatch.
func (q *DBQuerier) EnumInsideCompositeBatch(batch *pgx.Batch) {
	batch.Queue(enumInsideCompositeSQL)
}

// EnumInsideCompositeScan implements Querier.EnumInsideCompositeScan.
func (q *DBQuerier) EnumInsideCompositeScan(results pgx.BatchResults) (Device, error) {
	row := results.QueryRow()
	var item Device
	rowRow := newCompositeType(
		"device",
		[]string{"mac", "type"},
		&pgtype.Macaddr{},
		enumDecoderDeviceType,
	)
	if err := row.Scan(rowRow); err != nil {
		return item, fmt.Errorf("scan EnumInsideCompositeBatch row: %w", err)
	}
	if err := rowRow.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign EnumInsideComposite row: %w", err)
	}
	return item, nil
}
