package ali

// CreateTradeParam requests alipay.trade.create.
type CreateTradeParam struct {
	AppAuthToken         string         `json:"-"`
	NotifyURL            string         `json:"-"`
	OutTradeNo           string         `json:"out_trade_no,omitempty"`
	SellerID             string         `json:"seller_id,omitempty"`
	TotalAmount          string         `json:"total_amount"`
	DiscountableAmount   string         `json:"discountable_amount"`
	UndiscountableAmount string         `json:"undiscountable_amount"`
	BuyerLogonID         string         `json:"buyer_logon_id"`
	Subject              string         `json:"subject"`
	Body                 string         `json:"body"`
	BuyerID              string         `json:"buyer_id"`
	GoodsDetails         []*GoodsDetail `json:"goods_detail,omitempty"`
	OperatorID           string         `json:"operator_id"`
	StoreID              string         `json:"store_id"`
	TerminalID           string         `json:"terminal_id"`
	ExtendParams         *ExtendParam   `json:"extend_params,omitempty"`
	TimeoutExpress       string         `json:"timeout_express"`
	RoyaltyInfo          *RoyaltyInfo   `json:"royalty_info,omitempty"`
	AliPayStoreID        string         `json:"alipay_store_id"`
	SubMerchants         []SubMerchant  `json:"sub_merchant"`
	MerchantOrderNo      string         `json:"merchant_order_no"`
}

// URI returns the uri.
func (p CreateTradeParam) URI() string {
	return "alipay.trade.create"
}

// ExtraParams returns extra parameters.
func (p CreateTradeParam) ExtraParams() map[string]string {
	return map[string]string{
		"app_auth_token": p.AppAuthToken,
		"notify_url":     p.NotifyURL,
	}
}

// BizContent returns biz_content in JSON format.
func (p CreateTradeParam) BizContent() string {
	return marshalJSON(p)
}

// CreateTradeRsp responses alipay.trade.create.
type CreateTradeRsp struct {
	TradeCreateResponse struct {
		Code       string `json:"code"`
		Msg        string `json:"msg"`
		SubCode    string `json:"sub_code"`
		SubMsg     string `json:"sub_msg"`
		TradeNo    string `json:"trade_no"`
		OutTradeNo string `json:"out_trade_no"`
	} `json:"alipay_trade_create_response"`
	Sign string `json:"sign"`
}

// SubMerchant 二级商户信息.
type SubMerchant struct {
	MerchantID string `json:"merchant_id"`
}

// GoodsDetail 订单包含的商品列表信息.
type GoodsDetail struct {
	GoodsID       string `json:"goods_id"`
	AliPayGoodsID string `json:"alipay_goods_id"`
	GoodsName     string `json:"goods_name"`
	Quantity      string `json:"quantity"`
	Price         string `json:"price"`
	GoodsCategory string `json:"goods_category"`
	Body          string `json:"body"`
	ShowURL       string `json:"show_url"`
}

// ExtendParam 业务扩展参数.
type ExtendParam struct {
	SysServiceProviderID string `json:"sys_service_provider_id"`
	HbFqNum              string `json:"hb_fq_num"`
	HbFqSellerPercent    string `json:"hb_fq_seller_percent"`
	TimeoutExpress       string `json:"timeout_express"`
}

// RoyaltyInfo 描述分账信息.
type RoyaltyInfo struct {
	RoyaltyType        string               `json:"royalty_type"`
	RoyaltyDetailInfos []*RoyaltyDetailInfo `json:"royalty_detail_infos,omitempty"`
}

// RoyaltyDetailInfo 分账明细的信息.
type RoyaltyDetailInfo struct {
	SerialNo         string `json:"serial_no"`
	TransInType      string `json:"trans_in_type"`
	BatchNo          string `json:"batch_no"`
	OutRelationID    string `json:"out_relation_id"`
	TransOutType     string `json:"trans_out_type"`
	TransOut         string `json:"trans_out"`
	TransIn          string `json:"trans_in"`
	Amount           string `json:"amount"`
	Desc             string `json:"desc"`
	AmountPercentage string `json:"amount_percentage"`
	AliPayStoreID    string `json:"alipay_store_id"`
}
