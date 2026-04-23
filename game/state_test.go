package game

import (
	"testing"
	"time"
)

func TestNewGameStateDefaults(t *testing.T) {
	gs := NewGameState(3, "Nova")

	if gs.ID != 3 {
		t.Fatalf("expected slot 3, got %d", gs.ID)
	}
	if gs.PlayerName != "Nova" {
		t.Fatalf("expected player name Nova, got %q", gs.PlayerName)
	}
	if gs.CurrentRoomID != "village" {
		t.Fatalf("expected starting room village, got %q", gs.CurrentRoomID)
	}
	if gs.Health != 100 || gs.MaxHealth != 100 {
		t.Fatalf("unexpected health defaults: %+v", gs)
	}
	if gs.Level != 1 || gs.RequiredXP != 100 {
		t.Fatalf("unexpected progression defaults: %+v", gs)
	}
}

func TestResetRunPreservesSlotAndTutorialFlag(t *testing.T) {
	gs := NewGameState(2, "Hero")
	gs.HasSeenTutorial = true
	gs.World = CreateWorld()
	gs.Gold = 999
	gs.Health = 1
	gs.IsGameOver = true
	gs.PartyQuestKey = "monster-hunt"
	gs.PartyQuestName = "Monster Hunt"
	gs.PartyQuestStatus = "active"
	gs.PartyQuestPhase = "Trail the Beasts"
	gs.PartyQuestPhaseIndex = 1
	gs.PartyQuestPhaseGoal = 2
	gs.PartyQuestPhaseProgress = 1

	gs.ResetRun()

	if gs.ID != 2 || gs.PlayerName != "Hero" {
		t.Fatalf("identity changed after reset: %+v", gs)
	}
	if !gs.HasSeenTutorial {
		t.Fatal("expected tutorial flag to be preserved")
	}
	if gs.Health != 100 || gs.Gold != 50 || gs.CurrentRoomID != "village" {
		t.Fatalf("reset did not restore defaults: %+v", gs)
	}
	if gs.PartyQuestKey != "" || gs.PartyQuestStatus != "" || gs.PartyQuestPhase != "" {
		t.Fatalf("reset did not clear quest state: %+v", gs)
	}
	if gs.World == nil {
		t.Fatal("expected world to be preserved")
	}
}

func TestCheckLevelUp(t *testing.T) {
	gs := NewGameState(1, "Tester")
	gs.Experience = 100

	msg := gs.CheckLevelUp()

	if gs.Level != 2 {
		t.Fatalf("expected level 2, got %d", gs.Level)
	}
	if gs.Health != gs.MaxHealth {
		t.Fatalf("expected full heal on level up, got %d/%d", gs.Health, gs.MaxHealth)
	}
	if msg == "" {
		t.Fatal("expected level up message")
	}
}

func TestTransition(t *testing.T) {
	gs := NewGameState(1, "Tester")
	gs.World = CreateWorld()

	if err := gs.Transition("forest"); err != "" {
		t.Fatalf("unexpected transition error: %s", err)
	}
	if gs.CurrentRoomID != "forest" {
		t.Fatalf("expected forest, got %q", gs.CurrentRoomID)
	}
	if err := gs.Transition("missing"); err == "" {
		t.Fatal("expected error for missing room")
	}
}

func TestPartyQuestDefinitionsIncludeBranches(t *testing.T) {
	defs := partyQuestDefinitions()

	monster, ok := defs["monster-hunt"]
	if !ok {
		t.Fatal("expected monster-hunt quest to exist")
	}
	if len(monster.Phases) < 3 {
		t.Fatalf("expected monster-hunt to have multiple phases, got %d", len(monster.Phases))
	}
	if monster.RewardRelic == "" || monster.RewardGold <= 0 {
		t.Fatalf("expected monster-hunt reward metadata, got %+v", monster)
	}

	voidQuest, ok := defs["void-expedition"]
	if !ok {
		t.Fatal("expected void-expedition quest to exist")
	}
	if len(voidQuest.Phases) < 3 {
		t.Fatalf("expected void-expedition to have multiple phases, got %d", len(voidQuest.Phases))
	}
	if voidQuest.RewardRelic == "" || voidQuest.RewardGold <= 0 {
		t.Fatalf("expected void-expedition reward metadata, got %+v", voidQuest)
	}

	frostQuest, ok := defs["frost-pact"]
	if !ok {
		t.Fatal("expected frost-pact quest to exist")
	}
	if len(frostQuest.Phases) < 3 {
		t.Fatalf("expected frost-pact to have multiple phases, got %d", len(frostQuest.Phases))
	}
	if frostQuest.RewardRelic == "" || frostQuest.RewardGold <= 0 {
		t.Fatalf("expected frost-pact reward metadata, got %+v", frostQuest)
	}

	sunQuest, ok := defs["sun-covenant"]
	if !ok {
		t.Fatal("expected sun-covenant quest to exist")
	}
	if len(sunQuest.Phases) < 3 {
		t.Fatalf("expected sun-covenant to have multiple phases, got %d", len(sunQuest.Phases))
	}
	if sunQuest.RewardRelic == "" || sunQuest.RewardGold <= 0 {
		t.Fatalf("expected sun-covenant reward metadata, got %+v", sunQuest)
	}
}

func TestPartyQuestProgressSummary(t *testing.T) {
	phaseName, progress, goal := partyQuestProgressSummary("Monster Hunt", "Trail the Beasts", 1, 3)
	if phaseName != "Trail the Beasts" || progress != 1 || goal != 3 {
		t.Fatalf("unexpected quest summary: %q %d/%d", phaseName, progress, goal)
	}

	fallbackName, fallbackProgress, fallbackGoal := partyQuestProgressSummary("Void Expedition", "", 2, 4)
	if fallbackName != "Void Expedition" || fallbackProgress != 2 || fallbackGoal != 4 {
		t.Fatalf("expected fallback to quest name, got %q %d/%d", fallbackName, fallbackProgress, fallbackGoal)
	}
}

func TestPartyQuestCompletionRewardsDefaultToSafeValues(t *testing.T) {
	damageBonus, turns, guardReduction, guardTurns, gold, relic := partyQuestCompletionRewards(nil)
	if damageBonus != 10 || turns != 3 || guardReduction != 8 || guardTurns != 2 || gold != 150 || relic == "" {
		t.Fatalf("unexpected default completion rewards: %d %d %d %d %d %q", damageBonus, turns, guardReduction, guardTurns, gold, relic)
	}
}

func TestPartyQuestCooldownRemaining(t *testing.T) {
	if remaining := partyQuestCooldownRemaining(nil); remaining != 0 {
		t.Fatalf("expected zero cooldown for nil party, got %s", remaining)
	}

	party := &Party{QuestCooldownUntil: time.Now().Add(2 * time.Minute)}
	if remaining := partyQuestCooldownRemaining(party); remaining <= 0 {
		t.Fatalf("expected positive cooldown, got %s", remaining)
	}
}
