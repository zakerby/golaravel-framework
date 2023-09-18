package event

//go:generate mockery --name=Instance
type Instance interface {
	// Register event listeners to the application.
	Register(map[Event][]Listener)
	// Job create a new event task.
	Job(event Event, args []Arg) Task
	// GetEvents gets all registered events.
	GetEvents() map[Event][]Listener
}

type Event interface {
	// Handle the event.
	Handle(args []Arg) ([]Arg, error)
}

type Listener interface {
	// Signature returns the unique identifier for the listener.
	Signature() string
	// Queue configure the event queue options.
	Queue(args ...any) Queue
	// Handle the event.
	Handle(args ...any) error
}

//go:generate mockery --name=Task
type Task interface {
	// Dispatch an event and call the listeners.
	Dispatch() error
}

type Arg struct {
	Type  string
	Value any
}

type Queue struct {
	Enable     bool
	Connection string
	Queue      string
}
