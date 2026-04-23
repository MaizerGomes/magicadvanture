# Magic Adventure Tutorial

Created by Dante Gomes with assistance from Gemini and Codex.

This guide explains the core loop, progression, and the main systems that matter while playing.

## 0. Storage Model

- sqlite stores local save data, offline relic inventory, and other state the game needs to run without internet.
- Mongo stores the optional multiplayer layer: player presence, room chat, whispers, room events, party invitations, party follow/join actions, party heals, party rally buffs, party guards, party quests, party quest logs, party dungeon phases, quest cooldowns, and relic trading.
- Wise man AI settings are stored locally in sqlite after you configure them in-game. The setup flow can open the provider page in your browser, then wait for you to paste the key in the terminal.
- If Mongo is unreachable, the game automatically falls back to offline mode.

## 1. Start

- Launch the game.
- Pick one of the five save slots.
- If a slot already exists, continue it or overwrite it.
- On first run, the sqlite database is created automatically.

## 2. Core Loop

- Move between rooms using numbered actions.
- Explore biomes to find resources, bosses, and story triggers.
- Buy gear and healing items in the shop.
- Save progress happens automatically after actions and major transitions.
- If online mode is active, you can also see nearby players, chat in the room, whisper directly, broadcast room events, form a party, invite nearby players, share a party heal, rally the party for combat buffs, raise a party guard, start a party quest, open the party quest log, enter party dungeon phases, follow your party leader, trade relics, and open the notifications inbox for pending interactions.
- If the wise man is configured, you can ask him for clues in the Sage Hut. If he is not configured or the provider is unavailable, he falls back to a humorous offline reply.
- The game also remembers the last action result and shows it again on the next screen so you do not miss important text.

## 3. Combat

- Enemies have visible health bars.
- Your damage scales with level and gear.
- Steel Sword and Forsaken Blade dramatically increase attack output.
- Healing items restore HP up to your current maximum.

## 4. Progression

- XP fills your level bar.
- Leveling increases max health and base damage.
- Stronger zones are gated by items or prior victories.
- The Dragon at Mountain is the final test and expects you to be prepared.

## 5. Biomes

- Village: hub, shop, garden, and the path to the sage hut.
- Forest: major crossroads and the hidden trail to the old ruins.
- River: cave access, arctic bridge, and harbor route.
- Desert: Sun Amulet path.
- Arctic: Ice Crystal path.
- Harbor: gateway to the lighthouse, moonwell, and tide caves.
- Sage Hut: the wise man's home, the Moon Pearl offering point, and the Triptych Blessing turn-in.
- Settings: reconfigure the wise man AI provider or key at any time.
- Lighthouse: a clue-bearing lookout and a small reward stop.
- Tide Caves: the Moon Pearl side quest and the guardian fight.
- Old Ruins: hidden tablets that reveal one third of the Triptych quest.
- Moonwell: a glowing hidden spring that grants another Triptych sign.
- Observatory: a mountain lookout that completes the last Triptych sign.
- Void / Binary Sea: late-game glitch route.
- Party quest lines: Monster Hunt, Void Expedition, Frost Pact, and Sun Covenant each have multiple phases, distinct dungeon branches, shared rewards, quest-specific relic bonuses, and coordinated boss fights.
- Party quest boards enter a short cooldown after completion so the same reward chain cannot be farmed immediately.

## 6. Practical Advice

- Buy a sword early.
- Carry healing before entering dangerous zones.
- Use the garden when you want a safer way to build up.
- Do not rush the Dragon before level 4.
