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
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Body           string   `xml:"body"`
	Attach         string   `xml:"attach"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       string   `xml:"total_fee"`
	SpbillCreateIP string   `xml:"spbill_create_ip"`
	NotifyURL      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
}

func (req unifiedOrderReq) URI() string {
	return "https://api.mch.weixin.qq.com/pay/unifiedorder"
}

// UnifiedOrderRsp is the response returned by /pay/unifiedorder.
type UnifiedOrderRsp struct {
	XMLName     xml.Name `xml:"xml"`
	ReturnCode  string   `xml:"return_code"`
	ReturnMsg   string   `xml:"return_msg"`
	AppID       string   `xml:"appid"`
	MchID       string   `xml:"mch_id"`
	DeviceInfo  string   `xml:"device_info"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	ResultCode  string   `xml:"result_code"`
	ErrCode     string   `xml:"err_code"`
	ErrCodeDesc string   `xml:"err_code_des"`
	TradeType   string   `xml:"trade_type"`
	PrepayID    string   `xml:"prepay_id"`
	CodeURL     string   `xml:"code_url"`
}

type queryOrderReq struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid"`
	MchID         string   `xml:"mch_id"`
	TransactionID string   `xml:"transaction_id"`
	NonceStr      string   `xml:"nonce_str"`
}

func (req queryOrderReq) URI() string {
	return "https://api.mch.weixin.qq.com/pay/orderquery"
}

// QueryOrderRsp is the response returned by /pay/orderquery
type QueryOrderRsp struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`
	ReturnMsg      string   `xml:"return_msg"`
	AppID          string   `xml:"appid"`
	MchID          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	ResultCode     string   `xml:"result_code"`
	ErrCode        string   `xml:"err_code"`
	ErrCodeDesc    string   `xml:"err_code_des"`
	DeviceInfo     string   `xml:"device_info"`
	OpenID         string   `xml:"open_id"`
	IsSubscribe    string   `xml:"is_subscribe"`
	TradeType      string   `xml:"trade_type"`
	TradeState     string   `xml:"trade_state"`
	TradeStateDesc string   `xml:"trade_state_desc"`
	BankType       string   `xml:"bank_type"`
	TotalFee       string   `xml:"total_fee"`
	FeeType        string   `xml:"fee_type"`
	CashFee        string   `xml:"cash_fee"`
	CashFeeType    string   `xml:"cash_fee_type"`
	CouponFee      string   `xml:"coupon_fee"`
	CouponCount    string   `xml:"coupon_count"`
	TransactionID  string   `xml:"transaction_id"`
	OutTradeNo     string   `xml:"out_trade_no"`
	Attach         string   `xml:"attach"`
	TimeEnd        string   `xml:"time_end"`
}
