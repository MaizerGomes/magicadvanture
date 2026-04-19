package game

import "fmt"

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
		},
	}

	world["shop"] = &Room{
		ID:          "shop",
		Description: "The General Store. Everything you need for survival.",
		Actions: []Action{
			{ID: "buy_sword", Description: "Steel Sword (50g)", IsAvailable: func(gs *GameState) bool { return !gs.HasSteelSword }, Result: func(gs *GameState) string { if gs.Gold >= 50 { gs.Gold -= 50; gs.HasSteelSword = true; gs.BaseDamage += 20; return "Bought Sword!" }; return "No gold!" }},
			{ID: "buy_meat", Description: "Prime Meat (20g)", IsAvailable: func(gs *GameState) bool { return !gs.HasMeat }, Result: func(gs *GameState) string { if gs.Gold >= 20 { gs.Gold -= 20; gs.HasMeat = true; return "Bought Meat!" }; return "No gold!" }},
			{ID: "buy_garden", Description: "Garden Kit (30g)", IsAvailable: func(gs *GameState) bool { return !gs.HasGardenKit }, Result: func(gs *GameState) string { if gs.Gold >= 30 { gs.Gold -= 30; gs.HasGardenKit = true; return "Bought Kit!" }; return "No gold!" }},
			{ID: "buy_water", Description: "Water Skin (10g) - Required for Desert", IsAvailable: func(gs *GameState) bool { return !gs.HasWater }, Result: func(gs *GameState) string { if gs.Gold >= 10 { gs.Gold -= 10; gs.HasWater = true; return "Bought Water!" }; return "No gold!" }},
			{ID: "buy_fur", Description: "Fur Coat (40g) - Required for Arctic", IsAvailable: func(gs *GameState) bool { return !gs.HasFurCoat }, Result: func(gs *GameState) string { if gs.Gold >= 40 { gs.Gold -= 40; gs.HasFurCoat = true; return "Bought Fur Coat!" }; return "No gold!" }},
			{ID: "buy_potion", Description: "Minor Potion (20g) - Heals 30 HP", Result: func(gs *GameState) string {
				if gs.Gold >= 20 {
					gs.Gold -= 20
					gs.Health += 30
					if gs.Health > gs.MaxHealth { gs.Health = gs.MaxHealth }
					return "You drank a Minor Potion and feel better!"
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
			{ID: "work", Description: "Tend Garden (Needs Kit)", IsAvailable: func(gs *GameState) bool { return gs.HasGardenKit }, Result: func(gs *GameState) string { gs.SkillPoints += 15; if gs.SkillPoints >= gs.Level*40 { gs.Level++; gs.MaxHealth += 20; gs.Health = gs.MaxHealth; return "Level Up!" }; return "Tending garden... (+15 SP)" }},
			{ID: "back", Description: "Back to Village", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("village"); return "Leaving garden." }},
		},
	}

	// --- MAIN WORLD ---
	world["forest"] = &Room{
		ID:          "forest",
		Description: "A crossroad in the forest. Village to West, River North, Zoo South, and Desert East.",
		Actions: []Action{
			{ID: "go_village", Description: "Go West (Village)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("village"); return "Heading home." }},
			{ID: "go_north", Description: "Go North (River)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("river"); return "Water sound ahead." }},
			{ID: "go_south", Description: "Go South (Zoo)", Type: MoveAction, Result: func(gs *GameState) string { 
				gs.Transition("zoo"); 
				if !gs.IsBeastDead { gs.MonsterHealth = 120 }
				return "Roars ahead." 
			}},
			{ID: "go_east", Description: "Go East (Scorched Desert)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("desert"); return "The air gets hot." }},
		},
	}

	world["river"] = &Room{
		ID:          "river",
		Description: "A cold river. Forest South, Cave behind waterfall, Arctic North across the bridge.",
		Actions: []Action{
			{ID: "go_forest", Description: "South (Forest)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "Back to woods." }},
			{ID: "go_cave", Description: "Cave", Type: MoveAction, Result: func(gs *GameState) string { 
				gs.Transition("cave"); 
				if !gs.IsTrollDead { gs.MonsterHealth = 50 }
				return "Step inside." 
			}},
			{ID: "go_arctic", Description: "North (Arctic Peaks)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("arctic"); return "Crossing bridge..." }},
		},
	}

	// --- BIOME: DESERT ---
	world["desert"] = &Room{
		ID:          "desert",
		Description: "The Scorched Desert. Endless dunes and heat.",
		Actions: []Action{
			{ID: "go_forest", Description: "West (Back to Forest)", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "Escaping heat." }},
			{ID: "explore_desert", Description: "Search for the Sphinx", IsAvailable: func(gs *GameState) bool { return !gs.HasSunAmulet }, Result: func(gs *GameState) string {
				if !gs.HasWater { gs.Health -= 20; if gs.Health <= 0 { gs.IsGameOver = true; return "You died of thirst." }; return "You explore, but the heat is killing you! (-20 HP)" }
				gs.Transition("sphinx"); return "You find a massive stone lion with a human face."
			}},
			{ID: "go_mountain", Description: "Go to Mountain (Needs both Amulets)", IsAvailable: func(gs *GameState) bool { return gs.HasSunAmulet && gs.HasIceCrystal }, Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("mountain"); return "Ascending..." }},
		},
	}

	world["sphinx"] = &Room{
		ID:          "sphinx",
		Description: "The Sphinx speaks: 'I have no voice, but I can scream. I have no wings, but I can fly. What am I?'",
		Actions: []Action{
			{ID: "ans_wind", Description: "The Wind", Result: func(gs *GameState) string { gs.HasSunAmulet = true; gs.Gold += 200; gs.Transition("desert"); return "Correct! She hands you the Sun Amulet and 200 gold." }},
			{ID: "ans_shadow", Description: "A Shadow", Result: func(gs *GameState) string { 
				gs.Health -= 30; 
				if gs.Health <= 0 { gs.IsGameOver = true }; 
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
				if !gs.HasFurCoat { gs.Health -= 15; return "You are freezing! (-15 HP)" }
				gs.Transition("frost_giant"); 
				if !gs.IsFrostGiantDead { gs.MonsterHealth = 180 }
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
				dmg := 20 + gs.BaseDamage
				if gs.HasSteelSword { dmg += 20 }
				if gs.HasForsakenBlade { dmg += 40 }

				gs.MonsterHealth -= dmg
				DrawEnemyHealthBar("Frost Giant", gs.MonsterHealth, 180)

				if gs.MonsterHealth <= 0 { 
					gs.IsFrostGiantDead = true
					gs.HasIceCrystal = true; 
					gs.Gold += 300; 
					gs.Experience += 200;
					gs.Transition("frost_cliffs"); 
					return "Giant shattered! Found Ice Crystal and 200 XP." 
				}
				gs.Health -= 20; if gs.Health <= 0 { gs.IsGameOver = true }; return fmt.Sprintf("You deal %d, he hits you for 20!", dmg)
			}},
			{ID: "run", Description: "Run", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("frost_cliffs"); return "Run!" }},
		},
	}

	// --- END GAME ---
	world["mountain"] = &Room{
		ID:          "mountain",
		Description: "The Dragon's Peak. The Great Wyrm awaits.",
		Actions: []Action{
			{ID: "fight_dragon", Description: "Slay the Dragon", Result: func(gs *GameState) string {
				if gs.Level < 4 { return "The Dragon laughs. You are too weak! (Reach Lvl 4)" }
				gs.IsGameOver = true
				return "YOU DID IT! You slew the Dragon and became a legend of Oakhaven! THE END."
			}},
		},
	}

	// --- LEGACY & REBALANCED ---
	world["zoo"] = &Room{
		ID:          "zoo",
		Description: "A massive Beast blocks the entrance. It looks hungry and powerful.",
		Actions: []Action{
			{ID: "feed", Description: "Feed Meat", IsAvailable: func(gs *GameState) bool { return gs.HasMeat && !gs.IsBeastDead }, Result: func(gs *GameState) string { gs.HasMeat = false; gs.Gold += 150; gs.IsBeastDead = true; return "The beast eats the meat and falls asleep. You find 150g!" }},
			{ID: "fight_beast", Description: "Attack the Beast", IsAvailable: func(gs *GameState) bool { return !gs.IsBeastDead }, Result: func(gs *GameState) string {
				dmg := 20 + gs.BaseDamage
				if gs.HasSteelSword { dmg += 20 }
				if gs.HasForsakenBlade { dmg += 40 }
				
				gs.MonsterHealth -= dmg
				DrawEnemyHealthBar("Beast", gs.MonsterHealth, 120)
				
				if gs.MonsterHealth <= 0 {
					gs.IsBeastDead = true
					gs.Gold += 200
					gs.Experience += 150
					return fmt.Sprintf("You slew the Beast! Earned 150 XP and 200 gold.")
				}
				
				gs.Health -= 15
				if gs.Health <= 0 { gs.IsGameOver = true; return "The Beast crushed you. Game Over." }
				return fmt.Sprintf("You hit the Beast for %d damage! It bites you for 15 damage.", dmg)
			}},
			{ID: "retreat", Description: "Run back to Forest", Type: MoveAction, Result: func(gs *GameState) string { gs.Transition("forest"); return "You escaped." }},
		},
	}

	world["cave"] = &Room{
		ID:          "cave",
		Description: "Damp cave with a troll.",
		Actions: []Action{
			{ID: "hit", Description: "Hit Troll", IsAvailable: func(gs *GameState) bool { return !gs.IsTrollDead }, Result: func(gs *GameState) string {
				dmg := 20 + gs.BaseDamage
				if gs.HasSteelSword { dmg += 20 }
				if gs.HasForsakenBlade { dmg += 40 }

				gs.MonsterHealth -= dmg
				DrawEnemyHealthBar("Troll", gs.MonsterHealth, 50)

				if gs.MonsterHealth <= 0 { 
					gs.IsTrollDead = true; 
					gs.Gold += 100; 
					gs.Experience += 50;
					return "Troll dead! Found 100g and 50 XP." 
				}
				gs.Health -= 10; return fmt.Sprintf("You deal %d damage, it hit back for 10!", dmg)
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

	world["binary_sea"] = &Room{
		ID:          "binary_sea",
		Description: "A vast sea of data. A dark figure stands on a floating island of code.",
		Actions: []Action{
			{ID: "challenge", Description: "Challenge the Figure", IsAvailable: func(gs *GameState) bool { return !gs.Is1x1x1x1Dead }, Result: func(gs *GameState) string {
				gs.Transition("boss_1x1x1x1");
				gs.MonsterHealth = 200;
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
				dmg := 25 + gs.BaseDamage
				if gs.HasSteelSword { dmg += 20 }
				if gs.HasForsakenBlade { dmg += 40 }
				
				gs.MonsterHealth -= dmg
				DrawEnemyHealthBar("1x1x1x1", gs.MonsterHealth, 200)

				if gs.MonsterHealth <= 0 {
					gs.Is1x1x1x1Dead = true
					gs.HasForsakenBlade = true
					gs.Gold += 1000
					gs.Experience += 500
					gs.Transition("binary_sea")
					return "CRITICAL SYSTEM FAILURE! 1x1x1x1 shatters into raw data. You found the FORSAKEN BLADE and 500 XP!"
				}
				
				// Glitch attack
				if gs.MonsterHealth < 100 && gs.Health > 30 {
					gs.Health = 30
					return fmt.Sprintf("You deal %d damage. 1x1x1x1 GLITCHES your health to 30!", dmg)
				}
				
				gs.Health -= 25
				if gs.Health <= 0 { gs.IsGameOver = true }
				return fmt.Sprintf("You deal %d damage. 1x1x1x1 hits you for 25 damage!", dmg)
			}},
			{ID: "reboot", Description: "Try to reboot the world (Run)", Type: MoveAction, Result: func(gs *GameState) string {
				gs.Transition("void")
				return "You managed to reboot the local area and escaped to the Barrens!"
			}},
		},
	}

	return world
}
