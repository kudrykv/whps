package signals

import (
	"errors"
	"github.com/kudrykv/whps/whps/types"
	"sync"
	"time"
)

var sm = sync.Map{}

func Create(id string) (chan *types.Req, error) {
	if _, ok := sm.Load(id); ok {
		return nil, errors.New("id already exists")
	}

	ch := make(chan *types.Req)
	sm.Store(id, ch)

	go func() {
		timer := time.NewTimer(2 * time.Second)
		<-timer.C

		ch <- nil
		sm.Delete(id)
	}()

	return ch, nil
}

func Get(id string) (chan *types.Req, bool) {
	if value, ok := sm.Load(id); ok {
		if ch, ok2 := value.(chan *types.Req); !ok2 {
			return nil, ok2
		} else {
			return ch, ok2
		}
	} else {
		return nil, ok
	}
}
