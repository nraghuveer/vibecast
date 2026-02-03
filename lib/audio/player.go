package audio

import (
	"errors"
	"log"
	"os/exec"
	"sync"
)

const queueSize = 64

type Player struct {
	queue chan string
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
	select {
	case player.queue <- path:
		return nil
	default:
		return errors.New("audio queue full")
	}
}

func (p *Player) loop() {
	for path := range p.queue {
		if path == "" || afplayPath == "" {
			continue
		}
		cmd := exec.Command(afplayPath, path)
		if err := cmd.Run(); err != nil {
			log.Printf("audio playback failed: %v", err)
		}
	}
}
