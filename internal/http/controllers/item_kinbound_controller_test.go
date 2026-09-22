package controllers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EQEmuTools/spire/internal/auditlog"
	"github.com/EQEmuTools/spire/internal/database"
	"github.com/labstack/echo/v4"
	gocache "github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestItemKinboundPolicyMatchesSource(t *testing.T) {
	for _, tc := range []struct {
		name    string
		state   itemKinboundState
		blocked bool
	}{
		{"explicit no trade", itemKinboundState{Nodrop: 0}, false},
		{"ordinary tradeable", itemKinboundState{Nodrop: 1}, true},
		{"tradeable 255", itemKinboundState{Nodrop: 255}, true},
		{"no transfer", itemKinboundState{NoTransfer: 1}, true},
		{"known epic", itemKinboundState{Epic: true}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := itemKinboundBlockedReason(tc.state) != ""; got != tc.blocked {
				t.Fatalf("blocked = %v, want %v", got, tc.blocked)
			}
		})
	}
}

func TestItemKinboundHTTP(t *testing.T) {
	for _, tc := range []struct {
		name, method, payload                          string
		schema                                         int64
		current, nodrop, noTransfer                    int64
		epic                                           *int64
		auditFail                                      bool
		wantStatus, writes, audits, commits, rollbacks int
	}{
		{name: "old schema remains readable", method: "GET", schema: 0, wantStatus: 200},
		{name: "old schema cannot write", method: "PATCH", schema: 0, payload: `{"kinbound":1,"expected_kinbound":0}`, wantStatus: 409},
		{name: "read without catalog row", method: "GET", schema: 4, wantStatus: 200},
		{name: "explicit Always", method: "PATCH", schema: 4, payload: `{"kinbound":1,"expected_kinbound":0}`, wantStatus: 200, writes: 1, audits: 1, commits: 1},
		{name: "clear Never to unlisted", method: "PATCH", schema: 4, current: 2, payload: `{"kinbound":0,"expected_kinbound":2}`, wantStatus: 200, writes: 1, audits: 1, commits: 1},
		{name: "Never allowed on tradeable", method: "PATCH", schema: 4, nodrop: 255, payload: `{"kinbound":2,"expected_kinbound":0}`, wantStatus: 200, writes: 1, audits: 1, commits: 1},
		{name: "tradeable blocked", method: "PATCH", schema: 4, nodrop: 255, payload: `{"kinbound":1,"expected_kinbound":0}`, wantStatus: 400, rollbacks: 1},
		{name: "no transfer blocked", method: "PATCH", schema: 4, noTransfer: 1, payload: `{"kinbound":1,"expected_kinbound":0}`, wantStatus: 400, rollbacks: 1},
		{name: "known epic blocked", method: "PATCH", schema: 4, epic: kinboundInt64(1), payload: `{"kinbound":1,"expected_kinbound":0}`, wantStatus: 400, rollbacks: 1},
		{name: "concurrent Never preserved", method: "PATCH", schema: 4, current: 2, payload: `{"kinbound":1,"expected_kinbound":0}`, wantStatus: 409, rollbacks: 1},
		{name: "missing expected rejected", method: "PATCH", schema: 4, payload: `{"kinbound":1}`, wantStatus: 400},
		{name: "missing mode rejected", method: "PATCH", schema: 4, payload: `{"expected_kinbound":0}`, wantStatus: 400},
		{name: "unknown mode rejected", method: "PATCH", schema: 4, payload: `{"kinbound":3,"expected_kinbound":0}`, wantStatus: 400},
		{name: "fraction rejected", method: "PATCH", schema: 4, payload: `{"kinbound":1.5,"expected_kinbound":0}`, wantStatus: 400},
		{name: "audit failure prevents item write", method: "PATCH", schema: 4, auditFail: true, payload: `{"kinbound":1,"expected_kinbound":0}`, wantStatus: 500, audits: 1, rollbacks: 1},
		{name: "unchanged is idempotent", method: "PATCH", schema: 4, current: 1, payload: `{"kinbound":1,"expected_kinbound":1}`, wantStatus: 200, commits: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			record := &kinboundSQLRecord{schema: tc.schema, current: tc.current, nodrop: tc.nodrop, noTransfer: tc.noTransfer, epic: tc.epic, auditFail: tc.auditFail}
			pool := sql.OpenDB(kinboundSQLConnector{record})
			t.Cleanup(func() { _ = pool.Close() })
			db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true, SkipDefaultTransaction: true})
			if err != nil {
				t.Fatal(err)
			}
			cache := gocache.New(gocache.NoExpiration, 0)
			resolver := database.NewResolver(database.NewConnections(db, db, nil), nil, nil, cache)
			controller := NewItemKinboundController(resolver, auditlog.NewUserEvent(resolver, cache))
			e := echo.New()
			for _, route := range controller.Routes() {
				e.Add(route.Method(), "/"+route.Route(), route.Handler())
			}
			req := httptest.NewRequest(tc.method, "/item-kinbound/100", strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			response := httptest.NewRecorder()
			e.ServeHTTP(response, req)
			if response.Code != tc.wantStatus {
				t.Fatalf("status %d, want %d: %s", response.Code, tc.wantStatus, response.Body.String())
			}
			if record.writes != tc.writes || record.audits != tc.audits || record.commits != tc.commits || record.rollbacks != tc.rollbacks {
				t.Fatalf("writes/audits/commits/rollbacks = %d/%d/%d/%d, want %d/%d/%d/%d", record.writes, record.audits, record.commits, record.rollbacks, tc.writes, tc.audits, tc.commits, tc.rollbacks)
			}
			if response.Code == 200 {
				var state itemKinboundState
				if err := json.Unmarshal(response.Body.Bytes(), &state); err != nil {
					t.Fatal(err)
				}
				if state.Supported != (tc.schema == 4) {
					t.Fatalf("unexpected schema support: %+v", state)
				}
			}
			if tc.writes > 0 && !record.locked {
				t.Fatal("item mutation did not lock the latest saved item")
			}
		})
	}
}

func kinboundInt64(n int64) *int64 { return &n }

type kinboundSQLRecord struct {
	schema, current, nodrop, noTransfer int64
	epic                                *int64
	auditFail, locked                   bool
	writes, audits, commits, rollbacks  int
}
type kinboundSQLConnector struct{ record *kinboundSQLRecord }

func (c kinboundSQLConnector) Connect(context.Context) (driver.Conn, error) {
	return &kinboundSQLConn{c.record}, nil
}
func (c kinboundSQLConnector) Driver() driver.Driver { return kinboundSQLDriver{c.record} }

type kinboundSQLDriver struct{ record *kinboundSQLRecord }

func (d kinboundSQLDriver) Open(string) (driver.Conn, error) { return &kinboundSQLConn{d.record}, nil }

type kinboundSQLConn struct{ record *kinboundSQLRecord }

func (*kinboundSQLConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepared statement")
}
func (*kinboundSQLConn) Close() error                { return nil }
func (c *kinboundSQLConn) Begin() (driver.Tx, error) { return c, nil }
func (c *kinboundSQLConn) Commit() error             { c.record.commits++; return nil }
func (c *kinboundSQLConn) Rollback() error           { c.record.rollbacks++; return nil }
func (c *kinboundSQLConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	switch {
	case strings.Contains(query, "information_schema.columns"):
		return &kinboundSQLRows{columns: []string{"COUNT(*)"}, values: [][]driver.Value{{c.record.schema}}}, nil
	case strings.Contains(query, "FROM `items`"):
		c.record.locked = strings.Contains(query, "FOR UPDATE")
		return &kinboundSQLRows{columns: []string{"id", "kinbound", "nodrop", "notransfer"}, values: [][]driver.Value{{int64(100), c.record.current, c.record.nodrop, c.record.noTransfer}}}, nil
	case strings.Contains(query, "FROM `item_kinbound_policy`"):
		rows := &kinboundSQLRows{columns: []string{"epic"}}
		if c.record.epic != nil {
			rows.values = [][]driver.Value{{*c.record.epic}}
		}
		return rows, nil
	}
	return nil, errors.New("unexpected query: " + query)
}
func (c *kinboundSQLConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if strings.HasPrefix(query, "UPDATE `items`") {
		if len(args) != 5 || args[1].Value != int64(100) || args[2].Value != c.record.current || args[3].Value != c.record.nodrop || args[4].Value != c.record.noTransfer {
			return nil, errors.New("mutation is not scoped to the loaded item and mode")
		}
		c.record.writes++
		return driver.RowsAffected(1), nil
	}
	if strings.HasPrefix(query, "INSERT INTO `spire_user_event_log`") {
		c.record.audits++
		if c.record.auditFail {
			return nil, errors.New("audit unavailable")
		}
		return kinboundSQLResult{}, nil
	}
	return nil, errors.New("unexpected mutation: " + query)
}

type kinboundSQLResult struct{}

func (kinboundSQLResult) LastInsertId() (int64, error) { return 1, nil }
func (kinboundSQLResult) RowsAffected() (int64, error) { return 1, nil }

type kinboundSQLRows struct {
	columns []string
	values  [][]driver.Value
}

func (r *kinboundSQLRows) Columns() []string { return r.columns }
func (*kinboundSQLRows) Close() error        { return nil }
func (r *kinboundSQLRows) Next(dest []driver.Value) error {
	if len(r.values) == 0 {
		return io.EOF
	}
	copy(dest, r.values[0])
	r.values = r.values[1:]
	return nil
}
