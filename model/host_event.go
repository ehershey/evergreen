package model

import (
	"10gen.com/mci"
	"github.com/10gen-labs/slogger/v1"
	"time"
)

const (
	// resource type
	ResourceTypeHost = "HOST"

	// event types
	EventHostCreated            = "HOST_CREATED"
	EventHostStatusChanged      = "HOST_STATUS_CHANGED"
	EventHostDNSNameSet         = "HOST_DNS_NAME_SET"
	EventHostProvisionFailed    = "HOST_PROVISION_FAILED"
	EventHostProvisioned        = "HOST_PROVISIONED"
	EventHostRunningTaskSet     = "HOST_RUNNING_TASK_SET"
	EventHostRunningTaskCleared = "HOST_RUNNING_TASK_CLEARED"
	EventHostTaskPidSet         = "HOST_TASK_PID_SET"
)

// implements EventData
type HostEventData struct {
	// necessary for IsValid
	ResourceType string `bson:"r_type" json:"resource_type"`

	OldStatus string `bson:"o_s,omitempty" json:"old_status,omitempty"`
	NewStatus string `bson:"n_s,omitempty" json:"new_status,omitempty"`
	SetupLog  string `bson:"log,omitempty" json:"setup_log,omitempty"`
	Hostname  string `bson:"hn,omitempty" json:"hostname,omitempty"`
	TaskId    string `bson:"t_id,omitempty" json:"task_id,omitempty"`
	TaskPid   string `bson:"t_pid,omitempty" json:"task_pid,omitempty"`
}

func (self HostEventData) IsValid() bool {
	return self.ResourceType == ResourceTypeHost
}

func NewHostEventFinder() *EventFinder {
	return &EventFinder{resourceType: ResourceTypeHost}
}

func FindMostRecentHostEvents(hostId string, n int) ([]Event, error) {
	return NewHostEventFinder().FindMostRecentEvents(hostId, n)
}

func FindAllHostEventsInOrder(hostId string) ([]Event, error) {
	return NewHostEventFinder().FindAllEventsInOrder(hostId)
}

func LogHostEvent(hostId string, eventType string, eventData HostEventData) {
	eventData.ResourceType = ResourceTypeHost
	event := Event{
		Timestamp:  time.Now(),
		ResourceId: hostId,
		EventType:  eventType,
		Data:       EventDataWrapper{eventData},
	}

	logger := NewDBEventLogger(EventLogCollection)
	if err := logger.LogEvent(event); err != nil {
		mci.Logger.Errorf(slogger.ERROR, "Error logging host event: %v", err)
	}
}

func LogHostCreatedEvent(hostId string) {
	LogHostEvent(hostId, EventHostCreated, HostEventData{})
}

func LogHostStatusChangedEvent(hostId string, oldStatus string,
	newStatus string) {
	LogHostEvent(hostId, EventHostStatusChanged,
		HostEventData{OldStatus: oldStatus, NewStatus: newStatus})
}

func LogHostDNSNameSetEvent(hostId string, dnsName string) {
	LogHostEvent(hostId, EventHostDNSNameSet,
		HostEventData{Hostname: dnsName})
}

func LogHostProvisionedEvent(hostId string) {
	LogHostEvent(hostId, EventHostProvisioned, HostEventData{})
}

func LogHostRunningTaskSetEvent(hostId string, taskId string) {
	LogHostEvent(hostId, EventHostRunningTaskSet,
		HostEventData{TaskId: taskId})
}

func LogHostRunningTaskClearedEvent(hostId string, taskId string) {
	LogHostEvent(hostId, EventHostRunningTaskCleared,
		HostEventData{TaskId: taskId})
}

func LogHostTaskPidSetEvent(hostId string, taskPid string) {
	LogHostEvent(hostId, EventHostTaskPidSet,
		HostEventData{TaskPid: taskPid})
}

func LogProvisionFailedEvent(hostId string, setupLog string) {
	LogHostEvent(hostId, EventHostProvisionFailed, HostEventData{SetupLog: setupLog})
}
