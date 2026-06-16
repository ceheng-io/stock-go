package stock

import (
	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/internal/services"
)

// Client is the root SDK entry point.
type Client struct {
	config Config
	core   *core.Client

	Quotes      *services.QuoteService
	Calendar    *services.CalendarService
	Kline       *services.KlineService
	Indicator   *services.IndicatorService
	Board       *services.BoardService
	FundFlow    *services.FundFlowService
	Northbound  *services.NorthboundService
	DragonTiger *services.DragonTigerService
	BlockTrade  *services.BlockTradeService
	Margin      *services.MarginService
	Dividend    *services.DividendService
	Data        *services.DataService
	MarketEvent *services.MarketEventService
	Fund        *services.FundService
	Futures     *services.FuturesService
	Options     *services.OptionsService
}

// StockSDK preserves the TypeScript SDK root class name at the Go root package.
type StockSDK = Client

// New creates a Client with the provided options.
func New(options ...Option) *Client {
	config := defaultConfig()
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}
	coreClient := core.NewClient(core.Config{
		BaseURL:           config.BaseURL,
		TencentMinuteURL:  config.TencentMinuteURL,
		SearchBaseURL:     config.SearchBaseURL,
		CalendarURL:       config.CalendarURL,
		AShareListURL:     config.AShareListURL,
		USListURL:         config.USListURL,
		HKListURL:         config.HKListURL,
		FundListURL:       config.FundListURL,
		EastmoneyKlineURL: config.EastmoneyKlineURL,
		HTTPClient:        config.HTTPClient,
		ProxyPool: core.ProxyPoolConfig{
			URLs: append([]string(nil), config.ProxyPool.URLs...),
		},
		Timeout:         config.Timeout,
		UserAgent:       config.UserAgent,
		RotateUserAgent: config.RotateUserAgent,
		Headers:         cloneStringMap(config.Headers),
		Retry: core.RetryConfig{
			MaxRetries:           config.Retry.MaxRetries,
			BaseDelay:            config.Retry.BaseDelay,
			MaxDelay:             config.Retry.MaxDelay,
			BackoffMultiplier:    config.Retry.BackoffMultiplier,
			RetryableStatusCodes: append([]int(nil), config.Retry.RetryableStatusCodes...),
			RetryOnNetworkError:  cloneBoolPtr(config.Retry.RetryOnNetworkError),
			RetryOnTimeout:       cloneBoolPtr(config.Retry.RetryOnTimeout),
			OnRetry:              config.Retry.OnRetry,
		},
		RateLimiter:      newCoreRateLimiter(config.RateLimit),
		CircuitBreaker:   newCoreCircuitBreaker(config.CircuitBreaker),
		ProviderPolicies: newCoreProviderPolicies(config.ProviderPolicies),
		Hooks:            config.RequestHooks,
	})
	quoteService := services.NewQuoteService(coreClient, services.QuoteURLs{
		Minute:     config.TencentMinuteURL,
		Datacenter: config.EastmoneyDatacenterURL,
	})
	klineService := services.NewKlineService(coreClient, services.KlineURLs{
		CN:       config.EastmoneyKlineURL,
		CNTrends: config.EastmoneyTrendsURL,
		HK:       config.EastmoneyHKKlineURL,
		HKTrends: config.EastmoneyHKTrendsURL,
		US:       config.EastmoneyUSKlineURL,
		USTrends: config.EastmoneyUSTrendsURL,
	})
	return &Client{
		config:    config,
		core:      coreClient,
		Quotes:    quoteService,
		Calendar:  services.NewCalendarService(quoteService),
		Kline:     klineService,
		Indicator: services.NewIndicatorService(klineService, quoteService),
		Board: services.NewBoardService(coreClient, services.BoardURLs{
			IndustryList:         config.EastmoneyIndustryListURL,
			IndustrySpot:         config.EastmoneyIndustrySpotURL,
			IndustryConstituents: config.EastmoneyIndustryConstituentsURL,
			IndustryKline:        config.EastmoneyIndustryKlineURL,
			IndustryTrends:       config.EastmoneyIndustryTrendsURL,
			ConceptList:          config.EastmoneyConceptListURL,
			ConceptSpot:          config.EastmoneyConceptSpotURL,
			ConceptConstituents:  config.EastmoneyConceptConstituentsURL,
			ConceptKline:         config.EastmoneyConceptKlineURL,
			ConceptTrends:        config.EastmoneyConceptTrendsURL,
		}),
		FundFlow: services.NewFundFlowService(coreClient, services.FundFlowURLs{
			FFlow: config.EastmoneyFundFlowURL,
			Clist: config.EastmoneyClistURL,
		}),
		Northbound: services.NewNorthboundService(coreClient, services.NorthboundURLs{
			Minute:     config.EastmoneyNorthboundMinuteURL,
			Datacenter: config.EastmoneyDatacenterURL,
		}),
		DragonTiger: services.NewDragonTigerService(coreClient, config.EastmoneyDatacenterURL),
		BlockTrade:  services.NewBlockTradeService(coreClient, config.EastmoneyDatacenterURL),
		Margin:      services.NewMarginService(coreClient, config.EastmoneyDatacenterURL),
		Dividend:    services.NewDividendService(coreClient, config.EastmoneyDatacenterURL),
		Data: services.NewDataService(coreClient, services.DataServiceOptions{
			DatacenterURL:       config.EastmoneyDatacenterURL,
			F10BaseURL:          config.EastmoneyF10BaseURL,
			AnnouncementListURL: config.EastmoneyAnnouncementListURL,
			AnnouncementURL:     config.EastmoneyAnnouncementURL,
		}),
		MarketEvent: services.NewMarketEventService(coreClient, config.EastmoneyTopicURL, config.THSLimitUpPoolURL),
		Fund: services.NewFundService(coreClient, services.FundURLs{
			GZ:        config.EastmoneyFundGZURL,
			Pingzhong: config.EastmoneyFundPingzhongURL,
			DataIndex: config.EastmoneyFundDataIndexURL,
		}),
		Futures: services.NewFuturesService(coreClient, services.FuturesURLs{
			Kline:       config.EastmoneyFuturesKlineURL,
			GlobalSpot:  config.EastmoneyFuturesGlobalSpotURL,
			GlobalKline: config.EastmoneyFuturesGlobalKlineURL,
			Datacenter:  config.EastmoneyDatacenterURL,
		}),
		Options: services.NewOptionsService(coreClient, services.OptionsURLs{
			CFFEXQuotes:    config.EastmoneyOptionCFFEXURL,
			LHB:            config.EastmoneyOptionLHBURL,
			ETFMonths:      config.SinaETFOptionListURL,
			ETFExpire:      config.SinaETFOptionExpireURL,
			ETFMinute:      config.SinaETFOptionMinuteURL,
			ETFDaily:       config.SinaETFOptionDailyURL,
			ETF5Day:        config.SinaETFOption5DayURL,
			IndexSpot:      config.SinaIndexOptionSpotURL,
			IndexKline:     config.SinaIndexOptionKlineURL,
			CommoditySpot:  config.SinaCommodityOptionSpotURL,
			CommodityKline: config.SinaCommodityOptionKlineURL,
		}),
	}
}

func newCoreRateLimiter(options RateLimitOptions) core.RequestLimiter {
	if options.RequestsPerSecond <= 0 {
		return nil
	}
	return core.NewRateLimiter(core.RateLimiterOptions{
		RequestsPerSecond: options.RequestsPerSecond,
		MaxBurst:          options.MaxBurst,
	})
}

func newCoreCircuitBreaker(options CircuitBreakerOptions) core.RequestCircuitBreaker {
	if options.FailureThreshold <= 0 {
		return nil
	}
	return core.NewCircuitBreaker(core.CircuitBreakerOptions{
		FailureThreshold: options.FailureThreshold,
		ResetTimeout:     options.ResetTimeout,
		HalfOpenRequests: options.HalfOpenRequests,
		OnStateChange:    options.OnStateChange,
	})
}

func newCoreProviderPolicies(policies map[ProviderName]ProviderPolicy) map[core.ProviderName]core.ProviderPolicy {
	if len(policies) == 0 {
		return nil
	}
	converted := make(map[core.ProviderName]core.ProviderPolicy, len(policies))
	for provider, policy := range policies {
		converted[core.ProviderName(provider)] = core.ProviderPolicy{
			Timeout:         policy.Timeout,
			UserAgent:       policy.UserAgent,
			RotateUserAgent: cloneBoolPtr(policy.RotateUserAgent),
			Headers:         cloneStringMap(policy.Headers),
			Retry:           newCoreRetryConfig(policy.Retry),
			RateLimiter:     newCoreRateLimiterPtr(policy.RateLimit),
			CircuitBreaker:  newCoreCircuitBreakerPtr(policy.CircuitBreaker),
		}
	}
	return converted
}

func newCoreRetryConfig(options *RetryOptions) *core.RetryConfig {
	if options == nil {
		return nil
	}
	return &core.RetryConfig{
		MaxRetries:           options.MaxRetries,
		BaseDelay:            options.BaseDelay,
		MaxDelay:             options.MaxDelay,
		BackoffMultiplier:    options.BackoffMultiplier,
		RetryableStatusCodes: append([]int(nil), options.RetryableStatusCodes...),
		RetryOnNetworkError:  cloneBoolPtr(options.RetryOnNetworkError),
		RetryOnTimeout:       cloneBoolPtr(options.RetryOnTimeout),
		OnRetry:              options.OnRetry,
	}
}

func newCoreRateLimiterPtr(options *RateLimitOptions) core.RequestLimiter {
	if options == nil {
		return nil
	}
	return newCoreRateLimiter(*options)
}

func newCoreCircuitBreakerPtr(options *CircuitBreakerOptions) core.RequestCircuitBreaker {
	if options == nil {
		return nil
	}
	return newCoreCircuitBreaker(*options)
}
