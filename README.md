# 🌌 Magic Adventure 6.0: The Multiplayer & Cloud Update

Welcome to the most connected version of Magic Adventure! Version 6.0 transitions the game to a global cloud-based experience, allowing adventurers to see and interact with each other in real-time.

---

## 🚀 What's New in 6.0

- **🌐 Online Multiplayer**: Experience a living world. See online players and their locations as you explore.
- **🤝 Player Interactions**: Wave, challenge, or share stories with other players in the same room.
- **☁️ Cloud Saves (MongoDB)**: Your progress is now stored in the cloud. Access your character from anywhere!
- **📖 Interactive Tutorial**: New players are guided through a "How to Play" section to master the game mechanics.
- **❄️ Arctic Expansion**: The Frozen Tundra has been expanded with new locations: Arctic Entrance, Frost Cliffs, and the Ice Bridge.
- **📊 Refined Combat & Balance**: Monsters now have visible health bars, and combat has been rebalanced for a more rewarding progression.

---

## ✨ Core Features

- **🎮 Persistent World**: Every action updates your global presence.
- **🛡️ 5 Character Slots**: Manage multiple legends stored securely in MongoDB.
- **🎨 Dynamic ASCII Art**: Visual feedback for every biome and legendary boss.
- **🧩 Open-World Progression**: Collect the Sun Amulet and Ice Crystal to reach the final encounter at Dragon's Peak.

---

## 🛠️ Installation & Setup

1. **Prerequisites**: Ensure you have [Go](https://golang.org/dl/) installed.
2. **Environment Variable**: For security, set your MongoDB connection string:
   ```bash
   export MONGO_URI="your_mongodb_connection_string"
   ```
   *(If not set, it will default to the public development database)*.
3. **Build**:
   ```bash
   go build -o adventure main.go
   ```
4. **Run**:
   ```bash
   ./adventure
   ```

---

## 📖 Documentation

- **[TUTORIAL.md](./TUTORIAL.md)**: Details on XP, leveling, and biome survival.
- **[README.md](./README.md)**: Project overview and setup.

---

*“In a world of flickering code and shifting sands, you are never truly alone.”*
