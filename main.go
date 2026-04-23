package main

import (
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"magicadventure/game"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Load .env file if it exists
	godotenv.Load()

	store, err := game.InitDB()
	if err != nil {
		log.Fatalf("Could not initialize save store: %v", err)
	}
	defer store.Close()

	online, err := game.InitOnline()
	if err != nil {
		log.Printf("Online features unavailable: %v", err)
	}
	if online != nil {
		defer online.Close()
	}
	wiseman := game.NewWiseManService(store)

	reader := bufio.NewReader(os.Stdin)
	game.ClearScreen()
	lang := game.EnsureLanguageConfigured(reader, store)
	fmt.Println(game.ColorCyan + game.TranslateText(lang, "Welcome to Magic Adventure 6.0!") + game.ColorReset)
	fmt.Println(game.TranslateText(lang, "Created by Dante Gomes with assistance from Gemini and Codex."))
	if msg, err := wiseman.EnsureConfigured(reader); err != nil {
		log.Printf("Wise man setup unavailable: %v", err)
	} else if strings.TrimSpace(msg) != "" {
		fmt.Println(game.TranslateText(lang, msg))
		fmt.Println("\n" + game.TranslateText(lang, "Press Enter to continue..."))
		reader.ReadString('\n')
	}

	saves, err := store.GetAllSaves()
	if err != nil {
		log.Printf("Could not load saves: %v", err)
	}
	saveMap := make(map[int]game.GameState)
	for _, s := range saves {
		saveMap[s.ID] = s
	}

	unreadCounts := map[int]int{}
	if online != nil && online.Enabled() {
		for i, s := range saveMap {
			tmp := s
			if notifications, err := online.GetNotifications(&tmp, tmp.LastNotificationSeen, 100); err == nil {
				unreadCounts[i] = len(notifications)
			}
		}
	}

	fmt.Println("\n" + game.TranslateText(lang, "--- CHARACTER SLOTS ---"))
	for i := 1; i <= 5; i++ {
		if s, ok := saveMap[i]; ok {
			if unreadCounts[i] > 0 {
				fmt.Printf("[%d] %s (%s %d) - %s | %s: %d %s\n", i, s.PlayerName, game.TranslateText(lang, "Lvl"), s.Level, s.CurrentRoomID, game.TranslateText(lang, "Inbox"), unreadCounts[i], game.TranslateText(lang, "unread"))
			} else {
				fmt.Printf("[%d] %s (%s %d) - %s\n", i, s.PlayerName, game.TranslateText(lang, "Lvl"), s.Level, s.CurrentRoomID)
			}
		} else {
			fmt.Printf("[%d] %s\n", i, game.TranslateText(lang, "-- EMPTY SLOT --"))
		}
	}

	var gs *game.GameState
	var slotID int

	for {
		fmt.Print("\n" + game.TranslateText(lang, "Select a slot (1-5): "))
		input, _ := reader.ReadString('\n')
		slotID, err = strconv.Atoi(strings.TrimSpace(input))
		if err == nil && slotID >= 1 && slotID <= 5 {
			break
		}
		fmt.Println(game.TranslateText(lang, "Invalid slot. Please enter 1-5."))
	}

	if s, ok := saveMap[slotID]; ok {
		fmt.Printf("\n%s: %s. \n1) %s\n2) %s\n%s", game.TranslateText(lang, "Selected"), s.PlayerName, game.TranslateText(lang, "Continue"), game.TranslateText(lang, "Overwrite (New Game)"), game.TranslateText(lang, "Select: "))
		choice, _ := reader.ReadString('\n')
		if strings.TrimSpace(choice) == "1" {
			gs = &s
			lang = game.NormalizeLanguage(s.Language)
		}
	}

	if gs == nil {
		gs = createNewCharacter(reader, slotID, lang)
	}
	gs.Language = lang

	gs.World = game.CreateWorld()
	if online != nil && online.Enabled() {
		_ = online.SeedStarterRelics(gs)
		_ = online.SyncPresence(gs)
	}

	engine := game.NewEngine(gs, store, online, wiseman)
	engine.Run()
}

func createNewCharacter(reader *bufio.Reader, slotID int, lang string) *game.GameState {
	fmt.Print("\n" + game.TranslateText(lang, "Enter your character's name: "))
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Hero"
	}

	return game.NewGameState(slotID, name)
}
