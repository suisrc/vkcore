package voska_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/suisrc/vkcore/voska"

	"github.com/stretchr/testify/assert"
)

// go test ./voska -v -run TestSolveMp3
// go test ./voska -v -run TestSolveWav

func TestSolveMp3(t *testing.T) {
	voska.VOSK_MODEL = "../data/vosk"
	bts, _ := os.ReadFile("../data/dist/audio.mp3")
	txt, err := voska.GetTextByMp3(bts, true)
	assert.Nil(t, err)
	fmt.Println("content: ", txt)

	// assert.NotNil(t, nil)
}

func TestSolveWav(t *testing.T) {
	voska.VOSK_MODEL = "../data/vosk"
	bts, _ := os.ReadFile("../data/dist/audio.wav")
	txt, err := voska.GetTextByWav(bts, true)
	assert.Nil(t, err)
	fmt.Println("content: ", txt)

	// assert.NotNil(t, nil)
}
