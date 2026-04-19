package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"magicadventure/game"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	client, err := game.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		client.Disconnect(ctx)
	}()

	reader := bufio.NewReader(os.Stdin)
	game.ClearScreen()
	fmt.Println(game.ColorCyan + "Welcome to Magic Adventure 6.0!" + game.ColorReset)

	saves, _ := game.GetAllSaves()
	saveMap := make(map[int]game.GameState)
	for _, s := range saves {
		saveMap[s.ID] = s
	}

	fmt.Println("\n--- CHARACTER SLOTS ---")
	for i := 1; i <= 5; i++ {
		if s, ok := saveMap[i]; ok {
			fmt.Printf("[%d] %s (Lvl %d) - %s\n", i, s.PlayerName, s.Level, s.CurrentRoomID)
		} else {
			fmt.Printf("[%d] -- EMPTY SLOT --\n", i)
		}
	}

	var gs *game.GameState
	var slotID int

	for {
		fmt.Print("\nSelect a slot (1-5): ")
		input, _ := reader.ReadString('\n')
		slotID, err = strconv.Atoi(strings.TrimSpace(input))
		if err == nil && slotID >= 1 && slotID <= 5 {
			break
		}
		fmt.Println("Invalid slot. Please enter 1-5.")
	}

	if s, ok := saveMap[slotID]; ok {
		fmt.Printf("\nSelected: %s. \n1) Continue\n2) Overwrite (New Game)\nSelect: ", s.PlayerName)
		choice, _ := reader.ReadString('\n')
		if strings.TrimSpace(choice) == "1" {
			gs = &s
		}
	}

	if gs == nil {
		gs = createNewCharacter(reader, slotID)
	}

	gs.World = game.CreateWorld()
	engine := game.NewEngine(gs, client)
	engine.Run()
}

func createNewCharacter(reader *bufio.Reader, slotID int) *game.GameState {
	fmt.Print("\nEnter your character's name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Hero"
	}

	return &game.GameState{
		ID:            slotID,
		PlayerName:    name,
		CurrentRoomID: "village",
		Health:        100,
		MaxHealth:     100,
		IsGameOver:    false,
		Gold:          50,
		MonsterHealth: 50,
		Level:         1,
		Experience:    0,
		RequiredXP:    100,
		SkillPoints:   0,
		BaseDamage:    5,
	}
}
