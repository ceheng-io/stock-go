package stock

// WithTencentMinuteURL sets the Tencent minute timeline URL.
func WithTencentMinuteURL(minuteURL string) Option {
	return func(config *Config) {
		if minuteURL != "" {
			config.TencentMinuteURL = minuteURL
		}
	}
}

// WithAShareListURL sets the A-share code list URL.
func WithAShareListURL(listURL string) Option {
	return func(config *Config) {
		if listURL != "" {
			config.AShareListURL = listURL
		}
	}
}

// WithUSListURL sets the US code list URL.
func WithUSListURL(listURL string) Option {
	return func(config *Config) {
		if listURL != "" {
			config.USListURL = listURL
		}
	}
}

// WithHKListURL sets the HK code list URL.
func WithHKListURL(listURL string) Option {
	return func(config *Config) {
		if listURL != "" {
			config.HKListURL = listURL
		}
	}
}

// WithFundListURL sets the fund code list URL.
func WithFundListURL(listURL string) Option {
	return func(config *Config) {
		if listURL != "" {
			config.FundListURL = listURL
		}
	}
}

// WithCalendarURL sets the A-share trading calendar URL.
func WithCalendarURL(calendarURL string) Option {
	return func(config *Config) {
		if calendarURL != "" {
			config.CalendarURL = calendarURL
		}
	}
}

// WithBaseURL sets the base URL used by the Tencent quote client.
func WithBaseURL(baseURL string) Option {
	return func(config *Config) {
		if baseURL != "" {
			config.BaseURL = baseURL
		}
	}
}

// WithSearchBaseURL sets the Tencent Smartbox search URL.
func WithSearchBaseURL(baseURL string) Option {
	return func(config *Config) {
		if baseURL != "" {
			config.SearchBaseURL = baseURL
		}
	}
}
