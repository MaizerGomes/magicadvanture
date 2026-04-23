# Implementation Plan: Who Wants to Be a Millionaire Mini-Game

Add a new interactive mini-game "Who Wants to Be a Millionaire" (Quem Quer Ser um Milionário) with AI-generated challenges, persistent leaderboard, timed streaks, and multiplayer lifelines.

## Objectives
- Create a new room: `millionaire_hall`.
- Implement a persistent leaderboard (top 20) with 2-day streaks.
- Add automatic reward distribution and leaderboard reset after 2 days.
- Implement two lifelines:
  - **Ask the Audience**: Broadcast question to all online players and collect responses.
  - **Ask the Wise Man**: Use AI with Google Search integration (via the `WiseManService`).
- Track lifeline usage per streak.

## Key Files & Context
- `game/models.go`: Add `MillionairePoints`, `LifelineAudienceUsed`, `LifelineWiseManUsed`, and `MillionaireStreakID` to `GameState`.
- `game/db.go`: Update SQLite schema to persist new `GameState` fields.
- `game/online.go`: 
  - Define `MillionaireLeaderboard` and `MillionaireAudienceResponse` MongoDB collections.
  - Implement `SyncMillionaireScore`, `GetMillionaireLeaderboard`, `CheckMillionaireReset`, `BroadcastMillionaireQuestion`, and `SubmitAudienceResponse`.
- `game/wiseman.go`: Add `GenerateMillionaireChallenge` and `AskWiseManWithSearch`.
- `game/engine.go`: Add `RunMillionaireGame` and logic to handle audience/AI lifelines.
- `game/world.go`: Add `millionaire_hall` to the world map.
- `game/localization.go`: Add translations for all new UI elements and messages.

## Implementation Steps

### Phase 1: Data Models & Persistence
1. **Update `GameState`**:
   - `MillionairePoints int`: Total points in the current streak.
   - `LifelineAudienceUsed bool`: Whether "Ask the Audience" was used this streak.
   - `LifelineWiseManUsed bool`: Whether "Ask the Wise Man" was used this streak.
   - `MillionaireStreakID string`: ID of the current streak to detect resets.
2. **Update SQLite**: Add columns to the `saves` table and update `SaveGame`/`LoadSave` logic.
3. **Update MongoDB Models**:
   - `MillionaireStreak`: `{ _id, StartTime, EndTime, IsProcessed }`
   - `MillionaireScore`: `{ OnlineID, PlayerName, Score, StreakID }`
   - `AudienceResponse`: `{ QuestionID, PlayerID, Choice }`

### Phase 2: Logic & Service Updates
1. **WiseManService**:
   - `GenerateMillionaireChallenge(gs *GameState)`: Request a harder question from AI.
   - `AskWiseManWithSearch(question string)`: Integration for the Google search lifeline (leveraging tools if available or simulating via AI if not).
2. **OnlineService**:
   - `CheckAndResetMillionaireStreak()`: Background or on-check logic to process rewards (100g, 50g, 25g) and reset scores.
   - `BroadcastAudienceQuestion(gs, challenge)`: Send notifications to all active players.
   - `CollectAudienceResponses(questionID)`: Count results from other players.

### Phase 3: UI & Gameplay
1. **World Update**: Connect `millionaire_hall` from `knowledge_hall`.
2. **Engine Update**:
   - Create `RunMillionaireGame(reader)`: The main loop for the mini-game.
   - Implement "Lifeline" menu options.
   - Implement the Scoreboard display (Top 20 + Top 3 highlight).
3. **Localization**: Add all Portuguese and English strings.

## Verification & Testing
- **Persistence**: Verify scores and lifeline usage are saved across sessions.
- **Timed Reset**: Manually set a streak's end time to the past and verify rewards are granted and the board is cleared.
- **Audience Lifeline**:
  - Start a challenge in one terminal.
  - Check inbox and reply in another terminal.
  - Verify the first terminal sees the aggregated results.
- **Wise Man Lifeline**: Verify the AI provides a "searched" answer.
- **Multiplayer**: Ensure the leaderboard updates correctly for different players.
