package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const DefaultSQLitePath = "magicadventure.db"

type Store struct {
	path string
}

type sqliteSaveRecord struct {
	ID                   int    `json:"id"`
	PlayerName           string `json:"player_name"`
	CurrentRoomID        string `json:"current_room_id"`
	Health               int    `json:"health"`
	MaxHealth            int    `json:"max_health"`
	IsGameOver           int    `json:"is_game_over"`
	Gold                 int    `json:"gold"`
	MonsterHealth        int    `json:"monster_health"`
	IsTrollDead          int    `json:"is_troll_dead"`
	IsBeastDead          int    `json:"is_beast_dead"`
	IsFrostGiantDead     int    `json:"is_frost_giant_dead"`
	SkillPoints          int    `json:"skill_points"`
	Level                int    `json:"level"`
	Experience           int    `json:"experience"`
	RequiredXP           int    `json:"required_xp"`
	BaseDamage           int    `json:"base_damage"`
	HasGardenKit         int    `json:"has_garden_kit"`
	HasMeat              int    `json:"has_meat"`
	HasSteelSword        int    `json:"has_steel_sword"`
	HasWater             int    `json:"has_water"`
	HasFurCoat           int    `json:"has_fur_coat"`
	HasSunAmulet         int    `json:"has_sun_amulet"`
	HasIceCrystal        int    `json:"has_ice_crystal"`
	HasForsakenBlade     int    `json:"has_forsaken_blade"`
	HasMoonPearl         int    `json:"has_moon_pearl"`
	HasSageBlessing      int    `json:"has_sage_blessing"`
	HasRuinsToken        int    `json:"has_ruins_token"`
	HasMoonCharm         int    `json:"has_moon_charm"`
	HasStarMap           int    `json:"has_star_map"`
	HasTriptychBlessing  int    `json:"has_triptych_blessing"`
	Is1x1x1x1Dead        int    `json:"is_1x1x1x1_dead"`
	HasSeenTutorial      int    `json:"has_seen_tutorial"`
	ConsecutiveWrongAnswers int `json:"consecutive_wrong_answers"`
	MillionairePoints       int `json:"millionaire_points"`
	LifelineAudienceUsed    int `json:"lifeline_audience_used"`
	LifelineWiseManUsed     int `json:"lifeline_wiseman_used"`
	MillionaireStreakID     string `json:"millionaire_streak_id"`
	Language             string `json:"language"`
	PlayerCombatStatus   string `json:"player_combat_status"`
	PlayerCombatTurns    int    `json:"player_combat_turns"`
	MonsterCombatStatus  string `json:"monster_combat_status"`
	MonsterCombatTurns   int    `json:"monster_combat_turns"`
	OnlineID             string `json:"online_id"`
	OnlineRelics         string `json:"online_relics"`
	LastNotificationSeen int64  `json:"last_notification_seen"`
	LastActionMessage    string `json:"last_action_message"`
	LastSeen             int64  `json:"last_seen"`
}

func InitDB() (*Store, error) {
	path := os.Getenv("MAGIC_ADVENTURE_DB_PATH")
	if path == "" {
		path = DefaultSQLitePath
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	_ = file.Close()

	store := &Store{path: absPath}
	if err := store.ensureSchema(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) ensureSchema() error {
	schema := `
CREATE TABLE IF NOT EXISTS saves (
	id INTEGER PRIMARY KEY,
	player_name TEXT NOT NULL,
	current_room_id TEXT NOT NULL,
	health INTEGER NOT NULL,
	max_health INTEGER NOT NULL,
	is_game_over INTEGER NOT NULL,
	gold INTEGER NOT NULL,
	monster_health INTEGER NOT NULL,
	is_troll_dead INTEGER NOT NULL,
	is_beast_dead INTEGER NOT NULL,
	is_frost_giant_dead INTEGER NOT NULL,
	skill_points INTEGER NOT NULL,
	level INTEGER NOT NULL,
	experience INTEGER NOT NULL,
	required_xp INTEGER NOT NULL,
	base_damage INTEGER NOT NULL,
	has_garden_kit INTEGER NOT NULL,
	has_meat INTEGER NOT NULL,
	has_steel_sword INTEGER NOT NULL,
	has_water INTEGER NOT NULL,
	has_fur_coat INTEGER NOT NULL,
	has_sun_amulet INTEGER NOT NULL,
	has_ice_crystal INTEGER NOT NULL,
	has_forsaken_blade INTEGER NOT NULL,
	has_moon_pearl INTEGER NOT NULL DEFAULT 0,
	has_sage_blessing INTEGER NOT NULL DEFAULT 0,
	has_ruins_token INTEGER NOT NULL DEFAULT 0,
	has_moon_charm INTEGER NOT NULL DEFAULT 0,
	has_star_map INTEGER NOT NULL DEFAULT 0,
	has_triptych_blessing INTEGER NOT NULL DEFAULT 0,
	is_1x1x1x1_dead INTEGER NOT NULL,
	has_seen_tutorial INTEGER NOT NULL,
	consecutive_wrong_answers INTEGER NOT NULL DEFAULT 0,
	millionaire_points INTEGER NOT NULL DEFAULT 0,
	lifeline_audience_used INTEGER NOT NULL DEFAULT 0,
	lifeline_wiseman_used INTEGER NOT NULL DEFAULT 0,
	millionaire_streak_id TEXT NOT NULL DEFAULT '',
	language TEXT NOT NULL DEFAULT 'en',
	player_combat_status TEXT NOT NULL DEFAULT '',
	player_combat_turns INTEGER NOT NULL DEFAULT 0,
	monster_combat_status TEXT NOT NULL DEFAULT '',
	monster_combat_turns INTEGER NOT NULL DEFAULT 0,
	online_id TEXT NOT NULL DEFAULT '',
	online_relics TEXT NOT NULL DEFAULT '[]',
	last_notification_seen INTEGER NOT NULL DEFAULT 0,
	last_action_message TEXT NOT NULL DEFAULT '',
	last_seen INTEGER NOT NULL,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_saves_last_seen ON saves(last_seen);
CREATE INDEX IF NOT EXISTS idx_saves_updated_at ON saves(updated_at);
CREATE INDEX IF NOT EXISTS idx_saves_online_id ON saves(online_id);
CREATE TABLE IF NOT EXISTS app_settings (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL,
	updated_at INTEGER NOT NULL
);
`
	if err := s.execSQL(schema); err != nil {
		return err
	}
	if err := s.ensureColumn("online_id", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.ensureColumn("online_relics", "TEXT NOT NULL DEFAULT '[]'"); err != nil {
		return err
	}
	if err := s.ensureColumn("last_notification_seen", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("last_action_message", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.ensureColumn("player_combat_status", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.ensureColumn("player_combat_turns", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("monster_combat_status", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.ensureColumn("monster_combat_turns", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("has_moon_pearl", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("has_sage_blessing", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("has_ruins_token", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("has_moon_charm", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("has_star_map", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("has_triptych_blessing", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("consecutive_wrong_answers", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("millionaire_points", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("lifeline_audience_used", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("lifeline_wiseman_used", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if err := s.ensureColumn("millionaire_streak_id", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := s.ensureColumn("language", "TEXT NOT NULL DEFAULT 'en'"); err != nil {
		return err
	}
	if err := s.execSQL(`UPDATE saves SET online_id = 'slot-' || id || '-' || lower(hex(randomblob(8))) WHERE online_id = '' OR online_id IS NULL;`); err != nil {
		return err
	}
	return nil
}

func (s *Store) SetSetting(key, value string) error {
	if s == nil {
		return errors.New("sqlite store is not initialized")
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("setting key cannot be empty")
	}

	now := time.Now().Unix()
	sql := fmt.Sprintf(`
INSERT INTO app_settings (key, value, updated_at)
VALUES (%s, %s, %d)
ON CONFLICT(key) DO UPDATE SET
	value = excluded.value,
	updated_at = excluded.updated_at
`, sqliteString(key), sqliteString(value), now)
	return s.execSQL(sql)
}

func (s *Store) GetSetting(key string) (string, error) {
	if s == nil {
		return "", errors.New("sqlite store is not initialized")
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return "", errors.New("setting key cannot be empty")
	}

	rows, err := s.queryRows(fmt.Sprintf(`SELECT value FROM app_settings WHERE key = %s LIMIT 1`, sqliteString(key)))
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", nil
	}
	return rows[0]["value"], nil
}

func (s *Store) DeleteSetting(key string) error {
	if s == nil {
		return errors.New("sqlite store is not initialized")
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("setting key cannot be empty")
	}
	return s.execSQL(fmt.Sprintf(`DELETE FROM app_settings WHERE key = %s;`, sqliteString(key)))
}

func (s *Store) SaveGame(gs *GameState) error {
	if s == nil {
		return errors.New("sqlite store is not initialized")
	}

	if gs == nil {
		return errors.New("game state is nil")
	}
	gs.EnsureOnlineID()
	now := time.Now().Unix()
	gs.LastSeen = now

	sql := fmt.Sprintf(`
INSERT INTO saves (
	id, player_name, current_room_id, health, max_health, is_game_over, gold,
	monster_health, is_troll_dead, is_beast_dead, is_frost_giant_dead,
	skill_points, level, experience, required_xp, base_damage,
	has_garden_kit, has_meat, has_steel_sword, has_water, has_fur_coat,
	has_sun_amulet, has_ice_crystal, has_forsaken_blade, has_moon_pearl, has_sage_blessing,
	has_ruins_token, has_moon_charm, has_star_map, has_triptych_blessing, is_1x1x1x1_dead,
	has_seen_tutorial, consecutive_wrong_answers, millionaire_points, lifeline_audience_used, lifeline_wiseman_used, millionaire_streak_id, player_combat_status, player_combat_turns, monster_combat_status, monster_combat_turns,
	online_id, online_relics, last_seen, created_at, updated_at, last_notification_seen, last_action_message
	, language
)
VALUES (
	%d, %s, %s, %d, %d, %d, %d,
	%d, %d, %d, %d,
	%d, %d, %d, %d, %d,
	%d, %d, %d, %d, %d,
	%d, %d, %d, %d, %d,
	%d, %d, %d, %d, %d,
	%d, %d, %d, %d, %d, %s, %s, %d, %s, %d,
	%s, %s, %d, COALESCE((SELECT created_at FROM saves WHERE id = %d), %d), %d, %d, %s, %s
)
ON CONFLICT(id) DO UPDATE SET
	player_name = excluded.player_name,
	current_room_id = excluded.current_room_id,
	health = excluded.health,
	max_health = excluded.max_health,
	is_game_over = excluded.is_game_over,
	gold = excluded.gold,
	monster_health = excluded.monster_health,
	is_troll_dead = excluded.is_troll_dead,
	is_beast_dead = excluded.is_beast_dead,
	is_frost_giant_dead = excluded.is_frost_giant_dead,
	skill_points = excluded.skill_points,
	level = excluded.level,
	experience = excluded.experience,
	required_xp = excluded.required_xp,
	base_damage = excluded.base_damage,
	has_garden_kit = excluded.has_garden_kit,
	has_meat = excluded.has_meat,
	has_steel_sword = excluded.has_steel_sword,
	has_water = excluded.has_water,
	has_fur_coat = excluded.has_fur_coat,
	has_sun_amulet = excluded.has_sun_amulet,
	has_ice_crystal = excluded.has_ice_crystal,
	has_forsaken_blade = excluded.has_forsaken_blade,
	has_moon_pearl = excluded.has_moon_pearl,
	has_sage_blessing = excluded.has_sage_blessing,
	has_ruins_token = excluded.has_ruins_token,
	has_moon_charm = excluded.has_moon_charm,
	has_star_map = excluded.has_star_map,
	has_triptych_blessing = excluded.has_triptych_blessing,
	is_1x1x1x1_dead = excluded.is_1x1x1x1_dead,
	has_seen_tutorial = excluded.has_seen_tutorial,
	consecutive_wrong_answers = excluded.consecutive_wrong_answers,
	millionaire_points = excluded.millionaire_points,
	lifeline_audience_used = excluded.lifeline_audience_used,
	lifeline_wiseman_used = excluded.lifeline_wiseman_used,
	millionaire_streak_id = excluded.millionaire_streak_id,
	player_combat_status = excluded.player_combat_status,
	player_combat_turns = excluded.player_combat_turns,
	monster_combat_status = excluded.monster_combat_status,
	monster_combat_turns = excluded.monster_combat_turns,
	online_id = excluded.online_id,
	online_relics = excluded.online_relics,
	last_notification_seen = excluded.last_notification_seen,
	last_action_message = excluded.last_action_message,
	language = excluded.language,
	last_seen = excluded.last_seen,
	updated_at = excluded.updated_at
`,
		gs.ID,
		sqliteString(gs.PlayerName),
		sqliteString(gs.CurrentRoomID),
		gs.Health, gs.MaxHealth, boolToInt(gs.IsGameOver), gs.Gold,
		gs.MonsterHealth, boolToInt(gs.IsTrollDead), boolToInt(gs.IsBeastDead), boolToInt(gs.IsFrostGiantDead),
		gs.SkillPoints, gs.Level, gs.Experience, gs.RequiredXP, gs.BaseDamage,
		boolToInt(gs.HasGardenKit), boolToInt(gs.HasMeat), boolToInt(gs.HasSteelSword), boolToInt(gs.HasWater), boolToInt(gs.HasFurCoat),
		boolToInt(gs.HasSunAmulet), boolToInt(gs.HasIceCrystal), boolToInt(gs.HasForsakenBlade), boolToInt(gs.HasMoonPearl), boolToInt(gs.HasSageBlessing),
		boolToInt(gs.HasRuinsToken), boolToInt(gs.HasMoonCharm), boolToInt(gs.HasStarMap), boolToInt(gs.HasTriptychBlessing), boolToInt(gs.Is1x1x1x1Dead),
		boolToInt(gs.HasSeenTutorial), gs.ConsecutiveWrongAnswers, gs.MillionairePoints, boolToInt(gs.LifelineAudienceUsed), boolToInt(gs.LifelineWiseManUsed), sqliteString(gs.MillionaireStreakID), sqliteString(gs.PlayerCombatStatus), gs.PlayerCombatTurns, sqliteString(gs.MonsterCombatStatus), gs.MonsterCombatTurns,
		sqliteString(gs.OnlineID), sqliteString(mustJSON(gs.OnlineRelics)), gs.LastSeen, gs.ID, now, now, gs.LastNotificationSeen, sqliteString(gs.LastActionMessage), sqliteString(NormalizeLanguage(gs.Language)))
	return s.execSQL(sql)
}

func (s *Store) LoadSave(slotID int) (*GameState, error) {
	if s == nil {
		return nil, errors.New("sqlite store is not initialized")
	}

	records, err := s.queryRecords(fmt.Sprintf(`
SELECT id, player_name, current_room_id, health, max_health, is_game_over, gold,
	monster_health, is_troll_dead, is_beast_dead, is_frost_giant_dead,
	skill_points, level, experience, required_xp, base_damage,
	has_garden_kit, has_meat, has_steel_sword, has_water, has_fur_coat,
	has_sun_amulet, has_ice_crystal, has_forsaken_blade, has_moon_pearl, has_sage_blessing,
	has_ruins_token, has_moon_charm, has_star_map, has_triptych_blessing, is_1x1x1x1_dead,
	has_seen_tutorial, consecutive_wrong_answers, millionaire_points, lifeline_audience_used, lifeline_wiseman_used, millionaire_streak_id, language, player_combat_status, player_combat_turns, monster_combat_status, monster_combat_turns,
	online_id, online_relics, last_notification_seen, last_action_message, last_seen
FROM saves
WHERE id = %d`, slotID))
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return recordToGameState(records[0]), nil
}

func (s *Store) GetAllSaves() ([]GameState, error) {
	if s == nil {
		return nil, errors.New("sqlite store is not initialized")
	}

	records, err := s.queryRecords(`
SELECT id, player_name, current_room_id, health, max_health, is_game_over, gold,
	monster_health, is_troll_dead, is_beast_dead, is_frost_giant_dead,
	skill_points, level, experience, required_xp, base_damage,
	has_garden_kit, has_meat, has_steel_sword, has_water, has_fur_coat,
	has_sun_amulet, has_ice_crystal, has_forsaken_blade, has_moon_pearl, has_sage_blessing,
	has_ruins_token, has_moon_charm, has_star_map, has_triptych_blessing, is_1x1x1x1_dead,
	has_seen_tutorial, consecutive_wrong_answers, millionaire_points, lifeline_audience_used, lifeline_wiseman_used, millionaire_streak_id, language, player_combat_status, player_combat_turns, monster_combat_status, monster_combat_turns,
	online_id, online_relics, last_notification_seen, last_action_message, last_seen
FROM saves
ORDER BY last_seen DESC, id ASC`)
	if err != nil {
		return nil, err
	}

	saves := make([]GameState, 0, len(records))
	for _, record := range records {
		saves = append(saves, *recordToGameState(record))
	}
	return saves, nil
}

func (s *Store) GetOnlinePlayers(currentGS *GameState) ([]GameState, error) {
	if s == nil {
		return nil, errors.New("sqlite store is not initialized")
	}

	threshold := time.Now().Unix() - 120
	records, err := s.queryRecords(fmt.Sprintf(`
SELECT id, player_name, current_room_id, health, max_health, is_game_over, gold,
	monster_health, is_troll_dead, is_beast_dead, is_frost_giant_dead,
	skill_points, level, experience, required_xp, base_damage,
	has_garden_kit, has_meat, has_steel_sword, has_water, has_fur_coat,
	has_sun_amulet, has_ice_crystal, has_forsaken_blade, has_moon_pearl, has_sage_blessing,
	has_ruins_token, has_moon_charm, has_star_map, has_triptych_blessing, is_1x1x1x1_dead,
	has_seen_tutorial, language, online_id, online_relics, last_notification_seen, last_action_message, last_seen
FROM saves
WHERE last_seen > %d AND id != %d
ORDER BY last_seen DESC`, threshold, currentGS.ID))
	if err != nil {
		return nil, err
	}

	players := make([]GameState, 0, len(records))
	for _, record := range records {
		players = append(players, *recordToGameState(record))
	}
	return players, nil
}

func (s *Store) execSQL(sql string) error {
	cmd := exec.Command("sqlite3", s.path, sql)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return fmt.Errorf("sqlite3 exec failed: %w: %s", err, msg)
		}
		return fmt.Errorf("sqlite3 exec failed: %w", err)
	}
	return nil
}

func (s *Store) ensureColumn(name, definition string) error {
	rows, err := s.queryTableInfo()
	if err != nil {
		return err
	}
	for _, row := range rows {
		if strings.EqualFold(row.Name, name) {
			return nil
		}
	}
	return s.execSQL(fmt.Sprintf(`ALTER TABLE saves ADD COLUMN %s %s;`, name, definition))
}

func (s *Store) queryRecords(sql string) ([]sqliteSaveRecord, error) {
	cmd := exec.Command("sqlite3", "-json", s.path, sql)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return nil, fmt.Errorf("sqlite3 query failed: %w: %s", err, msg)
		}
		return nil, fmt.Errorf("sqlite3 query failed: %w", err)
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" || raw == "[]" {
		return nil, nil
	}

	var records []sqliteSaveRecord
	if err := json.Unmarshal(out, &records); err != nil {
		return nil, err
	}
	return records, nil
}

type sqliteTableColumn struct {
	Name string `json:"name"`
}

func (s *Store) queryTableInfo() ([]sqliteTableColumn, error) {
	cmd := exec.Command("sqlite3", "-json", s.path, `PRAGMA table_info(saves);`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return nil, fmt.Errorf("sqlite3 query failed: %w: %s", err, msg)
		}
		return nil, fmt.Errorf("sqlite3 query failed: %w", err)
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" || raw == "[]" {
		return nil, nil
	}

	var rows []sqliteTableColumn
	if err := json.Unmarshal(out, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Store) queryRows(sql string) ([]map[string]string, error) {
	cmd := exec.Command("sqlite3", "-json", s.path, sql)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return nil, fmt.Errorf("sqlite3 query failed: %w: %s", err, msg)
		}
		return nil, fmt.Errorf("sqlite3 query failed: %w", err)
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" || raw == "[]" {
		return nil, nil
	}

	var rows []map[string]string
	if err := json.Unmarshal(out, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func recordToGameState(record sqliteSaveRecord) *GameState {
	return &GameState{
		ID:                   record.ID,
		PlayerName:           record.PlayerName,
		CurrentRoomID:        record.CurrentRoomID,
		Health:               record.Health,
		MaxHealth:            record.MaxHealth,
		IsGameOver:           intToBool(record.IsGameOver),
		Gold:                 record.Gold,
		MonsterHealth:        record.MonsterHealth,
		IsTrollDead:          intToBool(record.IsTrollDead),
		IsBeastDead:          intToBool(record.IsBeastDead),
		IsFrostGiantDead:     intToBool(record.IsFrostGiantDead),
		SkillPoints:          record.SkillPoints,
		Level:                record.Level,
		Experience:           record.Experience,
		RequiredXP:           record.RequiredXP,
		BaseDamage:           record.BaseDamage,
		HasGardenKit:         intToBool(record.HasGardenKit),
		HasMeat:              intToBool(record.HasMeat),
		HasSteelSword:        intToBool(record.HasSteelSword),
		HasWater:             intToBool(record.HasWater),
		HasFurCoat:           intToBool(record.HasFurCoat),
		HasSunAmulet:         intToBool(record.HasSunAmulet),
		HasIceCrystal:        intToBool(record.HasIceCrystal),
		HasForsakenBlade:     intToBool(record.HasForsakenBlade),
		HasMoonPearl:         intToBool(record.HasMoonPearl),
		HasSageBlessing:      intToBool(record.HasSageBlessing),
		HasRuinsToken:        intToBool(record.HasRuinsToken),
		HasMoonCharm:         intToBool(record.HasMoonCharm),
		HasStarMap:           intToBool(record.HasStarMap),
		HasTriptychBlessing:  intToBool(record.HasTriptychBlessing),
		Is1x1x1x1Dead:        intToBool(record.Is1x1x1x1Dead),
		HasSeenTutorial:      intToBool(record.HasSeenTutorial),
		ConsecutiveWrongAnswers: record.ConsecutiveWrongAnswers,
		MillionairePoints:       record.MillionairePoints,
		LifelineAudienceUsed:    intToBool(record.LifelineAudienceUsed),
		LifelineWiseManUsed:     intToBool(record.LifelineWiseManUsed),
		MillionaireStreakID:     record.MillionaireStreakID,
		Language:             NormalizeLanguage(record.Language),
		PlayerCombatStatus:   record.PlayerCombatStatus,
		PlayerCombatTurns:    record.PlayerCombatTurns,
		MonsterCombatStatus:  record.MonsterCombatStatus,
		MonsterCombatTurns:   record.MonsterCombatTurns,
		OnlineID:             record.OnlineID,
		OnlineRelics:         parseJSONStrings(record.OnlineRelics),
		LastNotificationSeen: record.LastNotificationSeen,
		LastActionMessage:    record.LastActionMessage,
		LastSeen:             record.LastSeen,
	}
}

func sqliteString(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func intToBool(v int) bool {
	return v != 0
}

func mustJSON(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	data, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func parseJSONStrings(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "null" {
		return nil
	}

	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return nil
	}
	return values
}
