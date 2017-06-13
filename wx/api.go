package wx

import "encoding/xml"

// Payment returns to App.
type Payment struct {
	AppID     string
	PartnerID string
	PrepayID  string
	NonceStr  string
	Timestamp string
	Package   string
	Sign      string
}

type unifiedOrderReq struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`            // 应用ID
	MchID          string   `xml:"mch_id"`           // 商户号
	NonceStr       string   `xml:"nonce_str"`        // 随机字符串
	Body           string   `xml:"body"`             // 商品描述
	Attach         string   `xml:"attach"`           // 附加数据
	OutTradeNo     string   `xml:"out_trade_no"`     // 商户订单号
	TotalFee       string   `xml:"total_fee"`        // 总金额
	SpbillCreateIP string   `xml:"spbill_create_ip"` // 终端IP
	NotifyURL      string   `xml:"notify_url"`       // 通知地址
	TradeType      string   `xml:"trade_type"`       // 交易类型
}

func (req unifiedOrderReq) URI() string {
	return "https://api.mch.weixin.qq.com/pay/unifiedorder"
}

func (req unifiedOrderReq) SandBoxURI() string {
	return "https://api.mch.weixin.qq.com/sandboxnew/pay/unifiedorder"
}

// UnifiedOrderRsp is the response returned by /pay/unifiedorder.
type UnifiedOrderRsp struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"`  // 返回状态码
	ReturnMsg   string   `xml:"return_msg"`   // 返回信息
	AppID       string   `xml:"appid"`        // 应用APPID
	MchID       string   `xml:"mch_id"`       // 商户号
	DeviceInfo  string   `xml:"device_info"`  // 设备号
	NonceStr    string   `xml:"nonce_str"`    // 随机字符串
	Sign        string   `xml:"sign"`         // 签名
	ResultCode  string   `xml:"result_code"`  // 业务结果
	ErrCode     string   `xml:"err_code"`     // 错误代码
	ErrCodeDesc string   `xml:"err_code_des"` // 错误代码描述
	TradeType   string   `xml:"trade_type"`   // 交易类型
	PrepayID    string   `xml:"prepay_id"`    // 预支付交易会话标识
}

type queryOrderReq struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid"`          // 应用APPID
	MchID         string   `xml:"mch_id"`         // 商户号
	TransactionID string   `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	NonceStr      string   `xml:"nonce_str"`      // 随机字符串
}

func (req queryOrderReq) URI() string {
	return "https://api.mch.weixin.qq.com/pay/orderquery"
}

func (req queryOrderReq) SandBoxURI() string {
	return "https://api.mch.weixin.qq.com/sandboxnew/pay/orderquery"
}

// QueryOrderRsp is the response returned by /pay/orderquery
type QueryOrderRsp struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`      // 返回状态码
	ReturnMsg      string   `xml:"return_msg"`       // 返回信息
	AppID          string   `xml:"appid"`            // 应用APPID
	MchID          string   `xml:"mch_id"`           // 商户号
	NonceStr       string   `xml:"nonce_str"`        // 随机字符串
	Sign           string   `xml:"sign"`             // 签名
	ResultCode     string   `xml:"result_code"`      // 业务结果
	ErrCode        string   `xml:"err_code"`         // 错误代码
	ErrCodeDesc    string   `xml:"err_code_des"`     // 错误代码描述
	DeviceInfo     string   `xml:"device_info"`      // 设备号
	OpenID         string   `xml:"openid"`           // 用户标识
	IsSubscribe    string   `xml:"is_subscribe"`     // 是否关注公众账号
	TradeType      string   `xml:"trade_type"`       // 交易类型
	TradeState     string   `xml:"trade_state"`      // 交易状态
	BankType       string   `xml:"bank_type"`        // 付款银行
	TotalFee       string   `xml:"total_fee"`        // 总金额
	FeeType        string   `xml:"fee_type"`         // 货币种类
	CashFee        string   `xml:"cash_fee"`         // 现金支付金额
	CashFeeType    string   `xml:"cash_fee_type"`    // 现金支付货币类型
	CouponFee      string   `xml:"coupon_fee"`       // 代金券或立减优惠金额
	CouponCount    string   `xml:"coupon_count"`     // 代金券或立减优惠使用数量
	TransactionID  string   `xml:"transaction_id"`   // 微信支付订单号
	OutTradeNo     string   `xml:"out_trade_no"`     // 商户订单号
	Attach         string   `xml:"attach"`           // 附加数据
	TimeEnd        string   `xml:"time_end"`         // 支付完成时间
	TradeStateDesc string   `xml:"trade_state_desc"` // 交易状态描述
}

type refundOrderReq struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid"`          // 应用ID
	MchID         string   `xml:"mch_id"`         // 商户号
	NonceStr      string   `xml:"nonce_str"`      // 随机字符串
	TransactionID string   `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	OutRefundNo   string   `xml:"out_refund_no"`  // 商户退款单号
	TotalFee      string   `xml:"total_fee"`      // 总金额
	RefundFee     string   `xml:"refund_fee"`     // 退款金额
}

func (req refundOrderReq) URI() string {
	return "https://api.mch.weixin.qq.com/secapi/pay/refund"
}

func (req refundOrderReq) SandBoxURI() string {
	return "https://api.mch.weixin.qq.com/sandboxnew/secapi/pay/refund"
}

type RefundOrderRsp struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`    // 返回状态码
	ReturnMsg     string   `xml:"return_msg"`     // 返回信息
	ResultCode    string   `xml:"result_code"`    // 业务结果
	ErrCode       string   `xml:"err_code"`       // 错误代码
	ErrCodeDesc   string   `xml:"err_code_des"`   // 错误代码描述
	AppID         string   `xml:"appid"`          // 应用APPID
	MchID         string   `xml:"mch_id"`         // 商户号
	NonceStr      string   `xml:"nonce_str"`      // 随机字符串
	TransactionID string   `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	OutRefundNo   string   `xml:"out_refund_no"`  // 商户退款单号
	RefundID      string   `xml:"refund_id"`      // 微信退款单号
	RefundFee     string   `xml:"refund_fee"`     // 退款金额
	TotalFee      string   `xml:"total_fee"`      // 标价金额
	CashFee       string   `xml:"cash_fee"`       // 现金支付金额
}

type queryRefundReq struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid"`          // 应用ID
	MchID         string   `xml:"mch_id"`         // 商户号
	NonceStr      string   `xml:"nonce_str"`      // 随机字符串
	TransactionID string   `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	OutRefundNo   string   `xml:"out_refund_no"`  // 商户退款单号
	RefundID      string   `xml:"refund_id"`      // 微信退款单号
}

func (req queryRefundReq) URI() string {
	return "https://api.mch.weixin.qq.com/pay/refundquery"
}

func (req queryRefundReq) SandBoxURI() string {
	return "https://api.mch.weixin.qq.com/sandboxnew/pay/refundquery"
}

type QueryRefundRsp struct {
	ReturnCode    string `xml:"return_code"`    // 返回状态码
	ReturnMsg     string `xml:"return_msg"`     // 返回信息
	ResultCode    string `xml:"result_code"`    // 业务结果
	ErrCode       string `xml:"err_code"`       // 错误代码
	ErrCodeDesc   string `xml:"err_code_des"`   // 错误代码描述
	AppID         string `xml:"appid"`          // 应用APPID
	MchID         string `xml:"mch_id"`         // 商户号
	NonceStr      string `xml:"nonce_str"`      // 随机字符串
	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	TotalFee      string `xml:"total_fee"`      // 标价金额
	CashFee       string `xml:"cash_fee"`       // 现金支付金额
	RefundCount   string `xml:"refund_count"`   // 退款笔数
}

type getSandBoxSignKeyReq struct {
	MchID    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
}

func (req getSandBoxSignKeyReq) SandBoxURI() string {
	return "https://api.mch.weixin.qq.com/sandboxnew/pay/getsignkey"
}

// GetSandBoxSignKeyRsp is the response returned by
type GetSandBoxSignKeyRsp struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	MchID          string `xml:"mch_id"`
	SandBoxSignKey string `xml:"sandbox_signkey"`
}
