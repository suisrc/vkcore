package voska

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	vosk "github.com/alphacep/vosk-api/go"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
)

// 根据readme.md中的提示，安装vosk-model

func GetTextByMp3(mp3_data []byte, only_text bool) (string, error) {
	if len(mp3_data) == 0 {
		return "", fmt.Errorf("no mp3 data") // 没有音频数据
	}
	// mp3 -> wav

	// 解码 MP3 文件
	mp3_dec, err := mp3.NewDecoder(bytes.NewReader(mp3_data))
	if err != nil {
		return "", err
	}
	// defer mp3_dec.Close()
	// 获取MP3数据信息
	sample_rate := 44100 // 44100.0, CD音质
	bit_depth := 8
	num_chans := 1
	// 创建缓冲区
	wav_data := bytes.NewBuffer(nil)
	// 将缓冲区转换为 io.WriteSeeker 接口
	wseeker := &OffsetWriter{buffer: wav_data, position: 0}
	// 创建 WAV 编码器 // audioFormat: 1(PCM)
	wav_enc := wav.NewEncoder(wseeker, sample_rate, bit_depth, num_chans, 1)
	defer wav_enc.Close()
	// defer wav_enc.Close()
	// 将 MP3 数据解码PCM数据并写入 WAV 文件
	// 在 golang 种 byte 是 uint8 的别名
	buf := make([]uint8, 4096)
	for {
		num, err := mp3_dec.Read(buf)
		if num > 0 {
			if err := wav_enc.WriteFrame(buf[:num]); err != nil {
				return "", err // 写入失败
			}
		}
		if err == io.EOF {
			break // 解码完成, 结束遍历
		}
	}
	// 处理的音频数据， 无法预览， 但是 vosk 可以识别🤣
	// 对音频数据进行识别
	return GetTextByWav(wav_data.Bytes(), only_text)
}

var VOSK_MODEL = "data/vosk"

func GetTextByWav(wav_data []byte, only_text bool) (string, error) {
	if len(wav_data) == 0 {
		return "", fmt.Errorf("no wav data") // 没有音频数据
	}
	// fmt.Println("wav_data: ", len(wav_data))
	// os.WriteFile("../../dist/out.wav", wav_data, 0666)
	// 创建 Vosk 模型对象
	vosk.SetLogLevel(-1) // 关闭日志
	model, err := vosk.NewModel(VOSK_MODEL)
	if err != nil {
		return "", err
	}
	defer model.Free()

	// wav_dec := wav.NewDecoder(bytes.NewReader(wav_data))
	// wav_dec.ReadInfo() // float64(wav_dec.SampleRate)
	// 创建 Vosk 语音识别器对象, 这里就不看 wav_info, 直接使用CD音质的采样率
	rec, err := vosk.NewRecognizer(model, 44100.0)
	if err != nil {
		return "", err
	}
	defer rec.Free()

	// rec.SetWords(10)
	// rec.SetPartialWords(10)
	// 匹配可能性最高的10个结果
	// rec.SetMaxAlternatives(10)

	// 识别语音并输出识别结果
	if rec.AcceptWaveform(wav_data) != 0 {
		if only_text {
			// Unmarshal example for final result
			var jres map[string]interface{}
			json.Unmarshal([]byte(rec.Result()), &jres)
			return jres["text"].(string), nil
		} else {
			return rec.Result(), nil
		}
	}
	return "", nil // 没有音频数据
}

//=======================================================

// OffsetWriter is a io.WriteSeeker implementation that keeps track of the current offset.
type OffsetWriter struct {
	buffer   *bytes.Buffer
	position int64
}

func (w *OffsetWriter) Write(p []byte) (n int, err error) {
	n, err = w.buffer.Write(p)
	w.position += int64(n)
	return
}

func (w *OffsetWriter) Seek(offset int64, whence int) (int64, error) {
	var position int64

	switch whence {
	case io.SeekStart:
		position = offset
	case io.SeekCurrent:
		position = w.position + offset
	case io.SeekEnd:
		position = int64(w.buffer.Len()) + offset
	default:
		return w.position, os.ErrInvalid
	}

	if position < 0 {
		return w.position, os.ErrInvalid
	}

	w.position = position
	return w.position, nil
}
