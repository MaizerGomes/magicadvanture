package game

import (
	"fmt"
)

func CreateWorld() map[string]*Room {
	world := make(map[string]*Room)

	// --- VILLAGE & HUB ---
	world["village"] = &Room{
		ID:          "village",
		Description: "You are in the Village of Oakhaven. It's peaceful here. You can see a shop and a garden.",
		Actions: []Action{
			{ID: "go_forest", Description: "Go to the Forest", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "You head east." }},
			{ID: "go_shop", Description: "Enter Shop", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("shop"); return "Welcome!" }},
			{ID: "go_garden", Description: "Visit Zen Garden", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("garden"); return "Birds are chirping." }},
			{ID: "go_sage_hut", Description: "Visit the Wise Man", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("sage_hut"); return "The path grows quiet and thoughtful." }},
			{ID: "go_knowledge_hall", Description: "Go to Knowledge Hall", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("knowledge_hall"); return "You walk towards a majestic marble building." }},
		},
	}

	world["knowledge_hall"] = &Room{
		ID:          "knowledge_hall",
		Description: "You are in the Knowledge Hall. Large bookshelves line the walls. To the North is the Hall of Wisdom, but it is guarded.",
		Actions: []Action{
			{ID: "go_village", Description: "Back to Village", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("village"); return "Heading back." }},
			{ID: "go_guessing", Description: "Go to Guessing Table", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("guessing_game"); return "You see an old man with cups." }},
			{ID: "go_millionaire", Description: "Go to Millionaire Hall", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("millionaire_hall"); return "The lights get brighter and you hear dramatic music." }},
			{ID: "enter_wisdom_hall", Description: "Enter Hall of Wisdom (Challenge)", Type: CustomAction, Result: func(gs *GameState) string {
				// We'll handle the challenge in the engine for this CustomAction
				return "CHALLENGE_REQUIRED:wisdom_hall"
			}},
		},
	}

	world["millionaire_hall"] = &Room{
		ID:          "millionaire_hall",
		Description: "The Hall of Millionaires. Spotlights shine on a central stage with two chairs. A giant scoreboard hangs on the wall.",
		Actions: []Action{
			{ID: "play_millionaire", Description: "Start Millionaire Game", Type: CustomAction, Result: func(gs *GameState) string {
				return "CHALLENGE_REQUIRED:millionaire"
			}},
			{ID: "back_knowledge", Description: "Back to Knowledge Hall", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("knowledge_hall"); return "Leaving the stage." }},
		},
	}

	world["wisdom_hall"] = &Room{
		ID:          "wisdom_hall",
		Description: "The Hall of Wisdom. A golden pedestal stands in the center. You feel much smarter already!",
		Actions: []Action{
			{ID: "meditate", Description: "Meditate (+50 SP)", Result: func(gs *GameState) string {
				gs.SkillPoints += 50
				return "You meditated and gained 50 Skill Points!"
			}},
			{ID: "leave_wisdom", Description: "Leave to Knowledge Hall", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("knowledge_hall"); return "Stepping out." }},
		},
	}

	world["guessing_game"] = &Room{
		ID:          "guessing_game",
		Description: "A small table with three cups. An old man challenges you to guess where the pearl is.",
		Actions: []Action{
			{ID: "guess_1", Description: "Lift Cup 1 (10 gold)", Result: func(gs *GameState) string { return "MINI_GAME:guessing_game:1" }},
			{ID: "guess_2", Description: "Lift Cup 2 (10 gold)", Result: func(gs *GameState) string { return "MINI_GAME:guessing_game:2" }},
			{ID: "guess_3", Description: "Lift Cup 3 (10 gold)", Result: func(gs *GameState) string { return "MINI_GAME:guessing_game:3" }},
			{ID: "leave_guess", Description: "Back to Knowledge Hall", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("knowledge_hall"); return "Leaving the table." }},
		},
	}

	world["shop"] = &Room{
		ID:          "shop",
		Description: "The General Store. Everything you need for survival.",
		Actions: []Action{
			{ID: "buy_sword", Description: "Steel Sword (50g)", IsAvailable: func(gs *GameState) bool { return !gs.HasSteelSword }, Result: func(gs *GameState) string {
				if gs.Gold >= 50 {
					gs.Gold -= 50
					gs.HasSteelSword = true
					gs.BaseDamage += 20
					return "Bought Sword!"
				}
				return "No gold!"
			}},
			{ID: "buy_meat", Description: "Prime Meat (20g)", IsAvailable: func(gs *GameState) bool { return !gs.HasMeat }, Result: func(gs *GameState) string {
				if gs.Gold >= 20 {
					gs.Gold -= 20
					gs.HasMeat = true
					return "Bought Meat!"
				}
				return "No gold!"
			}},
			{ID: "buy_garden", Description: "Garden Kit (30g)", IsAvailable: func(gs *GameState) bool { return !gs.HasGardenKit }, Result: func(gs *GameState) string {
				if gs.Gold >= 30 {
					gs.Gold -= 30
					gs.HasGardenKit = true
					return "Bought Kit!"
				}
				return "No gold!"
			}},
			{ID: "buy_water", Description: "Water Skin (10g) - Required for Desert", IsAvailable: func(gs *GameState) bool { return !gs.HasWater }, Result: func(gs *GameState) string {
				if gs.Gold >= 10 {
					gs.Gold -= 10
					gs.HasWater = true
					return "Bought Water!"
				}
				return "No gold!"
			}},
			{ID: "buy_fur", Description: "Fur Coat (40g) - Required for Arctic", IsAvailable: func(gs *GameState) bool { return !gs.HasFurCoat }, Result: func(gs *GameState) string {
				if gs.Gold >= 40 {
					gs.Gold -= 40
					gs.HasFurCoat = true
					return "Bought Fur Coat!"
				}
				return "No gold!"
			}},
			{ID: "buy_potion", Description: "Minor Potion (20g) - Heals 30 HP", Result: func(gs *GameState) string {
				if gs.Gold >= 20 {
					gs.Gold -= 20
					restored := Heal(gs, 30)
					return fmt.Sprintf("You drank a Minor Potion and recovered %d HP.", restored)
				}
				return "You don't have enough gold!"
			}},
			{ID: "leave_shop", Description: "Leave Shop", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("village"); return "See ya!" }},
		},
	}

	world["garden"] = &Room{
		ID:          "garden",
		Description: "The tranquil Zen Garden.",
		Actions: []Action{
			{ID: "work", Description: "Tend Garden (Needs Kit)", IsAvailable: func(gs *GameState) bool { return gs.HasGardenKit }, Result: func(gs *GameState) string {
				gs.SkillPoints += 15
				if gs.SkillPoints >= gs.Level*40 {
					gs.Level++
					gs.MaxHealth += 20
					gs.Health = gs.MaxHealth
					return "Level Up!"
				}
				return "Tending garden... (+15 SP)"
			}},
			{ID: "back", Description: "Back to Village", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("village"); return "Leaving garden." }},
		},
	}

	// --- MAIN WORLD ---
	world["forest"] = &Room{
		ID:          "forest",
		Description: "A crossroad in the forest. Village to West, River North, Zoo South, Desert East, and ancient ruins off a hidden trail.",
		Actions: []Action{
			{ID: "go_village", Description: "Go West (Village)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("village"); return "Heading home." }},
			{ID: "go_north", Description: "Go North (River)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("river"); return "Water sound ahead." }},
			{ID: "go_south", Description: "Go South (Zoo)", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("zoo")
				if !gs.IsBeastDead {
					BeginCombat(gs, 120)
				}
				return "Roars ahead."
			}},
			{ID: "go_east", Description: "Go East (Scorched Desert)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("desert"); return "The air gets hot." }},
			{ID: "go_ruins", Description: "Take the Hidden Trail (Ruins)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("old_ruins"); return "Old stones await beyond the trees." }},
		},
	}

	world["river"] = &Room{
		ID:          "river",
		Description: "A cold river. Forest South, Cave behind waterfall, Arctic North across the bridge.",
		Actions: []Action{
			{ID: "go_forest", Description: "South (Forest)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "Back to woods." }},
			{ID: "go_cave", Description: "Cave", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("cave")
				if !gs.IsTrollDead {
					BeginCombat(gs, 50)
				}
				return "Step inside."
			}},
			{ID: "go_arctic", Description: "North (Arctic Peaks)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("arctic"); return "Crossing bridge..." }},
			{ID: "go_harbor", Description: "East (Harbor)", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("harbor")
				return "You follow the river until the salt air reaches you."
			}},
		},
	}

	world["sage_hut"] = &Room{
		ID:          "sage_hut",
		Description: "A candlelit hut stacked with scrolls, shells, and a thousand quiet answers.",
		Actions: []Action{
			{ID: "offer_pearl", Description: "Offer the Moon Pearl", IsAvailable: func(gs *GameState) bool { return gs.HasMoonPearl && !gs.HasSageBlessing }, Result: func(gs *GameState) string {
				if !gs.HasMoonPearl {
					return "The wise man says he cannot accept an empty palm."
				}
				gs.HasMoonPearl = false
				gs.HasSageBlessing = true
				gs.SkillPoints += 20
				return joinCombatMessages(
					"The wise man accepts the Moon Pearl and smiles.",
					AwardVictory(gs, 120, 120),
					"You receive the Sage's Blessing and a sharper mind.",
				)
			}},
			{ID: "complete_triptych", Description: "Offer the Three Signs", IsAvailable: func(gs *GameState) bool {
				return gs.HasRuinsToken && gs.HasMoonCharm && gs.HasStarMap && !gs.HasTriptychBlessing
			}, Result: func(gs *GameState) string {
				if gs.HasTriptychBlessing {
					return "The wise man already marked your triptych as complete."
				}
				if !gs.HasRuinsToken || !gs.HasMoonCharm || !gs.HasStarMap {
					return "The wise man says the three signs are not yet all in your hands."
				}
				gs.HasRuinsToken = false
				gs.HasMoonCharm = false
				gs.HasStarMap = false
				gs.HasTriptychBlessing = true
				gs.SkillPoints += 30
				gs.Gold += 150
				return joinCombatMessages(
					"The wise man studies the three signs and nods with approval.",
					AwardVictory(gs, 150, 150),
					"You gain the Triptych Blessing, 150 gold, and a sharper sense of hidden paths.",
				)
			}},
			{ID: "leave", Description: "Return to the Village", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("village"); return "You leave the sage's hut." }},
		},
	}

	world["harbor"] = &Room{
		ID:          "harbor",
		Description: "A busy harbor of salt wind, ropes, and creaking boats. An old lighthouse watches over the water.",
		Actions: []Action{
			{ID: "go_river", Description: "West (Back to River)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("river"); return "You head away from the harbor." }},
			{ID: "go_lighthouse", Description: "Climb the Lighthouse", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("lighthouse"); return "You climb toward the beacon." }},
			{ID: "go_tide_caves", Description: "Enter the Tide Caves", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("tide_caves"); return "The tide pulls you into the caves." }},
			{ID: "go_moonwell", Description: "Find the Moonwell", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("moonwell")
				return "You follow a lantern path to the Moonwell."
			}},
		},
	}

	world["lighthouse"] = &Room{
		ID:          "lighthouse",
		Description: "A tall lighthouse with a beacon that cuts through storm clouds.",
		Actions: []Action{
			{ID: "climb", Description: "Climb the beacon stairs", Result: func(gs *GameState) string {
				if gs.HasMoonPearl && gs.HasSageBlessing {
					gs.SkillPoints += 10
					return "From the beacon room you see every road you have walked. The wise man's words feel clearer now."
				}
				gs.Gold += 40
				return "You spot a hidden supply cache near the beacon and recover 40 gold."
			}},
			{ID: "return", Description: "Return to the Harbor", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("harbor"); return "You descend to the harbor." }},
		},
	}

	world["moonwell"] = &Room{
		ID:          "moonwell",
		Description: "A tide-fed Moonwell tucked behind the harbor. The water glows when you speak softly.",
		Actions: []Action{
			{ID: "draw_moonlight", Description: "Draw Moonlight Water", IsAvailable: func(gs *GameState) bool { return !gs.HasMoonCharm }, Result: func(gs *GameState) string {
				if gs.HasMoonCharm {
					return "The Moonwell is already calm."
				}
				gs.HasMoonCharm = true
				gs.SkillPoints += 10
				gs.Gold += 20
				return joinCombatMessages(
					"The Moonwell answers with a silver shimmer.",
					AwardVictory(gs, 60, 80),
					"You recover a Moon Charm and a little wisdom.",
				)
			}},
			{ID: "return", Description: "Return to the Harbor", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("harbor"); return "You leave the Moonwell behind." }},
		},
	}

	world["tide_caves"] = &Room{
		ID:          "tide_caves",
		Description: "Low caves carved by the tide. Something bright glints beneath the foam.",
		Actions: []Action{
			{ID: "recover_pearl", Description: "Recover the Moon Pearl", IsAvailable: func(gs *GameState) bool { return !gs.HasMoonPearl }, Result: func(gs *GameState) string {
				if gs.HasMoonPearl {
					return "You already recovered the Moon Pearl."
				}
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 90)
				}

				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The tide swallows your opening move.")...)
				}

				dmg := AttackDamage(gs, 24)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You strike the tide guardian for %d damage.", dmg))
				DrawEnemyHealthBar("Tide Guardian", gs.MonsterHealth, 90)

				if gs.MonsterHealth <= 0 {
					gs.HasMoonPearl = true
					gs.Transition("harbor")
					messages = append(messages, "The tide guardian collapses and reveals the Moon Pearl.")
					messages = append(messages, AwardVictory(gs, 100, 150))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Tide Guardian", 14, CombatStatusPoison, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "return", Description: "Return to the Harbor", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("harbor"); return "You escape the tide caves." }},
		},
	}

	world["old_ruins"] = &Room{
		ID:          "old_ruins",
		Description: "Collapsed stone arches and broken mosaics. The forest hides this place from casual travelers.",
		Actions: []Action{
			{ID: "study_ruins", Description: "Study the Ancient Tablets", IsAvailable: func(gs *GameState) bool { return !gs.HasRuinsToken }, Result: func(gs *GameState) string {
				if gs.HasRuinsToken {
					return "You already deciphered the tablets."
				}
				gs.HasRuinsToken = true
				gs.SkillPoints += 10
				gs.Gold += 15
				return joinCombatMessages(
					"The tablets describe a path of three hidden signs.",
					AwardVictory(gs, 60, 80),
					"You recover a Ruins Token and 15 gold.",
				)
			}},
			{ID: "return", Description: "Return to the Forest", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "You slip back into the trees." }},
		},
	}

	// --- BIOME: DESERT ---
	world["desert"] = &Room{
		ID:          "desert",
		Description: "The Scorched Desert. Endless dunes and heat.",
		Actions: []Action{
			{ID: "go_forest", Description: "West (Back to Forest)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "Escaping heat." }},
			{ID: "explore_desert", Description: "Search for the Sphinx", IsAvailable: func(gs *GameState) bool { return !gs.HasSunAmulet }, Result: func(gs *GameState) string {
				if !gs.HasWater {
					gs.Health -= 20
					if gs.Health <= 0 {
						gs.IsGameOver = true
						return "You died of thirst."
					}
					return "You explore, but the heat is killing you! (-20 HP)"
				}
				gs.Transition("sphinx")
				return "You find a massive stone lion with a human face."
			}},
			{ID: "go_mountain", Description: "Go to Mountain (Needs both Amulets)", IsAvailable: func(gs *GameState) bool { return gs.HasSunAmulet && gs.HasIceCrystal }, Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("mountain"); return "Ascending..." }},
		},
	}

	world["sphinx"] = &Room{
		ID:          "sphinx",
		Description: "The Sphinx speaks: 'I have no voice, but I can scream. I have no wings, but I can fly. What am I?'",
		Actions: []Action{
			{ID: "ans_wind", Description: "The Wind", Result: func(gs *GameState) string {
				gs.HasSunAmulet = true
				gs.Gold += 200
				gs.Transition("desert")
				return "Correct! She hands you the Sun Amulet and 200 gold."
			}},
			{ID: "ans_shadow", Description: "A Shadow", Result: func(gs *GameState) string {
				gs.Health -= 30
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				DrawEnemyHealthBar("Sphinx Aura", 100, 100) // Visual flair
				return "Wrong! She slaps you for 30 damage."
			}},
			{ID: "run", Description: "Run away", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("desert"); return "Cowardly retreat." }},
		},
	}

	// --- BIOME: ARCTIC ---
	world["arctic"] = &Room{
		ID:          "arctic",
		Description: "The Arctic Entrance. A frozen wasteland. To the North lie the Frost Cliffs, and to the East is an Ice Bridge.",
		Actions: []Action{
			{ID: "go_river", Description: "South (Back to River)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("river"); return "Leaving cold." }},
			{ID: "go_cliffs", Description: "North (Frost Cliffs)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("frost_cliffs"); return "Climbing the icy cliffs..." }},
			{ID: "go_bridge", Description: "East (Ice Bridge)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("ice_bridge"); return "Walking onto the slippery bridge..." }},
		},
	}

	world["frost_cliffs"] = &Room{
		ID:          "frost_cliffs",
		Description: "High Frost Cliffs. The wind howls here. A Frost Giant guards this area.",
		Actions: []Action{
			{ID: "go_arctic", Description: "South (Back to Entrance)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("arctic"); return "Descending." }},
			{ID: "fight_giant", Description: "Fight Frost Giant", IsAvailable: func(gs *GameState) bool { return !gs.IsFrostGiantDead }, Result: func(gs *GameState) string {
				if !gs.HasFurCoat {
					gs.Health -= 15
					if gs.Health <= 0 {
						gs.IsGameOver = true
					}
					return joinCombatMessages("You are freezing! (-15 HP)", ApplyPlayerCombatStatus(gs, CombatStatusFreeze, 1))
				}
				gs.Transition("frost_giant")
				if !gs.IsFrostGiantDead {
					BeginCombat(gs, 180)
				}
				return "A giant made of ice towers over you!"
			}},
		},
	}

	world["ice_bridge"] = &Room{
		ID:          "ice_bridge",
		Description: "A narrow Ice Bridge. You see a flickering Glitched Portal at the end of it.",
		Actions: []Action{
			{ID: "go_arctic", Description: "West (Back to Entrance)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("arctic"); return "Heading back." }},
			{ID: "enter_portal", Description: "Enter Glitched Portal", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("void"); return "010101... Dimension Shift... 010101" }},
		},
	}

	world["frost_giant"] = &Room{
		ID:          "frost_giant",
		Description: "The Frost Giant roars!",
		Actions: []Action{
			{ID: "fight", Description: "Attack Giant", Result: func(gs *GameState) string {
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "You are frozen solid and lose your attack.")...)
				}

				dmg := AttackDamage(gs, 20)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You deal %d damage to the Frost Giant.", dmg))
				DrawEnemyHealthBar("Frost Giant", gs.MonsterHealth, 180)

				if gs.MonsterHealth <= 0 {
					gs.IsFrostGiantDead = true
					gs.HasIceCrystal = true
					gs.Transition("frost_cliffs")
					messages = append(messages, "The Frost Giant shatters under the ice.")
					messages = append(messages, AwardVictory(gs, 200, 300))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "The Frost Giant", 20, CombatStatusFreeze, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "run", Description: "Run", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("frost_cliffs"); return "Run!" }},
		},
	}

	// --- END GAME ---
	world["mountain"] = &Room{
		ID:          "mountain",
		Description: "The Dragon's Peak. The Great Wyrm awaits, and an observatory sits higher up the trail.",
		Actions: []Action{
			{ID: "fight_dragon", Description: "Slay the Dragon", Result: func(gs *GameState) string {
				if gs.Level < 4 {
					return "The Dragon laughs. You are too weak! (Reach Lvl 4)"
				}
				gs.IsGameOver = true
				return "YOU DID IT! You slew the Dragon and became a legend of Oakhaven! THE END."
			}},
			{ID: "go_observatory", Description: "Climb to the Observatory", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("observatory")
				return "You follow the high ridge to the observatory."
			}},
		},
	}

	world["observatory"] = &Room{
		ID:          "observatory",
		Description: "A ruined observatory where the sky feels close enough to touch.",
		Actions: []Action{
			{ID: "study_stars", Description: "Align the Star Map", IsAvailable: func(gs *GameState) bool { return !gs.HasStarMap }, Result: func(gs *GameState) string {
				if gs.HasStarMap {
					return "The star map is already fixed in your notes."
				}
				gs.HasStarMap = true
				gs.SkillPoints += 10
				gs.Gold += 20
				return joinCombatMessages(
					"The observatory lens catches a pattern in the stars.",
					AwardVictory(gs, 60, 80),
					"You recover a Star Map and 20 gold.",
				)
			}},
			{ID: "return", Description: "Return to the Mountain", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("mountain"); return "You descend from the observatory." }},
		},
	}

	// --- LEGACY & REBALANCED ---
	world["zoo"] = &Room{
		ID:          "zoo",
		Description: "A massive Beast blocks the entrance. It looks hungry and powerful.",
		Actions: []Action{
			{ID: "feed", Description: "Feed Meat", IsAvailable: func(gs *GameState) bool { return gs.HasMeat && !gs.IsBeastDead }, Result: func(gs *GameState) string {
				gs.HasMeat = false
				gs.Gold += 150
				gs.IsBeastDead = true
				return "The beast eats the meat and falls asleep. You find 150g!"
			}},
			{ID: "fight_beast", Description: "Attack the Beast", IsAvailable: func(gs *GameState) bool { return !gs.IsBeastDead }, Result: func(gs *GameState) string {
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "You are stunned and cannot act.")...)
				}

				dmg := AttackDamage(gs, 20)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You hit the Beast for %d damage.", dmg))
				DrawEnemyHealthBar("Beast", gs.MonsterHealth, 120)

				if gs.MonsterHealth <= 0 {
					gs.IsBeastDead = true
					messages = append(messages, "The Beast collapses in a heap.")
					messages = append(messages, AwardVictory(gs, 150, 200))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "The Beast", 15, CombatStatusPoison, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Run back to Forest", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "You escaped." }},
		},
	}

	world["cave"] = &Room{
		ID:          "cave",
		Description: "Damp cave with a troll.",
		Actions: []Action{
			{ID: "hit", Description: "Hit Troll", IsAvailable: func(gs *GameState) bool { return !gs.IsTrollDead }, Result: func(gs *GameState) string {
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The troll's blow leaves you unable to swing.")...)
				}

				dmg := AttackDamage(gs, 20)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You deal %d damage to the Troll.", dmg))
				DrawEnemyHealthBar("Troll", gs.MonsterHealth, 50)

				if gs.MonsterHealth <= 0 {
					gs.IsTrollDead = true
					messages = append(messages, "The Troll staggers and falls.")
					messages = append(messages, AwardVictory(gs, 50, 100))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "The Troll", 10, CombatStatusPoison, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "leave", Description: "Leave", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("river"); return "Back." }},
		},
	}

	// --- BIOME: FORSAKEN DIMENSION ---
	world["void"] = &Room{
		ID:          "void",
		Description: "The Glitch Barrens. The reality here is flickering. You feel a strange presence.",
		Actions: []Action{
			{ID: "explore", Description: "Explore the flickering void", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("binary_sea"); return "The ground turns into 1s and 0s." }},
			{ID: "exit", Description: "Return through Portal", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("arctic"); return "Back to the cold peaks." }},
		},
	}

	world["party_hunt_gate"] = &Room{
		ID:          "party_hunt_gate",
		Description: "A beast-trail gate sealed by snarling wards. This is a party-only dungeon phase for Monster Hunt.",
		Actions: []Action{
			{ID: "clear_outer", Description: "Clear the outer sentries", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "monster-hunt" && gs.PartyQuestPhaseIndex == 0
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 60)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The outer sentries hold the line and you lose your opening.")...)
				}

				dmg := AttackDamage(gs, 18)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You cut down the outer sentries for %d damage.", dmg))
				DrawEnemyHealthBar("Outer Sentries", gs.MonsterHealth, 60)

				if gs.MonsterHealth <= 0 {
					gs.Transition("party_hunt_depths")
					BeginCombat(gs, 80)
					messages = append(messages, "The sentries collapse and the trail opens deeper into the lair.")
					messages = append(messages, AwardVictory(gs, 40, 60))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Outer Sentries", 10, CombatStatusPoison, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Retreat to the forest", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "The party falls back to safer ground." }},
		},
	}

	world["party_hunt_depths"] = &Room{
		ID:          "party_hunt_depths",
		Description: "The beast lair's lower chambers. Wave after wave of hunters and beasts fill the dark.",
		Actions: []Action{
			{ID: "purge_depths", Description: "Purge the lair waves", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "monster-hunt" && gs.PartyQuestPhaseIndex == 1
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 80)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The lair swarms you before you can swing.")...)
				}

				dmg := AttackDamage(gs, 22)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You tear through the lair waves for %d damage.", dmg))
				DrawEnemyHealthBar("Lair Waves", gs.MonsterHealth, 80)

				if gs.MonsterHealth <= 0 {
					messages = append(messages, "The lair waves collapse into silence.")
					messages = append(messages, AwardVictory(gs, 60, 80))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Lair Waves", 12, CombatStatusStun, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "descend", Description: "Descend to the alpha vault", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "monster-hunt" && gs.PartyQuestPhaseIndex >= 2
			}, Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_hunt_vault")
				BeginCombat(gs, 120)
				return "A final chamber opens before the alpha beast."
			}},
			{ID: "retreat", Description: "Withdraw to the gate", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_hunt_gate")
				BeginCombat(gs, 60)
				return "You retreat to the gate."
			}},
		},
	}

	world["party_hunt_vault"] = &Room{
		ID:          "party_hunt_vault",
		Description: "The alpha vault. The beast lord waits with its pack at its back.",
		Actions: []Action{
			{ID: "fight_alpha", Description: "Fight the alpha beast", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "monster-hunt" && gs.PartyQuestPhaseIndex >= 2
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 120)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The alpha beast swats away your opening.")...)
				}

				dmg := AttackDamage(gs, 30)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You strike the alpha beast for %d damage.", dmg))
				DrawEnemyHealthBar("Alpha Beast", gs.MonsterHealth, 120)

				if gs.MonsterHealth <= 0 {
					gs.Transition("forest")
					messages = append(messages, "The alpha beast collapses and the dungeon seals behind you.")
					messages = append(messages, AwardVictory(gs, 120, 180))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Alpha Beast", 16, CombatStatusBurn, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "coordinated_strike", Description: "Coordinated Strike (Party)", IsAvailable: func(gs *GameState) bool {
				return gs.PartySupport > 0 && gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "monster-hunt" && gs.PartyQuestPhaseIndex >= 2
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 120)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The party loses the opening and the alpha beast snaps back.")...)
				}

				dmg := AttackDamage(gs, 42)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("The party coordinates a strike for %d damage.", dmg))
				DrawEnemyHealthBar("Alpha Beast", gs.MonsterHealth, 120)
				if gs.PartySupport > 1 {
					messages = append(messages, ApplyMonsterCombatStatus(gs, CombatStatusStun, 1))
				}

				if gs.MonsterHealth <= 0 {
					gs.Transition("forest")
					messages = append(messages, "The alpha beast collapses under the coordinated strike.")
					messages = append(messages, AwardVictory(gs, 120, 180))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Alpha Beast", 14, CombatStatusBurn, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Retreat to the depths", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_hunt_depths")
				BeginCombat(gs, 80)
				return "You retreat from the alpha vault."
			}},
		},
	}

	world["party_void_gate"] = &Room{
		ID:          "party_void_gate",
		Description: "A fractured gate of shifting code. This is a party-only dungeon phase for Void Expedition.",
		Actions: []Action{
			{ID: "stabilize_entry", Description: "Stabilize the entry field", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "void-expedition" && gs.PartyQuestPhaseIndex == 0
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 70)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The static field overwhelms your opening.")...)
				}

				dmg := AttackDamage(gs, 20)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You stabilize the entry field for %d damage.", dmg))
				DrawEnemyHealthBar("Static Field", gs.MonsterHealth, 70)

				if gs.MonsterHealth <= 0 {
					gs.Transition("party_void_depths")
					BeginCombat(gs, 100)
					messages = append(messages, "The static field collapses and the route deepens.")
					messages = append(messages, AwardVictory(gs, 50, 70))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Static Field", 12, CombatStatusFreeze, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Return to the arctic", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("arctic"); return "The party escapes the glitch gate." }},
		},
	}

	world["party_void_depths"] = &Room{
		ID:          "party_void_depths",
		Description: "A lower chamber of flickering shards and broken rules.",
		Actions: []Action{
			{ID: "purge_shards", Description: "Purge the shard waves", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "void-expedition" && gs.PartyQuestPhaseIndex == 1
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 100)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The shard waves distort your swing.")...)
				}

				dmg := AttackDamage(gs, 24)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You tear through the shard waves for %d damage.", dmg))
				DrawEnemyHealthBar("Shard Waves", gs.MonsterHealth, 100)

				if gs.MonsterHealth <= 0 {
					messages = append(messages, "The shard waves collapse into fragments.")
					messages = append(messages, AwardVictory(gs, 80, 100))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Shard Waves", 14, CombatStatusBurn, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "descend", Description: "Descend to the core chamber", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "void-expedition" && gs.PartyQuestPhaseIndex >= 2
			}, Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_void_core")
				BeginCombat(gs, 140)
				return "The final core chamber hums ahead."
			}},
			{ID: "retreat", Description: "Return to the gate", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_void_gate")
				BeginCombat(gs, 70)
				return "You retreat to the glitch gate."
			}},
		},
	}

	world["party_void_core"] = &Room{
		ID:          "party_void_core",
		Description: "The void core chamber. The breach itself is trying to think.",
		Actions: []Action{
			{ID: "seal_core", Description: "Seal the void core", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "void-expedition" && gs.PartyQuestPhaseIndex >= 2
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 140)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The core rewrites your opening before you strike.")...)
				}

				dmg := AttackDamage(gs, 34)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You strike the void core for %d damage.", dmg))
				DrawEnemyHealthBar("Void Core", gs.MonsterHealth, 140)

				if gs.MonsterHealth <= 0 {
					gs.Transition("arctic")
					messages = append(messages, "The void core seals itself shut and the party escapes the chamber.")
					messages = append(messages, AwardVictory(gs, 140, 220))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Void Core", 18, CombatStatusBurn, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "coordinated_strike", Description: "Coordinated Strike (Party)", IsAvailable: func(gs *GameState) bool {
				return gs.PartySupport > 0 && gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "void-expedition" && gs.PartyQuestPhaseIndex >= 2
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 140)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The void core disrupts the coordinated strike.")...)
				}

				dmg := AttackDamage(gs, 46)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("The party coordinates a strike for %d damage.", dmg))
				DrawEnemyHealthBar("Void Core", gs.MonsterHealth, 140)
				if gs.PartySupport > 1 {
					messages = append(messages, ApplyMonsterCombatStatus(gs, CombatStatusStun, 1))
				}

				if gs.MonsterHealth <= 0 {
					gs.Transition("arctic")
					messages = append(messages, "The void core collapses under the coordinated strike.")
					messages = append(messages, AwardVictory(gs, 140, 220))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Void Core", 16, CombatStatusBurn, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Escape to the depths", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_void_depths")
				BeginCombat(gs, 100)
				return "You retreat from the core chamber."
			}},
		},
	}

	world["party_frost_gate"] = &Room{
		ID:          "party_frost_gate",
		Description: "A frozen gate sealed by ancient ice. This is a party-only dungeon phase for Frost Pact.",
		Actions: []Action{
			{ID: "shatter_ward", Description: "Shatter the outer ward", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "frost-pact" && gs.PartyQuestPhaseIndex == 0
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 70)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The ward freezes your opening.")...)
				}

				dmg := AttackDamage(gs, 18)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You shatter the outer ward for %d damage.", dmg))
				DrawEnemyHealthBar("Outer Ward", gs.MonsterHealth, 70)

				if gs.MonsterHealth <= 0 {
					gs.Transition("party_frost_depths")
					BeginCombat(gs, 90)
					messages = append(messages, "The ward shatters and the glacier opens.")
					messages = append(messages, AwardVictory(gs, 35, 55))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Outer Ward", 10, CombatStatusFreeze, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Retreat to the arctic", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("arctic"); return "You retreat from the frozen gate." }},
		},
	}

	world["party_frost_depths"] = &Room{
		ID:          "party_frost_depths",
		Description: "Glacier depths where the Frost Wyrm's breath freezes the air itself.",
		Actions: []Action{
			{ID: "break_glacier", Description: "Break the glacier ward", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "frost-pact" && gs.PartyQuestPhaseIndex == 1
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 90)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The glacier ward solidifies around your weapon.")...)
				}

				dmg := AttackDamage(gs, 22)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You break the glacier ward for %d damage.", dmg))
				DrawEnemyHealthBar("Glacier Ward", gs.MonsterHealth, 90)

				if gs.MonsterHealth <= 0 {
					messages = append(messages, "The glacier ward breaks apart.")
					messages = append(messages, AwardVictory(gs, 50, 75))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Glacier Ward", 12, CombatStatusFreeze, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "descend", Description: "Descend to the Frost Wyrm peak", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "frost-pact" && gs.PartyQuestPhaseIndex >= 2
			}, Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_frost_peak")
				BeginCombat(gs, 130)
				return "The party climbs toward the Frost Wyrm."
			}},
			{ID: "retreat", Description: "Return to the gate", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_frost_gate")
				BeginCombat(gs, 70)
				return "You retreat to the frozen gate."
			}},
		},
	}

	world["party_frost_peak"] = &Room{
		ID:          "party_frost_peak",
		Description: "The glacier peak. The Frost Wyrm coils around the summit.",
		Actions: []Action{
			{ID: "slay_wyrm", Description: "Slay the Frost Wyrm", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "frost-pact" && gs.PartyQuestPhaseIndex >= 2
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 130)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The Frost Wyrm coils around your opening.")...)
				}

				dmg := AttackDamage(gs, 32)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You strike the Frost Wyrm for %d damage.", dmg))
				DrawEnemyHealthBar("Frost Wyrm", gs.MonsterHealth, 130)

				if gs.MonsterHealth <= 0 {
					gs.Transition("arctic")
					messages = append(messages, "The Frost Wyrm shatters and the glacier calms.")
					messages = append(messages, AwardVictory(gs, 130, 200))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Frost Wyrm", 16, CombatStatusFreeze, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "coordinated_strike", Description: "Coordinated Strike (Party)", IsAvailable: func(gs *GameState) bool {
				return gs.PartySupport > 0 && gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "frost-pact" && gs.PartyQuestPhaseIndex >= 2
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 130)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The Frost Wyrm freezes the coordinated strike in place.")...)
				}

				dmg := AttackDamage(gs, 44)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("The party coordinates a strike for %d damage.", dmg))
				DrawEnemyHealthBar("Frost Wyrm", gs.MonsterHealth, 130)
				if gs.PartySupport > 1 {
					messages = append(messages, ApplyMonsterCombatStatus(gs, CombatStatusStun, 1))
				}

				if gs.MonsterHealth <= 0 {
					gs.Transition("arctic")
					messages = append(messages, "The Frost Wyrm shatters under the coordinated strike.")
					messages = append(messages, AwardVictory(gs, 130, 200))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Frost Wyrm", 14, CombatStatusFreeze, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Retreat to the depths", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_frost_depths")
				BeginCombat(gs, 90)
				return "You retreat from the Frost Wyrm peak."
			}},
		},
	}

	world["party_sun_gate"] = &Room{
		ID:          "party_sun_gate",
		Description: "A radiant temple gate lit by dawnfire. This is a party-only dungeon phase for Sun Covenant.",
		Actions: []Action{
			{ID: "ignite_gate", Description: "Ignite the dawn gate", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "sun-covenant" && gs.PartyQuestPhaseIndex == 0
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 80)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The dawn gate blinds your opening and you lose momentum.")...)
				}

				dmg := AttackDamage(gs, 20)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You ignite the dawn gate for %d damage.", dmg))
				DrawEnemyHealthBar("Dawn Gate", gs.MonsterHealth, 80)

				if gs.MonsterHealth <= 0 {
					gs.Transition("party_sun_depths")
					BeginCombat(gs, 110)
					messages = append(messages, "The dawn gate collapses and the temple opens.")
					messages = append(messages, AwardVictory(gs, 45, 70))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Dawn Gate", 11, CombatStatusBurn, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Retreat to the desert", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("desert"); return "The party falls back to the dunes." }},
		},
	}

	world["party_sun_depths"] = &Room{
		ID:          "party_sun_depths",
		Description: "Temple depths filled with burning dunes and radiant guardians.",
		Actions: []Action{
			{ID: "purge_dunes", Description: "Purge the burning dunes", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "sun-covenant" && gs.PartyQuestPhaseIndex == 1
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 110)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The burning dunes swallow your opening.")...)
				}

				dmg := AttackDamage(gs, 26)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You cut through the burning dunes for %d damage.", dmg))
				DrawEnemyHealthBar("Burning Dunes", gs.MonsterHealth, 110)

				if gs.MonsterHealth <= 0 {
					messages = append(messages, "The burning dunes collapse into calm sand.")
					messages = append(messages, AwardVictory(gs, 75, 100))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "Burning Dunes", 13, CombatStatusBurn, 1))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "descend", Description: "Descend to the crown chamber", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "sun-covenant" && gs.PartyQuestPhaseIndex >= 2
			}, Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_sun_crown")
				BeginCombat(gs, 150)
				return "The party climbs toward the Sunbound Colossus."
			}},
			{ID: "retreat", Description: "Return to the gate", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_sun_gate")
				BeginCombat(gs, 80)
				return "You retreat to the dawn gate."
			}},
		},
	}

	world["party_sun_crown"] = &Room{
		ID:          "party_sun_crown",
		Description: "The crown chamber. The Sunbound Colossus stands between the party and dawn.",
		Actions: []Action{
			{ID: "break_crown", Description: "Break the sun crown", IsAvailable: func(gs *GameState) bool {
				return gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "sun-covenant" && gs.PartyQuestPhaseIndex >= 2
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 150)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The colossus blinds you before you can strike.")...)
				}

				dmg := AttackDamage(gs, 34)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You crack the sun crown for %d damage.", dmg))
				DrawEnemyHealthBar("Sunbound Colossus", gs.MonsterHealth, 150)

				if gs.MonsterHealth <= 0 {
					gs.Transition("desert")
					messages = append(messages, "The sun crown shatters and the temple goes quiet.")
					messages = append(messages, AwardVictory(gs, 150, 240))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				if gs.MonsterHealth < 75 {
					messages = append(messages, "The colossus radiates a scorching pulse.")
				}
				messages = append(messages, MonsterCounterAttack(gs, "Sunbound Colossus", 18, CombatStatusBurn, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "coordinated_strike", Description: "Coordinated Strike (Party)", IsAvailable: func(gs *GameState) bool {
				return gs.PartySupport > 0 && gs.PartyQuestStatus == "active" && gs.PartyQuestKey == "sun-covenant" && gs.PartyQuestPhaseIndex >= 2
			}, Result: func(gs *GameState) string {
				if gs.MonsterHealth <= 0 {
					BeginCombat(gs, 150)
				}
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "The colossus blinds the party's opening strike.")...)
				}

				dmg := AttackDamage(gs, 48)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("The party drives the sun crown for %d damage.", dmg))
				DrawEnemyHealthBar("Sunbound Colossus", gs.MonsterHealth, 150)
				if gs.PartySupport > 1 {
					messages = append(messages, ApplyMonsterCombatStatus(gs, CombatStatusStun, 1))
				}

				if gs.MonsterHealth <= 0 {
					gs.Transition("desert")
					messages = append(messages, "The Sunbound Colossus kneels and the crown shatters.")
					messages = append(messages, AwardVictory(gs, 150, 240))
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				if gs.MonsterHealth < 75 {
					messages = append(messages, ApplyPlayerCombatStatus(gs, CombatStatusBurn, 1))
				}
				messages = append(messages, MonsterCounterAttack(gs, "Sunbound Colossus", 16, CombatStatusBurn, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "retreat", Description: "Retreat to the depths", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("party_sun_depths")
				BeginCombat(gs, 110)
				return "You retreat from the crown chamber."
			}},
		},
	}

	world["binary_sea"] = &Room{
		ID:          "binary_sea",
		Description: "A vast sea of data. A dark figure stands on a floating island of code.",
		Actions: []Action{
			{ID: "challenge", Description: "Challenge the Figure", IsAvailable: func(gs *GameState) bool { return !gs.Is1x1x1x1Dead }, Result: func(gs *GameState) string {
				gs.Transition("boss_1x1x1x1")
				BeginCombat(gs, 200)
				return "The figure turns around. It's 1x1x1x1! 'Null is the only truth.'"
			}},
			{ID: "search", Description: "Search for loot", IsAvailable: func(gs *GameState) bool { return gs.Is1x1x1x1Dead }, Result: func(gs *GameState) string {
				return "You found some Glitched Data (Worth 500 gold)!"
			}},
			{ID: "go_back", Description: "Go back to Void", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("void"); return "Reality stabilizes slightly." }},
		},
	}

	world["boss_1x1x1x1"] = &Room{
		ID:          "boss_1x1x1x1",
		Description: "1x1x1x1 is glitching through reality!",
		Actions: []Action{
			{ID: "attack", Description: "Strike with all your might", Result: func(gs *GameState) string {
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "Your body glitches and you lose your opening.")...)
				}

				dmg := AttackDamage(gs, 25)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("You strike 1x1x1x1 for %d damage.", dmg))
				DrawEnemyHealthBar("1x1x1x1", gs.MonsterHealth, 200)

				if gs.MonsterHealth <= 0 {
					gs.Is1x1x1x1Dead = true
					gs.HasForsakenBlade = true
					gs.Transition("binary_sea")
					messages = append(messages, "CRITICAL SYSTEM FAILURE! 1x1x1x1 shatters into raw data.")
					messages = append(messages, AwardVictory(gs, 500, 1000))
					messages = append(messages, "You found the FORSAKEN BLADE!")
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				if gs.MonsterHealth < 100 && gs.Health > 30 {
					gs.Health = 30
					messages = append(messages, fmt.Sprintf("You deal %d damage. 1x1x1x1 glitches your health to 30!", dmg))
					messages = append(messages, ApplyPlayerCombatStatus(gs, CombatStatusBurn, 2))
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "1x1x1x1", 25, CombatStatusBurn, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "coordinated_strike", Description: "Coordinated Strike (Party)", IsAvailable: func(gs *GameState) bool { return gs.PartySupport > 0 && !gs.Is1x1x1x1Dead }, Result: func(gs *GameState) string {
				blocked := gs.PlayerCombatStatus == CombatStatusFreeze || gs.PlayerCombatStatus == CombatStatusStun
				messages := TickCombatStatuses(gs)
				if blocked {
					return joinCombatMessages(append(messages, "Your party loses the opening and 1x1x1x1 slips away.")...)
				}

				dmg := AttackDamage(gs, 40)
				gs.MonsterHealth -= dmg
				messages = append(messages, fmt.Sprintf("Your party coordinates a strike for %d damage.", dmg))
				DrawEnemyHealthBar("1x1x1x1", gs.MonsterHealth, 200)
				if gs.PartySupport > 1 {
					messages = append(messages, ApplyMonsterCombatStatus(gs, CombatStatusStun, 1))
				}

				if gs.MonsterHealth <= 0 {
					gs.Is1x1x1x1Dead = true
					gs.HasForsakenBlade = true
					gs.Transition("binary_sea")
					messages = append(messages, "The party overwhelms the glitch core and it collapses.")
					messages = append(messages, AwardVictory(gs, 500, 1000))
					messages = append(messages, "You found the FORSAKEN BLADE!")
					ClearCombatStatus(gs)
					return joinCombatMessages(messages...)
				}

				messages = append(messages, MonsterCounterAttack(gs, "1x1x1x1", 20, CombatStatusBurn, 2))
				if gs.Health <= 0 {
					gs.IsGameOver = true
				}
				return joinCombatMessages(messages...)
			}},
			{ID: "reboot", Description: "Try to reboot the world (Run)", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("void")
				return "You managed to reboot the local area and escaped to the Barrens!"
			}},
		},
	}

	return world
}
