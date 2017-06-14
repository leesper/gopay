package ali

type RefundTradeParam struct {
	AppAuthToken string `json:"-"`              // 可选
	OutTradeNo   string `json:"out_trade_no"`   // 商户订单号 与TradeNo二选一
	TradeNo      string `json:"trade_no"`       // 支付宝交易号 与OutTradeNo二选一
	RefundAmount string `json:"refund_amount"`  // 必须 退款的金额
	RefundReason string `json:"refund_reason"`  // 可选 退款的原因说明
	OutRequestNo string `json:"out_request_no"` // 可选 标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传
	OperatorID   string `json:"operator_id"`    // 可选 商户的操作员编号
	StoreID      string `json:"store_id"`       // 可选 商户的门店编号
	TerminalID   string `json:"terminal_id"`    // 可选 商户的终端编号
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
		BuyerLogonID         string `json:"buyer_logon_id"`
		BuyerUserID          string `json:"buyer_user_id"`
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
