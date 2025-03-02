package main

const (
	SidebarLoad       = "sidebar-load"
	SidebarUpdateName = "sidebar-update-name"
	SidebarAdd        = "sidebar-add"
	SidebarDelete     = "sidebar-delete"

	HistoryChange           = "history-change"
	HistoryNetworkModify    = "history-network-modify"
	HistoryBusModify        = "history-bus-modify"
	HistoryNodeModify       = "history-node-modify"
	HistoryMessageModify    = "history-message-modify"
	HistorySignalModify     = "history-signal-modify"
	HistorySignalTypeModify = "history-signal-type-modify"
	HistorySignalUnitModify = "history-signal-unit-modify"
	HistorySignalEnumModify = "history-signal-enum-modify"

	BusAdded        = "bus-added"
	NodeAdded       = "node-added"
	MessageAdded    = "message-added"
	SignalAdded     = "signal-added"
	SignalTypeAdded = "signal-type-added"
	SignalUnitAdded = "signal-unit-added"
	SignalEnumAdded = "signal-enum-added"
)

type SidebarUpdateNameEvent struct {
	UpdatedID string `json:"updatedId"`
	Name      string `json:"name"`
}

type SidebarAddEvent struct {
	AddedItem SidebarItem `json:"addedItem"`
}

type SidebarDeleteEvent struct {
	DeletedID string `json:"deletedId"`
}
