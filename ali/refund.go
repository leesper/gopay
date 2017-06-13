package ali

type RefundTradeParam struct {
	AppAuthToken string `json:"-"`
	OutTradeNo   string `json:"out_trade_no"`
	TradeNo      string `json:"trade_no"`
	RefundAmount string `json:"refund_amount"`
	RefundReason string `json:"refund_reason"`
	OutRequestNo string `json:"out_request_no"`
	OperatorID   string `json:"operator_id"`
	StoreID      string `json:"store_id"`
	TerminalID   string `json:"terminal_id"`
}

func (p RefundTradeParam) URI() string {
	return "alipay.trade.refund"
}

func (p RefundTradeParam) ExtraParams() map[string]string {
	return map[string]string{
		"app_auth_token": p.AppAuthToken,
	}
}

func (p RefundTradeParam) BizContent() string {
	return marshalJSON(p)
}

type RefundTradeRsp struct {
	TradeRefundResponse struct {
		Code                 string `json:"code"`
		Msg                  string `json:"msg"`
		SubCode              string `json:"sub_code"`
		SubMsg               string `json:"sub_msg"`
		TradeNo              string `json:"trade_no"`
		OutTradeNo           string `json:"out_trade_no"`
		OpenID               string `json:"open_id"`
		BuyerLogonID         string `json:"buyer_logon_id"`
		FundChange           string `json:"fund_change"`
		RefundFee            string `json:"refund_fee"`
		GmtRefundPay         string `json:"gmt_refund_pay"`
		StoreName            string `store_name`
		RefundDetailItemList []struct {
			FundChannel string `json:"fund_channel"`
			Amount      string `json:"amount"`
			RealAmount  string `json:"real_amount"`
		} `json:"refund_detail_item_list"`
	} `json:"alipay_trade_refund_response"`
	Sign string `json:"sign"`
}

func (r *RefundTradeRsp) Success() bool {
	return r.TradeRefundResponse.Msg == "Success"
}

type QueryRefundParam struct {
	AppAuthToken string `json:"-"`
	OutTradeNo   string `json:"out_trade_no,omitempty"`
	TradeNo      string `json:"trade_no,omitempty"`
	OutRequestNo string `json:"out_request_no"`
}

func (p QueryRefundParam) URI() string {
	return "alipay.trade.fastpay.refund.query"
}

func (p QueryRefundParam) ExtraParams() map[string]string {
	return map[string]string{
		"app_auth_token": p.AppAuthToken,
	}
}

func (p QueryRefundParam) BizContent() string {
	return marshalJSON(p)
}

type QueryRefundRsp struct {
	RefundQueryResponse struct {
		Code         string `json:"code"`
		Msg          string `json:"msg"`
		SubCode      string `json:"sub_code"`
		SubMsg       string `json:"sub_msg"`
		TradeNo      string `json:"trade_no"`       // 支付宝交易号
		OutTradeNo   string `json:"out_trade_no"`   // 创建交易传入的商户订单号
		OutRequestNo string `json:"out_request_no"` // 本笔退款对应的退款请求号
		RefundReason string `json:"refund_reason"`  // 发起退款时，传入的退款原因
		TotalAmount  string `json:"total_amount"`   // 发该笔退款所对应的交易的订单金额
		RefundAmount string `json:"refund_amount"`  // 本次退款请求，对应的退款金额
	} `json:"alipay_trade_fastpay_refund_query_response"`
	Sign string `json:"sign"`
}
