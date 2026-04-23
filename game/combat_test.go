package game

import "testing"

func TestAttackDamageIncludesPartySupportAndBuffs(t *testing.T) {
	gs := NewGameState(1, "Hero")
	gs.BaseDamage = 10
	gs.PartySupport = 2
	gs.HasSteelSword = true
	gs.PartyBattleBuffTurns = 2
	gs.PartyBattleBuffBonus = 8

	got := AttackDamage(gs, 20)
	want := 20 + 10 + (2 * 5) + 8 + 20
	if got != want {
		t.Fatalf("expected %d damage, got %d", want, got)
	}
}

func TestAttackDamageIncludesQuestRelicBonus(t *testing.T) {
	gs := NewGameState(1, "Hero")
	gs.CurrentRoomID = "party_sun_crown"
	gs.OnlineRelics = []string{"Sunfire Prism"}

	got := AttackDamage(gs, 10)
	want := 10 + gs.BaseDamage + 10
	if got != want {
		t.Fatalf("expected %d damage, got %d", want, got)
	}
}

func TestMonsterCounterAttackRespectsPartyGuard(t *testing.T) {
	gs := NewGameState(1, "Hero")
	gs.Health = 100
	gs.PartyGuardTurns = 2
	gs.PartyGuardReduction = 7

	msg := MonsterCounterAttack(gs, "Beast", 10, CombatStatusPoison, 1)
	if gs.Health != 97 {
		t.Fatalf("expected guard to reduce damage to 3, got %d health", gs.Health)
	}
	if msg == "" {
		t.Fatal("expected counterattack message")
	}
	if gs.PlayerCombatStatus != CombatStatusPoison || gs.PlayerCombatTurns != 1 {
		t.Fatalf("expected inflicted status to be applied, got %+v", gs)
	}
}

func TestApplyPartyBuffHelpers(t *testing.T) {
	gs := NewGameState(1, "Hero")

	if msg := ApplyPartyBattleBuff(gs, 3, 9); msg == "" {
		t.Fatal("expected party battle buff message")
	}
	if gs.PartyBattleBuffTurns != 3 || gs.PartyBattleBuffBonus != 9 {
		t.Fatalf("unexpected battle buff state: %+v", gs)
	}

	if msg := ApplyPartyGuard(gs, 2, 6); msg == "" {
		t.Fatal("expected party guard message")
	}
	if gs.PartyGuardTurns != 2 || gs.PartyGuardReduction != 6 {
		t.Fatalf("unexpected guard state: %+v", gs)
	}

	msgs := TickCombatStatuses(gs)
	if len(msgs) == 0 {
		t.Fatal("expected tick messages for active party buffs")
	}
}
