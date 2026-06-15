package stock

// WithTHSLimitUpPoolURL sets the Tonghuashun limit-up pool URL.
func WithTHSLimitUpPoolURL(limitUpPoolURL string) Option {
	return func(config *Config) {
		if limitUpPoolURL != "" {
			config.THSLimitUpPoolURL = limitUpPoolURL
		}
	}
}
