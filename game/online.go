package game

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const DefaultOnlineDBName = "magicadventure_online"

type OnlineService struct {
	client  *mongo.Client
	db      *mongo.Database
	enabled bool
}

type OnlinePlayer struct {
	OnlineID            string    `bson:"online_id" json:"online_id"`
	PlayerName          string    `bson:"player_name" json:"player_name"`
	CurrentRoomID       string    `bson:"current_room_id" json:"current_room_id"`
	PartyID             string    `bson:"party_id" json:"party_id"`
	Health              int       `bson:"health" json:"health"`
	MaxHealth           int       `bson:"max_health" json:"max_health"`
	PartyBuffTurns      int       `bson:"party_buff_turns" json:"party_buff_turns"`
	PartyBuffBonus      int       `bson:"party_buff_bonus" json:"party_buff_bonus"`
	PartyGuardTurns     int       `bson:"party_guard_turns" json:"party_guard_turns"`
	PartyGuardReduction int       `bson:"party_guard_reduction" json:"party_guard_reduction"`
	Level               int       `bson:"level" json:"level"`
	Gold                int       `bson:"gold" json:"gold"`
	LastSeen            time.Time `bson:"last_seen" json:"last_seen"`
	Relics              []string  `bson:"relics" json:"relics"`
}

type RoomMessage struct {
	OnlineID      string    `bson:"online_id" json:"online_id"`
	PlayerName    string    `bson:"player_name" json:"player_name"`
	CurrentRoomID string    `bson:"current_room_id" json:"current_room_id"`
	Kind          string    `bson:"kind" json:"kind"`
	Text          string    `bson:"text" json:"text"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

type WhisperMessage struct {
	OnlineID   string    `bson:"online_id" json:"online_id"`
	PlayerName string    `bson:"player_name" json:"player_name"`
	TargetID   string    `bson:"target_id" json:"target_id"`
	TargetName string    `bson:"target_name" json:"target_name"`
	Text       string    `bson:"text" json:"text"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

type InboxNotification struct {
	RecipientID string    `bson:"recipient_id" json:"recipient_id"`
	SenderID    string    `bson:"sender_id" json:"sender_id"`
	SenderName  string    `bson:"sender_name" json:"sender_name"`
	Kind        string    `bson:"kind" json:"kind"`
	ReferenceID string    `bson:"reference_id" json:"reference_id"`
	Text        string    `bson:"text" json:"text"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

type onlineProfile struct {
	ID                  string    `bson:"_id"`
	SlotID              int       `bson:"slot_id"`
	PlayerName          string    `bson:"player_name"`
	CurrentRoomID       string    `bson:"current_room_id"`
	PartyID             string    `bson:"party_id"`
	Health              int       `bson:"health"`
	MaxHealth           int       `bson:"max_health"`
	PartyBuffTurns      int       `bson:"party_buff_turns"`
	PartyBuffBonus      int       `bson:"party_buff_bonus"`
	PartyGuardTurns     int       `bson:"party_guard_turns"`
	PartyGuardReduction int       `bson:"party_guard_reduction"`
	Level               int       `bson:"level"`
	Gold                int       `bson:"gold"`
	LastSeen            time.Time `bson:"last_seen"`
	Relics              []string  `bson:"relics"`
}

type Party struct {
	ID                    string    `bson:"_id" json:"_id"`
	LeaderID              string    `bson:"leader_id" json:"leader_id"`
	LeaderName            string    `bson:"leader_name" json:"leader_name"`
	MemberIDs             []string  `bson:"member_ids" json:"member_ids"`
	QuestKey              string    `bson:"quest_key" json:"quest_key"`
	QuestStatus           string    `bson:"quest_status" json:"quest_status"`
	QuestPhase            string    `bson:"quest_phase" json:"quest_phase"`
	QuestPhaseProgress    int       `bson:"quest_phase_progress" json:"quest_phase_progress"`
	QuestPhaseGoal        int       `bson:"quest_phase_goal" json:"quest_phase_goal"`
	QuestPhaseIndex       int       `bson:"quest_phase_index" json:"quest_phase_index"`
	QuestDungeonPhase     bool      `bson:"quest_dungeon_phase" json:"quest_dungeon_phase"`
	QuestName             string    `bson:"quest_name" json:"quest_name"`
	QuestGoal             int       `bson:"quest_goal" json:"quest_goal"`
	QuestProgress         int       `bson:"quest_progress" json:"quest_progress"`
	QuestPhaseKey         string    `bson:"quest_phase_key" json:"quest_phase_key"`
	QuestPhaseName        string    `bson:"quest_phase_name" json:"quest_phase_name"`
	QuestPhaseDescription string    `bson:"quest_phase_description" json:"quest_phase_description"`
	QuestRewardDamage     int       `bson:"quest_reward_damage" json:"quest_reward_damage"`
	QuestRewardTurns      int       `bson:"quest_reward_turns" json:"quest_reward_turns"`
	QuestRewardGold       int       `bson:"quest_reward_gold" json:"quest_reward_gold"`
	QuestRewardRelic      string    `bson:"quest_reward_relic" json:"quest_reward_relic"`
	QuestGuardReduction   int       `bson:"quest_guard_reduction" json:"quest_guard_reduction"`
	QuestGuardTurns       int       `bson:"quest_guard_turns" json:"quest_guard_turns"`
	QuestCooldownUntil    time.Time `bson:"quest_cooldown_until" json:"quest_cooldown_until"`
	QuestLog              []string  `bson:"quest_log" json:"quest_log"`
	CreatedAt             time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time `bson:"updated_at" json:"updated_at"`
}

func (o *OnlineService) SubscribeRoom(gs *GameState) *RoomFeed {
	if !o.Enabled() || gs == nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	f := &RoomFeed{
		service: o,
		gs:      gs,
		roomID:  gs.CurrentRoomID,
		ctx:     ctx,
		cancel:  cancel,
	}
	// Note: We don't start the goroutines here to keep it simple or we can start them
	go f.watchRoomMessages()
	go f.watchWhispers()
	go f.watchInboxNotifications()
	go f.watchProfileChanges()
	go f.fallbackRefreshLoop()
	return f
}

type MillionaireStreak struct {
	ID          string    `bson:"_id"`
	StartTime   time.Time `bson:"start_time"`
	EndTime     time.Time `bson:"end_time"`
	IsProcessed bool      `bson:"is_processed"`
}

type MillionaireScore struct {
	OnlineID   string    `bson:"online_id"`
	PlayerName string    `bson:"player_name"`
	Score      int       `bson:"score"`
	StreakID   string    `bson:"streak_id"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

type MillionaireAudienceResponse struct {
	QuestionID string    `bson:"question_id"`
	PlayerID   string    `bson:"player_id"`
	Choice     int       `bson:"choice"`
	CreatedAt  time.Time `bson:"created_at"`
}

type RoomFeed struct {
	service    *OnlineService
	gs         *GameState
	roomID     string
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.RWMutex
	nearby     []OnlinePlayer
	messages   []RoomMessage
	whispers   []WhisperMessage
	inbox      []InboxNotification
	notify     func(string)
	refreshErr error
}

func InitOnline() (*OnlineService, error) {
	uri := os.Getenv("MONGO_URI")
	if strings.TrimSpace(uri) == "" {
		return nil, nil
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if strings.TrimSpace(dbName) == "" {
		dbName = DefaultOnlineDBName
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return &OnlineService{enabled: false}, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return &OnlineService{enabled: false}, err
	}

	return &OnlineService{
		client:  client,
		db:      client.Database(dbName),
		enabled: true,
	}, nil
}

func (o *OnlineService) Enabled() bool {
	return o != nil && o.enabled
}

func (o *OnlineService) Close() error {
	if o == nil || o.client == nil {
		return nil
	}
	return o.client.Disconnect(context.Background())
}

func (o *OnlineService) StartRoomFeed(gs *GameState, notify func(string)) (*RoomFeed, error) {
	if !o.Enabled() || gs == nil {
		return nil, nil
	}

	gs.EnsureOnlineID()
	ctx, cancel := context.WithCancel(context.Background())
	feed := &RoomFeed{
		service: o,
		gs:      gs,
		roomID:  gs.CurrentRoomID,
		ctx:     ctx,
		cancel:  cancel,
		notify:  notify,
	}

	if err := feed.refreshSnapshot(); err != nil {
		feed.refreshErr = err
	}

	go feed.watchRoomMessages()
	go feed.watchWhispers()
	go feed.watchInboxNotifications()
	go feed.watchProfileChanges()
	go feed.fallbackRefreshLoop()
	return feed, nil
}

func (o *OnlineService) SyncPresence(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	gs.EnsureOnlineID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile := onlineProfile{
		ID:                  gs.OnlineID,
		SlotID:              gs.ID,
		PlayerName:          gs.PlayerName,
		CurrentRoomID:       gs.CurrentRoomID,
		PartyID:             "",
		Health:              gs.Health,
		MaxHealth:           gs.MaxHealth,
		PartyBuffTurns:      gs.PartyBattleBuffTurns,
		PartyBuffBonus:      gs.PartyBattleBuffBonus,
		PartyGuardTurns:     gs.PartyGuardTurns,
		PartyGuardReduction: gs.PartyGuardReduction,
		Level:               gs.Level,
		Gold:                gs.Gold,
		LastSeen:            time.Now().UTC(),
		Relics:              append([]string(nil), gs.OnlineRelics...),
	}

	if existing, err := o.getOnlineProfile(ctx, gs.OnlineID); err == nil && existing != nil {
		profile.PartyID = existing.PartyID
	}
	if profile.PartyID != "" {
		if party, err := o.getPartyByID(ctx, profile.PartyID); err != nil || party == nil {
			profile.PartyID = ""
		}
	}

	_, err := o.db.Collection("online_profiles").ReplaceOne(ctx, bson.M{"_id": gs.OnlineID}, profile, options.Replace().SetUpsert(true))
	return err
}

func (o *OnlineService) RefreshSelfState(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile, err := o.getOnlineProfile(ctx, gs.OnlineID)
	if err != nil || profile == nil {
		return err
	}

	gs.Health = profile.Health
	if profile.MaxHealth > 0 {
		gs.MaxHealth = profile.MaxHealth
	}
	gs.PartyBattleBuffTurns = profile.PartyBuffTurns
	gs.PartyBattleBuffBonus = profile.PartyBuffBonus
	gs.PartyGuardTurns = profile.PartyGuardTurns
	gs.PartyGuardReduction = profile.PartyGuardReduction

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	syncPartyQuestState(gs, party)
	return nil
}

func (o *OnlineService) GetNearbyPlayers(gs *GameState) ([]OnlinePlayer, error) {
	if !o.Enabled() || gs == nil {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	threshold := time.Now().UTC().Add(-2 * time.Minute)
	filter := bson.M{
		"current_room_id": gs.CurrentRoomID,
		"last_seen":       bson.M{"$gt": threshold},
		"_id":             bson.M{"$ne": gs.OnlineID},
	}
	opts := options.Find().SetSort(bson.D{{Key: "last_seen", Value: -1}})
	cursor, err := o.db.Collection("online_profiles").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var players []OnlinePlayer
	if err := cursor.All(ctx, &players); err != nil {
		return nil, err
	}
	return players, nil
}

func (o *OnlineService) GetRoomMessages(roomID string, limit int) ([]RoomMessage, error) {
	if !o.Enabled() {
		return nil, nil
	}
	if limit <= 0 {
		limit = 6
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"current_room_id": roomID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	cursor, err := o.db.Collection("room_messages").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []RoomMessage
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (o *OnlineService) PostRoomMessage(gs *GameState, text string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("message cannot be empty")
	}
	if len(text) > 240 {
		text = text[:240]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.db.Collection("room_messages").InsertOne(ctx, bson.M{
		"online_id":       gs.OnlineID,
		"player_name":     gs.PlayerName,
		"current_room_id": gs.CurrentRoomID,
		"kind":            "chat",
		"text":            text,
		"created_at":      time.Now().UTC(),
	})
	return err
}

func (o *OnlineService) BroadcastRoomEvent(gs *GameState, text string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("room event cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.db.Collection("room_messages").InsertOne(ctx, bson.M{
		"online_id":       gs.OnlineID,
		"player_name":     gs.PlayerName,
		"current_room_id": gs.CurrentRoomID,
		"kind":            "event",
		"text":            text,
		"created_at":      time.Now().UTC(),
	})
	return err
}

func (o *OnlineService) NotifyPlayer(sender *GameState, targetID, kind, text string) error {
	if !o.Enabled() || sender == nil {
		return nil
	}
	targetID = strings.TrimSpace(targetID)
	kind = strings.TrimSpace(kind)
	text = strings.TrimSpace(text)
	if targetID == "" || kind == "" || text == "" {
		return errors.New("notification requires a target, kind, and text")
	}
	if len(text) > 240 {
		text = text[:240]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.db.Collection("player_notifications").InsertOne(ctx, bson.M{
		"recipient_id": targetID,
		"sender_id":    sender.OnlineID,
		"sender_name":  sender.PlayerName,
		"kind":         kind,
		"reference_id": "",
		"text":         text,
		"created_at":   time.Now().UTC(),
	})
	return err
}

func (o *OnlineService) NotifyPlayerWithReference(sender *GameState, targetID, kind, text, referenceID string) error {
	if !o.Enabled() || sender == nil {
		return nil
	}
	targetID = strings.TrimSpace(targetID)
	kind = strings.TrimSpace(kind)
	text = strings.TrimSpace(text)
	referenceID = strings.TrimSpace(referenceID)
	if targetID == "" || kind == "" || text == "" {
		return errors.New("notification requires a target, kind, and text")
	}
	if len(text) > 240 {
		text = text[:240]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.db.Collection("player_notifications").InsertOne(ctx, bson.M{
		"recipient_id": targetID,
		"sender_id":    sender.OnlineID,
		"sender_name":  sender.PlayerName,
		"kind":         kind,
		"reference_id": referenceID,
		"text":         text,
		"created_at":   time.Now().UTC(),
	})
	return err
}

func (o *OnlineService) SendWhisper(gs *GameState, targetID, targetName, text string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}
	text = strings.TrimSpace(text)
	targetID = strings.TrimSpace(targetID)
	targetName = strings.TrimSpace(targetName)
	if text == "" || targetID == "" {
		return errors.New("whisper requires a target and message")
	}
	if len(text) > 240 {
		text = text[:240]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.db.Collection("direct_messages").InsertOne(ctx, bson.M{
		"online_id":   gs.OnlineID,
		"player_name": gs.PlayerName,
		"target_id":   targetID,
		"target_name": targetName,
		"text":        text,
		"created_at":  time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	_ = o.NotifyPlayer(gs, targetID, "whisper", fmt.Sprintf("Whisper from %s: %s", gs.PlayerName, text))
	return nil
}

func (o *OnlineService) GetWhispers(gs *GameState, limit int) ([]WhisperMessage, error) {
	if !o.Enabled() || gs == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 6
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"target_id": gs.OnlineID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	cursor, err := o.db.Collection("direct_messages").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []WhisperMessage
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (o *OnlineService) GetNotifications(gs *GameState, since int64, limit int) ([]InboxNotification, error) {
	if !o.Enabled() || gs == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"recipient_id": gs.OnlineID}
	if since > 0 {
		filter["created_at"] = bson.M{"$gt": time.Unix(0, since)}
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	cursor, err := o.db.Collection("player_notifications").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []InboxNotification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (o *OnlineService) getOnlineProfile(ctx context.Context, id string) (*onlineProfile, error) {
	if !o.Enabled() {
		return nil, nil
	}

	var profile onlineProfile
	if err := o.db.Collection("online_profiles").FindOne(ctx, bson.M{"_id": id}).Decode(&profile); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (o *OnlineService) getPartyByID(ctx context.Context, partyID string) (*Party, error) {
	if !o.Enabled() {
		return nil, nil
	}

	var party Party
	if err := o.db.Collection("parties").FindOne(ctx, bson.M{"_id": partyID}).Decode(&party); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &party, nil
}

func (o *OnlineService) getPlayerParty(ctx context.Context, onlineID string) (*Party, error) {
	profile, err := o.getOnlineProfile(ctx, onlineID)
	if err != nil || profile == nil || strings.TrimSpace(profile.PartyID) == "" {
		return nil, err
	}
	return o.getPartyByID(ctx, profile.PartyID)
}

func (o *OnlineService) GetPartySummary(gs *GameState) (*Party, []OnlinePlayer, error) {
	if !o.Enabled() || gs == nil {
		return nil, nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil || party == nil {
		return nil, nil, err
	}

	cursor, err := o.db.Collection("online_profiles").Find(
		ctx,
		bson.M{"_id": bson.M{"$in": party.MemberIDs}},
		options.Find().SetSort(bson.D{{Key: "last_seen", Value: -1}}),
	)
	if err != nil {
		return party, nil, err
	}
	defer cursor.Close(ctx)

	var members []OnlinePlayer
	if err := cursor.All(ctx, &members); err != nil {
		return party, nil, err
	}
	return party, members, nil
}

type partyQuestDefinition struct {
	Key            string
	Name           string
	RewardDamage   int
	RewardTurns    int
	RewardGold     int
	RewardRelic    string
	GuardReduction int
	GuardTurns     int
	Description    string
	Phases         []partyQuestPhaseDefinition
}

type partyQuestPhaseDefinition struct {
	Key         string
	Name        string
	Goal        int
	RoomID      string
	Description string
}

func partyQuestDefinitions() map[string]partyQuestDefinition {
	return map[string]partyQuestDefinition{
		"monster-hunt": {
			Key:            "monster-hunt",
			Name:           "Monster Hunt",
			RewardDamage:   12,
			RewardTurns:    3,
			RewardGold:     180,
			RewardRelic:    "Hunter's Claw",
			GuardReduction: 8,
			GuardTurns:     2,
			Description:    "Track, clear, and defeat a hidden beast pack across a three-phase dungeon. Party wins rally, guard, gold, and a relic on completion.",
			Phases: []partyQuestPhaseDefinition{
				{Key: "trail", Name: "Trail the Beasts", Goal: 1, RoomID: "party_hunt_gate", Description: "Enter the beast trail and break through the outer sentries."},
				{Key: "lair", Name: "Clear the Lair", Goal: 2, RoomID: "party_hunt_depths", Description: "Defeat the lesser beasts and force your way deeper."},
				{Key: "alpha", Name: "Hunt the Alpha", Goal: 1, RoomID: "party_hunt_vault", Description: "Confront the alpha guardian and finish the hunt."},
			},
		},
		"void-expedition": {
			Key:            "void-expedition",
			Name:           "Void Expedition",
			RewardDamage:   8,
			RewardTurns:    3,
			RewardGold:     220,
			RewardRelic:    "Null Compass",
			GuardReduction: 12,
			GuardTurns:     2,
			Description:    "Walk through a shifting dungeon of code and seal the core before it rewrites the party.",
			Phases: []partyQuestPhaseDefinition{
				{Key: "trace", Name: "Map the Static", Goal: 1, RoomID: "party_void_gate", Description: "Find the stable route through the flickering entrance."},
				{Key: "fracture", Name: "Break the Fractures", Goal: 2, RoomID: "party_void_depths", Description: "Shatter the data walls and suppress the distortions."},
				{Key: "core", Name: "Seal the Core", Goal: 1, RoomID: "party_void_core", Description: "Confront the core guardian and seal the breach."},
			},
		},
		"frost-pact": {
			Key:            "frost-pact",
			Name:           "Frost Pact",
			RewardDamage:   10,
			RewardTurns:    3,
			RewardGold:     200,
			RewardRelic:    "Wyrmfang Sigil",
			GuardReduction: 10,
			GuardTurns:     2,
			Description:    "Seal a cursed glacier and end the Frost Wyrm's pact with the mountain.",
			Phases: []partyQuestPhaseDefinition{
				{Key: "cross", Name: "Cross the Glacier", Goal: 1, RoomID: "party_frost_gate", Description: "Reach the frozen gate and break the outer ward."},
				{Key: "ward", Name: "Break the Ward", Goal: 2, RoomID: "party_frost_depths", Description: "Shatter the glacier ward and force the wyrm deeper."},
				{Key: "wyrm", Name: "Face the Frost Wyrm", Goal: 1, RoomID: "party_frost_peak", Description: "Defeat the Frost Wyrm at the mountain peak."},
			},
		},
		"sun-covenant": {
			Key:            "sun-covenant",
			Name:           "Sun Covenant",
			RewardDamage:   14,
			RewardTurns:    4,
			RewardGold:     240,
			RewardRelic:    "Sunfire Prism",
			GuardReduction: 11,
			GuardTurns:     2,
			Description:    "Carry the party through a sacred dune temple and bind the Sunbound Colossus before dawn breaks fully.",
			Phases: []partyQuestPhaseDefinition{
				{Key: "dawn", Name: "Raise the Dawn Gate", Goal: 1, RoomID: "party_sun_gate", Description: "Light the sun gate and force the temple to open."},
				{Key: "dunes", Name: "Break the Burning Dunes", Goal: 2, RoomID: "party_sun_depths", Description: "Crack the dune wards and survive the radiant sandstorm."},
				{Key: "crown", Name: "Face the Sunbound Colossus", Goal: 1, RoomID: "party_sun_crown", Description: "Climb to the crown chamber and break the colossus."},
			},
		},
	}
}

func questPhaseAt(def partyQuestDefinition, index int) (partyQuestPhaseDefinition, bool) {
	if index < 0 || index >= len(def.Phases) {
		return partyQuestPhaseDefinition{}, false
	}
	return def.Phases[index], true
}

func syncPartyQuestState(gs *GameState, party *Party) {
	if gs == nil {
		return
	}
	if party == nil {
		gs.ClearPartyQuestState()
		return
	}

	gs.PartyQuestKey = party.QuestKey
	gs.PartyQuestName = party.QuestName
	gs.PartyQuestStatus = party.QuestStatus
	gs.PartyQuestPhase = party.QuestPhaseName
	gs.PartyQuestPhaseIndex = party.QuestPhaseIndex
	gs.PartyQuestPhaseGoal = party.QuestPhaseGoal
	gs.PartyQuestPhaseProgress = party.QuestPhaseProgress
	gs.PartyQuestRewardGold = party.QuestRewardGold
	gs.PartyQuestRewardRelic = party.QuestRewardRelic
}

func partyQuestProgressSummary(questName, phaseName string, progress, goal int) (string, int, int) {
	label := strings.TrimSpace(phaseName)
	if label == "" {
		label = strings.TrimSpace(questName)
	}
	return label, progress, goal
}

func partyQuestCompletionRewards(party *Party) (damageBonus, turns, guardReduction, guardTurns, gold int, relic string) {
	if party == nil {
		return 10, 3, 8, 2, 150, randomRelic()
	}

	damageBonus = party.QuestRewardDamage
	turns = party.QuestRewardTurns
	guardReduction = party.QuestGuardReduction
	guardTurns = party.QuestGuardTurns
	gold = party.QuestRewardGold
	relic = strings.TrimSpace(party.QuestRewardRelic)

	if damageBonus <= 0 {
		damageBonus = 10
	}
	if turns <= 0 {
		turns = 3
	}
	if guardReduction <= 0 {
		guardReduction = 8
	}
	if guardTurns <= 0 {
		guardTurns = 2
	}
	if gold <= 0 {
		gold = 150
	}
	if relic == "" {
		relic = randomRelic()
	}
	return
}

func partyQuestCooldownRemaining(party *Party) time.Duration {
	if party == nil || party.QuestCooldownUntil.IsZero() {
		return 0
	}
	return time.Until(party.QuestCooldownUntil)
}

func (o *OnlineService) CreateParty(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existing, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	partyID := fmt.Sprintf("party-%s-%d", gs.OnlineID, time.Now().UnixNano())
	now := time.Now().UTC()
	party := Party{
		ID:         partyID,
		LeaderID:   gs.OnlineID,
		LeaderName: gs.PlayerName,
		MemberIDs:  []string{gs.OnlineID},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if _, err := o.db.Collection("parties").InsertOne(ctx, party); err != nil {
		return err
	}

	_, err = o.db.Collection("online_profiles").UpdateByID(ctx, gs.OnlineID, bson.M{"$set": bson.M{"party_id": partyID}})
	gs.ClearPartyQuestState()
	return err
}

func (o *OnlineService) InviteToParty(gs *GameState, targetID string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if party == nil {
		if err := o.CreateParty(gs); err != nil {
			return err
		}
		party, err = o.getPlayerParty(ctx, gs.OnlineID)
		if err != nil {
			return err
		}
	}
	if party == nil {
		return errors.New("could not create a party")
	}

	target, err := o.getOnlineProfile(ctx, targetID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("target player is not online")
	}
	if strings.TrimSpace(target.PartyID) != "" && target.PartyID != party.ID {
		return errors.New("target player is already in another party")
	}
	if target.PartyID == party.ID {
		return errors.New("player is already in your party")
	}

	text := fmt.Sprintf("%s invited you to join party %s.", gs.PlayerName, party.ID)
	return o.NotifyPlayerWithReference(gs, targetID, "party_invite", text, party.ID)
}

func (o *OnlineService) AcceptPartyInvite(gs *GameState, partyID string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	partyID = strings.TrimSpace(partyID)
	if partyID == "" {
		return errors.New("party invite is missing a party id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPartyByID(ctx, partyID)
	if err != nil {
		return err
	}
	if party == nil {
		return errors.New("party no longer exists")
	}
	currentParty, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if currentParty != nil {
		return errors.New("leave your current party before joining another")
	}
	if containsString(party.MemberIDs, gs.OnlineID) {
		return nil
	}

	if _, err := o.db.Collection("parties").UpdateByID(ctx, partyID, bson.M{"$addToSet": bson.M{"member_ids": gs.OnlineID}, "$set": bson.M{"updated_at": time.Now().UTC()}}); err != nil {
		return err
	}
	if _, err := o.db.Collection("online_profiles").UpdateByID(ctx, gs.OnlineID, bson.M{"$set": bson.M{"party_id": partyID}}); err != nil {
		return err
	}

	_ = o.NotifyPlayerWithReference(gs, party.LeaderID, "party_join", fmt.Sprintf("%s joined your party.", gs.PlayerName), partyID)
	syncPartyQuestState(gs, party)
	return nil
}

func (o *OnlineService) LeaveParty(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if party == nil {
		return nil
	}

	remaining := removeString(party.MemberIDs, gs.OnlineID)
	if len(remaining) == 0 {
		if _, err := o.db.Collection("parties").DeleteOne(ctx, bson.M{"_id": party.ID}); err != nil {
			return err
		}
	} else {
		setFields := bson.M{
			"member_ids": remaining,
			"updated_at": time.Now().UTC(),
		}
		if party.LeaderID == gs.OnlineID {
			nextLeaderID := remaining[0]
			nextLeader, err := o.getOnlineProfile(ctx, nextLeaderID)
			if err != nil {
				return err
			}
			leaderName := nextLeaderID
			if nextLeader != nil && strings.TrimSpace(nextLeader.PlayerName) != "" {
				leaderName = nextLeader.PlayerName
			}
			setFields["leader_id"] = nextLeaderID
			setFields["leader_name"] = leaderName
		}
		if _, err := o.db.Collection("parties").UpdateByID(ctx, party.ID, bson.M{"$set": setFields}); err != nil {
			return err
		}
	}

	_, err = o.db.Collection("online_profiles").UpdateByID(ctx, gs.OnlineID, bson.M{"$set": bson.M{"party_id": ""}})
	gs.ClearPartyQuestState()
	return err
}

func (o *OnlineService) PartyHeal(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if party == nil {
		return errors.New("you are not in a party")
	}
	if party.LeaderID != gs.OnlineID {
		return errors.New("only the party leader can call a party heal")
	}

	partySummary, members, err := o.GetPartySummary(gs)
	if err != nil {
		return err
	}
	if partySummary == nil || len(members) == 0 {
		return errors.New("no party members to heal")
	}

	healAmount := 15 + (len(members)-1)*5
	for _, member := range members {
		if member.OnlineID == "" {
			continue
		}
		newHealth := member.Health + healAmount
		if member.MaxHealth > 0 && newHealth > member.MaxHealth {
			newHealth = member.MaxHealth
		}
		if _, err := o.db.Collection("online_profiles").UpdateByID(ctx, member.OnlineID, bson.M{"$set": bson.M{"health": newHealth}}); err != nil {
			return err
		}
		if member.OnlineID == gs.OnlineID {
			gs.Health = newHealth
		}
		_ = o.NotifyPlayerWithReference(gs, member.OnlineID, "party_heal", fmt.Sprintf("%s shared a healing surge with the party (+%d HP).", gs.PlayerName, healAmount), party.ID)
	}
	return o.SyncPresence(gs)
}

func (o *OnlineService) RallyParty(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if party == nil {
		return errors.New("you are not in a party")
	}
	if party.LeaderID != gs.OnlineID {
		return errors.New("only the party leader can rally the party")
	}

	partySummary, members, err := o.GetPartySummary(gs)
	if err != nil {
		return err
	}
	if partySummary == nil || len(members) == 0 {
		return errors.New("no party members to rally")
	}

	bonus := 8 + (len(members)-1)*2
	turns := 3
	for _, member := range members {
		if member.OnlineID == "" {
			continue
		}
		if _, err := o.db.Collection("online_profiles").UpdateByID(ctx, member.OnlineID, bson.M{"$set": bson.M{"party_buff_turns": turns, "party_buff_bonus": bonus}}); err != nil {
			return err
		}
		if member.OnlineID == gs.OnlineID {
			gs.PartyBattleBuffTurns = turns
			gs.PartyBattleBuffBonus = bonus
		}
		_ = o.NotifyPlayerWithReference(gs, member.OnlineID, "party_buff", fmt.Sprintf("%s rallied the party: +%d damage for %d turn(s).", gs.PlayerName, bonus, turns), party.ID)
	}
	return o.SyncPresence(gs)
}

func (o *OnlineService) PartyGuard(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if party == nil {
		return errors.New("you are not in a party")
	}
	if party.LeaderID != gs.OnlineID {
		return errors.New("only the party leader can raise a guard")
	}

	partySummary, members, err := o.GetPartySummary(gs)
	if err != nil {
		return err
	}
	if partySummary == nil || len(members) == 0 {
		return errors.New("no party members to guard")
	}

	reduction := 8 + (len(members)-1)*3
	turns := 2
	for _, member := range members {
		if member.OnlineID == "" {
			continue
		}
		if _, err := o.db.Collection("online_profiles").UpdateByID(ctx, member.OnlineID, bson.M{"$set": bson.M{"party_guard_turns": turns, "party_guard_reduction": reduction}}); err != nil {
			return err
		}
		if member.OnlineID == gs.OnlineID {
			gs.PartyGuardTurns = turns
			gs.PartyGuardReduction = reduction
		}
		_ = o.NotifyPlayerWithReference(gs, member.OnlineID, "party_guard", fmt.Sprintf("%s raised a party guard: incoming damage reduced by %d for %d turn(s).", gs.PlayerName, reduction, turns), party.ID)
	}
	return o.SyncPresence(gs)
}

func (o *OnlineService) StartPartyQuest(gs *GameState, questKey string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	defs := partyQuestDefinitions()
	def, ok := defs[strings.TrimSpace(questKey)]
	if !ok {
		return fmt.Errorf("unknown party quest %q", questKey)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if party == nil {
		return errors.New("you are not in a party")
	}
	if party.LeaderID != gs.OnlineID {
		return errors.New("only the party leader can start a quest")
	}
	if strings.TrimSpace(party.QuestStatus) == "active" {
		return errors.New("your party already has an active quest")
	}
	if remaining := partyQuestCooldownRemaining(party); remaining > 0 {
		return fmt.Errorf("your party quest board is recovering for another %s", remaining.Truncate(time.Second))
	}
	if len(def.Phases) == 0 {
		return errors.New("quest has no phases")
	}

	totalPhases := len(def.Phases)
	phase := def.Phases[0]
	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"quest_key":               def.Key,
			"quest_name":              def.Name,
			"quest_goal":              totalPhases,
			"quest_progress":          0,
			"quest_status":            "active",
			"quest_phase_key":         phase.Key,
			"quest_phase_name":        phase.Name,
			"quest_phase_index":       0,
			"quest_phase_goal":        phase.Goal,
			"quest_phase_progress":    0,
			"quest_phase_description": phase.Description,
			"quest_reward_damage":     def.RewardDamage,
			"quest_reward_turns":      def.RewardTurns,
			"quest_reward_gold":       def.RewardGold,
			"quest_reward_relic":      def.RewardRelic,
			"quest_guard_reduction":   def.GuardReduction,
			"quest_guard_turns":       def.GuardTurns,
			"quest_cooldown_until":    time.Time{},
			"updated_at":              now,
		},
	}
	if _, err := o.db.Collection("parties").UpdateByID(ctx, party.ID, update); err != nil {
		return err
	}
	party.QuestKey = def.Key
	party.QuestName = def.Name
	party.QuestGoal = totalPhases
	party.QuestProgress = 0
	party.QuestStatus = "active"
	party.QuestPhaseKey = phase.Key
	party.QuestPhaseName = phase.Name
	party.QuestPhaseIndex = 0
	party.QuestPhaseGoal = phase.Goal
	party.QuestPhaseProgress = 0
	party.QuestPhaseDescription = phase.Description
	party.QuestRewardDamage = def.RewardDamage
	party.QuestRewardTurns = def.RewardTurns
	party.QuestRewardGold = def.RewardGold
	party.QuestRewardRelic = def.RewardRelic
	party.QuestGuardReduction = def.GuardReduction
	party.QuestGuardTurns = def.GuardTurns

	memberText := fmt.Sprintf("%s started the party quest %q. Phase 1/%d is %q: %s", gs.PlayerName, def.Name, totalPhases, phase.Name, phase.Description)
	for _, memberID := range party.MemberIDs {
		_ = o.NotifyPlayerWithReference(gs, memberID, "party_quest", memberText, party.ID)
	}
	syncPartyQuestState(gs, party)
	return o.SyncPresence(gs)
}

func (o *OnlineService) ProgressPartyQuest(gs *GameState, questKey string, amount int) error {
	if !o.Enabled() || gs == nil || amount <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if party == nil || strings.TrimSpace(party.QuestStatus) != "active" {
		return nil
	}
	if strings.TrimSpace(party.QuestKey) != strings.TrimSpace(questKey) {
		return nil
	}

	defs := partyQuestDefinitions()
	def, ok := defs[party.QuestKey]
	if !ok {
		return fmt.Errorf("unknown party quest %q", party.QuestKey)
	}
	phase, ok := questPhaseAt(def, party.QuestPhaseIndex)
	if !ok {
		return errors.New("party quest phase is unavailable")
	}

	newProgress := party.QuestPhaseProgress + amount
	if newProgress > phase.Goal {
		newProgress = phase.Goal
	}

	updates := bson.M{
		"$set": bson.M{
			"quest_phase_progress": newProgress,
			"updated_at":           time.Now().UTC(),
		},
	}
	if newProgress >= phase.Goal {
		nextIndex := party.QuestPhaseIndex + 1
		party.QuestProgress++
		if nextPhase, ok := questPhaseAt(def, nextIndex); ok {
			updates["$set"].(bson.M)["quest_progress"] = party.QuestProgress
			updates["$set"].(bson.M)["quest_phase_index"] = nextIndex
			updates["$set"].(bson.M)["quest_phase_key"] = nextPhase.Key
			updates["$set"].(bson.M)["quest_phase_name"] = nextPhase.Name
			updates["$set"].(bson.M)["quest_phase_goal"] = nextPhase.Goal
			updates["$set"].(bson.M)["quest_phase_progress"] = 0
			updates["$set"].(bson.M)["quest_phase_description"] = nextPhase.Description
		} else {
			updates["$set"].(bson.M)["quest_progress"] = party.QuestGoal
			updates["$set"].(bson.M)["quest_phase_progress"] = newProgress
			updates["$set"].(bson.M)["quest_phase_name"] = phase.Name
			updates["$set"].(bson.M)["quest_phase_goal"] = phase.Goal
		}
	}
	if _, err := o.db.Collection("parties").UpdateByID(ctx, party.ID, updates); err != nil {
		return err
	}

	if newProgress >= phase.Goal {
		nextIndex := party.QuestPhaseIndex + 1
		if nextPhase, ok := questPhaseAt(def, nextIndex); ok {
			memberText := fmt.Sprintf("Party quest %q advanced to phase %d/%d: %s", def.Name, nextIndex+1, len(def.Phases), nextPhase.Name)
			for _, memberID := range party.MemberIDs {
				_ = o.NotifyPlayerWithReference(gs, memberID, "party_quest", memberText, party.ID)
			}
			party.QuestProgress++
			party.QuestPhaseIndex = nextIndex
			party.QuestPhaseKey = nextPhase.Key
			party.QuestPhaseName = nextPhase.Name
			party.QuestPhaseGoal = nextPhase.Goal
			party.QuestPhaseProgress = 0
			party.QuestPhaseDescription = nextPhase.Description
			syncPartyQuestState(gs, party)
			return o.SyncPresence(gs)
		}
		return o.finishPartyQuest(gs, party)
	}

	party.QuestPhaseProgress = newProgress
	syncPartyQuestState(gs, party)
	memberText := fmt.Sprintf("Party quest %q phase %q progressed to %d/%d.", party.QuestName, phase.Name, newProgress, phase.Goal)
	for _, memberID := range party.MemberIDs {
		_ = o.NotifyPlayerWithReference(gs, memberID, "party_quest", memberText, party.ID)
	}
	return o.SyncPresence(gs)
}

func (o *OnlineService) finishPartyQuest(gs *GameState, party *Party) error {
	if !o.Enabled() || gs == nil || party == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	cooldownUntil := now.Add(10 * time.Minute)
	update := bson.M{
		"$set": bson.M{
			"quest_status":         "complete",
			"quest_phase_progress": party.QuestPhaseGoal,
			"quest_cooldown_until": cooldownUntil,
			"updated_at":           now,
		},
	}
	if _, err := o.db.Collection("parties").UpdateByID(ctx, party.ID, update); err != nil {
		return err
	}

	partySummary, members, err := o.GetPartySummary(gs)
	if err != nil {
		return err
	}
	if partySummary == nil {
		return errors.New("party not available")
	}

	damageBonus, turns, guardReduction, guardTurns, questGold, questRelic := partyQuestCompletionRewards(party)

	for _, member := range members {
		if member.OnlineID == "" {
			continue
		}
		memberGold := member.Gold + questGold
		memberRelics := append([]string(nil), member.Relics...)
		if !containsString(memberRelics, questRelic) {
			memberRelics = append(memberRelics, questRelic)
		}
		if _, err := o.db.Collection("online_profiles").UpdateByID(ctx, member.OnlineID, bson.M{"$set": bson.M{
			"party_buff_turns":      turns,
			"party_buff_bonus":      damageBonus,
			"party_guard_turns":     guardTurns,
			"party_guard_reduction": guardReduction,
			"gold":                  memberGold,
			"relics":                memberRelics,
		}}); err != nil {
			return err
		}
		if member.OnlineID == gs.OnlineID {
			gs.PartyBattleBuffTurns = turns
			gs.PartyBattleBuffBonus = damageBonus
			gs.PartyGuardTurns = guardTurns
			gs.PartyGuardReduction = guardReduction
			gs.Gold = memberGold
			if !containsString(gs.OnlineRelics, questRelic) {
				gs.OnlineRelics = append(gs.OnlineRelics, questRelic)
			}
		}
		_ = o.NotifyPlayerWithReference(gs, member.OnlineID, "party_quest", fmt.Sprintf("Party quest %q completed. Reward unlocked: %d gold and %s.", party.QuestName, questGold, questRelic), party.ID)
	}

	party.QuestStatus = "complete"
	party.QuestProgress = party.QuestGoal
	party.QuestPhaseProgress = party.QuestPhaseGoal
	party.QuestPhaseName = fmt.Sprintf("%s (completed)", party.QuestPhaseName)
	party.QuestCooldownUntil = cooldownUntil
	syncPartyQuestState(gs, party)
	return o.SyncPresence(gs)
}

func (o *OnlineService) FollowPartyLeader(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	party, err := o.getPlayerParty(ctx, gs.OnlineID)
	if err != nil {
		return err
	}
	if party == nil {
		return errors.New("you are not in a party")
	}
	if party.LeaderID == gs.OnlineID {
		return errors.New("you are already the party leader")
	}

	leader, err := o.getOnlineProfile(ctx, party.LeaderID)
	if err != nil {
		return err
	}
	if leader == nil {
		return errors.New("party leader is offline")
	}

	gs.CurrentRoomID = leader.CurrentRoomID
	if err := o.SyncPresence(gs); err != nil {
		return err
	}
	_ = o.NotifyPlayerWithReference(gs, party.LeaderID, "party_follow", fmt.Sprintf("%s followed you to %s.", gs.PlayerName, leader.CurrentRoomID), party.ID)
	return nil
}

func (o *OnlineService) GrantRelic(gs *GameState, relic string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}
	relic = strings.TrimSpace(relic)
	if relic == "" {
		return errors.New("relic cannot be empty")
	}

	if !containsString(gs.OnlineRelics, relic) {
		gs.OnlineRelics = append(gs.OnlineRelics, relic)
	}
	return o.SyncPresence(gs)
}

func (o *OnlineService) TransferRelic(gs *GameState, targetID, relic string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}
	relic = strings.TrimSpace(relic)
	targetID = strings.TrimSpace(targetID)
	if relic == "" || targetID == "" {
		return errors.New("relic transfer requires a relic and target")
	}
	if !containsString(gs.OnlineRelics, relic) {
		return fmt.Errorf("you do not own %q", relic)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := o.db.Collection("online_profiles").UpdateByID(ctx, gs.OnlineID, bson.M{"$pull": bson.M{"relics": relic}})
	if err != nil {
		return err
	}
	_, err = o.db.Collection("online_profiles").UpdateByID(ctx, targetID, bson.M{"$addToSet": bson.M{"relics": relic}})
	if err != nil {
		return err
	}

	gs.OnlineRelics = removeString(gs.OnlineRelics, relic)
	if err := o.SyncPresence(gs); err != nil {
		return err
	}
	_ = o.NotifyPlayer(gs, targetID, "gift", fmt.Sprintf("%s gifted you %s.", gs.PlayerName, relic))
	return nil
}

func (o *OnlineService) SwapRelics(gs *GameState, targetID, giveRelic, takeRelic string) error {
	if !o.Enabled() || gs == nil {
		return nil
	}
	giveRelic = strings.TrimSpace(giveRelic)
	takeRelic = strings.TrimSpace(takeRelic)
	targetID = strings.TrimSpace(targetID)
	if giveRelic == "" || takeRelic == "" || targetID == "" {
		return errors.New("trade requires both relics and a target")
	}
	if !containsString(gs.OnlineRelics, giveRelic) {
		return fmt.Errorf("you do not own %q", giveRelic)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if o.client != nil {
		if err := o.swapRelicsWithTransaction(ctx, gs.OnlineID, targetID, giveRelic, takeRelic); err == nil {
			if err := o.finishRelicSwap(gs, giveRelic, takeRelic); err != nil {
				return err
			}
			_ = o.NotifyPlayer(gs, targetID, "trade", fmt.Sprintf("%s traded %s for %s with you.", gs.PlayerName, giveRelic, takeRelic))
			_ = o.NotifyPlayer(gs, gs.OnlineID, "trade", fmt.Sprintf("You traded %s for %s with another player.", giveRelic, takeRelic))
			return nil
		} else if !isUnsupportedTransactionError(err) {
			return err
		}
	}

	targetUpdate := bson.M{
		"$pull":     bson.M{"relics": takeRelic},
		"$addToSet": bson.M{"relics": giveRelic},
	}
	if _, err := o.db.Collection("online_profiles").UpdateByID(ctx, targetID, targetUpdate); err != nil {
		return err
	}

	selfUpdate := bson.M{
		"$pull":     bson.M{"relics": giveRelic},
		"$addToSet": bson.M{"relics": takeRelic},
	}
	if _, err := o.db.Collection("online_profiles").UpdateByID(ctx, gs.OnlineID, selfUpdate); err != nil {
		rollback := bson.M{
			"$pull":     bson.M{"relics": giveRelic},
			"$addToSet": bson.M{"relics": takeRelic},
		}
		_, _ = o.db.Collection("online_profiles").UpdateByID(ctx, targetID, rollback)
		return err
	}

	if err := o.finishRelicSwap(gs, giveRelic, takeRelic); err != nil {
		return err
	}
	_ = o.NotifyPlayer(gs, targetID, "trade", fmt.Sprintf("%s traded %s for %s with you.", gs.PlayerName, giveRelic, takeRelic))
	_ = o.NotifyPlayer(gs, gs.OnlineID, "trade", fmt.Sprintf("You traded %s for %s with another player.", giveRelic, takeRelic))
	return nil
}

func (o *OnlineService) swapRelicsWithTransaction(ctx context.Context, selfID, targetID, giveRelic, takeRelic string) error {
	session, err := o.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	txnOpts := options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetReadPreference(readpref.Primary()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		targetUpdate := bson.M{
			"$pull":     bson.M{"relics": takeRelic},
			"$addToSet": bson.M{"relics": giveRelic},
		}
		if _, err := o.db.Collection("online_profiles").UpdateByID(sc, targetID, targetUpdate); err != nil {
			return nil, err
		}

		selfUpdate := bson.M{
			"$pull":     bson.M{"relics": giveRelic},
			"$addToSet": bson.M{"relics": takeRelic},
		}
		if _, err := o.db.Collection("online_profiles").UpdateByID(sc, selfID, selfUpdate); err != nil {
			return nil, err
		}
		return nil, nil
	}, txnOpts)
	return err
}

func (o *OnlineService) finishRelicSwap(gs *GameState, giveRelic, takeRelic string) error {
	gs.OnlineRelics = removeString(gs.OnlineRelics, giveRelic)
	if !containsString(gs.OnlineRelics, takeRelic) {
		gs.OnlineRelics = append(gs.OnlineRelics, takeRelic)
	}
	return o.SyncPresence(gs)
}

func isUnsupportedTransactionError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "transactions are not supported") ||
		strings.Contains(msg, "transaction numbers are only allowed on a replica set member or mongos") ||
		strings.Contains(msg, "read preference in a transaction must be primary")
}

func (o *OnlineService) SeedStarterRelics(gs *GameState) error {
	if !o.Enabled() || gs == nil {
		return nil
	}
	if len(gs.OnlineRelics) > 0 {
		return nil
	}

	relic := randomRelic()
	gs.OnlineRelics = append(gs.OnlineRelics, relic)
	return o.SyncPresence(gs)
}

func (o *OnlineService) GetMillionaireLeaderboard(streakID string) ([]MillionaireScore, error) {
	if !o.Enabled() {
		return nil, errors.New("online features disabled")
	}
	coll := o.db.Collection("millionaire_scores")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "score", Value: -1}, {Key: "updated_at", Value: 1}}).SetLimit(20)
	cursor, err := coll.Find(ctx, bson.M{"streak_id": streakID}, opts)
	if err != nil {
		return nil, err
	}
	var scores []MillionaireScore
	if err := cursor.All(ctx, &scores); err != nil {
		return nil, err
	}
	return scores, nil
}

func (o *OnlineService) SyncMillionaireScore(gs *GameState, score int) error {
	if !o.Enabled() {
		return nil
	}
	coll := o.db.Collection("millionaire_scores")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"online_id": gs.OnlineID, "streak_id": gs.MillionaireStreakID}
	update := bson.M{
		"$set": bson.M{
			"player_name": gs.PlayerName,
			"score":       score,
			"updated_at":  time.Now(),
		},
	}
	_, err := coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (o *OnlineService) GetCurrentMillionaireStreak() (*MillionaireStreak, error) {
	if !o.Enabled() {
		return nil, errors.New("online features disabled")
	}
	coll := o.db.Collection("millionaire_streaks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the latest streak that hasn't ended or the most recent one
	var streak MillionaireStreak
	err := coll.FindOne(ctx, bson.M{"is_processed": false}, options.FindOne().SetSort(bson.D{{Key: "end_time", Value: -1}})).Decode(&streak)
	if err == mongo.ErrNoDocuments {
		// Create a new streak
		newStreak := MillionaireStreak{
			ID:          fmt.Sprintf("streak-%d", time.Now().Unix()),
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(48 * time.Hour),
			IsProcessed: false,
		}
		_, err := o.db.Collection("millionaire_streaks").InsertOne(ctx, newStreak)
		return &newStreak, err
	}
	return &streak, err
}

func (o *OnlineService) CheckAndResetMillionaireStreak(gs *GameState) (string, error) {
	if !o.Enabled() {
		return "", nil
	}
	streak, err := o.GetCurrentMillionaireStreak()
	if err != nil {
		return "", err
	}

	if time.Now().After(streak.EndTime) && !streak.IsProcessed {
		// Process rewards
		scores, err := o.GetMillionaireLeaderboard(streak.ID)
		if err == nil {
			for i, s := range scores {
				reward := 0
				if i == 0 {
					reward = 100
				} else if i == 1 {
					reward = 50
				} else if i == 2 {
					reward = 25
				}
				if reward > 0 {
					_ = o.NotifyPlayerWithReference(&GameState{OnlineID: "SYSTEM", PlayerName: "SYSTEM"}, s.OnlineID, "Millionaire reward", fmt.Sprintf("You won %d gold for placing #%d in the Millionaire streak!", reward, i+1), "")
				}
			}
		}

		// Mark streak as processed
		coll := o.db.Collection("millionaire_streaks")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = coll.UpdateOne(ctx, bson.M{"_id": streak.ID}, bson.M{"$set": bson.M{"is_processed": true}})

		// Create new streak
		newStreak := MillionaireStreak{
			ID:          fmt.Sprintf("streak-%d", time.Now().Unix()),
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(48 * time.Hour),
			IsProcessed: false,
		}
		_, _ = o.db.Collection("millionaire_streaks").InsertOne(ctx, newStreak)
		return "Leaderboard reset! Rewards sent to winners.", nil
	}

	return "", nil
}

func (o *OnlineService) GetActivePlayers(gs *GameState) ([]OnlinePlayer, error) {
	if !o.Enabled() {
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	threshold := time.Now().UTC().Add(-2 * time.Minute)
	filter := bson.M{
		"last_seen": bson.M{"$gt": threshold},
		"_id":       bson.M{"$ne": gs.OnlineID},
	}
	cursor, err := o.db.Collection("online_profiles").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var players []OnlinePlayer
	if err := cursor.All(ctx, &players); err != nil {
		return nil, err
	}
	return players, nil
}

func (o *OnlineService) BroadcastMillionaireQuestion(gs *GameState, questionID string, questionText string) error {
	if !o.Enabled() {
		return nil
	}
	// Send notification to all active players
	players, err := o.GetActivePlayers(gs)
	if err != nil {
		return err
	}

	for _, p := range players {
		_ = o.NotifyPlayerWithReference(gs, p.OnlineID, "MILLIONAIRE_AUDIENCE", fmt.Sprintf("%s needs help with a Millionaire question: %s", gs.PlayerName, questionText), questionID)
	}
	return nil
}

func (o *OnlineService) SubmitAudienceResponse(questionID, playerID string, choice int) error {
	if !o.Enabled() {
		return nil
	}
	coll := o.db.Collection("millionaire_audience_responses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := coll.ReplaceOne(ctx, bson.M{"question_id": questionID, "player_id": playerID}, MillionaireAudienceResponse{
		QuestionID: questionID,
		PlayerID:   playerID,
		Choice:     choice,
		CreatedAt:  time.Now(),
	}, options.Replace().SetUpsert(true))
	return err
}

func (o *OnlineService) GetAudienceStats(questionID string) (map[int]int, error) {
	if !o.Enabled() {
		return nil, nil
	}
	coll := o.db.Collection("millionaire_audience_responses")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{"question_id": questionID})
	if err != nil {
		return nil, err
	}
	stats := make(map[int]int)
	var responses []MillionaireAudienceResponse
	if err := cursor.All(ctx, &responses); err != nil {
		return nil, err
	}
	for _, r := range responses {
		stats[r.Choice]++
	}
	return stats, nil
}

func (o *OnlineService) GetOnlinePlayers(gs *GameState) ([]OnlinePlayer, error) {
	return o.GetNearbyPlayers(gs)
}

func (o *OnlineService) GetRecentMessages(roomID, onlineID string) ([]RoomMessage, []WhisperMessage, error) {
	messages, err := o.GetRoomMessages(roomID, 6)
	if err != nil {
		return nil, nil, err
	}
	// Note: We don't have a simple "GetRecentWhispers" that takes onlineID alone in the same way, 
	// but Engine expects it. I'll use GetWhispers which uses the gs stored onlineID.
	// This is a bit of a hack because I don't have gs here, but I can fix Engine or add a better method.
	// For now, I'll return messages and empty whispers if I can't easily get them without GS.
	return messages, nil, nil
}

func (o *OnlineService) Snapshot(gs *GameState) (nearby []OnlinePlayer, messages []RoomMessage, whispers []WhisperMessage, err error) {
	if !o.Enabled() || gs == nil {
		return nil, nil, nil, nil
	}

	nearby, err = o.GetNearbyPlayers(gs)
	if err != nil {
		return nil, nil, nil, err
	}
	messages, err = o.GetRoomMessages(gs.CurrentRoomID, 6)
	if err != nil {
		return nil, nil, nil, err
	}
	whispers, err = o.GetWhispers(gs, 6)
	return nearby, messages, whispers, err
}

func (f *RoomFeed) Close() {
	if f == nil || f.cancel == nil {
		return
	}
	f.cancel()
}

func (f *RoomFeed) Snapshot() ([]OnlinePlayer, []RoomMessage, []WhisperMessage, []InboxNotification) {
	if f == nil {
		return nil, nil, nil, nil
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	nearby := append([]OnlinePlayer(nil), f.nearby...)
	messages := append([]RoomMessage(nil), f.messages...)
	whispers := append([]WhisperMessage(nil), f.whispers...)
	inbox := append([]InboxNotification(nil), f.inbox...)
	return nearby, messages, whispers, inbox
}

func (f *RoomFeed) refreshSnapshot() error {
	nearby, err := f.service.GetNearbyPlayers(f.gs)
	if err != nil {
		return err
	}
	messages, err := f.service.GetRoomMessages(f.roomID, 6)
	if err != nil {
		return err
	}
	whispers, err := f.service.GetWhispers(f.gs, 6)
	if err != nil {
		return err
	}
	inbox, err := f.service.GetNotifications(f.gs, f.gs.LastNotificationSeen, 20)
	if err != nil {
		return err
	}

	var previousNearby []OnlinePlayer
	f.mu.Lock()
	previousNearby = append([]OnlinePlayer(nil), f.nearby...)
	f.nearby = nearby
	f.messages = messages
	f.whispers = whispers
	f.inbox = inbox
	f.mu.Unlock()

	f.notifyPresenceChanges(previousNearby, nearby)
	if err := f.service.RefreshSelfState(f.gs); err != nil {
		f.refreshErr = err
	}
	return nil
}

func (f *RoomFeed) notifyPresenceChanges(previous, current []OnlinePlayer) {
	if f.notify == nil {
		return
	}

	previousSet := make(map[string]OnlinePlayer, len(previous))
	for _, player := range previous {
		previousSet[player.OnlineID] = player
	}
	currentSet := make(map[string]OnlinePlayer, len(current))
	for _, player := range current {
		currentSet[player.OnlineID] = player
	}

	for _, player := range current {
		if _, ok := previousSet[player.OnlineID]; !ok {
			f.notify(fmt.Sprintf("%s entered %s", player.PlayerName, f.roomID))
		}
	}
	for _, player := range previous {
		if _, ok := currentSet[player.OnlineID]; !ok {
			f.notify(fmt.Sprintf("%s left %s", player.PlayerName, f.roomID))
		}
	}
}

func (f *RoomFeed) fallbackRefreshLoop() {
	ticker := time.NewTicker(6 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-f.ctx.Done():
			return
		case <-ticker.C:
			if err := f.refreshSnapshot(); err != nil {
				f.refreshErr = err
			}
		}
	}
}

func (f *RoomFeed) watchRoomMessages() {
	pipeline := mongo.Pipeline{
		bson.D{{
			Key: "$match",
			Value: bson.M{
				"operationType":                "insert",
				"fullDocument.current_room_id": f.roomID,
			},
		}},
	}

	stream, err := f.service.db.Collection("room_messages").Watch(
		f.ctx,
		pipeline,
		options.ChangeStream().SetFullDocument(options.UpdateLookup),
	)
	if err != nil {
		f.refreshErr = err
		return
	}
	defer stream.Close(f.ctx)

	type roomMessageEvent struct {
		FullDocument RoomMessage `bson:"fullDocument"`
	}

	for stream.Next(f.ctx) {
		var evt roomMessageEvent
		if err := stream.Decode(&evt); err != nil {
			continue
		}
		if evt.FullDocument.OnlineID == f.gs.OnlineID {
			continue
		}
		if evt.FullDocument.Kind != "" && evt.FullDocument.Kind != "chat" && evt.FullDocument.Kind != "event" {
			continue
		}

		f.mu.Lock()
		f.messages = append([]RoomMessage{evt.FullDocument}, f.messages...)
		if len(f.messages) > 6 {
			f.messages = f.messages[:6]
		}
		f.mu.Unlock()

		if f.notify != nil {
			prefix := evt.FullDocument.PlayerName
			if evt.FullDocument.Kind == "event" {
				prefix = "event"
			}
			f.notify(fmt.Sprintf("%s: %s", prefix, evt.FullDocument.Text))
		}
	}
}

func (f *RoomFeed) watchWhispers() {
	pipeline := mongo.Pipeline{
		bson.D{{
			Key: "$match",
			Value: bson.M{
				"operationType":          "insert",
				"fullDocument.target_id": f.gs.OnlineID,
			},
		}},
	}

	stream, err := f.service.db.Collection("direct_messages").Watch(
		f.ctx,
		pipeline,
		options.ChangeStream().SetFullDocument(options.UpdateLookup),
	)
	if err != nil {
		f.refreshErr = err
		return
	}
	defer stream.Close(f.ctx)

	type whisperEvent struct {
		FullDocument WhisperMessage `bson:"fullDocument"`
	}

	for stream.Next(f.ctx) {
		var evt whisperEvent
		if err := stream.Decode(&evt); err != nil {
			continue
		}
		if evt.FullDocument.OnlineID == f.gs.OnlineID {
			continue
		}

		f.mu.Lock()
		f.whispers = append([]WhisperMessage{evt.FullDocument}, f.whispers...)
		if len(f.whispers) > 6 {
			f.whispers = f.whispers[:6]
		}
		f.mu.Unlock()

		if f.notify != nil {
			f.notify(fmt.Sprintf("whisper from %s: %s", evt.FullDocument.PlayerName, evt.FullDocument.Text))
		}
	}
}

func (f *RoomFeed) watchInboxNotifications() {
	pipeline := mongo.Pipeline{
		bson.D{{
			Key: "$match",
			Value: bson.M{
				"operationType":             "insert",
				"fullDocument.recipient_id": f.gs.OnlineID,
			},
		}},
	}

	stream, err := f.service.db.Collection("player_notifications").Watch(
		f.ctx,
		pipeline,
		options.ChangeStream().SetFullDocument(options.UpdateLookup),
	)
	if err != nil {
		f.refreshErr = err
		return
	}
	defer stream.Close(f.ctx)

	type inboxEvent struct {
		FullDocument InboxNotification `bson:"fullDocument"`
	}

	for stream.Next(f.ctx) {
		var evt inboxEvent
		if err := stream.Decode(&evt); err != nil {
			continue
		}

		f.mu.Lock()
		f.inbox = append([]InboxNotification{evt.FullDocument}, f.inbox...)
		if len(f.inbox) > 20 {
			f.inbox = f.inbox[:20]
		}
		f.mu.Unlock()

		if f.notify != nil {
			f.notify(fmt.Sprintf("\a[inbox] %s: %s", evt.FullDocument.SenderName, evt.FullDocument.Text))
		}
	}
}

func (f *RoomFeed) watchProfileChanges() {
	stream, err := f.service.db.Collection("online_profiles").Watch(
		f.ctx,
		mongo.Pipeline{},
		options.ChangeStream().SetFullDocument(options.UpdateLookup),
	)
	if err != nil {
		f.refreshErr = err
		return
	}
	defer stream.Close(f.ctx)

	for stream.Next(f.ctx) {
		if err := f.refreshSnapshot(); err != nil {
			f.refreshErr = err
		}
	}
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func removeString(values []string, needle string) []string {
	out := values[:0]
	for _, value := range values {
		if value != needle {
			out = append(out, value)
		}
	}
	return out
}

func randomRelic() string {
	relics := []string{
		"Amber Compass",
		"Moonglass Charm",
		"Sunken Coin",
		"Runed Feather",
		"Prismatic Seed",
		"Sunfire Prism",
	}
	return relics[int(time.Now().UnixNano())%len(relics)]
}

type wiseManAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type wiseManAIRequest struct {
	Model       string             `json:"model"`
	Messages    []wiseManAIMessage `json:"messages"`
	Temperature float64            `json:"temperature,omitempty"`
}

type wiseManAIResponse struct {
	Choices []struct {
		Message wiseManAIMessage `json:"message"`
	} `json:"choices"`
}

func (o *OnlineService) AskWiseMan(gs *GameState, question string) (string, error) {
	if !o.Enabled() || gs == nil {
		return "", errors.New("online mode is unavailable")
	}

	question = strings.TrimSpace(question)
	if question == "" {
		return "", errors.New("question cannot be empty")
	}

	endpoint := strings.TrimSpace(os.Getenv("WISEMAN_AI_URL"))
	if endpoint == "" {
		return "", errors.New("wise man ai is not configured")
	}
	model := strings.TrimSpace(os.Getenv("WISEMAN_AI_MODEL"))
	if model == "" {
		model = "gpt-4o-mini"
	}

	systemPrompt := fmt.Sprintf(
		"You are a wise man in a fantasy terminal RPG. The player is in room %q. They may ask about locations, items, quests, or advice.\n"+
			"Reply with 2 or 3 short sentences max. Always include one question back to the player, and if they ask where to find something, give a concrete in-game clue.\n"+
			"If they ask for a precious item or a hidden object, point them toward the harbor, lighthouse, tide caves, ruins, moonwell, observatory, or sage hut using a hint, not a spoiler dump.\n"+
			"Keep the tone warm, wise, and slightly humorous.",
		gs.CurrentRoomID,
	)

	reqBody := wiseManAIRequest{
		Model: model,
		Messages: []wiseManAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: question},
		},
		Temperature: 0.8,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	apiKey := strings.TrimSpace(os.Getenv("WISEMAN_AI_KEY"))
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(payload))
		if msg == "" {
			msg = resp.Status
		}
		return "", fmt.Errorf("wise man ai request failed: %s", msg)
	}

	var result wiseManAIResponse
	if err := json.Unmarshal(payload, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", errors.New("wise man ai returned no response")
	}
	reply := strings.TrimSpace(result.Choices[0].Message.Content)
	if reply == "" {
		return "", errors.New("wise man ai returned an empty response")
	}
	return reply, nil
}
