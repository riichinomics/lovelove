package engine

import (
	"context"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	lovelove "hanafuda.moe/lovelove/proto"
)

type playerMeta struct {
	connections      []chan proto.Message
	id               string
	position         lovelove.PlayerPosition
	cancelDisconnect func()
}

type gameContextActivityState struct {
	mutex              sync.Mutex
	connectedPlayers   int
	activeRequests     int
	cleanupCancelation func()
}

type gameContext struct {
	id               string
	GameState        *gameState
	players          map[string]*playerMeta
	requestQueue     chan func()
	cleanupRequested chan context.Context

	activityState gameContextActivityState
}

func (context *gameContext) BroadcastUpdates(gameUpdates map[string][]*lovelove.GameStateUpdatePart) {
	for playerId, updates := range gameUpdates {
		player, ok := context.players[playerId]
		if !ok {
			continue
		}

		for _, listener := range player.connections {
			listener <- &lovelove.GameStateUpdate{
				Updates: updates,
			}
		}
	}
}

func (gameContext *gameContext) ChangeConnectionStatus(userId string, connected bool) {
	player, playerExists := gameContext.players[userId]
	if !playerExists {
		return
	}

	connectionStatusUpdates := make(map[string][]*lovelove.GameStateUpdatePart)
	for id, _ := range gameContext.players {
		if id == userId {
			continue
		}

		connectionStatusUpdates[id] = []*lovelove.GameStateUpdatePart{
			{
				ConnectionStatusUpdate: &lovelove.ConnectionStatusUpdate{
					Player:    player.position,
					Connected: connected,
				},
			},
		}
	}

	gameContext.BroadcastUpdates(connectionStatusUpdates)
}

func (activityState *gameContextActivityState) cancelCleanup() {
	if activityState.cleanupCancelation == nil {
		return
	}
	activityState.cleanupCancelation()
	activityState.cleanupCancelation = nil
}

func (gameContext *gameContext) StartRequest() {
	activityState := &gameContext.activityState
	activityState.mutex.Lock()
	activityState.activeRequests++
	activityState.cancelCleanup()
	defer activityState.mutex.Unlock()
}

func (gameContext *gameContext) PlayerConnected() {
	activityState := &gameContext.activityState
	activityState.mutex.Lock()
	activityState.connectedPlayers++
	activityState.cancelCleanup()
	defer activityState.mutex.Unlock()
}

func (gameContext *gameContext) EndRequest() {
	activityState := &gameContext.activityState
	activityState.mutex.Lock()
	activityState.activeRequests--
	if activityState.cleanupCancelation == nil && activityState.activeRequests == 0 && activityState.connectedPlayers == 0 {
		activityState.cleanupCancelation = gameContext.ScheduleGameCleanup(5)
	}
	defer activityState.mutex.Unlock()
}

func (gameContext *gameContext) PlayerDisconnected() {
	activityState := &gameContext.activityState
	activityState.mutex.Lock()
	activityState.connectedPlayers--
	if activityState.cleanupCancelation == nil && activityState.activeRequests == 0 && activityState.connectedPlayers == 0 {
		activityState.cleanupCancelation = gameContext.ScheduleGameCleanup(30)
	}
	defer activityState.mutex.Unlock()
}

func (gameContext *gameContext) PlayerLeftRoom(player *playerMeta, connection chan proto.Message) {
	for i, listener := range player.connections {
		if listener == connection {
			player.connections = append(player.connections[:i], player.connections[i+1:]...)
			break
		}
	}

	gameContext.PlayerDisconnected()

	if len(player.connections) != 0 {
		return
	}

	disconnectedContext, cancelDisconnect := context.WithCancel(context.Background())
	player.cancelDisconnect = cancelDisconnect
	go func() {
		defer func() {
			player.cancelDisconnect = nil
		}()

		select {
		case <-disconnectedContext.Done():
			return
		case <-time.After(5 * time.Second):
			cancelDisconnect()
			gameContext.requestQueue <- func() {
				if len(player.connections) != 0 {
					return
				}

				for _, gamePlayer := range gameContext.players {
					if len(gamePlayer.connections) > 0 {
						gameContext.ChangeConnectionStatus(player.id, false)
						return
					}
				}
			}
		}
	}()
}

func (gameContext *gameContext) ScheduleGameCleanup(delay int) func() {
	cleanupContext, cancelCleanup := context.WithCancel(context.Background())
	go func() {
		select {
		case <-cleanupContext.Done():
			return
		case <-time.After(time.Duration(delay) * time.Second):
			gameContext.cleanupRequested <- cleanupContext
		}
	}()
	return cancelCleanup
}
