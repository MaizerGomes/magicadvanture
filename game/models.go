package game

import "fmt"

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
	ID             int    `db:"id" bson:"_id"`
	PlayerName     string `db:"player_name" bson:"player_name"`
	CurrentRoomID  string `db:"current_room_id" bson:"current_room_id"`
	Health         int    `db:"health" bson:"health"`
	MaxHealth      int    `db:"max_health" bson:"max_health"`
	IsGameOver     bool   `db:"is_game_over" bson:"is_game_over"`
	Gold           int    `db:"gold" bson:"gold"`
	MonsterHealth    int    `db:"monster_health" bson:"monster_health"`
	IsTrollDead      bool   `db:"is_troll_dead" bson:"is_troll_dead"`
	IsBeastDead      bool   `db:"is_beast_dead" bson:"is_beast_dead"`
	IsFrostGiantDead bool   `db:"is_frost_giant_dead" bson:"is_frost_giant_dead"`
	SkillPoints      int    `db:"skill_points" bson:"skill_points"`
	Level          int    `db:"level" bson:"level"`
	Experience     int    `db:"experience" bson:"experience"`
	RequiredXP     int    `db:"required_xp" bson:"required_xp"`
	BaseDamage     int    `db:"base_damage" bson:"base_damage"`
	HasGardenKit   bool   `db:"has_garden_kit" bson:"has_garden_kit"`
	HasMeat        bool   `db:"has_meat" bson:"has_meat"`
	HasSteelSword  bool   `db:"has_steel_sword" bson:"has_steel_sword"`
	HasWater       bool   `db:"has_water" bson:"has_water"`
	HasFurCoat     bool   `db:"has_fur_coat" bson:"has_fur_coat"`
	HasSunAmulet   bool   `db:"has_sun_amulet" bson:"has_sun_amulet"`
	HasIceCrystal  bool   `db:"has_ice_crystal" bson:"has_ice_crystal"`
	HasForsakenBlade bool `db:"has_forsaken_blade" bson:"has_forsaken_blade"`
	Is1x1x1x1Dead    bool   `db:"is_1x1x11_dead" bson:"is_1x1x11_dead"`
	HasSeenTutorial  bool   `db:"has_seen_tutorial" bson:"has_seen_tutorial"`
	LastSeen         int64  `db:"last_seen" bson:"last_seen"` // Unix timestamp
	World            map[string]*Room `db:"-" bson:"-"`
	Inventory      []string `db:"-" bson:"-"`
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
		return ""
	}
	return fmt.Sprintf("Error: Room %s not found", roomID)
}
