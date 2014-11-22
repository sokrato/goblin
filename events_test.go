package goblin

import (
    "testing"
)

type MyListener struct {
    smiled bool
    cried bool
}

func (ml *MyListener)HandleEvent(evt string, args ...interface{}) {
    if evt == "smiled" {
        ml.smiled = true
    } else if evt == "cried" {
        ml.cried = true
    }
}

func TestEmitter(t *testing.T) {
    ml := &MyListener{}
    em := make(EventEmitter, 0)
    em.On("cried", ml)
    em.Once("smiled", ml)
    em.Emit("cried")
    em.Emit("smiled")
    if !(ml.cried && ml.smiled) {
        t.Fail()
    }
    ml.smiled = false
    em.Emit("smiled")
    if ml.smiled {
        t.Fail()
    }
}
