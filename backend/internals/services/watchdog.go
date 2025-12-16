package services

import "github.com/gorilla/websocket"

//a routine tha frequently checks the health of various services and reports them if they are unresponsive.

type WatchdogService struct {
	Containers map[string]string // containers to watch
	Interval   int               // check interval in seconds
	Notify     func()
	conn       *websocket.Conn // alert configurations
}

func NewWatchdogService(interval int, notify func()) *WatchdogService {
	return &WatchdogService{
		Containers: make(map[string]string),
		Interval:   interval,
		Notify:     notify,
	}
}

func NewWatchdogServiceWithConnection(interval int, notify func(), conn *websocket.Conn) *WatchdogService {
	return &WatchdogService{
		Containers: make(map[string]string),
		Interval:   interval,
		Notify:     notify,
		conn:       conn,
	}
}

func (w *WatchdogService) MonitorServices() error {

	return nil
}

func (w *WatchdogService) ContainerWatch(containerID string) {
	// fetcth container runnnig
	// cross check with stored status
	// add new to map if not present notify
	// anything missess report notify
}

func (w *WatchdogService) ResourceWatch(notify func()) {
	// monitor resource usage
	// if threshold crossed notify
	// STREAM RESOURCE USAGE LOGS
}
