package screens

import (
	"math"
	"strings"
)

type WaveModel struct {
	config WaveConfig
}

type WaveConfig struct {
	Phase     float64
	Frequency float64
	Amplitude float64
}

func NewWaveModel(phase, frequency, amplitude float64) WaveModel {
	return WaveModel{
		config: WaveConfig{
			Phase:     phase,
			Frequency: frequency,
			Amplitude: amplitude,
		},
	}
}

func (m *WaveModel) SetFrequency(freq float64) {
	m.config.Frequency = freq
}

func (m *WaveModel) AdvancePhase(delta float64) {
	m.config.Phase += delta
}

func (m *WaveModel) Render(width, height int) []string {
	if width <= 0 {
		return []string{}
	}

	var lines []string
	for y := 0; y < height; y++ {
		var line strings.Builder
		for x := 0; x < width; x++ {
			wave := math.Sin((float64(x)*m.config.Frequency + m.config.Phase + float64(y)*0.1))
			if math.Abs(wave) > 0.7 {
				line.WriteString("▀")
			} else if math.Abs(wave) > 0.3 {
				line.WriteString("▄")
			} else if math.Abs(wave) > 0.1 {
				line.WriteString("░")
			} else {
				line.WriteString(" ")
			}
		}
		lines = append(lines, line.String())
	}
	return lines
}
