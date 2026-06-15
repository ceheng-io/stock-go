package stock

// WithEastmoneyIndustryListURL sets the Eastmoney industry board list URL.
func WithEastmoneyIndustryListURL(listURL string) Option {
	return func(config *Config) {
		if listURL != "" {
			config.EastmoneyIndustryListURL = listURL
		}
	}
}

// WithEastmoneyIndustrySpotURL sets the Eastmoney industry board spot URL.
func WithEastmoneyIndustrySpotURL(spotURL string) Option {
	return func(config *Config) {
		if spotURL != "" {
			config.EastmoneyIndustrySpotURL = spotURL
		}
	}
}

// WithEastmoneyIndustryConstituentsURL sets the Eastmoney industry board constituents URL.
func WithEastmoneyIndustryConstituentsURL(consURL string) Option {
	return func(config *Config) {
		if consURL != "" {
			config.EastmoneyIndustryConstituentsURL = consURL
		}
	}
}

// WithEastmoneyIndustryKlineURL sets the Eastmoney industry board K-line URL.
func WithEastmoneyIndustryKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.EastmoneyIndustryKlineURL = klineURL
		}
	}
}

// WithEastmoneyIndustryTrendsURL sets the Eastmoney industry board trends URL.
func WithEastmoneyIndustryTrendsURL(trendsURL string) Option {
	return func(config *Config) {
		if trendsURL != "" {
			config.EastmoneyIndustryTrendsURL = trendsURL
		}
	}
}

// WithEastmoneyConceptListURL sets the Eastmoney concept board list URL.
func WithEastmoneyConceptListURL(listURL string) Option {
	return func(config *Config) {
		if listURL != "" {
			config.EastmoneyConceptListURL = listURL
		}
	}
}

// WithEastmoneyConceptSpotURL sets the Eastmoney concept board spot URL.
func WithEastmoneyConceptSpotURL(spotURL string) Option {
	return func(config *Config) {
		if spotURL != "" {
			config.EastmoneyConceptSpotURL = spotURL
		}
	}
}

// WithEastmoneyConceptConstituentsURL sets the Eastmoney concept board constituents URL.
func WithEastmoneyConceptConstituentsURL(consURL string) Option {
	return func(config *Config) {
		if consURL != "" {
			config.EastmoneyConceptConstituentsURL = consURL
		}
	}
}

// WithEastmoneyConceptKlineURL sets the Eastmoney concept board K-line URL.
func WithEastmoneyConceptKlineURL(klineURL string) Option {
	return func(config *Config) {
		if klineURL != "" {
			config.EastmoneyConceptKlineURL = klineURL
		}
	}
}

// WithEastmoneyConceptTrendsURL sets the Eastmoney concept board trends URL.
func WithEastmoneyConceptTrendsURL(trendsURL string) Option {
	return func(config *Config) {
		if trendsURL != "" {
			config.EastmoneyConceptTrendsURL = trendsURL
		}
	}
}
