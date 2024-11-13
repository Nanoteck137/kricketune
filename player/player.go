package player

import (
	"fmt"
	"sync"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
	"github.com/kr/pretty"
	"github.com/nanoteck137/kricketune/core/log"
)

type Track struct {
	Name   string
	Artist string
	Uri    string
}

type Player struct {
	playbin *gst.Element

	queueMutex sync.Mutex
	index      int
	tracks     []Track
}

func New() (*Player, error) {
	playbin, err := createPlayer()
	if err != nil {
		return nil, err
	}

	return &Player{
		playbin: playbin,
	}, nil
}

func (p *Player) PrepareChange() {
	p.playbin.SetState(gst.StateReady)
}

func (p *Player) SetURI(uri string) {
	p.playbin.Set("uri", uri)
	p.playbin.SetState(gst.StatePlaying)
}

func (p *Player) PlayTrack(track Track) {
	p.SetURI(track.Uri)
}

// TODO(patrik): Rename?
func (p *Player) Reset() {
	p.PrepareChange()
	p.PlayTrack(p.CurrentTrack())
}

func (p *Player) Play() {
	p.playbin.SetState(gst.StatePlaying)
}

func (p *Player) Pause() {
	p.playbin.SetState(gst.StatePaused)
}

func (p *Player) NextTrack() {
	p.queueMutex.Lock()
	defer p.queueMutex.Unlock()

	p.index++
	if p.index >= len(p.tracks) {
		p.index = 0
	}

	p.PlayTrack(p.CurrentTrack())
}

func (p *Player) PrevTrack() {
	p.queueMutex.Lock()
	defer p.queueMutex.Unlock()

	p.index--
	if p.index < 0 {
		p.index = 0
	}

	p.PlayTrack(p.CurrentTrack())
}

func (p *Player) CurrentTrack() Track {
	return p.tracks[p.index]
}

func (p *Player) AddTrack(track Track) {
	p.tracks = append(p.tracks, track)
}

func Launch(player *Player) error {
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	playbin := player.playbin

	bus := playbin.GetBus()

	playbin.Connect("about-to-finish", func(element *gst.Element) {
		player.NextTrack()
	})

	bus.AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageError:
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			pretty.Println(err)
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG")
				fmt.Println(debug)
			}
			return false
		default:
			// fmt.Printf("msg.String(): %v\n", msg.String())
		}
		return true
	})

	go func() {
		err := mainLoop.RunError()
		if err != nil {
			log.Fatal("Failed to run loop", "err", err)
		}
	}()

	return nil
}

func createOutputs() (*gst.Bin, error) {
	outputs := gst.NewBin("outputs")

	tee, err := gst.NewElement("tee")
	if err != nil {
		return nil, err
	}

	outputs.Add(tee)

	teeSink := tee.GetStaticPad("sink")
	ghostPad := gst.NewGhostPad("sink", teeSink)
	outputs.AddPad(ghostPad.Pad)

	// out := "autoaudiosink"
	out := "audioresample ! audioconvert ! audio/x-raw,rate=48000,channels=2,format=S16LE ! autoaudiosink"
	output, err := gst.NewBinFromString(out, true)
	if err != nil {
		return nil, err
	}

	outputs.Add(output.Element)

	queue, err := gst.NewElement("queue")
	if err != nil {
		return nil, err
	}

	outputs.Add(queue)

	queue.Link(output.Element)
	tee.Link(queue)

	return outputs, nil
}

func createPlayer() (*gst.Element, error) {
	playbin, err := gst.NewElement("playbin")
	if err != nil {
		return nil, err
	}

	playbin.Set("flags", gst.StreamTypeAudio)
	playbin.Set("buffer-size", 5<<20)

	// playbin audio sink = bin -> queue -> volume -> output
	// output = tee -> queue -> parsed bin

	audioSink := gst.NewBin("audio-sink")

	queue, err := gst.NewElement("queue")
	if err != nil {
		return nil, err
	}

	volume, err := gst.NewElement("volume")
	if err != nil {
		return nil, err
	}

	volume.Set("volume", 0.01)

	outputs, err := createOutputs()
	if err != nil {
		return nil, err
	}

	audioSink.AddMany(queue, outputs.Element, volume)

	queue.Link(volume)
	volume.Link(outputs.Element)

	queueSink := queue.GetStaticPad("sink")
	ghostPad := gst.NewGhostPad("sink", queueSink)
	audioSink.AddPad(ghostPad.Pad)

	playbin.Set("audio-sink", audioSink)

	return playbin, nil
}
