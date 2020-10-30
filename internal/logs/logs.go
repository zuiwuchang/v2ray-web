package logs

import (
	"strings"
	"sync"
)

// MaxLogsCount .
const MaxLogsCount = 50

var single = _Logs{
	listeners: make(map[int64]logsListenerFunc),
	items:     make([]string, MaxLogsCount),
}

type logsListenerFunc func(string)
type _Logs struct {
	listeners map[int64]logsListenerFunc
	id        int64

	text  string
	items []string
	dirty bool
	index int
	count int
	sync.Mutex
}

func (l *_Logs) AddListener(listener logsListenerFunc) (id int64) {
	l.Lock()
	id = l.id
	l.id++
	l.listeners[id] = listener
	text := l.unsafeGetText()
	if text != "" {
		listener(text)
	}
	l.Unlock()
	return
}
func (l *_Logs) RemoveListener(id int64) {
	l.Lock()
	delete(l.listeners, id)
	l.Unlock()
	return
}
func (l *_Logs) Push(v string) {
	v = strings.TrimSpace(v)
	if v == "" {
		return
	}
	l.Lock()
	l.items[l.index] = v
	l.index++

	count := len(l.items)
	if l.index == count {
		l.index = 0
	}
	if l.count < count {
		l.count++
	}

	l.dirty = true
	l.notify(v)
	l.Unlock()
}
func (l *_Logs) unsafeGetText() string {
	if !l.dirty {
		return l.text
	}
	if l.count == 0 {
		l.text = ""
	} else if l.count == 1 {
		index := l.index - 1
		if index < 0 {
			index += len(l.items)
		}
		l.text = l.items[index]
	} else {
		arrs := make([]string, l.count)
		var index int
		count := len(l.items)
		for i := 0; i < l.count; i++ {
			index = l.index - 1 - i
			if index < 0 {
				index += count
			}
			arrs[l.count-1-i] = l.items[index]
		}
		l.text = strings.Join(arrs, "\n")
	}
	l.dirty = false
	return l.text
}
func (l *_Logs) notify(str string) {
	for _, f := range l.listeners {
		f(str)
	}
}

// AddListener 添加監聽器
func AddListener(listener logsListenerFunc) (id int64) {
	return single.AddListener(listener)
}

// RemoveListener 移除監聽器
func RemoveListener(id int64) {
	single.RemoveListener(id)
}

// Push 壓入日誌
func Push(v string) {
	single.Push(v)
}
