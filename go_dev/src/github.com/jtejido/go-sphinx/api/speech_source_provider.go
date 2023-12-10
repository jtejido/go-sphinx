package api

type SpeechSourceProvider struct{}

// use portaudio
func (ssp *SpeechSourceProvider) GetMicrophone() (*Microphone, error) {
	return NewMicrophone(16000, 16)
}
