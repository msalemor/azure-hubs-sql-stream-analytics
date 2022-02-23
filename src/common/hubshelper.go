package common

import (
	"context"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

type EventHubsHelper struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Hub    *eventhub.Hub
}

func (*EventHubsHelper) NewEventHubsHelper() EventHubsHelper {
	return EventHubsHelper{}
}

func SendMessage(helper *EventHubsHelper, event string) {
	helper.Hub.Send(helper.Ctx, eventhub.NewEventFromString(event))
}

func EventHandler() {

}
