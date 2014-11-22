package goblin


const defaultListenerSize = 4

type EventListener interface {
    HandleEvent(evt string, args ...interface{})
}

type oneTimeListener struct {
    emitter EventEmitter
    delegate EventListener
}

func (otl *oneTimeListener) HandleEvent(evt string, args ...interface{}) {
    otl.delegate.HandleEvent(evt, args ...)
    otl.emitter.RemoveListener(evt, otl)
}

type EventEmitter map[string][]EventListener

func (em EventEmitter) Emit(evt string, args ...interface{}) {
    listeners, ok := em[evt]
    if !ok || len(listeners) < 1 {
        return
    }
    for _, lsn := range listeners {
        lsn.HandleEvent(evt, args...)
    }
}

func (em EventEmitter) On(evt string, listener EventListener) {
    _, ok := em[evt]
    if !ok {
        em[evt] = make([]EventListener, 0, defaultListenerSize)
    }
    em[evt] = append(em[evt], listener)
}

func (em EventEmitter) Once(evt string, listener EventListener) {
    em.On(evt, &oneTimeListener{
        emitter: em,
        delegate: listener,
    })
}

func (em EventEmitter) RemoveListener(evt string, listener EventListener) {
    listeners, ok := em[evt]
    if !ok {
        return
    }
    lsnCopy := make([]EventListener, 0, defaultListenerSize)
    for _, lsn := range listeners {
        if lsn == listener  || lsn == nil{
            continue
        }
        lsnCopy = append(lsnCopy, lsn)
    }
    em[evt] = lsnCopy
}
