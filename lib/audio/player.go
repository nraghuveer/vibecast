package audio

import (
	"errors"
	"log"
	"os/exec"
	"sync"
)

const queueSize = 64

type Player struct {
	queue   chan string
	mu      sync.Mutex
	cond    *sync.Cond
	pending int
}

var (
	playerOnce sync.Once
	player     *Player
	afplayPath string
	afplayErr  error
)

// Start initializes the background audio player.
// On macOS it uses afplay. If unavailable, audio is disabled.
func Start() *Player {
	playerOnce.Do(func() {
		path, err := exec.LookPath("afplay")
		if err != nil {
			afplayErr = err
			log.Printf("audio disabled: %v", err)
		} else {
			afplayPath = path
		}
		player = &Player{queue: make(chan string, queueSize)}
		player.cond = sync.NewCond(&player.mu)
		go player.loop()
	})
	return player
}

// Enqueue schedules an audio file for playback.
// It is non-blocking; when the queue is full it returns an error.
func Enqueue(path string) error {
	Start()
	if afplayErr != nil {
		return afplayErr
	}
	player.mu.Lock()
	player.pending++
	player.mu.Unlock()
	select {
	case player.queue <- path:
		return nil
	default:
		player.decrementPending()
		return errors.New("audio queue full")
	}
}

// Drain blocks until all queued audio has finished playing.
func Drain() {
	Start()
	if afplayErr != nil {
		return
	}
	player.mu.Lock()
	for player.pending > 0 {
		player.cond.Wait()
	}
	player.mu.Unlock()
}

func (p *Player) loop() {
	for path := range p.queue {
		if path == "" || afplayPath == "" {
			p.decrementPending()
			continue
		}
		cmd := exec.Command(afplayPath, path)
		if err := cmd.Run(); err != nil {
			log.Printf("audio playback failed: %v", err)
		}
		p.decrementPending()
	}
}

func (p *Player) decrementPending() {
	p.mu.Lock()
	if p.pending > 0 {
		p.pending--
	}
	if p.pending == 0 {
		p.cond.Broadcast()
	}
	p.mu.Unlock()
}
