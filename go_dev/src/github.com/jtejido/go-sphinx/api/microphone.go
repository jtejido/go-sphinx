package api

import (
	"fmt"

	"github.com/gen2brain/malgo"
)

// Microphone struct represents a simple microphone interface
type Microphone struct {
	ctx    malgo.Context
	device *malgo.Device
}

// NewMicrophone creates a new Microphone instance
func NewMicrophone(sampleRate int, channels int) (*Microphone, error) {
	mic := &Microphone{}

	// Configure the microphone context
	config := malgo.DefaultDeviceConfig(malgo.Capture)
	config.Capture.Format = malgo.FormatS16
	config.Capture.Channels = uint32(channels)
	config.SampleRate = uint32(sampleRate)

	// Initialize the context
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		return nil, err
	}
	mic.ctx = ctx.Context

	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)

	sizeInBytes := uint32(malgo.SampleSizeInBytes(config.Capture.Format))
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {

		sampleCount := framecount * config.Capture.Channels * sizeInBytes

		newCapturedSampleCount := capturedSampleCount + sampleCount

		pCapturedSamples = append(pCapturedSamples, pSample...)

		capturedSampleCount = newCapturedSampleCount

	}
	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	mic.device, err = malgo.InitDevice(mic.ctx, config, captureCallbacks)
	if err != nil {
		return nil, err
	}

	return mic, nil
}

// StartRecording starts recording from the microphone
func (mic *Microphone) StartRecording() error {
	return mic.device.Start()
}

// StopRecording stops recording from the microphone
func (mic *Microphone) StopRecording() error {
	return mic.device.Stop()
}

// CloseConnection closes the connection to the microphone
func (mic *Microphone) CloseConnection() {
	mic.device.Uninit()
}
