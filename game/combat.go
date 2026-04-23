package game

import (
	"fmt"
	"strings"
)

const (
	CombatStatusNone   = ""
	CombatStatusPoison = "poison"
	CombatStatusBurn   = "burn"
	CombatStatusFreeze = "freeze"
	CombatStatusStun   = "stun"
)

func AttackDamage(gs *GameState, base int) int {
	if gs == nil {
		return base
	}

	damage := base + gs.BaseDamage
	if gs.PartySupport > 0 {
		damage += gs.PartySupport * 5
	}
	if gs.PartyBattleBuffTurns > 0 {
		damage += gs.PartyBattleBuffBonus
	}
	damage += questRelicDamageBonus(gs)
	if gs.HasSteelSword {
		damage += 20
	}
	if gs.HasForsakenBlade {
		damage += 40
	}
	return damage
}

func Heal(gs *GameState, amount int) int {
	if gs == nil || amount <= 0 {
		return 0
	}

	before := gs.Health
	gs.Health += amount
	if gs.Health > gs.MaxHealth {
		gs.Health = gs.MaxHealth
	}
	return gs.Health - before
}

func AwardVictory(gs *GameState, xp, gold int) string {
	if gs == nil {
		return ""
	}

	gs.Experience += xp
	gs.Gold += gold

	switch {
	case xp > 0 && gold > 0:
		return fmt.Sprintf("You earned %d XP and %d gold.", xp, gold)
	case xp > 0:
		return fmt.Sprintf("You earned %d XP.", xp)
	case gold > 0:
		return fmt.Sprintf("You earned %d gold.", gold)
	default:
		return ""
	}
}

func BeginCombat(gs *GameState, monsterHealth int) {
	if gs == nil {
		return
	}
	gs.MonsterHealth = monsterHealth
	gs.PlayerCombatStatus = ""
	gs.PlayerCombatTurns = 0
	gs.MonsterCombatStatus = ""
	gs.MonsterCombatTurns = 0
}

func ApplyPlayerCombatStatus(gs *GameState, status string, turns int) string {
	return applyCombatStatus(gs, true, status, turns)
}

func ApplyMonsterCombatStatus(gs *GameState, status string, turns int) string {
	return applyCombatStatus(gs, false, status, turns)
}

func applyCombatStatus(gs *GameState, player bool, status string, turns int) string {
	if gs == nil {
		return ""
	}
	status = normalizeCombatStatus(status)
	if status == CombatStatusNone || turns <= 0 {
		return ""
	}

	if player {
		gs.PlayerCombatStatus = status
		gs.PlayerCombatTurns = turns
	} else {
		gs.MonsterCombatStatus = status
		gs.MonsterCombatTurns = turns
	}
	return fmt.Sprintf("%s is affected by %s for %d turn(s).", combatTargetLabel(player), status, turns)
}

func TickCombatStatuses(gs *GameState) []string {
	if gs == nil {
		return nil
	}

	var messages []string
	if msg := tickPartyBattleBuff(gs); msg != "" {
		messages = append(messages, msg)
	}
	if msg := tickPartyGuard(gs); msg != "" {
		messages = append(messages, msg)
	}
	if msg := tickCombatStatus(gs, true); msg != "" {
		messages = append(messages, msg)
	}
	if msg := tickCombatStatus(gs, false); msg != "" {
		messages = append(messages, msg)
	}
	return messages
}

func tickCombatStatus(gs *GameState, player bool) string {
	var status *string
	var turns *int
	target := "you"
	if player {
		status = &gs.PlayerCombatStatus
		turns = &gs.PlayerCombatTurns
	} else {
		status = &gs.MonsterCombatStatus
		turns = &gs.MonsterCombatTurns
		target = "the monster"
	}

	if status == nil || turns == nil || *status == CombatStatusNone || *turns <= 0 {
		return ""
	}

	message := ""
	switch *status {
	case CombatStatusPoison:
		amount := 6
		if player {
			gs.Health -= amount
			message = fmt.Sprintf("%s suffers %d poison damage.", target, amount)
		} else {
			gs.MonsterHealth -= amount
			message = fmt.Sprintf("%s suffers %d poison damage.", target, amount)
		}
	case CombatStatusBurn:
		amount := 10
		if player {
			gs.Health -= amount
		} else {
			gs.MonsterHealth -= amount
		}
		message = fmt.Sprintf("%s suffers %d burn damage.", target, amount)
	case CombatStatusFreeze:
		message = fmt.Sprintf("%s is frozen and loses momentum.", target)
	case CombatStatusStun:
		message = fmt.Sprintf("%s is stunned and struggles to react.", target)
	default:
		message = ""
	}

	*turns--
	if *turns <= 0 {
		*status = CombatStatusNone
	}
	return message
}

func MonsterCounterAttack(gs *GameState, monsterName string, damage int, inflictedStatus string, inflictedTurns int) string {
	if gs == nil {
		return ""
	}
	if gs.MonsterCombatStatus == CombatStatusFreeze || gs.MonsterCombatStatus == CombatStatusStun {
		return fmt.Sprintf("%s cannot counter this turn.", monsterName)
	}
	if gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun {
		return fmt.Sprintf("%s cannot counter this turn.", monsterName)
	}

	actualDamage := damage
	msg := fmt.Sprintf("%s hits you for %d damage.", monsterName, damage)
	if gs.PartyGuardTurns > 0 && gs.PartyGuardReduction > 0 {
		reduced := gs.PartyGuardReduction
		if reduced > actualDamage {
			reduced = actualDamage
		}
		actualDamage -= reduced
		msg += fmt.Sprintf(" Your party guard reduces the hit by %d.", reduced)
	}
	if actualDamage < 0 {
		actualDamage = 0
	}
	gs.Health -= actualDamage
	if actualDamage != damage {
		msg = fmt.Sprintf("%s hits you for %d damage, reduced to %d by your party.", monsterName, damage, actualDamage)
	}
	if normalized := normalizeCombatStatus(inflictedStatus); normalized != CombatStatusNone && inflictedTurns > 0 {
		gs.PlayerCombatStatus = normalized
		gs.PlayerCombatTurns = inflictedTurns
		msg += fmt.Sprintf(" You are afflicted with %s.", normalized)
	}
	if gs.Health <= 0 {
		gs.IsGameOver = true
	}
	return msg
}

func ClearCombatStatus(gs *GameState) {
	if gs == nil {
		return
	}
	gs.PlayerCombatStatus = CombatStatusNone
	gs.PlayerCombatTurns = 0
	gs.MonsterCombatStatus = CombatStatusNone
	gs.MonsterCombatTurns = 0
}

func ApplyPartyBattleBuff(gs *GameState, turns, bonus int) string {
	if gs == nil || turns <= 0 || bonus <= 0 {
		return ""
	}

	gs.PartyBattleBuffTurns = turns
	gs.PartyBattleBuffBonus = bonus
	return fmt.Sprintf("Your party gains a battle buff: +%d damage for %d turn(s).", bonus, turns)
}

func ApplyPartyGuard(gs *GameState, turns, reduction int) string {
	if gs == nil || turns <= 0 || reduction <= 0 {
		return ""
	}

	gs.PartyGuardTurns = turns
	gs.PartyGuardReduction = reduction
	return fmt.Sprintf("Your party raises a guard: incoming damage reduced by %d for %d turn(s).", reduction, turns)
}

func questRelicDamageBonus(gs *GameState) int {
	if gs == nil {
		return 0
	}

	roomID := strings.TrimSpace(gs.CurrentRoomID)
	questKey := strings.TrimSpace(gs.PartyQuestKey)
	bonus := 0

	if containsString(gs.OnlineRelics, "Hunter's Claw") && (questKey == "monster-hunt" || strings.HasPrefix(roomID, "party_hunt") || roomID == "zoo") {
		bonus += 8
	}
	if containsString(gs.OnlineRelics, "Null Compass") && (questKey == "void-expedition" || strings.HasPrefix(roomID, "party_void") || roomID == "binary_sea") {
		bonus += 8
	}
	if containsString(gs.OnlineRelics, "Wyrmfang Sigil") && (questKey == "frost-pact" || strings.HasPrefix(roomID, "party_frost") || roomID == "frost_giant") {
		bonus += 8
	}
	if containsString(gs.OnlineRelics, "Sunfire Prism") && (questKey == "sun-covenant" || strings.HasPrefix(roomID, "party_sun") || roomID == "desert" || roomID == "sphinx") {
		bonus += 10
	}

	return bonus
}

func tickPartyBattleBuff(gs *GameState) string {
	if gs == nil || gs.PartyBattleBuffTurns <= 0 || gs.PartyBattleBuffBonus <= 0 {
		return ""
	}

	gs.PartyBattleBuffTurns--
	if gs.PartyBattleBuffTurns <= 0 {
		bonus := gs.PartyBattleBuffBonus
		gs.PartyBattleBuffBonus = 0
		return fmt.Sprintf("Your party battle buff fades after granting +%d damage.", bonus)
	}
	return fmt.Sprintf("Your party battle buff remains active for %d turn(s).", gs.PartyBattleBuffTurns)
}

func tickPartyGuard(gs *GameState) string {
	if gs == nil || gs.PartyGuardTurns <= 0 || gs.PartyGuardReduction <= 0 {
		return ""
	}

	gs.PartyGuardTurns--
	if gs.PartyGuardTurns <= 0 {
		reduction := gs.PartyGuardReduction
		gs.PartyGuardReduction = 0
		return fmt.Sprintf("Your party guard fades after shielding %d damage.", reduction)
	}
	return fmt.Sprintf("Your party guard remains active for %d turn(s).", gs.PartyGuardTurns)
}

func combatTargetLabel(player bool) string {
	if player {
		return "You"
	}
	return "The monster"
}

func normalizeCombatStatus(status string) string {
	switch status {
	case CombatStatusPoison, CombatStatusBurn, CombatStatusFreeze, CombatStatusStun:
		return status
	default:
		return CombatStatusNone
	}
}

func joinCombatMessages(parts ...string) string {
	var out []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return strings.Join(out, "\n")
}
