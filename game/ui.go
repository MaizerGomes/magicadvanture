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
	default:
		return ""
	}
}

func DrawEnemyHealthBar(name string, current, max int) {
	DrawHealthBar(fmt.Sprintf("%-14s", name), current, max, ColorRed)
}

func ShowTutorial() {
	ClearScreen()
	fmt.Println(ColorCyan + ColorBold + "=== HOW TO PLAY ===" + ColorReset)
	fmt.Println("\nWelcome, adventurer! Magic Adventure is a multiplayer text RPG.")
	fmt.Println("\n1. " + ColorBold + "Navigation" + ColorReset + ": Choose numbered options to move between locations.")
	fmt.Println("2. " + ColorBold + "Combat" + ColorReset + ": Fight monsters to earn XP and Gold. Your stats increase automatically.")
	fmt.Println("3. " + ColorBold + "Multiplayer" + ColorReset + ": You can see other online players and their locations.")
	fmt.Println("4. " + ColorBold + "Interaction" + ColorReset + ": When in the same room as another player, you can interact with them!")
	fmt.Println("5. " + ColorBold + "Progression" + ColorReset + ": Buy items in the Shop and complete biomes to find the Dragon.")
	fmt.Println("\n" + ColorYellow + "Tip: Keep an eye on your Health! Use Potions to stay alive." + ColorReset)
}

func PrintHeader(gs *GameState) {
	fmt.Printf("%s=== %s | Lvl %d ===%s\n", ColorCyan+ColorBold, gs.PlayerName, gs.Level, ColorReset)
	DrawHealthBar("Health", gs.Health, gs.MaxHealth, ColorGreen)
	DrawHealthBar("Exp   ", gs.Experience, gs.RequiredXP, ColorBlue)
	
	items := []string{}
	if gs.HasSteelSword { items = append(items, "Steel Sword") }
	if gs.HasGardenKit { items = append(items, "Garden Kit") }
	if gs.HasMeat { items = append(items, "Meat") }
	if gs.HasWater { items = append(items, "Water") }
	if gs.HasFurCoat { items = append(items, "Fur Coat") }
	if gs.HasSunAmulet { items = append(items, "Sun Amulet") }
	if gs.HasIceCrystal { items = append(items, "Ice Crystal") }
	if gs.HasForsakenBlade { items = append(items, "Forsaken Blade") }
	
	inventoryStr := strings.Join(items, ", ")
	if inventoryStr == "" { inventoryStr = "None" }

	fmt.Printf("Gold: %s%d%s | SP: %s%d%s | Items: [%s]\n", ColorYellow, gs.Gold, ColorReset, ColorCyan, gs.SkillPoints, ColorReset, inventoryStr)
	fmt.Printf("Location: %s%s%s\n", ColorBlue, gs.CurrentRoomID, ColorReset)
	fmt.Println(strings.Repeat("-", 40))
}
