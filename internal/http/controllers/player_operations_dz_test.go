package controllers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestPlayerOperationsDZOptionalSchema(t *testing.T) {
	for _, scenario := range []string{"missing", "partial tables", "partial columns"} {
		t.Run(scenario, func(t *testing.T) {
			fixture := &playerOperationsDZSchemaFixture{scenario: scenario}
			sqlDB := sql.OpenDB(fixture)
			t.Cleanup(func() { _ = sqlDB.Close() })
			db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true})
			if err != nil {
				t.Fatal(err)
			}
			result, err := loadPlayerOperationsSavedExpeditions(db, 37)
			if err != nil {
				t.Fatal(err)
			}
			if result.Available || len(result.Entries) != 0 {
				t.Fatalf("incomplete schema reported available: %#v", result)
			}
			if scenario == "missing" && !strings.Contains(result.Message, "not installed") {
				t.Fatal(result.Message)
			}
			if scenario != "missing" && !strings.Contains(result.Message, "incomplete") {
				t.Fatal(result.Message)
			}
		})
	}
}

// A strict read-only fixture fails on any attempted feature-table query before
// optional schema readiness. This prevents older connected servers breaking.
type playerOperationsDZSchemaFixture struct{ scenario string }

func (f *playerOperationsDZSchemaFixture) Connect(context.Context) (driver.Conn, error) {
	return f, nil
}
func (f *playerOperationsDZSchemaFixture) Driver() driver.Driver            { return f }
func (f *playerOperationsDZSchemaFixture) Open(string) (driver.Conn, error) { return f, nil }
func (*playerOperationsDZSchemaFixture) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepared statement")
}
func (*playerOperationsDZSchemaFixture) Begin() (driver.Tx, error) {
	return nil, errors.New("unexpected transaction")
}
func (*playerOperationsDZSchemaFixture) Close() error { return nil }
func (f *playerOperationsDZSchemaFixture) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(query, "FROM `character_data`") {
		return &playerOperationsDZFixtureRows{columns: []string{"id", "zone_instance"}, values: [][]driver.Value{{int64(37), int64(0)}}}, nil
	}
	if strings.Contains(query, "information_schema.tables") {
		rows := &playerOperationsDZFixtureRows{columns: []string{"name", "engine"}}
		if f.scenario != "missing" {
			for _, table := range []string{"dynamic_zone_resume", "dynamic_zone_resume_history", "dynamic_zones", "dynamic_zone_members", "instance_list", "instance_list_player", "dynamic_zone_lockouts", "character_expedition_lockouts"} {
				rows.values = append(rows.values, []driver.Value{table, "InnoDB"})
			}
			if f.scenario == "partial columns" {
				rows.values = append(rows.values, []driver.Value{"character_dynamic_zone_resume", "InnoDB"})
			}
		}
		return rows, nil
	}
	if strings.Contains(query, "information_schema.columns") && f.scenario == "partial columns" {
		return &playerOperationsDZFixtureRows{columns: []string{"table_name", "column_name"}, values: [][]driver.Value{{"dynamic_zone_resume", "dz_id"}}}, nil
	}
	return nil, errors.New("unexpected optional data query: " + query)
}

type playerOperationsDZFixtureRows struct {
	columns []string
	values  [][]driver.Value
}

func (r *playerOperationsDZFixtureRows) Columns() []string { return r.columns }
func (*playerOperationsDZFixtureRows) Close() error        { return nil }
func (r *playerOperationsDZFixtureRows) Next(dest []driver.Value) error {
	if len(r.values) == 0 {
		return io.EOF
	}
	copy(dest, r.values[0])
	r.values = r.values[1:]
	return nil
}

func TestPlayerOperationsDZStatus(t *testing.T) {
	entry := playerOperationsSavedExpedition{DZID: 42, ZoneID: 202, InstanceID: 101, UUID: "01234567-0123-4567-89ab-0123456789ab", MaxMembers: 6, Expires: 200}
	state := playerOperationsSavedExpeditions{ConfiguredEnabled: true, ServerTime: 100}
	for _, test := range []struct {
		name   string
		change func(*playerOperationsSavedExpedition, *playerOperationsDZCharacter, *playerOperationsSavedExpeditions)
		want   string
	}{
		{"eligible still needs world", func(*playerOperationsSavedExpedition, *playerOperationsDZCharacter, *playerOperationsSavedExpeditions) {
		}, "server_check"},
		{"corrupt uuid", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			e.UUID = "wrong"
		}, "invalid"},
		{"invalid instance", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			e.InstanceID = 65536
		}, "invalid"},
		{"feature disabled", func(_ *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, s *playerOperationsSavedExpeditions) {
			s.ConfiguredEnabled = false
		}, "disabled"},
		{"exact expiry", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			e.Expires = 100
		}, "expired"},
		{"already member precedes locked", func(e *playerOperationsSavedExpedition, c *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			c.CurrentDZ = 42
			e.Locked = true
		}, "already_member"},
		{"another expedition", func(_ *playerOperationsSavedExpedition, c *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			c.CurrentDZ = 43
		}, "in_expedition"},
		{"saved instance location", func(_ *playerOperationsSavedExpedition, c *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			c.ZoneInstance = 17
		}, "in_instance"},
		{"locked", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			e.Locked = true
		}, "locked"},
		{"full", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			e.Members = 6
		}, "full"},
		{"over capacity", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			e.Members = 7
		}, "full"},
		{"season mismatch", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			e.SeasonID = 4
		}, "season"},
		{"matching season", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, s *playerOperationsSavedExpeditions) {
			e.SeasonID = 4
			s.CharacterSeason = 4
		}, "server_check"},
		{"lockout conflict", func(e *playerOperationsSavedExpedition, _ *playerOperationsDZCharacter, _ *playerOperationsSavedExpeditions) {
			e.HasLockout = true
		}, "lockout"},
	} {
		t.Run(test.name, func(t *testing.T) {
			e, c, s := entry, playerOperationsDZCharacter{}, state
			test.change(&e, &c, &s)
			if got := playerOperationsDZStatus(e, c, s); got != test.want {
				t.Fatalf("status = %q, want %q", got, test.want)
			}
		})
	}
}

func TestPlayerOperationsSavedExpeditionsQueryPreservesSourceEligibility(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: "user:pass@tcp(localhost:3306)/eqemu", SkipInitializeWithVersion: true,
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	query := db.Raw(playerOperationsSavedExpeditionsSQL, 37, 37, int64(100))
	if len(query.Statement.Vars) != 3 || query.Statement.Vars[0] != 37 || query.Statement.Vars[1] != 37 {
		t.Fatalf("character constraints not bound: %#v", query.Statement.Vars)
	}
	sql := query.Statement.SQL.String()
	for _, fragment := range []string{
		"h.uuid = r.uuid", "d.uuid = r.uuid", "d.instance_id = r.instance_id",
		"l.character_id = h.character_id", "l.dz_id = r.dz_id", "l.uuid = r.uuid",
		"l.season_id = r.season_id", "l.zone_id = i.zone", "l.zone_version = i.version", "l.expedition_name = d.name",
		"h.character_id = ?", "h.can_resume = 1", "r.revoked = 0", "d.type = 1",
		"i.never_expires = 0", "(i.start_time + i.duration) > ?", "LIMIT 16",
		"cl.expire_time > NOW()", "COALESCE(cl.from_expedition_uuid, '') <> d.uuid", "dl.id IS NULL",
	} {
		if !strings.Contains(sql, fragment) {
			t.Errorf("missing required eligibility constraint %q", fragment)
		}
	}
	// Historical access must not require current roster membership: kicked
	// members can still return, and empty retained DZs remain discoverable.
	if strings.Contains(sql, "JOIN dynamic_zone_members") {
		t.Fatal("candidate list requires active roster membership")
	}
}

func TestPlayerOperationsDZRuleBoolMatchesSource(t *testing.T) {
	for _, test := range []struct {
		value string
		want  bool
	}{
		{"true", true}, {"false", false}, {"1", true}, {"0", false}, {"2", true},
		{"yes", true}, {"enabled", true}, {"on", true}, {"TRUE", false}, {"", false},
	} {
		if got := playerOperationsDZRuleBool(test.value); got != test.want {
			t.Errorf("%q = %v, want %v", test.value, got, test.want)
		}
	}
}
