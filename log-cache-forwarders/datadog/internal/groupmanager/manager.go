package groupmanager

import (
	"context"
	"time"
)

// Manager syncs the source IDs from the GroupProvider with the GroupUpdater.
// It does so at the configured interval.
type Manager struct {
	groupName string
	ticker    <-chan time.Time
	gp        GroupProvider
	gu        GroupUpdater
}

// GroupProvider returns the desired SourceIDs.
type GroupProvider interface {
	SourceIDs() []string
}

// GroupUpdater is used to add (and keep alive) the source IDs for a group.
type GroupUpdater interface {
	// SetShardGroup adds source IDs to the LogCache sub-groups.
	SetShardGroup(ctx context.Context, name string, sourceIDs ...string) error
}

// Start creates and starts a Manager.
func Start(groupName string, ticker <-chan time.Time, gp GroupProvider, gu GroupUpdater) {
	m := &Manager{
		groupName: groupName,
		ticker:    ticker,
		gp:        gp,
		gu:        gu,
	}

	go m.run()
}

func (m *Manager) run() {
	for range m.ticker {
		sourceIDs := m.gp.SourceIDs()
		m.updateSourceIDs(sourceIDs)
	}
}

func (m *Manager) updateSourceIDs(sourceIDs []string) {
	for _, id := range sourceIDs {
		m.gu.SetShardGroup(context.Background(), m.groupName, id)
	}
}
