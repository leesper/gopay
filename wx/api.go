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

type unifiedOrderRsp struct {
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

type queryOrderReq struct{}

type queryOrderRsp struct{}
