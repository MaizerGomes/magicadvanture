# Magic Adventure

Magic Adventure is a terminal RPG built in Go. It combines slot-based saves, a branching world, turn-style combat, local sqlite persistence, and optional Mongo-backed multiplayer.
Created by Dante Gomes with assistance from Gemini and Codex.

## What It Does Well

- Five character slots
- Auto-created sqlite save database on first run for offline progress
- Optional Mongo-backed multiplayer when `MONGO_URI` is available
- Open-world progression across village, forest, river, harbor, lighthouse, tide caves, old ruins, moonwell, observatory, desert, arctic, sage hut, and glitch zones
- Boss fights with visible health bars
- Tutorial flow for new players

## Install

### Homebrew (macOS/Linux)

```bash
brew tap MaizerGomes/magicadvanture
brew install magicadventure
```

### Direct Download

See the [Releases](https://github.com/MaizerGomes/magicadvanture/releases) page for binaries.

## Run It

```bash
go run .
```

The first run creates `magicadventure.db` in the project directory unless `MAGIC_ADVENTURE_DB_PATH` is set.
If `MONGO_URI` is set and reachable, the game also enables online presence, room chat, whispers, room events, party invites, party following, party healing, party rally buffs, party guards, party quests, party quest logs, party dungeon phases, quest cooldowns, relic trading, and a persistent notifications inbox. The wise man is configured separately inside the game and can use Gemini or Cloudflare as an optional online assistant. The current party quest lines are Monster Hunt, Void Expedition, Frost Pact, and Sun Covenant. A separate side quest lets you collect the Ruins Token, Moon Charm, and Star Map from new locations and turn them in at the sage hut for the Triptych Blessing. Quest relics also provide thematic combat bonuses in their matching zones. If Mongo is unavailable, the game stays fully playable offline.

## Configuration

- `MAGIC_ADVENTURE_DB_PATH`: Optional path to the sqlite database file
- `MONGO_URI`: Optional MongoDB connection string for online multiplayer
- `MONGO_DB_NAME`: Optional Mongo database name, defaults to `magicadventure_online`
- `WISEMAN_AI_PROVIDER`: Optional bootstrap provider for the wise man: `gemini`, `cloudflare`, or `legacy`
- `WISEMAN_AI_KEY`: Optional bootstrap API key for the wise man
- `WISEMAN_AI_MODEL`: Optional bootstrap model name
- `WISEMAN_AI_ACCOUNT_ID`: Optional bootstrap Cloudflare account ID
- `WISEMAN_AI_URL`: Optional legacy chat-completions endpoint if you still want to use the old OpenAI-compatible flow
- `.env`: Loaded automatically if present

The recommended setup flow is in-game: at startup, the game can open the provider page in your browser, wait for you to paste the key, and save the wise man settings in sqlite. You can change the provider or key later from the Settings menu.

## Build

```bash
go build -o adventure .
```

## Notes

- Save data is always local by default.
- Online features are best effort and never block offline play.
- The sqlite database file and release binaries are ignored by git.
- See `TUTORIAL.md` for the game rules and progression notes.
