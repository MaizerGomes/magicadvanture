package game

import (
	"path/filepath"
	"testing"
)

func TestWiseManConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "saves.db")
	t.Setenv("MAGIC_ADVENTURE_DB_PATH", dbPath)

	store, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer store.Close()

	service := NewWiseManService(store)
	cfg := &WiseManConfig{
		Enabled:  true,
		Provider: WiseManProviderGemini,
		Model:    "gemini-2.5-flash",
		APIKey:   "test-key",
	}
	if err := service.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loaded, err := service.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected wise man config to load")
	}
	if loaded.Provider != WiseManProviderGemini || loaded.Model != "gemini-2.5-flash" || loaded.APIKey != "test-key" {
		t.Fatalf("unexpected config round-trip: %+v", loaded)
	}
}

func TestWiseManConfigLoadsFromEnvFallback(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "saves.db")
	t.Setenv("MAGIC_ADVENTURE_DB_PATH", dbPath)
	t.Setenv("WISEMAN_AI_PROVIDER", WiseManProviderCloudflare)
	t.Setenv("WISEMAN_AI_KEY", "cloudflare-token")
	t.Setenv("WISEMAN_AI_ACCOUNT_ID", "acct-123")
	t.Setenv("WISEMAN_AI_MODEL", DefaultCloudflareModel)

	store, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer store.Close()

	service := NewWiseManService(store)
	loaded, err := service.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected env wise man config to load")
	}
	if loaded.Provider != WiseManProviderCloudflare || loaded.AccountID != "acct-123" || loaded.APIKey != "cloudflare-token" {
		t.Fatalf("unexpected env fallback config: %+v", loaded)
	}
}
