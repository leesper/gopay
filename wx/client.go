package wx

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// constants for response.
const (
	Success = "SUCCESS"
)

// Config contains all configuration info.
type Config struct {
	AppID     string
	AppKey    string
	MchID     string
	NotifyURL string
	TradeType string
	SandBox   bool
}

// Client handles all transactions.
type Client struct {
	config    Config
	tlsClient http.Client
}

// NewClient returns a *Client ready to use.
func NewClient(cfg Config) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := http.Client{Transport: tr}
	return &Client{
		config:    cfg,
		tlsClient: client,
	}
}

// UnifiedOrder creates new order from Weixin.
func (c *Client) UnifiedOrder(totalFee int, desc, orderID, clientIP string) (*UnifiedOrderRsp, error) {
	req := unifiedOrderReq{
		AppID:          c.config.AppID,
		MchID:          c.config.MchID,
		NonceStr:       generateNonceStr(),
		Body:           desc,
		Attach:         "optional",
		OutTradeNo:     orderID,
		TotalFee:       fmt.Sprintf("%d", totalFee),
		SpbillCreateIP: clientIP,
		NotifyURL:      c.config.NotifyURL,
		TradeType:      c.config.TradeType,
	}

	reqMap, err := toMap(req)
	if err != nil {
		return nil, err
	}

	reqMap["sign"] = signature(reqMap, c.config.AppKey)
	xmlStr := toXMLStr(reqMap)

	uri := req.URI()
	if c.config.SandBox {
		uri = req.SandBoxURI()
	}

	data, err := c.doHTTPRequest(uri, xmlStr)
	if err != nil {
		return nil, err
	}

	rsp := &UnifiedOrderRsp{}
	if err = xml.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	rspMap, err := toMap(rsp)
	if err != nil {
		return nil, err
	}

	rspSign := signature(rspMap, c.config.AppKey)
	if rspSign != rspMap["sign"] {
		return nil, fmt.Errorf("signature failed, expected %s, got %s, response %#v, rsp map %v", rspSign, rspMap["sign"], rsp, rspMap)
	}

	if rsp.ReturnCode != Success {
		return nil, fmt.Errorf("return code %s, return msg %s", rsp.ReturnCode, rsp.ReturnMsg)
	}

	if rsp.ResultCode != Success {
		return nil, fmt.Errorf("err code %s, err code desc %s", rsp.ErrCode, rsp.ErrCodeDesc)
	}

	return rsp, nil
}

// ToPayment returns Payment from prePayID.
func (c *Client) ToPayment(prePayID string) Payment {
	nonceStr := generateNonceStr()
	timestampStr := generateTimestampStr()
	params := map[string]string{
		"appid":     c.config.AppID,
		"partnerid": c.config.MchID,
		"prepayid":  prePayID,
		"noncestr":  nonceStr,
		"timestamp": timestampStr,
		"package":   "Sign=WXPay",
	}

	return Payment{
		AppID:     c.config.AppID,
		PartnerID: c.config.MchID,
		PrepayID:  prePayID,
		NonceStr:  nonceStr,
		Timestamp: timestampStr,
		Package:   "Sign=WXPay",
		Sign:      signature(params, c.config.AppKey),
	}
}

// QueryOrder queries order info from Weixin.
func (c *Client) QueryOrder(transID string) (*QueryOrderRsp, error) {
	req := queryOrderReq{
		AppID:         c.config.AppID,
		MchID:         c.config.MchID,
		TransactionID: transID,
		NonceStr:      generateNonceStr(),
	}

	reqMap, err := toMap(req)
	if err != nil {
		return nil, err
	}

	reqMap["sign"] = signature(reqMap, c.config.AppKey)
	xmlStr := toXMLStr(reqMap)

	uri := req.URI()
	if c.config.SandBox {
		uri = req.SandBoxURI()
	}

	data, err := c.doHTTPRequest(uri, xmlStr)
	if err != nil {
		return nil, err
	}

	rsp := &QueryOrderRsp{}
	if err = xml.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	if rsp.ReturnCode != Success {
		return nil, fmt.Errorf("return code %s, return msg %s", rsp.ReturnCode, rsp.ReturnMsg)
	}

	rspMap, err := toMap(rsp)
	if err != nil {
		return nil, err
	}

	rspSign := signature(rspMap, c.config.AppKey)
	if rspSign != rspMap["sign"] {
		return nil, fmt.Errorf("signature failed, expected %s, got %s", rspSign, rspMap["sign"])
	}

	if rsp.ResultCode != Success {
		return nil, fmt.Errorf("err code %s, err code desc %s", rsp.ErrCode, rsp.ErrCodeDesc)
	}

	return rsp, nil
}

// AsyncNotification retrieves the asynchronous notification from Weixin.
func (c *Client) AsyncNotification(req *http.Request) (*AsyncNotificationResult, error) {
	defer req.Body.Close()
	result := &AsyncNotificationResult{}
	if err := xml.NewDecoder(req.Body).Decode(result); err != nil {
		return nil, err
	}

	if result.ReturnCode != Success {
		return nil, fmt.Errorf("return code %s, return msg %s", result.ReturnCode, result.ReturnMsg)
	}

	rspMap, err := toMap(result)
	if err != nil {
		return nil, err
	}

	rspSign := signature(rspMap, c.config.AppKey)
	if rspSign != rspMap["sign"] {
		return nil, fmt.Errorf("signature failed, expected %s, got %s", rspSign, rspMap["sign"])
	}

	if result.ResultCode != Success {
		return nil, fmt.Errorf("err code %s, err code desc %s", result.ErrCode, result.ErrCodeDesc)
	}

	return result, nil
}

// AnswerAsyncNotify returns a xml in string answering Weixin asynchronous notification.
func (c *Client) AnswerAsyncNotify(returnCode, returnMsg string) string {
	retMap := map[string]string{
		"return_code": returnCode,
		"return_msg":  returnMsg,
	}
	return toXMLStr(retMap)
}

// GetSandBoxSignKey gets sandox sign key from Weixin.
func (c *Client) GetSandBoxSignKey() (*GetSandBoxSignKeyRsp, error) {
	req := getSandBoxSignKeyReq{
		MchID:    c.config.MchID,
		NonceStr: generateNonceStr(),
	}

	reqMap, err := toMap(req)
	if err != nil {
		return nil, err
	}

	reqMap["sign"] = signature(reqMap, c.config.AppKey)
	xmlStr := toXMLStr(reqMap)

	data, err := c.doHTTPRequest(req.SandBoxURI(), xmlStr)
	if err != nil {
		return nil, err
	}

	rsp := &GetSandBoxSignKeyRsp{}
	if err = xml.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	if rsp.ReturnCode != Success {
		return nil, fmt.Errorf("return code %s, return msg %s", rsp.ReturnCode, rsp.ReturnMsg)
	}

	return rsp, nil
}

// AsyncNotificationResult is the result return from Weixin.
type AsyncNotificationResult struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	AppID          string `xml:"appid"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	ResultCode     string `xml:"result_code"`
	ErrCode        string `xml:"err_code"`
	ErrCodeDesc    string `xml:"err_code_des"`
	DeviceInfo     string `xml:"device_info"`
	OpenID         string `xml:"open_id"`
	IsSubscribe    string `xml:"is_subscribe"`
	TradeType      string `xml:"trade_type"`
	BankType       string `xml:"bank_type"`
	TotalFee       string `xml:"total_fee"`
	FeeType        string `xml:"fee_type"`
	CashFee        string `xml:"cash_fee"`
	CashFeeType    string `xml:"cash_fee_type"`
	CouponFee      string `xml:"coupon_fee"`
	CouponCount    string `xml:"coupon_count"`
	TransactionID  string `xml:"transaction_id"`
	OutTradeNo     string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	TimeEnd        string `xml:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc"`
}

func (c *Client) doHTTPRequest(uri string, xmlStr string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader([]byte(xmlStr)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	rsp, err := c.tlsClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
