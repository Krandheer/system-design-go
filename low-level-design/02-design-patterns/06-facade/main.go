package main

import "fmt"

// --- The Complex Subsystem Components ---
// These are the individual parts of our system. A client would normally
// have to interact with all of them directly to get something done.

type Amplifier struct{}

func (a *Amplifier) On()             { fmt.Println("Amplifier is on") }
func (a *Amplifier) SetVolume(v int) { fmt.Printf("Amplifier volume set to %d\n", v) }
func (a *Amplifier) Off()            { fmt.Println("Amplifier is off") }

type DvdPlayer struct{}

func (d *DvdPlayer) On()   { fmt.Println("DVD Player is on") }
func (d *DvdPlayer) Play() { fmt.Println("DVD Player is playing") }
func (d *DvdPlayer) Off()  { fmt.Println("DVD Player is off") }

type Projector struct{}

func (p *Projector) On()   { fmt.Println("Projector is on") }
func (p *Projector) Off()  { fmt.Println("Projector is off") }

type Screen struct{}

func (s *Screen) Down() { fmt.Println("Screen is down") }
func (s *Screen) Up()   { fmt.Println("Screen is up") }

// --- The Facade ---
// HomeTheaterFacade provides a simplified, unified interface to the subsystem.
type HomeTheaterFacade struct {
	amp       *Amplifier
	dvd       *DvdPlayer
	projector *Projector
	screen    *Screen
}

// NewHomeTheaterFacade is a constructor for our facade.
// It initializes all the subsystem components.
func NewHomeTheaterFacade() *HomeTheaterFacade {
	return &HomeTheaterFacade{
		amp:       &Amplifier{},
		dvd:       &DvdPlayer{},
		projector: &Projector{},
		screen:    &Screen{},
	}
}

// WatchMovie is the simplified method. It hides the complexity of the
// sequence of operations needed to watch a movie.
func (htf *HomeTheaterFacade) WatchMovie() {
	fmt.Println("Get ready to watch a movie...")
	htf.screen.Down()
	htf.projector.On()
	htf.amp.On()
	htf.amp.SetVolume(11)
	htf.dvd.On()
	htf.dvd.Play()
}

// EndMovie provides a similar simplified interface for shutting down.
func (htf *HomeTheaterFacade) EndMovie() {
	fmt.Println("\nShutting movie theater down...")
	htf.dvd.Off()
	htf.amp.Off()
	htf.projector.Off()
	htf.screen.Up()
}

func main() {
	// The client code interacts with the simple facade, not the complex subsystem.
	homeTheater := NewHomeTheaterFacade()
	
	// With one simple call, the client can start the movie.
	homeTheater.WatchMovie()
	
	// And with one simple call, the client can end it.
	homeTheater.EndMovie()
}
