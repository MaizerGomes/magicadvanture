package game

import (
	"fmt"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorMagenta = "\033[35m"
	ColorBold   = "\033[1m"
)

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func DrawHealthBar(label string, current, max int, color string) {
	width := 20
	filled := 0
	if max > 0 {
		filled = (current * width) / max
	}
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("-", width-filled)
	fmt.Printf("%s: [%s%s%s] %d/%d\n", label, color, bar, ColorReset, current, max)
}

func GetArt(roomID string) string {
	switch roomID {
	case "desert":
		return `
      .    .
   .  |  .
    \ | /
  -- ( ) --
    / | \
   '  |  '
  ~~~~~~~~~~~~~~~~
    /  \  /  \
   /    \/    \
  /            \`
	case "arctic":
		return `
     *    *  *
  *    *    *
     *   *
    / \ / \ / \
   /   \   /   \
  /_____\ /_____\`
	case "void":
		return `
     [ GLITCH ]
   01011100101
   10 glitch 10
   00111000110
     [ ERROR ]`
	case "binary_sea":
		return `
    010101010101
    101010101010
    011100110101
    ~~~~~~~~~~~~`
	case "boss_1x1x1x1":
		return `
    [ 1x1x1x1 ]
       _|_
      /0 0\
     |  ^  |
      \___/
     /|   |\
    / |___| \`
	case "mountain":
		return `
       /\
      /  \
     / /\ \
    / /  \ \
   / /    \ \`
	case "village":
		return `
      _      _      _
    _| |_  _| |_  _| |_
   |     ||     ||     |
   |  _  ||  _  ||  _  |
   | | | || | | || | | |
   |_|_|_||_|_|_||_|_|_|`
	case "shop":
		return `
    _________________
   |  $   SHOP   $  |
   |  [ ] [ ] [ ]  |
   |  [ ] [ ] [ ]  |
   |_______________|`
	case "garden":
		return `
      \ | /
    '-.ooo.-'
   -- oooo --
    .-'ooo'-.
      / | \
     \ | /
      \|/`
	case "forest":
		return `
      v .   ._, |_  .,
   -._\/  .  \ /    |/_
       \\  _\, y | \//
 _\_.___\\, \\/ -.\||
   ` + "`" + `7-,--.` + "`" + `._||  / / ,
   /'     ` + "`" + `  ` + "`" + ` y  , y | /
  /            |  /_ L ,
  |            |   |
  |   _    _   |   |
  |  |_|  |_|  |   |
  |____________|___|`
	case "river":
		return `
   ~~~~~~~~~~~~~~~
  ~~~~~~~~~~~~~~~~~
 ~~~~~~~~~~~~~~~~~~~
       _ ___ _
     _(_)_(_)_
    (_) (_) (_)
   ~~~~~~~~~~~~~~~`
	case "cave":
		return `
      __________
     /          \
    /   ______   \
   /   /      \   \
  /   /        \   \
 /___/          \___\`
	case "zoo":
		return `
    _______________
   |  [ MONSTER ]  |
   |   (O)   (O)   |
   |      <        |
   |    \___/      |
   |_______________|`
	case "party_sun_gate":
		return `
      \  |  /
   ---  \ | /  ---
        .-*-.
   ---  / | \  ---
      /  |  \`
	case "party_sun_depths":
		return `
     .-^^^^^-.
    /  SUN   \
   /  DUNES   \
   \   ***    /
    '._   _.-'
       'v'`
	case "party_sun_crown":
		return `
      /\  /\
     /  \/  \
    |  SUN   |
    | CROWN  |
     \      /
      \____/`
	case "sage_hut":
		return `
      /\____/\
     /  /\  \ \
    /__/  \__\ \
    \  \__/  / /
     \______/_/`
	case "harbor":
		return `
    ~~~  __  ~~~
   ~~~  /  \  ~~~
      _/____\_
  ___/  /\   \___
   \____/  \____/`
	case "lighthouse":
		return `
       /\ 
      /  \
     / || \
    /  ||  \
   /___||___\
      /__\`
	case "tide_caves":
		return `
   ~~~~~~~~~~~~~
  /  /\    /\  \
 /__/  \__/  \__\
 \  \   ____   /
  \__\_/____\_/`
	case "old_ruins":
		return `
    .-^-.
   /_/_\_\
  |  _   |
  | |_|  |
  |_____ |
    |_|`
	case "moonwell":
		return `
     .-.
   _(   )_
  (_______)
    \   /
     \_/
    ~~~~~`
	case "observatory":
		return `
      /\ 
     /  \
    /____\
   /| __ |\
  /_|_||_|_\
    /_/\_\
		`
	default:
		return ""
	}
}

func DrawEnemyHealthBar(name string, current, max int) {
	DrawHealthBar(fmt.Sprintf("%-14s", name), current, max, ColorRed)
}

func ShowTutorial(lang string) {
	ClearScreen()
	fmt.Println(ColorCyan + ColorBold + TranslateText(lang, "=== HOW TO PLAY ===") + ColorReset)
	fmt.Println("\n" + TranslateText(lang, "Welcome, adventurer! Magic Adventure is a multiplayer text RPG."))
	fmt.Println(TranslateText(lang, "Created by Dante Gomes with assistance from Gemini and Codex."))
	fmt.Println("\n1. " + ColorBold + TranslateText(lang, "Navigation") + ColorReset + ": " + TranslateText(lang, "Choose numbered options to move between locations."))
	fmt.Println("2. " + ColorBold + TranslateText(lang, "Combat") + ColorReset + ": " + TranslateText(lang, "Fight monsters to earn XP and Gold. Your stats increase automatically."))
	fmt.Println("3. " + ColorBold + TranslateText(lang, "Multiplayer") + ColorReset + ": " + TranslateText(lang, "You can see other online players and their locations."))
	fmt.Println("4. " + ColorBold + TranslateText(lang, "Party Play") + ColorReset + ": " + TranslateText(lang, "Leaders can heal, rally, and guard the party for real combat advantages."))
	fmt.Println("5. " + ColorBold + TranslateText(lang, "Party Quests") + ColorReset + ": " + TranslateText(lang, "Parties can start cooperative quests that reward stronger shared buffs."))
	fmt.Println("6. " + ColorBold + TranslateText(lang, "Exploration") + ColorReset + ": " + TranslateText(lang, "The village, forest, river, harbor, lighthouse, tide caves, ruins, observatory, and sage hut all hide different rewards."))
	fmt.Println("7. " + ColorBold + TranslateText(lang, "Wise Man") + ColorReset + ": " + TranslateText(lang, "You can wire him to Gemini or Cloudflare and reconfigure him later from the Settings menu."))
	fmt.Println("8. " + ColorBold + TranslateText(lang, "Progression") + ColorReset + ": " + TranslateText(lang, "Buy items in the Shop and complete biomes to find the Dragon."))
	fmt.Println("\n" + ColorYellow + TranslateText(lang, "Tip: Keep an eye on your Health! Use Potions to stay alive.") + ColorReset)
}

func PrintHeader(gs *GameState, unreadNotifications int, lang string, wiseManStatus string) {
	fmt.Printf("%s=== %s %d | %s | %s %d ===%s\n", ColorCyan+ColorBold, TranslateText(lang, "Slot"), gs.ID, gs.PlayerName, TranslateText(lang, "Lvl"), gs.Level, ColorReset)
	DrawHealthBar(TranslateText(lang, "Health"), gs.Health, gs.MaxHealth, ColorGreen)
	DrawHealthBar(TranslateText(lang, "Exp   "), gs.Experience, gs.RequiredXP, ColorBlue)

	items := []string{}
	if gs.HasSteelSword {
		items = append(items, TranslateText(lang, "Steel Sword"))
	}
	if gs.HasGardenKit {
		items = append(items, TranslateText(lang, "Garden Kit"))
	}
	if gs.HasMeat {
		items = append(items, TranslateText(lang, "Meat"))
	}
	if gs.HasWater {
		items = append(items, TranslateText(lang, "Water"))
	}
	if gs.HasFurCoat {
		items = append(items, TranslateText(lang, "Fur Coat"))
	}
	if gs.HasSunAmulet {
		items = append(items, TranslateText(lang, "Sun Amulet"))
	}
	if gs.HasIceCrystal {
		items = append(items, TranslateText(lang, "Ice Crystal"))
	}
	if gs.HasForsakenBlade {
		items = append(items, TranslateText(lang, "Forsaken Blade"))
	}
	if gs.HasMoonPearl {
		items = append(items, TranslateText(lang, "Moon Pearl"))
	}
	if gs.HasSageBlessing {
		items = append(items, TranslateText(lang, "Sage Blessing"))
	}
	if gs.HasRuinsToken {
		items = append(items, TranslateText(lang, "Ruins Token"))
	}
	if gs.HasMoonCharm {
		items = append(items, TranslateText(lang, "Moon Charm"))
	}
	if gs.HasStarMap {
		items = append(items, TranslateText(lang, "Star Map"))
	}
	if gs.HasTriptychBlessing {
		items = append(items, TranslateText(lang, "Triptych Blessing"))
	}

	inventoryStr := strings.Join(items, ", ")
	if inventoryStr == "" {
		inventoryStr = TranslateText(lang, "None")
	}

	fmt.Printf("%s: %s%d%s | %s: %s%d%s | %s: [%s]\n", TranslateText(lang, "Gold"), ColorYellow, gs.Gold, ColorReset, TranslateText(lang, "SP"), ColorCyan, gs.SkillPoints, ColorReset, TranslateText(lang, "Items"), inventoryStr)
	fmt.Printf("%s: %s%s%s\n", TranslateText(lang, "Location"), ColorBlue, TranslateRoomDescription(lang, gs.CurrentRoomID, gs.CurrentRoomID), ColorReset)
	if unreadNotifications > 0 {
		fmt.Printf("%s%s:%s %d %s\n", ColorCyan, TranslateText(lang, "Inbox"), ColorReset, unreadNotifications, TranslateText(lang, "unread notification(s)"))
	} else {
		fmt.Printf("%s%s:%s %s\n", ColorCyan, TranslateText(lang, "Inbox"), ColorReset, TranslateText(lang, "empty"))
	}
	if gs.PartySupport > 0 {
		fmt.Printf("%s%s:%s +%d %s\n", ColorBlue, TranslateText(lang, "Party bonus"), ColorReset, gs.PartySupport*5, TranslateText(lang, "damage per hit"))
	}
	if strings.TrimSpace(wiseManStatus) != "" {
		fmt.Printf("%s%s:%s %s\n", ColorCyan, TranslateText(lang, "Wise Man"), ColorReset, TranslateText(lang, wiseManStatus))
	}
	if strings.TrimSpace(gs.PartyQuestName) != "" {
		fmt.Printf("%s%s:%s %s [%s]\n", ColorBlue, TranslateText(lang, "Quest"), ColorReset, TranslateText(lang, gs.PartyQuestName), TranslateText(lang, gs.PartyQuestStatus))
		if phaseName, phaseProgress, phaseGoal := partyQuestProgressSummary(gs.PartyQuestName, gs.PartyQuestPhase, gs.PartyQuestPhaseProgress, gs.PartyQuestPhaseGoal); strings.TrimSpace(phaseName) != "" {
			fmt.Printf("%s%s:%s %s (%d/%d)\n", ColorBlue, TranslateText(lang, "Phase"), ColorReset, TranslateText(lang, phaseName), phaseProgress, phaseGoal)
		}
	}
	fmt.Println(strings.Repeat("-", 40))
}

func PrintRecentAction(gs *GameState) {
	if gs == nil {
		return
	}

	msg := strings.TrimSpace(gs.LastActionMessage)
	if msg == "" {
		return
	}

	fmt.Printf("%s%s:%s %s\n", ColorYellow, TranslateText(gs.Language, "Last result"), ColorReset, TranslateText(gs.Language, msg))
	fmt.Println(strings.Repeat("-", 40))
}
