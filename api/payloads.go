package oandaapi

type GetInstrumentsReq struct {
	Instruments []string // query param - csv format
}

type GetInstrumentsRes struct {
	Instruments       []*Instrument `json:"instruments"`
	LastTransactionID string        `json:"lastTransactionID"`
}

type GetCandlesReq struct {
	Instrument        string
	Granularity       string
	Price             string
	Count             int64
	From              string
	To                string
	Smooth            bool
	IncludeFirst      bool
	DailyAlignment    int64
	AlignmentTimezone string
	WeeklyAlignment   string
}

type GetCandlesRes struct {
	Instrument  string         `json:"instrument"`
	Granularity string         `json:"granularity"`
	Candles     []*Candlestick `json:"candles"`
}

type GetOrderBookReq struct {
	Instrument string
}

type GetOrderBookRes struct {
	Instrument  string `json:"instrument"`
	Time        string `json:"time"`
	Price       string `json:"price"`
	BucketWidth string `json:"bucketWidth"`
	Buckets     []struct {
		Price             string `json:"price"`
		LongCountPercent  string `json:"longCountPercent"`
		shortCountPercent string `json:"shortCountPercent"`
	} `json:"buckets"`
}

type GetOrdersReq struct {
	IDs        []string
	State      string
	Instrument string
	Count      int64
	BeforeID   string
}

type GetOrdersRes struct {
	Orders []struct {
		ID               string `json:"id"`
		CreateTime       string `json:"createTime"`
		State            string `json:"state"`
		ClientExtensions struct {
			ClientID string `json:"clientID"`
			Tag      string `json:"clientTag"`
			Comment  string `json:"clientComment"`
		} `json:"clientExtensions"`
	} `json:"orders"`
	LastTransactionID string `json:"lastTransactionID"`
}

type MarketOrder struct {
}

type StreamRes struct {
}
