package event_test
import (
	"log"
	"fmt"
	"gopkg.in/orivil/event.v0"
)

// test event name
const (
	GetUp = "GetUp"
	GoToWork = "GoToWork"
	BackHome = "BackHome"
)

// test listener
type People struct {
	name string
	priority int
}

func (p *People) GetUp(saySomething string) {
	fmt.Printf("%s get up and say: %s\n", p.name, saySomething)
}

func (p *People) GoToWork(saySomething string) {
	fmt.Printf("%s go to work and say: %s\n", p.name, saySomething)
}

func (p *People) BackHome(saySomething string) {
	fmt.Printf("%s back home and say: %s\n", p.name, saySomething)
}

// implement event.Listener interface
func (p *People) GetSubscribe() (name string, subscribes []event.Subscribe) {
	// listener name 要有唯一性
	name = p.name

	// 订阅事件
	subscribes = []event.Subscribe{
		// priority 可以不写， 即为默认值 0
		{Name: GoToWork, Priority: p.priority},
		{Name: GetUp, Priority: p.priority},
		{Name: BackHome, Priority: p.priority},
	}
	return
}

// test events
var testEvents = []*event.Event{
	{
		Name: GetUp,
		Call: func(listener interface{}, param ...interface{}) {
			listener.(*People).GetUp(param[0].(string))
		},
	},

	{
		Name: GoToWork,
		Call: func(listener interface{}, param ...interface{}) {
			listener.(*People).GoToWork(param[0].(string))
		},
	},

	{
		Name: BackHome,
		Call: func(listener interface{}, param ...interface{}) {
			listener.(*People).BackHome(param[0].(string))
		},
	},
}

func ExampleDispatcher()  {
	// 1. new dispatcher
	dispatcher := event.NewDispatcher()

	// 2. add events
	err := dispatcher.AddEvents(testEvents)
	if err != nil {
		log.Fatalln(err)
	}

	// new listener
	foo := &People{name: "foo", priority: 2}
	bar := &People{name: "bar"} // use default priority 0

	// 3. add listener
	err = dispatcher.AddListener(foo)
	if err == event.ErrListenerExist {
		log.Fatalf("people %s exist", foo.name)
	}
	err = dispatcher.AddListener(bar)
	if err == event.ErrListenerExist {
		log.Fatalf("people %s exist", bar.name)
	}

	// 4. trigger event
	say := "what a beautiful day"
	dispatcher.Trigger(GetUp, say)

	say = "I like my work"
	dispatcher.Trigger(GoToWork, say)

	say = "kiss my little baby girl"
	dispatcher.Trigger(BackHome, say)
	// Output:
	// foo get up and say: what a beautiful day
	// bar get up and say: what a beautiful day
	// foo go to work and say: I like my work
	// bar go to work and say: I like my work
	// foo back home and say: kiss my little baby girl
	// bar back home and say: kiss my little baby girl
}
