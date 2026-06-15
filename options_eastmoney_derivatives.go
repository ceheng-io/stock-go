package stock

// WithEastmoneyFuturesKlineURL sets the Eastmoney futures K-line URL.
func WithEastmoneyFuturesKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.EastmoneyFuturesKlineURL = klineURL
		}
	}
}

// WithEastmoneyFuturesGlobalSpotURL sets the Eastmoney global futures spot URL.
func WithEastmoneyFuturesGlobalSpotURL(spotURL string) Option {
	return func(config *Config) {
		if spotURL != "" {
			config.EastmoneyFuturesGlobalSpotURL = spotURL
		}
	}
}

// WithEastmoneyFuturesGlobalKlineURL sets the Eastmoney global futures K-line URL.
func WithEastmoneyFuturesGlobalKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.EastmoneyFuturesGlobalKlineURL = klineURL
		}
	}
}

// WithEastmoneyOptionCFFEXURL sets the Eastmoney CFFEX option quotes URL.
func WithEastmoneyOptionCFFEXURL(optionURL string) Option {
	return func(config *Config) {
		if optionURL != "" {
			config.EastmoneyOptionCFFEXURL = optionURL
		}
	}
}

// WithEastmoneyOptionLHBURL sets the Eastmoney option LHB URL.
func WithEastmoneyOptionLHBURL(lhbURL string) Option {
	return func(config *Config) {
		if lhbURL != "" {
			config.EastmoneyOptionLHBURL = lhbURL
		}
	}
}
