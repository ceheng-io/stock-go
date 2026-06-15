package stock

// WithSinaETFOptionListURL sets the Sina ETF option month list URL.
func WithSinaETFOptionListURL(listURL string) Option {
	return func(config *Config) {
		if listURL != "" {
			config.SinaETFOptionListURL = listURL
		}
	}
}

// WithSinaETFOptionExpireURL sets the Sina ETF option expire-day URL.
func WithSinaETFOptionExpireURL(expireURL string) Option {
	return func(config *Config) {
		if expireURL != "" {
			config.SinaETFOptionExpireURL = expireURL
		}
	}
}

// WithSinaETFOptionMinuteURL sets the Sina ETF option minute URL.
func WithSinaETFOptionMinuteURL(minuteURL string) Option {
	return func(config *Config) {
		if minuteURL != "" {
			config.SinaETFOptionMinuteURL = minuteURL
		}
	}
}

// WithSinaETFOptionDailyURL sets the Sina ETF option daily K-line URL.
func WithSinaETFOptionDailyURL(dailyURL string) Option {
	return func(config *Config) {
		if dailyURL != "" {
			config.SinaETFOptionDailyURL = dailyURL
		}
	}
}

// WithSinaETFOption5DayURL sets the Sina ETF option 5-day minute URL.
func WithSinaETFOption5DayURL(fiveDayURL string) Option {
	return func(config *Config) {
		if fiveDayURL != "" {
			config.SinaETFOption5DayURL = fiveDayURL
		}
	}
}

// WithSinaIndexOptionSpotURL sets the Sina index option T-quote URL.
func WithSinaIndexOptionSpotURL(spotURL string) Option {
	return func(config *Config) {
		if spotURL != "" {
			config.SinaIndexOptionSpotURL = spotURL
		}
	}
}

// WithSinaIndexOptionKlineURL sets the Sina index option daily K-line URL.
func WithSinaIndexOptionKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.SinaIndexOptionKlineURL = klineURL
		}
	}
}

// WithSinaCommodityOptionSpotURL sets the Sina commodity option T-quote URL.
func WithSinaCommodityOptionSpotURL(spotURL string) Option {
	return func(config *Config) {
		if spotURL != "" {
			config.SinaCommodityOptionSpotURL = spotURL
		}
	}
}

// WithSinaCommodityOptionKlineURL sets the Sina commodity option daily K-line URL.
func WithSinaCommodityOptionKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.SinaCommodityOptionKlineURL = klineURL
		}
	}
}
