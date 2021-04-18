package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"sort"
	"strconv"
	"text/template"
	"unicode"
)

// BEGIN CONFIGURATION

const (
	targetTypeName               = "EventDispatcher"
	targetFilePath               = "pkg/events/event_dispatcher.go"
	eventLoopMethodName          = "RunEventLoop"
	initialPriorityQueueCapacity = 100000
)

var supportedEvents = []EventType{
	//NewEventType("TimeTick", 1),
	NewEventType("StateChanged", 2),
	NewEventType("GameCreate", 2),
	NewEventType("UserJoined", 2),
	NewEventType("UserLeft", 2),
	NewEventType("CardPlayed", 2),
	NewEventType("CardChanged", 2),
	NewEventType("Sing", 2),
	//NewEventType("UserInput", 3),
	//NewEventType("ScoreSent", 3),
}

// END CONFIGURATION

func pascalCaseToCamel(str string) string {
	out := []rune(str)
	out[0] = unicode.ToLower(out[0])
	return string(out)
}

type Priority struct {
	Number uint
}

// sort.Interface implementation

type Priorites []Priority

func (slice Priorites) Len() int {
	return len(slice)
}

func (slice Priorites) Less(i, j int) bool {
	return slice[i].Number < slice[j].Number
}

func (slice Priorites) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// sort.Interface implementation end

func (priority *Priority) QueueName() string {
	return "priority" + strconv.Itoa(int(priority.Number)) + "EventsQueue"
}

type EventType struct {
	eventTypeName string
	priority      Priority
}

func NewEventType(typeName string, priority uint) EventType {
	return EventType{
		eventTypeName: typeName,
		priority:      Priority{priority},
	}
}

func (et *EventType) TypeName() string {
	return et.eventTypeName
}

func (et *EventType) ListenerTypeName() string {
	return et.eventTypeName + "Listener"
}

func (et *EventType) HandlerTypeName() string {
	return pascalCaseToCamel(et.eventTypeName + "Handler")
}

func (et *EventType) ListenerListName() string {
	return pascalCaseToCamel(et.eventTypeName + "Listeners")
}

func (et *EventType) ListenerHandleMethodName() string {
	return "Handle" + et.eventTypeName
}

func (et *EventType) RegisterMethodName() string {
	return "Register" + et.eventTypeName + "Listener"
}

func (et *EventType) FireMethodName() string {
	return "Fire" + et.eventTypeName
}

func (et *EventType) EventsQueueName() string {
	return et.priority.QueueName()
}

type Metadata struct {
	TypeName                     string
	EventLoopMethodName          string
	InitialPriorityQueueCapacity uint
	EventHandlerInterfaceName    string
	EventTypes                   []EventType
}

func (et Metadata) OrderedPriorities() []Priority {
	prioritiesSet := make(map[Priority]bool)

	for _, eventType := range et.EventTypes {
		prioritiesSet[eventType.priority] = true
	}

	uniqPriorities := make(Priorites, len(prioritiesSet))
	i := 0
	for priority := range prioritiesSet {
		uniqPriorities[i] = priority
		i++
	}
	sort.Sort(uniqPriorities)

	return uniqPriorities
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	codeTemplate := `
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
		// 2. Register your listeners using {{ .TypeName }}.Register<event type name>Listener methods.
		// 3. Run event loop by calling {{ .TypeName }}.{{ .EventLoopMethodName }}() method.
		// 4. Trigger events using {{ .TypeName }}.Fire<event type name> methods.

		// LISTENER INTERFACES

		{{ range .EventTypes }}
			type {{ .ListenerTypeName }} interface {
				{{ .ListenerHandleMethodName }}(*{{ .TypeName }})
			}
		{{ end }}

		// ##############################
		// END OF INTERFACE DOCUMENTATION
		// ##############################

		const (
			eventQueuesCapacity = {{ .InitialPriorityQueueCapacity }}
			idleDispatcherSleepTime time.Duration = 5 * time.Millisecond
			registeringListenerWhileRunningErrorMessage = "Tried to register listener while running event loop. Registering listeners is not thread safe therefore prohibited after starting event loop."
		)

		// PRIVATE EVENT HANDLERS

		type {{ .EventHandlerInterfaceName }} interface {
			handle()
		}

		{{ range .EventTypes }}
			type {{ .HandlerTypeName }} struct {
				event *{{ .TypeName }}
				eventListeners []{{ .ListenerTypeName }}
			}

			func (handler *{{ .HandlerTypeName }}) handle() {
				for _, listener := range handler.eventListeners {
					listener.{{ .ListenerHandleMethodName }}(handler.event)
				}
			}
		{{ end }}

		// EVENT DISPATCHER

		type {{ .TypeName }} struct {
			running bool

			// EVENT QUEUES

			{{ range .OrderedPriorities }}
				{{ .QueueName }} chan {{ $.EventHandlerInterfaceName }}
			{{ end }}

			// LISTENER LISTS

			{{ range .EventTypes }}
				{{ .ListenerListName }} []{{ .ListenerTypeName }}
			{{ end }}
		}

		// EVENT DISPATCHER CONSTRUCTOR

		func New{{ .TypeName }}() *{{ .TypeName }} {
			return &{{ .TypeName }}{
				running: false,

				// EVENT QUEUES

				{{ range .OrderedPriorities }}
					{{ .QueueName }}: make(chan {{ $.EventHandlerInterfaceName }}, eventQueuesCapacity),
				{{ end }}

				// LISTENER LISTS

				{{ range .EventTypes }}
					{{ .ListenerListName }}: []{{ .ListenerTypeName }}{},
				{{ end }}
			}
		}

		// MAIN EVENT LOOP

		func (dispatcher *{{ .TypeName }}) {{ .EventLoopMethodName }}() {
			dispatcher.running = true

			for {
				select {
					{{ range .OrderedPriorities }}
						case handler := <-dispatcher.{{ .QueueName }}:
							handler.handle()
					{{ end }}
				default:
					time.Sleep(idleDispatcherSleepTime)
				}
			}
		}

		func (dispatcher *{{ $.TypeName }}) panicWhenEventLoopRunning() {
			if(dispatcher.running) {
				panic(registeringListenerWhileRunningErrorMessage)
			}
		}

		// PUBLIC EVENT DISPATCHER METHODS

		type QueueFilling struct {
			CurrentLength int
			Capacity int
		}

		func (dispatcher *{{ $.TypeName }}) QueuesFilling() map[int]QueueFilling {
			filling := make(map[int]QueueFilling)

			{{ range .OrderedPriorities }}
				filling[{{ .Number }}] = QueueFilling{len(dispatcher.{{ .QueueName }}), eventQueuesCapacity}
			{{ end }}

			return filling
		}

		{{ range .EventTypes }}
			// {{ .TypeName }}

			func (dispatcher *{{ $.TypeName }}) {{ .RegisterMethodName }}(listener {{ .ListenerTypeName }}) {
				dispatcher.panicWhenEventLoopRunning()

				dispatcher.{{ .ListenerListName }} = append(dispatcher.{{ .ListenerListName }}, listener)
			}

			func (dispatcher *{{ $.TypeName }}) {{ .FireMethodName }}(event *{{ .TypeName }}) {
				handler := &{{ .HandlerTypeName }}{
					event: event,
					eventListeners: dispatcher.{{ .ListenerListName }},
				}

				dispatcher.{{ .EventsQueueName }} <- handler
			}
		{{ end }}
	`

	tpl, err := template.New("event_dispatcher").Parse(codeTemplate)
	checkError(err)

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, Metadata{
		TypeName:                     targetTypeName,
		EventLoopMethodName:          eventLoopMethodName,
		InitialPriorityQueueCapacity: initialPriorityQueueCapacity,
		EventHandlerInterfaceName:    "eventHandler",
		EventTypes:                   supportedEvents,
	})
	checkError(err)

	formatted_code, err := format.Source(buffer.Bytes())

	err = ioutil.WriteFile(targetFilePath, formatted_code, 0644)
	checkError(err)
}
