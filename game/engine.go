package game

import (
	"bufio"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Engine struct {
	State *GameState
	DB    *mongo.Client
}

func NewEngine(initialState *GameState, db *mongo.Client) *Engine {
	return &Engine{State: initialState, DB: db}
}

func (e *Engine) Run() {
	reader := bufio.NewReader(os.Stdin)

	if !e.State.HasSeenTutorial {
		ShowTutorial()
		e.State.HasSeenTutorial = true
		SaveGame(e.State)
		fmt.Println("\nPress Enter to start your adventure...")
		reader.ReadString('\n')
	}

	for !e.State.IsGameOver {
		ClearScreen()
		PrintHeader(e.State)

		// Multiplayer: Show online players
		onlinePlayers, _ := GetOnlinePlayers(e.State)
		if len(onlinePlayers) > 0 {
			fmt.Printf("%sOnline Players:%s ", ColorYellow, ColorReset)
			for i, p := range onlinePlayers {
				fmt.Printf("%s (%s)", p.PlayerName, p.CurrentRoomID)
				if i < len(onlinePlayers)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println()
		}

		art := GetArt(e.State.CurrentRoomID)
		if art != "" {
			fmt.Println(ColorYellow + art + ColorReset)
		}

		currentRoom := e.State.World[e.State.CurrentRoomID]
		fmt.Printf("\n%s\n\n", currentRoom.Description)

		// Filter available actions
		var availableActions []Action
		for _, action := range currentRoom.Actions {
			if action.IsAvailable == nil || action.IsAvailable(e.State) {
				availableActions = append(availableActions, action)
			}
		}

		// Multiplayer: Dynamic interactions with players in the same room
		for _, p := range onlinePlayers {
			if p.CurrentRoomID == e.State.CurrentRoomID {
				target := p // capture for closure
				availableActions = append(availableActions, Action{
					ID:          "interact_" + target.PlayerName,
					Description: fmt.Sprintf("%sInteract with %s%s", ColorCyan, target.PlayerName, ColorReset),
					Result: func(gs *GameState) string {
						interactions := []string{
							"You wave at %s. They look busy, but they smile back!",
							"You challenge %s to a duel! They laugh and say 'Not today, rookie'.",
							"You try to pickpocket %s... but they notice and slap your hand!",
							"You share a story about the Dragon with %s. They look impressed.",
							"You and %s perform a synchronized dance. The birds stop to watch.",
						}
						rand.Seed(time.Now().UnixNano())
						res := interactions[rand.Intn(len(interactions))]
						return fmt.Sprintf(res, target.PlayerName)
					},
				})
			}
		}

		fmt.Println("What do you want to do?")
		for i, action := range availableActions {
			fmt.Printf("%d) %s\n", i+1, action.Description)
		}

		fmt.Print("\nSelect an option: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			SaveGame(e.State)
			fmt.Println("Game Saved. Goodbye!")
			return
		}

		choice, err := strconv.Atoi(input)
		if err == nil && choice > 0 && choice <= len(availableActions) {
			action := availableActions[choice-1]
			resultMsg := action.Result(e.State)
			
			// Check Level Up after action
			levelUpMsg := e.State.CheckLevelUp()
			resultMsg += levelUpMsg
			
			// Show result briefly
			ClearScreen()
			PrintHeader(e.State)
			fmt.Printf("\n%s\n", resultMsg)
			
			fmt.Println("\nPress Enter to continue...")
			reader.ReadString('\n')
			
			SaveGame(e.State)
		} else {
			fmt.Println("Invalid selection. Press Enter to try again.")
			reader.ReadString('\n')
		}
	}

	if e.State.IsGameOver {
		ClearScreen()
		fmt.Println(ColorRed + "G A M E   O V E R" + ColorReset)
		fmt.Println("\nTry again? (yes/no)")
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) == "yes" {
			e.restart()
			e.Run()
		}
	}
}

func (e *Engine) restart() {
	e.State.IsGameOver = false
	e.State.CurrentRoomID = "village"
	e.State.Health = 100
	e.State.MaxHealth = 100
	e.State.Level = 1
	e.State.Experience = 0
	e.State.RequiredXP = 100
	e.State.Gold = 50
	e.State.SkillPoints = 0
	e.State.BaseDamage = 5
	e.State.HasGardenKit = false
	e.State.HasMeat = false
	e.State.HasSteelSword = false
	e.State.IsTrollDead = false
	e.State.IsBeastDead = false
	e.State.IsFrostGiantDead = false
	e.State.Is1x1x1x1Dead = false
	SaveGame(e.State)
}
