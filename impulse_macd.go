package imacd

// ImpulseMACD represents the Impulse MACD indicator
type ImpulseMACD struct {
	lengthMA     int
	lengthSignal int

	// Internal state for SMMA calculations
	smmaHigh *SMMA
	smmaLow  *SMMA

	// Internal state for ZLEMA calculation
	zlema *ZLEMA

	// Internal state for signal SMA
	signalSMA *SMA

	// Historical values for calculations
	values []ImpulseValue
}

// ImpulseValue represents a single calculation result
type ImpulseValue struct {
	MD    float64 // Main difference
	SB    float64 // Signal
	SH    float64 // Histogram (MD - SB)
	Color string  // Color indication
}

// SMMA (Smoothed Moving Average) helper
type SMMA struct {
	length int
	value  float64
	isInit bool
}

// ZLEMA (Zero Lag EMA) helper
type ZLEMA struct {
	length int
	ema1   *EMA
	ema2   *EMA
}

// EMA (Exponential Moving Average) helper
type EMA struct {
	length     int
	multiplier float64
	value      float64
	isInit     bool
}

// SMA (Simple Moving Average) helper
type SMA struct {
	length int
	values []float64
	sum    float64
}

// NewImpulseMACD creates a new Impulse MACD indicator
func NewImpulseMACD(lengthMA, lengthSignal int) *ImpulseMACD {
	return &ImpulseMACD{
		lengthMA:     lengthMA,
		lengthSignal: lengthSignal,
		smmaHigh:     NewSMMA(lengthMA),
		smmaLow:      NewSMMA(lengthMA),
		zlema:        NewZLEMA(lengthMA),
		signalSMA:    NewSMA(lengthSignal),
		values:       make([]ImpulseValue, 0),
	}
}

// Update processes new price data (high, low, close)
func (im *ImpulseMACD) Update(high, low, close float64) ImpulseValue {
	// Calculate HLC3 (typical price)
	hlc3 := (high + low + close) / 3.0

	// Update SMMA for high and low
	hi := im.smmaHigh.Update(high)
	lo := im.smmaLow.Update(low)

	// Update ZLEMA for HLC3
	mi := im.zlema.Update(hlc3)

	// Calculate main difference (md)
	var md float64
	if mi > hi {
		md = mi - hi
	} else if mi < lo {
		md = mi - lo
	} else {
		md = 0
	}

	// Calculate signal (sb)
	sb := im.signalSMA.Update(md)

	// Calculate histogram (sh)
	sh := md - sb

	// Determine color
	var color string
	if hlc3 > mi {
		if hlc3 > hi {
			color = "lime"
		} else {
			color = "green"
		}
	} else {
		if hlc3 < lo {
			color = "red"
		} else {
			color = "orange"
		}
	}

	value := ImpulseValue{
		MD:    md,
		SB:    sb,
		SH:    sh,
		Color: color,
	}

	im.values = append(im.values, value)
	return value
}

// GetValues returns all calculated values
func (im *ImpulseMACD) GetValues() []ImpulseValue {
	return im.values
}

// GetLatest returns the most recent calculation
func (im *ImpulseMACD) GetLatest() *ImpulseValue {
	if len(im.values) == 0 {
		return nil
	}
	return &im.values[len(im.values)-1]
}

// SMMA implementation
func NewSMMA(length int) *SMMA {
	return &SMMA{
		length: length,
		isInit: false,
	}
}

func (s *SMMA) Update(value float64) float64 {
	if !s.isInit {
		s.value = value // First value acts as SMA base
		s.isInit = true
	} else {
		s.value = (s.value*float64(s.length-1) + value) / float64(s.length)
	}
	return s.value
}

// ZLEMA implementation
func NewZLEMA(length int) *ZLEMA {
	return &ZLEMA{
		length: length,
		ema1:   NewEMA(length),
		ema2:   NewEMA(length),
	}
}

func (z *ZLEMA) Update(value float64) float64 {
	ema1 := z.ema1.Update(value)
	ema2 := z.ema2.Update(ema1)
	d := ema1 - ema2
	return ema1 + d
}

// EMA implementation
func NewEMA(length int) *EMA {
	multiplier := 2.0 / (float64(length) + 1.0)
	return &EMA{
		length:     length,
		multiplier: multiplier,
		isInit:     false,
	}
}

func (e *EMA) Update(value float64) float64 {
	if !e.isInit {
		e.value = value
		e.isInit = true
	} else {
		e.value = (value * e.multiplier) + (e.value * (1.0 - e.multiplier))
	}
	return e.value
}

// SMA implementation
func NewSMA(length int) *SMA {
	return &SMA{
		length: length,
		values: make([]float64, 0, length),
		sum:    0,
	}
}

func (s *SMA) Update(value float64) float64 {
	if len(s.values) < s.length {
		s.values = append(s.values, value)
		s.sum += value
	} else {
		s.sum -= s.values[0]
		copy(s.values, s.values[1:])
		s.values[s.length-1] = value
		s.sum += value
	}

	return s.sum / float64(len(s.values))
}

// Helper function to create default Impulse MACD (34, 9)
func NewDefaultImpulseMACD() *ImpulseMACD {
	return NewImpulseMACD(34, 9)
}

// BatchUpdate processes multiple price bars at once
func (im *ImpulseMACD) BatchUpdate(bars []PriceBar) []ImpulseValue {
	results := make([]ImpulseValue, len(bars))
	for i, bar := range bars {
		results[i] = im.Update(bar.High, bar.Low, bar.Close)
	}
	return results
}

// PriceBar re																																																																																																																										ents a price bar with OHLC data
type PriceBar struct {
	High  float64
	Low   float64
	Close float64
}

// Reset clears all internal state
func (im *ImpulseMACD) Reset() {
	im.smmaHigh = NewSMMA(im.lengthMA)
	im.smmaLow = NewSMMA(im.lengthMA)
	im.zlema = NewZLEMA(im.lengthMA)
	im.signalSMA = NewSMA(im.lengthSignal)
	im.values = make([]ImpulseValue, 0)
}