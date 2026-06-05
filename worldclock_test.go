package main

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestZoneData(t *testing.T) {
	names := zoneNames()
	if len(names) < 100 {
		t.Fatalf("expected hundreds of zones, got %d", len(names))
	}
	hel, err := loadZone("Europe/Helsinki")
	if err != nil {
		t.Fatalf("load Helsinki: %v", err)
	}
	la, err := loadZone("America/Los_Angeles")
	if err != nil {
		t.Fatalf("load Los_Angeles: %v", err)
	}
	now := time.Date(2026, 6, 5, 12, 0, 0, 0, time.UTC)
	_, ho := now.In(hel).Zone()
	_, lo := now.In(la).Zone()
	if ho == lo {
		t.Fatalf("expected different offsets, both %d", ho)
	}
	t.Logf("zones=%d helsinki=%+d la=%+d (hours)", len(names), ho/3600, lo/3600)
}

func TestRenderClock(t *testing.T) {
	out := renderClock(time.Date(2026, 6, 5, 10, 10, 30, 0, time.UTC), 40, 20)
	lines := strings.Split(out, "\n")
	if len(lines) < 5 {
		t.Fatalf("clock too small: %d lines", len(lines))
	}
	t.Logf("rendered %dx%d clock:\n%s", len(lines[0]), len(lines), out)
}

func TestModelView(t *testing.T) {
	m := initialModel()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = updated.(model)
	out := m.View()
	if !strings.Contains(out, "Europe/Helsinki") || !strings.Contains(out, "America/Los_Angeles") {
		t.Fatalf("view missing zone labels:\n%s", out)
	}
	t.Logf("full view:\n%s", out)
}
