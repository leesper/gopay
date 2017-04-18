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
		SubCode        string `json:"sub_code"`
		SubMsg         string `json:"sub_msg"`
		TradeNo        string `json:"trade_no"`       // 支付宝交易号
		OutTradeNo     string `json:"out_trade_no"`   // 商家订单号
		OpenID         string `json:"open_id"`        // 买家支付宝用户号，该字段将废弃，不要使用
		BuyerLogonID   string `json:"buyer_logon_id"` // 买家支付宝账号
		TradeStatus    string `json:"trade_status"`   // 交易状态
		TotalAmount    string `json:"total_amount"`   // 交易的订单金额
		ReceiptAmount  string `json:"receipt_amount"`
		BuyerPayAmount string `json:"buyer_pay_amount"` // 买家实付金额，单位为元，两位小数。
		PointAmount    string `json:"point_amount"`     // 积分支付的金额，单位为元，两位小数。
		InvoiceAmount  string `json:"invoice_amount"`   // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
		SendPayDate    string `json:"send_pay_date"`    // 本次交易打款给卖家的时间
		AliPayStoreID  string `json:"alipay_store_id"`  // 支付宝店铺编号
		StoreID        string `json:"store_id"`         // 商户门店编号
		TerminalID     string `json:"terminal_id"`      // 商户机具终端编号
		FundBillList   []struct {
			FundChannel string `json:"fund_channel"` // 交易使用的资金渠道，详见 支付渠道列表
			Amount      string `json:"amount"`       // 该支付工具类型所使用的金额
			RealAmount  string `json:"real_amount"`  // 渠道实际付款金额
		} `json:"fund_bill_list"` // 交易支付使用的资金渠道
		StoreName           string          `json:"store_name"`            // 请求交易支付中的商户店铺的名称
		BuyerUserID         string          `json:"buyer_user_id"`         // 买家在支付宝的用户id
		DiscountGoodsDetail string          `json:"discount_goods_detail"` // 本次交易支付所使用的单品券优惠的商品优惠信息
		IndustrySepcDetail  string          `json:"industry_sepc_detail"`  // 行业特殊信息（例如在医保卡支付业务中，向用户返回医疗信息）。
		VoucherDetailList   []VoucherDetail `json:"voucher_detail_list"`   // 本交易支付时使用的所有优惠券信息
	} `json:"alipay_trade_query_response"`
	Sign string `json:"sign"`
}

// VoucherDetail 本交易支付时使用的所有优惠券信息
type VoucherDetail struct {
	ID                 string `json:"id"`                  // 券id
	Name               string `json:"name"`                // 券名称
	Type               string `json:"type"`                // 当前有三种类型： ALIPAY_FIX_VOUCHER - 全场代金券, ALIPAY_DISCOUNT_VOUCHER - 折扣券, ALIPAY_ITEM_VOUCHER - 单品优惠
	Amount             string `json:"amount"`              // 优惠券面额，它应该会等于商家出资加上其他出资方出资
	MerchantContribute string `json:"merchant_contribute"` // 商家出资（特指发起交易的商家出资金额）
	OtherContribute    string `json:"other_contribute"`    // 其他出资方出资金额，可能是支付宝，可能是品牌商，或者其他方，也可能是他们的一起出资
	Memo               string `json:"memo"`                // 优惠券备注信息
}
