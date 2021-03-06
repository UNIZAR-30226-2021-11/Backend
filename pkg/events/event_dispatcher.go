// THIS IS AN AUTOMATICALLY GENERATED CODE! DO NOT EDIT THIS FILE!
// ADD YOUR EVENT TO 'generate_event_dispatcher.go' AND RUN 'go generate'

package events

import (
	"time"
)

// #######################
// INTERFACE DOCUMENTATION
// #######################

// 1. Create EventDispatcher using NewEventDispatcher() function.
// 2. Register your listeners using EventDispatcher.Register<event type name>Listener methods.
// 3. Run event loop by calling EventDispatcher.RunEventLoop() method.
// 4. Trigger events using EventDispatcher.Fire<event type name> methods.

// LISTENER INTERFACES

type StateChangedListener interface {
	HandleStateChanged(*StateChanged)
}

type GameCreateListener interface {
	HandleGameCreate(*GameCreate)
}

type SingleGameCreateListener interface {
	HandleSingleGameCreate(*SingleGameCreate)
}

type GamePauseListener interface {
	HandleGamePause(*GamePause)
}

type VotePauseListener interface {
	HandleVotePause(*VotePause)
}

type UserJoinedListener interface {
	HandleUserJoined(*UserJoined)
}

type UserLeftListener interface {
	HandleUserLeft(*UserLeft)
}

type CardPlayedListener interface {
	HandleCardPlayed(*CardPlayed)
}

type CardChangedListener interface {
	HandleCardChanged(*CardChanged)
}

type SingListener interface {
	HandleSing(*Sing)
}

// ##############################
// END OF INTERFACE DOCUMENTATION
// ##############################

const (
	eventQueuesCapacity                                       = 100000
	idleDispatcherSleepTime                     time.Duration = 5 * time.Millisecond
	registeringListenerWhileRunningErrorMessage               = "Tried to register listener while running event loop. Registering listeners is not thread safe therefore prohibited after starting event loop."
)

// PRIVATE EVENT HANDLERS

type eventHandler interface {
	handle()
}

type stateChangedHandler struct {
	event          *StateChanged
	eventListeners []StateChangedListener
}

func (handler *stateChangedHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleStateChanged(handler.event)
	}
}

type gameCreateHandler struct {
	event          *GameCreate
	eventListeners []GameCreateListener
}

func (handler *gameCreateHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleGameCreate(handler.event)
	}
}

type singleGameCreateHandler struct {
	event          *SingleGameCreate
	eventListeners []SingleGameCreateListener
}

func (handler *singleGameCreateHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleSingleGameCreate(handler.event)
	}
}

type gamePauseHandler struct {
	event          *GamePause
	eventListeners []GamePauseListener
}

func (handler *gamePauseHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleGamePause(handler.event)
	}
}

type votePauseHandler struct {
	event          *VotePause
	eventListeners []VotePauseListener
}

func (handler *votePauseHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleVotePause(handler.event)
	}
}

type userJoinedHandler struct {
	event          *UserJoined
	eventListeners []UserJoinedListener
}

func (handler *userJoinedHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleUserJoined(handler.event)
	}
}

type userLeftHandler struct {
	event          *UserLeft
	eventListeners []UserLeftListener
}

func (handler *userLeftHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleUserLeft(handler.event)
	}
}

type cardPlayedHandler struct {
	event          *CardPlayed
	eventListeners []CardPlayedListener
}

func (handler *cardPlayedHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleCardPlayed(handler.event)
	}
}

type cardChangedHandler struct {
	event          *CardChanged
	eventListeners []CardChangedListener
}

func (handler *cardChangedHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleCardChanged(handler.event)
	}
}

type singHandler struct {
	event          *Sing
	eventListeners []SingListener
}

func (handler *singHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleSing(handler.event)
	}
}

// EVENT DISPATCHER

type EventDispatcher struct {
	running bool

	// EVENT QUEUES

	priority2EventsQueue chan eventHandler

	// LISTENER LISTS

	stateChangedListeners []StateChangedListener

	gameCreateListeners []GameCreateListener

	singleGameCreateListeners []SingleGameCreateListener

	gamePauseListeners []GamePauseListener

	votePauseListeners []VotePauseListener

	userJoinedListeners []UserJoinedListener

	userLeftListeners []UserLeftListener

	cardPlayedListeners []CardPlayedListener

	cardChangedListeners []CardChangedListener

	singListeners []SingListener
}

// EVENT DISPATCHER CONSTRUCTOR

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		running: false,

		// EVENT QUEUES

		priority2EventsQueue: make(chan eventHandler, eventQueuesCapacity),

		// LISTENER LISTS

		stateChangedListeners: []StateChangedListener{},

		gameCreateListeners: []GameCreateListener{},

		singleGameCreateListeners: []SingleGameCreateListener{},

		gamePauseListeners: []GamePauseListener{},

		votePauseListeners: []VotePauseListener{},

		userJoinedListeners: []UserJoinedListener{},

		userLeftListeners: []UserLeftListener{},

		cardPlayedListeners: []CardPlayedListener{},

		cardChangedListeners: []CardChangedListener{},

		singListeners: []SingListener{},
	}
}

// MAIN EVENT LOOP

func (dispatcher *EventDispatcher) RunEventLoop() {
	dispatcher.running = true

	for {
		select {

		case handler := <-dispatcher.priority2EventsQueue:
			handler.handle()

		default:
			time.Sleep(idleDispatcherSleepTime)
		}
	}
}

func (dispatcher *EventDispatcher) panicWhenEventLoopRunning() {
	if dispatcher.running {
		panic(registeringListenerWhileRunningErrorMessage)
	}
}

// PUBLIC EVENT DISPATCHER METHODS

type QueueFilling struct {
	CurrentLength int
	Capacity      int
}

func (dispatcher *EventDispatcher) QueuesFilling() map[int]QueueFilling {
	filling := make(map[int]QueueFilling)

	filling[2] = QueueFilling{len(dispatcher.priority2EventsQueue), eventQueuesCapacity}

	return filling
}

// StateChanged

func (dispatcher *EventDispatcher) RegisterStateChangedListener(listener StateChangedListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.stateChangedListeners = append(dispatcher.stateChangedListeners, listener)
}

func (dispatcher *EventDispatcher) FireStateChanged(event *StateChanged) {
	handler := &stateChangedHandler{
		event:          event,
		eventListeners: dispatcher.stateChangedListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// GameCreate

func (dispatcher *EventDispatcher) RegisterGameCreateListener(listener GameCreateListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.gameCreateListeners = append(dispatcher.gameCreateListeners, listener)
}

func (dispatcher *EventDispatcher) FireGameCreate(event *GameCreate) {
	handler := &gameCreateHandler{
		event:          event,
		eventListeners: dispatcher.gameCreateListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// SingleGameCreate

func (dispatcher *EventDispatcher) RegisterSingleGameCreateListener(listener SingleGameCreateListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.singleGameCreateListeners = append(dispatcher.singleGameCreateListeners, listener)
}

func (dispatcher *EventDispatcher) FireSingleGameCreate(event *SingleGameCreate) {
	handler := &singleGameCreateHandler{
		event:          event,
		eventListeners: dispatcher.singleGameCreateListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// GamePause

func (dispatcher *EventDispatcher) RegisterGamePauseListener(listener GamePauseListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.gamePauseListeners = append(dispatcher.gamePauseListeners, listener)
}

func (dispatcher *EventDispatcher) FireGamePause(event *GamePause) {
	handler := &gamePauseHandler{
		event:          event,
		eventListeners: dispatcher.gamePauseListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// VotePause

func (dispatcher *EventDispatcher) RegisterVotePauseListener(listener VotePauseListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.votePauseListeners = append(dispatcher.votePauseListeners, listener)
}

func (dispatcher *EventDispatcher) FireVotePause(event *VotePause) {
	handler := &votePauseHandler{
		event:          event,
		eventListeners: dispatcher.votePauseListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// UserJoined

func (dispatcher *EventDispatcher) RegisterUserJoinedListener(listener UserJoinedListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.userJoinedListeners = append(dispatcher.userJoinedListeners, listener)
}

func (dispatcher *EventDispatcher) FireUserJoined(event *UserJoined) {
	handler := &userJoinedHandler{
		event:          event,
		eventListeners: dispatcher.userJoinedListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// UserLeft

func (dispatcher *EventDispatcher) RegisterUserLeftListener(listener UserLeftListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.userLeftListeners = append(dispatcher.userLeftListeners, listener)
}

func (dispatcher *EventDispatcher) FireUserLeft(event *UserLeft) {
	handler := &userLeftHandler{
		event:          event,
		eventListeners: dispatcher.userLeftListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// CardPlayed

func (dispatcher *EventDispatcher) RegisterCardPlayedListener(listener CardPlayedListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.cardPlayedListeners = append(dispatcher.cardPlayedListeners, listener)
}

func (dispatcher *EventDispatcher) FireCardPlayed(event *CardPlayed) {
	handler := &cardPlayedHandler{
		event:          event,
		eventListeners: dispatcher.cardPlayedListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// CardChanged

func (dispatcher *EventDispatcher) RegisterCardChangedListener(listener CardChangedListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.cardChangedListeners = append(dispatcher.cardChangedListeners, listener)
}

func (dispatcher *EventDispatcher) FireCardChanged(event *CardChanged) {
	handler := &cardChangedHandler{
		event:          event,
		eventListeners: dispatcher.cardChangedListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}

// Sing

func (dispatcher *EventDispatcher) RegisterSingListener(listener SingListener) {
	dispatcher.panicWhenEventLoopRunning()

	dispatcher.singListeners = append(dispatcher.singListeners, listener)
}

func (dispatcher *EventDispatcher) FireSing(event *Sing) {
	handler := &singHandler{
		event:          event,
		eventListeners: dispatcher.singListeners,
	}

	dispatcher.priority2EventsQueue <- handler
}
