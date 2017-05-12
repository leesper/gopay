package ali

// QueryTradeParam requests alipay.trade.query.
type QueryTradeParam struct {
	AppAuthToken string `json:"-"`
	OutTradeNo   string `json:"out_trade_no,omitempty"`
	TradeNo      string `json:"trade_no,omitempty"`
}

// URI returns the uri.
func (p QueryTradeParam) URI() string {
	return "alipay.trade.query"
}

// ExtraParams returns extra parameters.
func (p QueryTradeParam) ExtraParams() map[string]string {
	return map[string]string{
		"app_auth_token": p.AppAuthToken,
	}
}

// BizContent returns biz_content in JSON format.
func (p QueryTradeParam) BizContent() string {
	return marshalJSON(p)
}

// QueryTradeRsp responses alipay.trade.query.
type QueryTradeRsp struct {
	TradeQueryResponse struct {
		Code           string `json:"code"`
		Msg            string `json:"msg"`
		BuyerLogonID   string `json:"buyer_logon_id"`   // 买家支付宝账号
		BuyerPayAmount string `json:"buyer_pay_amount"` // 买家实付金额，单位为元，两位小数。
		BuyerUserID    string `json:"buyer_user_id"`    // 买家在支付宝的用户id
		InvoiceAmount  string `json:"invoice_amount"`   // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
		OpenID         string `json:"open_id"`          // 买家支付宝用户号，该字段将废弃，不要使用
		OutTradeNo     string `json:"out_trade_no"`     // 商家订单号
		PointAmount    string `json:"point_amount"`     // 积分支付的金额，单位为元，两位小数。
		ReceiptAmount  string `json:"receipt_amount"`
		SendPayDate    string `json:"send_pay_date"` // 本次交易打款给卖家的时间
		TotalAmount    string `json:"total_amount"`  // 交易的订单金额
		TradeNo        string `json:"trade_no"`      // 支付宝交易号
		TradeStatus    string `json:"trade_status"`  // 交易状态
	} `json:"alipay_trade_query_response"`
	Sign string `json:"sign"`
}

// VoucherDetail 本交易支付时使用的所有优惠券信息
// type VoucherDetail struct {
// 	ID                 string `json:"id"`                  // 券id
// 	Name               string `json:"name"`                // 券名称
// 	Type               string `json:"type"`                // 当前有三种类型： ALIPAY_FIX_VOUCHER - 全场代金券, ALIPAY_DISCOUNT_VOUCHER - 折扣券, ALIPAY_ITEM_VOUCHER - 单品优惠
// 	Amount             string `json:"amount"`              // 优惠券面额，它应该会等于商家出资加上其他出资方出资
// 	MerchantContribute string `json:"merchant_contribute"` // 商家出资（特指发起交易的商家出资金额）
// 	OtherContribute    string `json:"other_contribute"`    // 其他出资方出资金额，可能是支付宝，可能是品牌商，或者其他方，也可能是他们的一起出资
// 	Memo               string `json:"memo"`                // 优惠券备注信息
// }
