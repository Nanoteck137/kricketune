package player

import (
	"fmt"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
	"github.com/kr/pretty"
	"github.com/nanoteck137/kricketune/core/log"
)

type Queue interface {
	Next()
	Prev()
	CurrentTrack() (Track, bool)
}

type Track struct {
	Name   string
	Artist string
	Album  string
	Uri    string
}

type Player struct {
	playbin       *gst.Element
	volumeControl *gst.Element

	volume float32
	mute   bool

	// queueMutex sync.Mutex
	// index      int
	// tracks     []Track
	queue Queue
}

func New(audioOutput string) (*Player, error) {
	elements, err := createPlayer(audioOutput)
	if err != nil {
		return nil, err
	}

	return &Player{
		playbin:       elements.playbin,
		volumeControl: elements.volume,
	}, nil
}

func (p *Player) SetQueue(queue Queue) {
	p.queue = queue
}

func (p *Player) SetVolume(vol float32) {
	p.volumeControl.Set("volume", vol)
	p.volume = vol
}

func (p *Player) GetVolume() float32 {
	return p.volume
}

func (p *Player) SetMute(mute bool) {
	p.volumeControl.Set("mute", mute)
	p.mute = mute
}

func (p *Player) GetMute() bool {
	return p.mute
}

func (p *Player) IsPlaying() bool {
	return p.playbin.GetCurrentState() == gst.StatePlaying
}

func (p *Player) PrepareChange() {
	p.playbin.SetState(gst.StateReady)
}

func (p *Player) SetURI(uri string) {
	p.playbin.Set("uri", uri)
}

func (p *Player) PlayTrack(track Track) {
	log.Info("Now Playing", "name", track.Name, "artist", track.Artist, "album", track.Album)

	p.SetURI(track.Uri)
	p.Play()
}

func (p *Player) Start() {
	p.PrepareChange()

	track, ok := p.queue.CurrentTrack()
	if ok {
		p.PlayTrack(track)
	}
}

func (p *Player) Stop() {
	p.PrepareChange()
}

func (p *Player) Play() {
	p.playbin.SetState(gst.StatePlaying)
}

func (p *Player) Pause() {
	p.playbin.SetState(gst.StatePaused)
}

func (p *Player) NextTrack() {
	p.queue.Next()
	track, ok := p.queue.CurrentTrack()
	if ok {
		p.PlayTrack(track)
	}
}

func (p *Player) PrevTrack() {
	p.queue.Prev()
	track, ok := p.queue.CurrentTrack()
	if ok {
		p.PlayTrack(track)
	}
}

func (p *Player) RewindTrack() {
	p.playbin.SeekTime(0, gst.SeekFlagFlush)
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

	player.Start()

	go func() {
		err := mainLoop.RunError()
		if err != nil {
			log.Fatal("Failed to run loop", "err", err)
		}
	}()

	return nil
}

func createOutputs(audioOutput string) (*gst.Bin, error) {
	outputs := gst.NewBin("outputs")

	tee, err := gst.NewElement("tee")
	if err != nil {
		return nil, err
	}

	outputs.Add(tee)

	teeSink := tee.GetStaticPad("sink")
	ghostPad := gst.NewGhostPad("sink", teeSink)
	outputs.AddPad(ghostPad.Pad)

	output, err := gst.NewBinFromString(audioOutput, true)
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

type Elements struct {
	playbin *gst.Element
	volume  *gst.Element
}

func createPlayer(audioOutput string) (Elements, error) {
	playbin, err := gst.NewElement("playbin")
	if err != nil {
		return Elements{}, err
	}

	playbin.Set("flags", gst.StreamTypeAudio)
	playbin.Set("buffer-size", 5<<20)

	// playbin audio sink = bin -> queue -> volume -> output
	// output = tee -> queue -> parsed bin

	audioSink := gst.NewBin("audio-sink")

	queue, err := gst.NewElement("queue")
	if err != nil {
		return Elements{}, err
	}

	volume, err := gst.NewElement("volume")
	if err != nil {
		return Elements{}, err
	}

	outputs, err := createOutputs(audioOutput)
	if err != nil {
		return Elements{}, err
	}

	audioSink.AddMany(queue, outputs.Element, volume)

	queue.Link(volume)
	volume.Link(outputs.Element)

	queueSink := queue.GetStaticPad("sink")
	ghostPad := gst.NewGhostPad("sink", queueSink)
	audioSink.AddPad(ghostPad.Pad)

	playbin.Set("audio-sink", audioSink)

	return Elements{
		playbin: playbin,
		volume:  volume,
	}, nil
}
