package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Engine struct {
	State      *GameState
	Store      *Store
	Online     *OnlineService
	WiseMan    *WiseManService
	feed       *RoomFeed
	feedRoomID string
}

type MenuSection struct {
	Title   string
	Actions []Action
}

func NewEngine(initialState *GameState, store *Store, online *OnlineService, wiseMan *WiseManService) *Engine {
	return &Engine{State: initialState, Store: store, Online: online, WiseMan: wiseMan}
}

func (e *Engine) RunMillionaireGame(reader *bufio.Reader) {
	// 1. Sync Streak
	streak, _ := e.Online.GetCurrentMillionaireStreak()
	if streak != nil {
		if e.State.MillionaireStreakID != streak.ID {
			e.State.MillionairePoints = 0
			e.State.LifelineAudienceUsed = false
			e.State.LifelineWiseManUsed = false
			e.State.MillionaireStreakID = streak.ID
			e.persist()
		}
		resetMsg, _ := e.Online.CheckAndResetMillionaireStreak(e.State)
		if resetMsg != "" {
			fmt.Println(ColorCyan + TranslateText(e.State.Language, resetMsg) + ColorReset)
			time.Sleep(2 * time.Second)
			// Refresh streak info after reset
			streak, _ = e.Online.GetCurrentMillionaireStreak()
		}
	}

	for {
		ClearScreen()
		fmt.Println(ColorYellow + "=== WHO WANTS TO BE A MILLIONAIRE ===" + ColorReset)
		if streak != nil {
			fmt.Printf("Streak Ends: %s\n", streak.EndTime.Format("2006-01-02 15:04"))
		}
		fmt.Printf("Current Streak Points: %d\n", e.State.MillionairePoints)
		fmt.Println()

		// Show Leaderboard
		if streak != nil {
			scores, _ := e.Online.GetMillionaireLeaderboard(streak.ID)
			fmt.Println("--- TOP 20 LEADERS ---")
			for i, s := range scores {
				color := ""
				if i == 0 {
					color = ColorCyan + "[1st] "
				} else if i == 1 {
					color = ColorYellow + "[2nd] "
				} else if i == 2 {
					color = ColorMagenta + "[3rd] "
				} else {
					color = fmt.Sprintf("[%d] ", i+1)
				}
				fmt.Printf("%s%s - %d pts%s\n", color, s.PlayerName, s.Score, ColorReset)
			}
			fmt.Println()
		}

		fmt.Println("1) Answer a Question")
		fmt.Println("2) Exit Millionaire Hall")
		fmt.Print("\nSelection: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "2" || input == "exit" {
			return
		}

		if input == "1" {
			for {
				if !e.handleMillionaireChallenge(reader) {
					break
				}
			}
		}
	}
}

func (e *Engine) handleMillionaireChallenge(reader *bufio.Reader) bool {
	challenge, err := e.WiseMan.GenerateMillionaireChallenge(e.State)
	if err != nil {
		fmt.Println(ColorRed + "AI error: " + err.Error() + ColorReset)
		time.Sleep(2 * time.Second)
		return false
	}

	questionID := fmt.Sprintf("q-%d", time.Now().UnixNano())

	for {
		ClearScreen()
		fmt.Println(ColorYellow + "--- QUESTION ---" + ColorReset)
		fmt.Println(challenge.Question)
		for i, opt := range challenge.Options {
			fmt.Printf("%d) %s\n", i+1, opt)
		}
		fmt.Println()

		fmt.Println("Lifelines:")
		if !e.State.LifelineAudienceUsed {
			fmt.Println("A) Ask the Audience (Online Players)")
		} else {
			fmt.Println("A) [USED]")
		}
		if !e.State.LifelineWiseManUsed {
			fmt.Println("W) Ask the Wise Man (AI + Web Search)")
		} else {
			fmt.Println("W) [USED]")
		}

		fmt.Print("\nYour answer (1-4) or Lifeline (A/W): ")
		input, _ := reader.ReadString('\n')
		input = strings.ToUpper(strings.TrimSpace(input))

		if input == "A" && !e.State.LifelineAudienceUsed {
			e.State.LifelineAudienceUsed = true
			e.persist()
			_ = e.Online.BroadcastMillionaireQuestion(e.State, questionID, challenge.Question)
			fmt.Println("Broadcasting question to all online players...")
			fmt.Println("Waiting 15 seconds for responses...")
			time.Sleep(15 * time.Second)
			stats, _ := e.Online.GetAudienceStats(questionID)
			fmt.Println("\nAudience Results:")
			for i := 1; i <= 4; i++ {
				fmt.Printf("Option %d: %d votes\n", i, stats[i])
			}
			fmt.Println("\nPress Enter to return to the question.")
			reader.ReadString('\n')
			continue
		}

		if input == "W" && !e.State.LifelineWiseManUsed {
			e.State.LifelineWiseManUsed = true
			e.persist()
			fmt.Println("Consulting the Wise Man (searching the web)...")
			answer, _ := e.WiseMan.AskWiseManWithSearch(e.State, challenge.Question)
			fmt.Printf("\nWise Man says: %s\n", answer)
			fmt.Println("\nPress Enter to return to the question.")
			reader.ReadString('\n')
			continue
		}

		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 1 && choice <= 4 {
			if choice == challenge.CorrectIndex+1 {
				points := 100 // Harder questions give more points
				e.State.MillionairePoints += points
				e.State.Gold += 50
				_ = e.Online.SyncMillionaireScore(e.State, e.State.MillionairePoints)
				fmt.Println(ColorGreen + "CORRECT! +100 points and 50 gold." + ColorReset)
				e.persist()
				time.Sleep(2 * time.Second)
				return true
			} else {
				e.State.Gold -= 50
				if e.State.Gold < 0 {
					e.State.Gold = 0
				}
				fmt.Println(ColorRed + "WRONG! -50 gold." + ColorReset)
				fmt.Printf("The correct answer was: %s\n", challenge.CorrectAnswer)
				e.persist()
				time.Sleep(3 * time.Second)
				return false
			}
		}
	}
}

func (e *Engine) RunMiniGame(gameName, choice string) string {
	if gameName == "guessing_game" {
		if e.State.Gold < 10 {
			return TranslateText(e.State.Language, "You don't have enough gold!")
		}
		e.State.Gold -= 10
		winningCup := strconv.Itoa(rand.Intn(3) + 1)
		if choice == winningCup {
			e.State.Gold += 30
			return fmt.Sprintf("%s! %s.",
				TranslateText(e.State.Language, "Correct"),
				TranslateText(e.State.Language, "You guessed correctly! You won 30 gold!"))
		}
		return fmt.Sprintf("%s. %s",
			TranslateText(e.State.Language, "Wrong cup! Better luck next time."),
			fmt.Sprintf(TranslateText(e.State.Language, "The pearl was in cup %s."), winningCup))
	}
	return "Unknown mini-game"
}

func (e *Engine) RunEducationalChallenge(reader *bufio.Reader) bool {
	ClearScreen()
	fmt.Println(TranslateText(e.State.Language, "Consulting the Wise Man for a challenge..."))
	challenge, err := e.WiseMan.GenerateEducationalChallenge(e.State)
	if err != nil {
		fmt.Println(ColorRed + "Error connecting to the Wise Man: " + err.Error() + ColorReset)
		time.Sleep(2 * time.Second)
		// We'll proceed with a simple default math question if AI fails completely and no fallback list
		challenge = EducationalChallenge{
			Question:      "What is 5 + 5?",
			Options:       []string{"8", "10", "12", "15"},
			CorrectIndex:  1,
			CorrectAnswer: "10",
		}
	}

	ClearScreen()
	fmt.Println(ColorCyan + TranslateText(e.State.Language, "The Wise Man appears with a challenge!") + ColorReset)
	fmt.Printf("\n%s\n", challenge.Question)
	for i, opt := range challenge.Options {
		fmt.Printf("%d) %s\n", i+1, opt)
	}

	fmt.Print("\n" + TranslateText(e.State.Language, "Choose the correct answer: "))
	input, _ := reader.ReadString('\n')
	choice, err := strconv.Atoi(strings.TrimSpace(input))

	if err == nil && choice == challenge.CorrectIndex+1 {
		e.State.Gold += 50
		e.State.ConsecutiveWrongAnswers = 0
		fmt.Println(ColorGreen + TranslateText(e.State.Language, "Correct! You earned 50 gold.") + ColorReset)
		time.Sleep(2 * time.Second)
		return true
	}

	// Wrong answer
	e.State.Gold -= 50
	if e.State.Gold < 0 {
		e.State.Gold = 0
	}
	e.State.ConsecutiveWrongAnswers++
	fmt.Println(ColorRed + TranslateText(e.State.Language, "Wrong answer! You lost 50 gold.") + ColorReset)
	fmt.Printf("%s: %s\n", TranslateText(e.State.Language, "The correct answer was"), challenge.CorrectAnswer)

	if e.State.ConsecutiveWrongAnswers >= 3 {
		e.State.Health -= 10
		e.State.ConsecutiveWrongAnswers = 0
		fmt.Println(ColorRed + TranslateText(e.State.Language, "That's 3 wrong in a row! You lost 10 Health.") + ColorReset)
		if e.State.Health <= 0 {
			e.State.Health = 0
			e.State.IsGameOver = true
		}
	}

	time.Sleep(3 * time.Second)
	return false
}

func (e *Engine) Run() {
	reader := bufio.NewReader(os.Stdin)
	rand.Seed(time.Now().UnixNano())

	if !e.State.HasSeenTutorial {
		ShowTutorial(e.State.Language)
		e.State.HasSeenTutorial = true
		e.persist()
		fmt.Println("\n" + TranslateText(e.State.Language, "Press Enter to start your adventure..."))
		reader.ReadString('\n')
	}

	for !e.State.IsGameOver {
		e.ensureFeed()
		ClearScreen()

		nearbyPlayers, messages, whispers, inbox := e.feedSnapshot()
		party, partyMembers := e.partySnapshot()
		if len(partyMembers) > 1 {
			e.State.PartySupport = len(partyMembers) - 1
		} else {
			e.State.PartySupport = 0
		}
		PrintHeader(e.State, len(inbox), e.State.Language, e.wiseManStatus())
		e.printConnectivity()
		e.printPartySummary(party, partyMembers)
		if len(nearbyPlayers) > 0 {
			fmt.Printf("%sOnline Players:%s ", ColorYellow, ColorReset)
			for i, p := range nearbyPlayers {
				fmt.Printf("%s (%s)", p.PlayerName, p.CurrentRoomID)
				if i < len(nearbyPlayers)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println()
		}
		e.printRoomMessages(messages)
		e.printWhispers(whispers)
		e.printInboxSummary(inbox)
		PrintRecentAction(e.State)

		art := GetArt(e.State.CurrentRoomID)
		if art != "" {
			fmt.Println(ColorYellow + art + ColorReset)
		}

		currentRoom := e.State.World[e.State.CurrentRoomID]
		fmt.Printf("\n%s\n\n", TranslateRoomDescription(e.State.Language, currentRoom.ID, currentRoom.Description))

		var movementActions []Action
		var roomActions []Action
		for _, action := range currentRoom.Actions {
			if action.IsAvailable != nil && !action.IsAvailable(e.State) {
				continue
			}
			if action.Type == MoveAction {
				movementActions = append(movementActions, action)
				continue
			}
			roomActions = append(roomActions, action)
		}

		sections := []MenuSection{}
		if len(movementActions) > 0 {
			sections = append(sections, MenuSection{Title: "Movement", Actions: movementActions})
		}
		if len(roomActions) > 0 {
			sections = append(sections, MenuSection{Title: "Room", Actions: roomActions})
		}
		sections = append(sections, e.buildOnlineSections(reader, nearbyPlayers, messages, whispers, inbox, party, partyMembers)...)
		sections = append(sections, e.buildSupportSections(reader)...)
		if guidance := e.buildWiseManGuidanceSection(reader, sections); guidance != nil {
			sections = append(sections, *guidance)
		}

		availableActions := e.printMenuSections(sections, e.State.Language)
		if len(availableActions) == 0 {
			fmt.Println(TranslateText(e.State.Language, "No actions are available right now."))
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Print("\n" + TranslateText(e.State.Language, "Select an option: "))
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		input = strings.TrimSpace(input)

		if input == "i" || input == "inbox" {
			action := e.findActionByID(availableActions, "open_inbox")
			if action == nil {
				fmt.Println(TranslateText(e.State.Language, "Inbox is not available right now."))
				time.Sleep(2 * time.Second)
				continue
			}
			resultMsg := action.Result(e.State)
			levelUpMsg := e.State.CheckLevelUp()
			resultMsg += levelUpMsg
			e.postActionOnlineEvent(*action, resultMsg)
			inbox = nil

			ClearScreen()
			PrintHeader(e.State, len(inbox), e.State.Language, e.wiseManStatus())
			e.printConnectivity()
			e.printPartySummary(party, partyMembers)
			e.printRoomMessages(messages)
			e.printWhispers(whispers)
			e.printInboxSummary(inbox)
			PrintRecentAction(e.State)
			fmt.Printf("\n%s\n", TranslateText(e.State.Language, resultMsg))
			e.State.LastActionMessage = resultMsg
			e.persist()
			time.Sleep(2 * time.Second)
			continue
		}

		if input == "exit" || input == "quit" {
			e.persist()
			fmt.Println(TranslateText(e.State.Language, "Game Saved. Goodbye!"))
			return
		}

		choice, err := strconv.Atoi(input)
		if err == nil && choice > 0 && choice <= len(availableActions) {
			action := availableActions[choice-1]
			resultMsg := action.Result(e.State)

			if strings.HasPrefix(resultMsg, "CHALLENGE_REQUIRED:") {
				targetRoom := strings.TrimPrefix(resultMsg, "CHALLENGE_REQUIRED:")
				if targetRoom == "millionaire" {
					e.RunMillionaireGame(reader)
					resultMsg = "" // RunMillionaireGame handles its own UI
				} else if e.RunEducationalChallenge(reader) {
					e.State.Transition(targetRoom)
					resultMsg = TranslateText(e.State.Language, "You passed the challenge and entered!")
				} else {
					resultMsg = TranslateText(e.State.Language, "You failed the Wise Man's challenge.")
				}
			}

			if strings.HasPrefix(resultMsg, "MINI_GAME:") {
				parts := strings.Split(resultMsg, ":")
				if len(parts) >= 3 {
					gameName := parts[1]
					choice := parts[2]
					resultMsg = e.RunMiniGame(gameName, choice)
				}
			}

			// Check Level Up after action
			levelUpMsg := e.State.CheckLevelUp()
			resultMsg += levelUpMsg
			e.postActionOnlineEvent(action, resultMsg)
			if action.ID == "open_inbox" {
				inbox = nil
			}

			// Show result briefly
			ClearScreen()
			if len(nearbyPlayers) > 0 {
				fmt.Printf("%sOnline Players:%s ", ColorYellow, ColorReset)
				for i, p := range nearbyPlayers {
					fmt.Printf("%s (%s)", p.PlayerName, p.CurrentRoomID)
					if i < len(nearbyPlayers)-1 {
						fmt.Print(", ")
					}
				}
				fmt.Println()
			}
			PrintHeader(e.State, len(inbox), e.State.Language, e.wiseManStatus())
			e.printConnectivity()
			e.printPartySummary(party, partyMembers)
			e.printRoomMessages(messages)
			e.printWhispers(whispers)
			e.printInboxSummary(inbox)
			PrintRecentAction(e.State)
			fmt.Printf("\n%s\n", resultMsg)
			e.State.LastActionMessage = resultMsg
			e.persist()
			time.Sleep(2 * time.Second)
		} else {
			fmt.Println(TranslateText(e.State.Language, "Invalid selection."))
			time.Sleep(2 * time.Second)
		}
	}

	if e.State.IsGameOver {
		ClearScreen()
		fmt.Println(ColorRed + TranslateText(e.State.Language, "G A M E   O V E R") + ColorReset)
		if speech := e.wiseManDeathSpeech(); speech != "" {
			fmt.Println("\n" + TranslateText(e.State.Language, speech))
		}
		fmt.Print("\n" + TranslateText(e.State.Language, "Try again? (yes/no): "))
		input, _ := reader.ReadString('\n')
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(input)), "y") {
			e.State.ResetRun()
			e.persist()
			e.Run()
		}
	}
}

func (e *Engine) persist() {
	if e.Store != nil {
		_ = e.Store.SaveGame(e.State)
	}
	if e.Online != nil && e.Online.Enabled() {
		_ = e.Online.SyncPresence(e.State)
	}
}

func (e *Engine) ensureFeed() {
	if e.Online == nil || !e.Online.Enabled() {
		return
	}
	if e.feed == nil || e.feedRoomID != e.State.CurrentRoomID {
		if e.feed != nil {
			e.feed.Close()
		}
		e.feed = e.Online.SubscribeRoom(e.State)
		e.feedRoomID = e.State.CurrentRoomID
	}
}

func (e *Engine) feedSnapshot() (nearby []OnlinePlayer, messages []RoomMessage, whispers []WhisperMessage, inbox []InboxNotification) {
	if e.Online == nil || !e.Online.Enabled() {
		return nil, nil, nil, nil
	}
	nearby, _ = e.Online.GetOnlinePlayers(e.State)
	messages, whispers, _ = e.Online.GetRecentMessages(e.State.CurrentRoomID, e.State.OnlineID)
	inbox, _ = e.Online.GetNotifications(e.State, e.State.LastNotificationSeen, 50)
	return
}

func (e *Engine) partySnapshot() (*Party, []OnlinePlayer) {
	if e.Online == nil || !e.Online.Enabled() || e.State.PartyQuestKey == "" {
		return nil, nil
	}
	// Simplified party info for summary
	return nil, nil
}

func (e *Engine) wiseManStatus() string {
	if e.WiseMan == nil || !e.WiseMan.Configured() {
		return "Offline"
	}
	return "Online"
}

func (e *Engine) printConnectivity() {
	if e.Online == nil || !e.Online.Enabled() {
		fmt.Println(ColorRed + "Connection: Offline (Play alone)" + ColorReset)
		return
	}
	fmt.Println(ColorGreen + "Connection: Online" + ColorReset)
}

func (e *Engine) printPartySummary(party *Party, members []OnlinePlayer) {
	// Logic to print party status
}

func (e *Engine) printRoomMessages(messages []RoomMessage) {
	if len(messages) == 0 {
		return
	}
	fmt.Println(ColorBold + "--- Chat ---" + ColorReset)
	for _, m := range messages {
		fmt.Printf("%s: %s\n", m.PlayerName, m.Text)
	}
	fmt.Println()
}

func (e *Engine) printWhispers(whispers []WhisperMessage) {
	if len(whispers) == 0 {
		return
	}
	fmt.Println(ColorMagenta + "--- Private ---" + ColorReset)
	for _, w := range whispers {
		fmt.Printf("From %s: %s\n", w.PlayerName, w.Text)
	}
	fmt.Println()
}

func (e *Engine) printInboxSummary(inbox []InboxNotification) {
	if len(inbox) == 0 {
		return
	}
	fmt.Printf("%sInbox: %d unread (Type 'i' to open)%s\n", ColorBold+ColorCyan, len(inbox), ColorReset)
}

func (e *Engine) buildOnlineSections(reader *bufio.Reader, nearby []OnlinePlayer, messages []RoomMessage, whispers []WhisperMessage, inbox []InboxNotification, party *Party, members []OnlinePlayer) []MenuSection {
	if e.Online == nil || !e.Online.Enabled() {
		return nil
	}
	// Logic to build online sections like Chat, Party, etc.
	return nil
}

func (e *Engine) buildSupportSections(reader *bufio.Reader) []MenuSection {
	// Settings, etc.
	return nil
}

func (e *Engine) buildWiseManGuidanceSection(reader *bufio.Reader, sections []MenuSection) *MenuSection {
	// Wise man talk
	return nil
}

func (e *Engine) printMenuSections(sections []MenuSection, lang string) []Action {
	var allActions []Action
	for _, s := range sections {
		fmt.Printf("\n--- %s ---\n", TranslateText(lang, s.Title))
		for _, a := range s.Actions {
			allActions = append(allActions, a)
			fmt.Printf("%d) %s\n", len(allActions), TranslateActionDescription(lang, a.Description))
		}
	}
	return allActions
}

func (e *Engine) findActionByID(actions []Action, id string) *Action {
	for _, a := range actions {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

func (e *Engine) postActionOnlineEvent(action Action, result string) {
	// Sync events to mongo
}

func (e *Engine) wiseManDeathSpeech() string {
	if e.WiseMan == nil {
		return ""
	}
	s, _ := e.WiseMan.Farewell(e.State)
	return s
}
