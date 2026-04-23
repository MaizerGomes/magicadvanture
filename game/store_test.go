package game

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSQLiteStoreCreatesDatabaseAndRoundTripsSave(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "saves.db")
	t.Setenv("MAGIC_ADVENTURE_DB_PATH", dbPath)

	store, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer store.Close()

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected sqlite file to exist: %v", err)
	}

	gs := NewGameState(4, "Astra")
	gs.World = CreateWorld()
	gs.HasSeenTutorial = true
	gs.Gold = 321
	gs.Experience = 77
	gs.OnlineRelics = []string{"Amber Compass", "Sunken Coin"}
	gs.LastNotificationSeen = 123456789

	if err := store.SaveGame(gs); err != nil {
		t.Fatalf("SaveGame failed: %v", err)
	}

	loaded, err := store.LoadSave(4)
	if err != nil {
		t.Fatalf("LoadSave failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected saved game to load")
	}

	if loaded.PlayerName != "Astra" || loaded.Gold != 321 || !loaded.HasSeenTutorial {
		t.Fatalf("loaded state mismatch: %+v", loaded)
	}
	if loaded.LastNotificationSeen != 123456789 {
		t.Fatalf("expected notification cursor to round-trip, got %d", loaded.LastNotificationSeen)
	}
	if loaded.OnlineID == "" {
		t.Fatal("expected online id to be persisted")
	}
	if len(loaded.OnlineRelics) != 2 || loaded.OnlineRelics[0] != "Amber Compass" {
		t.Fatalf("expected relics to round-trip, got %+v", loaded.OnlineRelics)
	}
}

func TestGetOnlinePlayersExcludesCurrentSlot(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "saves.db")
	t.Setenv("MAGIC_ADVENTURE_DB_PATH", dbPath)

	store, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer store.Close()

	current := NewGameState(1, "Current")
	current.World = CreateWorld()
	other := NewGameState(2, "Other")
	other.World = CreateWorld()

	if err := store.SaveGame(current); err != nil {
		t.Fatalf("SaveGame current failed: %v", err)
	}
	if err := store.SaveGame(other); err != nil {
		t.Fatalf("SaveGame other failed: %v", err)
	}

	players, err := store.GetOnlinePlayers(current)
	if err != nil {
		t.Fatalf("GetOnlinePlayers failed: %v", err)
	}
	if len(players) != 1 {
		t.Fatalf("expected one other player, got %d", len(players))
	}
	if players[0].ID != 2 {
		t.Fatalf("expected slot 2, got %d", players[0].ID)
	}
}
