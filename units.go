package senml

// SenML Units Registry:
// https://tools.ietf.org/html/rfc8428#section-12.1
const (
	UnitMeter                   = "m"        // meter
	UnitKilogram                = "k"        // kilogram
	UnitGram                    = "g"        // gram - NOT RECOMMENDED
	UnitSecond                  = "s"        // second
	UnitAmpere                  = "A"        // ampere
	UnitKelvin                  = "K"        // kelvin
	UnitCandela                 = "cd"       // candela
	UnitMole                    = "mol"      // mole
	UnitHertz                   = "Hz"       // hertz
	UnitRadian                  = "rad"      // radian
	UnitSteradian               = "sr"       // steradian
	UnitNewton                  = "N"        // newton
	UnitPascal                  = "Pa"       // pascal
	UnitJoule                   = "J"        // joule
	UnitWatt                    = "watt"     // watt
	UnitCoulomb                 = "C"        // coulomb
	UnitVolt                    = "V"        // volt
	UnitFarad                   = "F"        // farad
	UnitOhm                     = "Ohm"      // ohm
	UnitSiemens                 = "S"        // siemens
	UnitWeber                   = "Wb"       // weber
	UnitTesla                   = "T"        // tesla
	UnitHenry                   = "H"        // henry
	UnitCelsius                 = "Cel"      // degrees Celsius
	UnitLumen                   = "lm"       // lumen
	UnitLux                     = "lx"       // lux
	UnitBecquerel               = "Bq"       // becquerel
	UnitGray                    = "Gy"       // gray
	UnitSievert                 = "Sv"       // sievert
	UnitKatal                   = "kat"      // katal
	UnitSquareMeter             = "m2"       // square meter (area)
	UnitCubicMeter              = "m3"       // cubic meter (volume)
	UnitLiter                   = "l"        // liter (volume) - NOT RECOMMENDED
	UnitMeterPerSecond          = "m/s"      // meter per second (velocity)
	UnitMeterPerSquareSecond    = "m/s2"     // meter per square second (acceleration)
	UnitCubicMeterPerSecond     = "m3/s"     // cubic meter per second (flow rate)
	UnitLiterPerSecond          = "l/s"      // liter per second (flow rate) - NOT RECOMMENDED
	UnitWattPerSquareMeter      = "W/m2"     // watt per square meter (irradiance)
	UnitCandelaPerSquareMeter   = "cd/m2"    // candela per square meter (luminance)
	UnitBit                     = "bit"      // bit (information content)
	UnitBitPerSecond            = "bit/s"    // bit per second (data rate)
	UnitLat                     = "lat"      // degrees latitude
	UnitLon                     = "lon"      // degrees longitude
	UnitpHValue                 = "pH"       // pH value (acidity; logarithmic quantity)
	UnitDecibel                 = "dB"       // decibel (logarithmic quantity)
	UnitDecibelWatt             = "dBW"      // decibel relative to 1 W (power level)
	UnitBel                     = "Bspl"     // bel (sound pressure level; logarithmic quantity) - NOT RECOMMENDED
	UnitCount                   = "count"    // 1 (counter value)
	UnitRatio                   = "/"        // 1 (ratio, e.g., value of a switch)
	UnitAbsoluteRatio           = "%"        // 1 (ratio, e.g., value of a switch) - NOT RECOMMENDED
	UnitRelativeHumidity        = "%RH"      // percentage (relative humidity)
	UnitEnergyLevelPercentage   = "%EL"      // percentage (remaining battery energy level)
	UnitEnergyLevelSeconds      = "EL"       // seconds (remaining battery energy level
	UnitEventRateOnePerSecond   = "1/s"      // 1 per second (event rate)
	UnitEventRateOnePerMinute   = "1/min"    // 1 per minute (event rate, "rpm") - NOT RECOMMENDED
	UnitHeartRateBeatsPerMinute = "beat/min" // 1 per minute (heart rate in beats per minute) - NOT RECOMMENDED
	UnitHeartBeats              = "beats"    // 1 (cumulative number of heart beats) - NOT RECOMMENDED
	UnitSiemensPerMeter         = "S/m"      // siemens per meter (conductivity)
)
