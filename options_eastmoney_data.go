package stock

// WithEastmoneyFundFlowURL sets the Eastmoney fund-flow URL.
func WithEastmoneyFundFlowURL(fflowURL string) Option {
	return func(config *Config) {
		if fflowURL != "" {
			config.EastmoneyFundFlowURL = fflowURL
		}
	}
}

// WithEastmoneyClistURL sets the Eastmoney clist URL.
func WithEastmoneyClistURL(clistURL string) Option {
	return func(config *Config) {
		if clistURL != "" {
			config.EastmoneyClistURL = clistURL
		}
	}
}

// WithEastmoneyNorthboundMinuteURL sets the Eastmoney northbound minute URL.
func WithEastmoneyNorthboundMinuteURL(minuteURL string) Option {
	return func(config *Config) {
		if minuteURL != "" {
			config.EastmoneyNorthboundMinuteURL = minuteURL
		}
	}
}

// WithEastmoneyDatacenterURL sets the Eastmoney datacenter URL.
func WithEastmoneyDatacenterURL(datacenterURL string) Option {
	return func(config *Config) {
		if datacenterURL != "" {
			config.EastmoneyDatacenterURL = datacenterURL
		}
	}
}

// WithEastmoneyTopicURL sets the Eastmoney topic data base URL.
func WithEastmoneyTopicURL(topicURL string) Option {
	return func(config *Config) {
		if topicURL != "" {
			config.EastmoneyTopicURL = topicURL
		}
	}
}

// WithEastmoneyFundGZURL sets the Eastmoney fund estimate URL.
func WithEastmoneyFundGZURL(fundGZURL string) Option {
	return func(config *Config) {
		if fundGZURL != "" {
			config.EastmoneyFundGZURL = fundGZURL
		}
	}
}

// WithEastmoneyFundPingzhongURL sets the Eastmoney fund pingzhongdata URL.
func WithEastmoneyFundPingzhongURL(pingzhongURL string) Option {
	return func(config *Config) {
		if pingzhongURL != "" {
			config.EastmoneyFundPingzhongURL = pingzhongURL
		}
	}
}

// WithEastmoneyFundDataIndexURL sets the Eastmoney fund data index URL.
func WithEastmoneyFundDataIndexURL(dataIndexURL string) Option {
	return func(config *Config) {
		if dataIndexURL != "" {
			config.EastmoneyFundDataIndexURL = dataIndexURL
		}
	}
}
