package game

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type ActionType int

const (
	MoveAction ActionType = iota
	ChoiceAction
	CustomAction
)

type Action struct {
	ID          string
	Description string
	Trigger     string // e.g., "1", "2"
	Type        ActionType
	Result      func(*GameState) string
	IsAvailable func(*GameState) bool // Optional: Defaults to true
}

type Room struct {
	ID          string
	Description string
	Actions     []Action
}

type GameState struct {
	ID                      int              `db:"id"`
	PlayerName              string           `db:"player_name"`
	CurrentRoomID           string           `db:"current_room_id"`
	Health                  int              `db:"health"`
	MaxHealth               int              `db:"max_health"`
	IsGameOver              bool             `db:"is_game_over"`
	Gold                    int              `db:"gold"`
	MonsterHealth           int              `db:"monster_health"`
	IsTrollDead             bool             `db:"is_troll_dead"`
	IsBeastDead             bool             `db:"is_beast_dead"`
	IsFrostGiantDead        bool             `db:"is_frost_giant_dead"`
	SkillPoints             int              `db:"skill_points"`
	Level                   int              `db:"level"`
	Experience              int              `db:"experience"`
	RequiredXP              int              `db:"required_xp"`
	BaseDamage              int              `db:"base_damage"`
	HasGardenKit            bool             `db:"has_garden_kit"`
	HasMeat                 bool             `db:"has_meat"`
	HasSteelSword           bool             `db:"has_steel_sword"`
	HasWater                bool             `db:"has_water"`
	HasFurCoat              bool             `db:"has_fur_coat"`
	HasSunAmulet            bool             `db:"has_sun_amulet"`
	HasIceCrystal           bool             `db:"has_ice_crystal"`
	HasForsakenBlade        bool             `db:"has_forsaken_blade"`
	HasMoonPearl            bool             `db:"has_moon_pearl"`
	HasSageBlessing         bool             `db:"has_sage_blessing"`
	HasRuinsToken           bool             `db:"has_ruins_token"`
	HasMoonCharm            bool             `db:"has_moon_charm"`
	HasStarMap              bool             `db:"has_star_map"`
	HasTriptychBlessing     bool             `db:"has_triptych_blessing"`
	Is1x1x1x1Dead           bool             `db:"is_1x1x1x1_dead"`
	HasSeenTutorial         bool             `db:"has_seen_tutorial"`
	ConsecutiveWrongAnswers int              `db:"consecutive_wrong_answers"`
	MillionairePoints       int              `db:"millionaire_points"`
	LifelineAudienceUsed    bool             `db:"lifeline_audience_used"`
	LifelineWiseManUsed     bool             `db:"lifeline_wiseman_used"`
	MillionaireStreakID     string           `db:"millionaire_streak_id"`
	Language                string           `db:"language"`
	PlayerCombatStatus      string           `db:"player_combat_status"`
	PlayerCombatTurns       int              `db:"player_combat_turns"`
	MonsterCombatStatus     string           `db:"monster_combat_status"`
	MonsterCombatTurns      int              `db:"monster_combat_turns"`
	OnlineID                string           `db:"online_id"`
	OnlineRelics            []string         `db:"online_relics"`
	LastNotificationSeen    int64            `db:"last_notification_seen"`
	LastActionMessage       string           `db:"last_action_message"`
	LastSeen                int64            `db:"last_seen"`
	PartySupport            int              `db:"-"`
	PartyBattleBuffTurns    int              `db:"-"`
	PartyBattleBuffBonus    int              `db:"-"`
	PartyGuardTurns         int              `db:"-"`
	PartyGuardReduction     int              `db:"-"`
	PartyQuestKey           string           `db:"-"`
	PartyQuestName          string           `db:"-"`
	PartyQuestStatus        string           `db:"-"`
	PartyQuestPhase         string           `db:"-"`
	PartyQuestPhaseIndex    int              `db:"-"`
	PartyQuestPhaseGoal     int              `db:"-"`
	PartyQuestPhaseProgress int              `db:"-"`
	PartyQuestRewardGold    int              `db:"-"`
	PartyQuestRewardRelic   string           `db:"-"`
	World                   map[string]*Room `db:"-"`
	Inventory               []string         `db:"-"`
}

func NewGameState(id int, playerName string) *GameState {
	gs := &GameState{
		ID:            id,
		PlayerName:    playerName,
		CurrentRoomID: "village",
		Health:        100,
		MaxHealth:     100,
		IsGameOver:    false,
		Gold:          50,
		MonsterHealth: 50,
		Level:         1,
		Experience:    0,
		RequiredXP:    100,
		SkillPoints:   0,
		BaseDamage:              5,
		Language:                "en",
		MillionairePoints:       0,
		LifelineAudienceUsed:    false,
		LifelineWiseManUsed:     false,
		MillionaireStreakID:     "",
		}
		gs.EnsureOnlineID()

	return gs
}

func (gs *GameState) EnsureOnlineID() {
	if gs == nil || gs.OnlineID != "" {
		return
	}

	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		gs.OnlineID = fmt.Sprintf("slot-%d-%d", gs.ID, time.Now().UnixNano())
		return
	}
	gs.OnlineID = "slot-" + hex.EncodeToString(buf[:])
}

func (gs *GameState) ResetRun() {
	if gs == nil {
		return
	}

	world := gs.World
	seenTutorial := gs.HasSeenTutorial
	gs.CurrentRoomID = "village"
	gs.Health = 100
	gs.MaxHealth = 100
	gs.IsGameOver = false
	gs.Gold = 50
	gs.MonsterHealth = 50
	gs.IsTrollDead = false
	gs.IsBeastDead = false
	gs.IsFrostGiantDead = false
	gs.SkillPoints = 0
	gs.Level = 1
	gs.Experience = 0
	gs.RequiredXP = 100
	gs.BaseDamage = 5
	gs.HasGardenKit = false
	gs.HasMeat = false
	gs.HasSteelSword = false
	gs.HasWater = false
	gs.HasFurCoat = false
	gs.HasSunAmulet = false
	gs.HasIceCrystal = false
	gs.HasForsakenBlade = false
	gs.HasMoonPearl = false
	gs.HasSageBlessing = false
	gs.HasRuinsToken = false
	gs.HasMoonCharm = false
	gs.HasStarMap = false
	gs.HasTriptychBlessing = false
	gs.Is1x1x1x1Dead = false
	gs.HasSeenTutorial = seenTutorial
	gs.ConsecutiveWrongAnswers = 0
	gs.MillionairePoints = 0
	gs.LifelineAudienceUsed = false
	gs.LifelineWiseManUsed = false
	gs.MillionaireStreakID = ""
	gs.Language = NormalizeLanguage(gs.Language)
	gs.ClearCombatState()
	gs.PartyBattleBuffTurns = 0
	gs.PartyBattleBuffBonus = 0
	gs.PartyGuardTurns = 0
	gs.PartyGuardReduction = 0
	gs.ClearPartyQuestState()
	gs.EnsureOnlineID()
	gs.LastSeen = 0
	gs.LastActionMessage = ""
	gs.World = world
	gs.Inventory = nil
}

func NormalizeLanguage(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "pt", "pt-br", "pt_pt", "portuguese":
		return "pt"
	default:
		return "en"
	}
}

func (gs *GameState) CheckLevelUp() string {
	msg := ""
	for gs.Experience >= gs.RequiredXP {
		gs.Level++
		gs.Experience -= gs.RequiredXP
		gs.RequiredXP = int(float64(gs.RequiredXP) * 1.6)
		gs.MaxHealth += 25
		gs.Health = gs.MaxHealth
		gs.BaseDamage += 5
		msg += fmt.Sprintf("\n%s*** LEVEL UP! You are now Level %d! ***%s\nMax Health and Damage increased!\n", ColorCyan+ColorBold, gs.Level, ColorReset)
	}
	return msg
}

func (gs *GameState) Transition(roomID string) string {
	if _, ok := gs.World[roomID]; ok {
		gs.CurrentRoomID = roomID
		gs.ClearCombatState()
		return ""
	}
	return fmt.Sprintf("Error: Room %s not found", roomID)
}

func (gs *GameState) ClearCombatState() {
	if gs == nil {
		return
	}
	gs.PlayerCombatStatus = ""
	gs.PlayerCombatTurns = 0
	gs.MonsterCombatStatus = ""
	gs.MonsterCombatTurns = 0
}

func (gs *GameState) ClearPartyQuestState() {
	if gs == nil {
		return
	}
	gs.PartyQuestKey = ""
	gs.PartyQuestName = ""
	gs.PartyQuestStatus = ""
	gs.PartyQuestPhase = ""
	gs.PartyQuestPhaseIndex = 0
	gs.PartyQuestPhaseGoal = 0
	gs.PartyQuestPhaseProgress = 0
	gs.PartyQuestRewardGold = 0
	gs.PartyQuestRewardRelic = ""
}
