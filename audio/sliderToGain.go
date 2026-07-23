package audio

import "math"

func SliderToGain(value float64) float32 {
	var db float64

	if value <= 0 {
		return 0
	}

	if value <= 0.5 {
		// 0.0 to 0.5 = -60dB to 0dB
		db = -60 + (value/0.5)*60
	} else {
		// 0.5 to 1.0 = 0dB to +6dB
		db = ((value - 0.5) / 0.5) * 6
	}

	return float32(math.Pow(10, db/20))
}

func GainToSlider(gain float32) float64 {
	if gain <= 0 {
		return 0
	}

	// Convert gain back to dB
	db := 20 * math.Log10(float64(gain))

	if db <= 0 {
		// -60dB to 0dB maps to 0.0 - 0.5
		if db < -60 {
			db = -60
		}

		return (db + 60) / 60 * 0.5
	}

	// 0dB to +6dB maps to 0.5 - 1.0
	if db > 6 {
		db = 6
	}

	return 0.5 + (db/6)*0.5
}
