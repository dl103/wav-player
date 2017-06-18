package wavplayer

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/dl103/wav"
	"github.com/gordonklaus/portaudio"
)

type Player struct {
	wavPath   string
	wavReader *wav.Reader
}

func NewPlayer(wavPath string) (p Player) {
	p = Player{wavPath: wavPath}
	p.loadReader()
	return
}

// Currently only plays from the default output device. Can change this later
func (p *Player) Play() {
	out := make([]int16, 1)
	// TODO: get the sample rate from the file
	stream, err := portaudio.OpenDefaultStream(0, 1, 48000, len(out), &out)
	chk(err)
	defer stream.Close()
	chk(stream.Start())

	// Reset before playing so we can re-use the player
	p.wavReader.Reset()

	buf := new(bytes.Buffer)
readLoop:
	for {
		s, err := p.wavReader.ReadRawSample()
		if err == io.EOF {
			break readLoop
		}
		chk(binary.Write(buf, binary.LittleEndian, s))
		chk(binary.Read(buf, binary.LittleEndian, out))
		chk(stream.Write())
	}
	return
}

func (p *Player) loadReader() {
	wavInfo, err := os.Stat(p.wavPath)
	chk(err)
	wavFile, err := os.Open(p.wavPath)
	chk(err)
	p.wavReader, err = wav.NewReader(wavFile, wavInfo.Size())
	chk(err)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
