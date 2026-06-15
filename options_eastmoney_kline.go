package stock

// WithEastmoneyKlineURL sets the Eastmoney K-line URL.
func WithEastmoneyKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.EastmoneyKlineURL = klineURL
		}
	}
}

// WithEastmoneyTrendsURL sets the Eastmoney CN trends URL.
func WithEastmoneyTrendsURL(trendsURL string) Option {
	return func(config *Config) {
		if trendsURL != "" {
			config.EastmoneyTrendsURL = trendsURL
		}
	}
}

// WithEastmoneyHKKlineURL sets the Eastmoney HK K-line URL.
func WithEastmoneyHKKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.EastmoneyHKKlineURL = klineURL
		}
	}
}

// WithEastmoneyHKTrendsURL sets the Eastmoney HK trends URL.
func WithEastmoneyHKTrendsURL(trendsURL string) Option {
	return func(config *Config) {
		if trendsURL != "" {
			config.EastmoneyHKTrendsURL = trendsURL
		}
	}
}

// WithEastmoneyUSKlineURL sets the Eastmoney US K-line URL.
func WithEastmoneyUSKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.EastmoneyUSKlineURL = klineURL
		}
	}
}

// WithEastmoneyUSTrendsURL sets the Eastmoney US trends URL.
func WithEastmoneyUSTrendsURL(trendsURL string) Option {
	return func(config *Config) {
		if trendsURL != "" {
			config.EastmoneyUSTrendsURL = trendsURL
		}
	}
}
