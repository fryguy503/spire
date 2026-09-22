package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/EQEmuTools/spire/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type playerOperationsSavedExpedition struct {
	DZID             uint32 `json:"dz_id"`
	UUID             string `json:"uuid"`
	InstanceID       uint32 `json:"instance_id"`
	Name             string `json:"name"`
	Leader           string `json:"leader"`
	ZoneID           uint32 `json:"zone_id"`
	ZoneVersion      uint32 `json:"zone_version"`
	SeasonID         int    `json:"season_id"`
	Members          uint32 `json:"members"`
	MaxMembers       uint32 `json:"max_members"`
	Expires          int64  `json:"expires"`
	Locked           bool   `json:"locked"`
	HasLockout       bool   `json:"-"`
	DatabaseEligible bool   `json:"database_eligible"`
	Status           string `json:"status"`
}

type playerOperationsSavedExpeditions struct {
	Available         bool                              `json:"available"`
	Message           string                            `json:"message"`
	ServerTime        int64                             `json:"server_time"`
	ConfiguredEnabled bool                              `json:"configured_enabled"`
	CharacterSeason   int                               `json:"character_season"`
	Entries           []playerOperationsSavedExpedition `json:"entries"`
}

type playerOperationsDZCharacter struct {
	ID           int
	ZoneInstance int
	CurrentDZ    uint32
}

func (p *PlayerOperationsController) getSavedExpeditions(c echo.Context) error {
	id, err := playerOperationsPositiveParam(c, "id", "Character")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	var result playerOperationsSavedExpeditions
	if p == nil || p.db == nil {
		return c.JSON(http.StatusServiceUnavailable, echo.Map{"error": "Player operations database is unavailable"})
	}
	db := p.db.Get(models.CharacterDatum{}, c)
	if db == nil {
		return c.JSON(http.StatusServiceUnavailable, echo.Map{"error": "Player operations database is unavailable"})
	}
	// The authenticated character connection owns all resume, rule and lockout
	// tables. A read-only snapshot avoids mixing rosters from different moments.
	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		var loadErr error
		result, loadErr = loadPlayerOperationsSavedExpeditions(tx, id)
		return loadErr
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return playerOperationsLoadError(c, "Character", err)
	}
	return c.JSON(http.StatusOK, result)
}

func loadPlayerOperationsSavedExpeditions(db *gorm.DB, characterID int) (playerOperationsSavedExpeditions, error) {
	result := playerOperationsSavedExpeditions{Entries: make([]playerOperationsSavedExpedition, 0)}
	var character playerOperationsDZCharacter
	if err := db.Table("character_data").Select("id, zone_instance").Where("id = ?", characterID).Take(&character).Error; err != nil {
		return result, err
	}
	var tables []struct {
		Name   string
		Engine string
	}
	if err := db.Raw(`SELECT table_name AS name, engine FROM information_schema.tables
		WHERE table_schema = DATABASE() AND table_name IN (
		'dynamic_zone_resume', 'dynamic_zone_resume_history', 'character_dynamic_zone_resume',
		'dynamic_zones', 'dynamic_zone_members', 'instance_list', 'instance_list_player',
		'dynamic_zone_lockouts', 'character_expedition_lockouts')`).Scan(&tables).Error; err != nil {
		return result, err
	}
	installed := make(map[string]bool)
	for _, table := range tables {
		installed[table.Name] = true
	}
	if !installed["dynamic_zone_resume"] && !installed["dynamic_zone_resume_history"] && !installed["character_dynamic_zone_resume"] {
		result.Message = "Saved expedition history is not installed on this database."
		return result, nil
	}
	if len(tables) != 9 {
		result.Message = "Saved expedition history is incomplete. Install the matching DZManager migrations."
		return result, nil
	}
	for _, table := range tables {
		if !strings.EqualFold(table.Engine, "InnoDB") {
			result.Message = "DZManager requires InnoDB storage for its expedition and history tables."
			return result, nil
		}
	}
	columnsReady, err := playerOperationsDZColumnsReady(db)
	if err != nil {
		return result, err
	}
	if !columnsReady {
		result.Message = "Saved expedition history is incomplete. Install the matching DZManager migrations."
		return result, nil
	}
	if err := db.Raw("SELECT UNIX_TIMESTAMP()").Scan(&result.ServerTime).Error; err != nil {
		return result, err
	}
	activeSeason, err := loadPlayerOperationsDZRules(db, &result)
	if err != nil {
		return result, err
	}
	if activeSeason > 0 {
		var bucket string
		if err := db.Table("data_buckets").Select("value").Where("character_id = ? AND `key` = ? AND (expires = 0 OR expires > ?)",
			characterID, "SeasonalCharacter", result.ServerTime).Order("id DESC").Limit(1).Scan(&bucket).Error; err != nil {
			return result, err
		}
		if season, parseErr := strconv.Atoi(bucket); parseErr == nil && season == activeSeason {
			result.CharacterSeason = season
		}
	}
	if err := db.Table("dynamic_zone_members m").Select("d.id").
		Joins("JOIN dynamic_zones d ON d.id = m.dynamic_zone_id").
		Where("m.character_id = ? AND d.type = 1", characterID).Order("d.id").Limit(1).Scan(&character.CurrentDZ).Error; err != nil {
		return result, err
	}
	if err := db.Raw(playerOperationsSavedExpeditionsSQL, characterID, characterID, result.ServerTime).Scan(&result.Entries).Error; err != nil {
		return result, err
	}
	for index := range result.Entries {
		entry := &result.Entries[index]
		entry.Status = playerOperationsDZStatus(*entry, character, result)
		entry.DatabaseEligible = entry.Status == "server_check"
	}
	result.Available = true
	result.Message = "Database snapshot only. DZManager checks live rules, pending invitations and current location again when the player rejoins."
	return result, nil
}

func playerOperationsDZColumnsReady(db *gorm.DB) (bool, error) {
	required := map[string][]string{
		"dynamic_zone_resume":           {"dz_id", "uuid", "instance_id", "season_id", "revoked"},
		"dynamic_zone_resume_history":   {"dz_id", "uuid", "character_id", "can_resume"},
		"character_dynamic_zone_resume": {"character_id", "zone_id", "zone_version", "expedition_name", "season_id", "dz_id", "uuid"},
		"dynamic_zones":                 {"id", "uuid", "instance_id", "name", "leader_id", "max_players", "is_locked", "type"},
		"dynamic_zone_members":          {"dynamic_zone_id", "character_id"},
		"instance_list":                 {"id", "zone", "version", "start_time", "duration", "never_expires"},
		"dynamic_zone_lockouts":         {"id", "dynamic_zone_id", "event_name"},
		"character_expedition_lockouts": {"character_id", "expedition_name", "event_name", "expire_time", "from_expedition_uuid"},
	}
	names := make([]string, 0, len(required))
	for name := range required {
		names = append(names, name)
	}
	var columns []struct {
		TableName  string
		ColumnName string
	}
	if err := db.Raw(`SELECT table_name, column_name FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name IN ?`, names).Scan(&columns).Error; err != nil {
		return false, err
	}
	present := make(map[string]bool)
	for _, column := range columns {
		present[column.TableName+"."+column.ColumnName] = true
	}
	for table, names := range required {
		for _, column := range names {
			if !present[table+"."+column] {
				return false, nil
			}
		}
	}
	return true, nil
}

// Mirrors world/DynamicZoneManager::HandleResumeRequest. History survives a
// kick, but only the exact latest DZ UUID for this content/season can return.
// An expired or deleted newer copy must never reopen an older saved copy.
const playerOperationsSavedExpeditionsSQL = `SELECT
	 d.id AS dz_id, d.uuid, d.instance_id, d.name, COALESCE(c.name, '') AS leader,
	 i.zone AS zone_id, i.version AS zone_version, r.season_id,
	 (i.start_time + i.duration) AS expires, d.is_locked AS locked, d.max_players AS max_members,
	 (SELECT COUNT(*) FROM dynamic_zone_members m WHERE m.dynamic_zone_id = d.id) AS members,
	 EXISTS(SELECT 1 FROM character_expedition_lockouts cl
		LEFT JOIN dynamic_zone_lockouts dl ON dl.dynamic_zone_id = d.id AND dl.event_name = cl.event_name
		WHERE cl.character_id = ? AND cl.expedition_name = d.name AND cl.expire_time > NOW()
		AND ((cl.event_name = 'Replay Timer' AND COALESCE(cl.from_expedition_uuid, '') <> d.uuid)
		OR (COALESCE(cl.event_name, '') <> 'Replay Timer' AND dl.id IS NULL))) AS has_lockout
	FROM dynamic_zone_resume r
	JOIN dynamic_zone_resume_history h ON h.dz_id = r.dz_id AND h.uuid = r.uuid
	JOIN dynamic_zones d ON d.id = r.dz_id AND d.uuid = r.uuid AND d.instance_id = r.instance_id
	JOIN instance_list i ON i.id = r.instance_id
	JOIN character_dynamic_zone_resume l ON l.character_id = h.character_id AND l.dz_id = r.dz_id AND l.uuid = r.uuid
		AND l.season_id = r.season_id AND l.zone_id = i.zone AND l.zone_version = i.version AND l.expedition_name = d.name
	LEFT JOIN character_data c ON c.id = d.leader_id
	WHERE h.character_id = ? AND h.can_resume = 1 AND r.revoked = 0 AND d.type = 1
		AND i.never_expires = 0 AND (i.start_time + i.duration) > ?
	ORDER BY (i.start_time + i.duration), d.id LIMIT 16`

func loadPlayerOperationsDZRules(db *gorm.DB, result *playerOperationsSavedExpeditions) (int, error) {
	var configured string
	if err := db.Table("variables").Select("value").Where("varname = ?", "RuleSet").Limit(1).Scan(&configured).Error; err != nil {
		return 0, err
	}
	if configured == "" {
		configured = "default"
	}
	var sets []struct {
		RulesetID uint `gorm:"column:ruleset_id"`
		Name      string
	}
	if err := db.Table("rule_sets").Select("ruleset_id, name").Where("name IN ?", []string{"default", configured}).Scan(&sets).Error; err != nil {
		return 0, err
	}
	ids := map[string]uint{}
	for _, set := range sets {
		ids[set.Name] = set.RulesetID
	}
	if _, ok := ids[configured]; !ok {
		return 0, errors.New("configured world ruleset does not exist")
	}
	if _, ok := ids["default"]; !ok {
		return 0, errors.New("default world ruleset does not exist")
	}
	var rules []struct {
		RulesetID uint `gorm:"column:ruleset_id"`
		RuleName  string
		RuleValue string
	}
	if err := db.Table("rule_values").Select("ruleset_id, rule_name, rule_value").
		Where("ruleset_id IN ? AND rule_name IN ?", []uint{ids["default"], ids[configured]},
			[]string{"Expedition:EnableDzResume", "Seasons:EnableSeasonalCharacters"}).Scan(&rules).Error; err != nil {
		return 0, err
	}
	activeSeason := 0
	// World loads default first, then overlays the selected ruleset. Missing
	// rules keep their compiled defaults: resume disabled and season zero.
	for _, id := range []uint{ids["default"], ids[configured]} {
		for _, rule := range rules {
			if rule.RulesetID != id {
				continue
			}
			if rule.RuleName == "Expedition:EnableDzResume" {
				result.ConfiguredEnabled = playerOperationsDZRuleBool(rule.RuleValue)
			} else {
				activeSeason, _ = strconv.Atoi(strings.TrimSpace(rule.RuleValue))
			}
		}
	}
	return activeSeason, nil
}

// Match Source's Strings::ToBool, including its case-sensitive substring
// checks, instead of interpreting a stored rule differently from world.
func playerOperationsDZRuleBool(value string) bool {
	for _, token := range []string{"true", "y", "on", "enable"} {
		if strings.Contains(value, token) {
			return true
		}
	}
	numeric, err := strconv.ParseInt(value, 10, 32)
	return err == nil && numeric != 0
}

var playerOperationsDZUUID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func playerOperationsDZStatus(entry playerOperationsSavedExpedition, character playerOperationsDZCharacter, state playerOperationsSavedExpeditions) string {
	switch {
	case entry.DZID == 0 || entry.ZoneID == 0 || entry.InstanceID == 0 || entry.InstanceID > 65535 || !playerOperationsDZUUID.MatchString(entry.UUID):
		return "invalid"
	case !state.ConfiguredEnabled:
		return "disabled"
	case entry.Expires <= state.ServerTime:
		return "expired"
	case character.CurrentDZ == entry.DZID:
		return "already_member"
	case character.CurrentDZ != 0:
		return "in_expedition"
	case character.ZoneInstance != 0:
		return "in_instance"
	case entry.Locked:
		return "locked"
	case entry.Members >= entry.MaxMembers:
		return "full"
	case state.CharacterSeason != entry.SeasonID:
		return "season"
	case entry.HasLockout:
		return "lockout"
	default:
		return "server_check"
	}
}
