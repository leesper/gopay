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

	if rsp.ReturnCode != Success {
		return nil, fmt.Errorf("return code %s, return msg %s", rsp.ReturnCode, rsp.ReturnMsg)
	}

	if rsp.ResultCode != Success {
		return nil, fmt.Errorf("err code %s, err code desc %s", rsp.ErrCode, rsp.ErrCodeDesc)
	}

	rspMap, err := toMap(rsp)
	if err != nil {
		return nil, err
	}
	rspSign := signature(rspMap, c.config.AppKey)
	if rspSign != rspMap["sign"] {
		return nil, fmt.Errorf("signature failed, expected %s, got %s", rspSign, rspMap["sign"])
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
func (c *Client) QueryOrder(transID string, tradeNo string) (*QueryOrderRsp, error) {
	req := queryOrderReq{
		AppID:         c.config.AppID,
		MchID:         c.config.MchID,
		TransactionID: transID,
		OutTradeNo:    tradeNo,
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

	if rsp.ResultCode != Success {
		return nil, fmt.Errorf("err code %s, err code desc %s", rsp.ErrCode, rsp.ErrCodeDesc)
	}

	rspMap, err := toMap(rsp)
	if err != nil {
		return nil, err
	}

	rspSign := signature(rspMap, c.config.AppKey)
	if rspSign != rspMap["sign"] {
		return nil, fmt.Errorf("signature failed, expected %s, got %s", rspSign, rspMap["sign"])
	}

	return rsp, nil
}

func (c *Client) RefundOrder(transID, tradeNo, refundNo string, totalFee, refundFee int) (*RefundOrderRsp, error) {
	req := refundOrderReq{
		AppID:         c.config.AppID,
		MchID:         c.config.MchID,
		NonceStr:      generateNonceStr(),
		TransactionID: transID,
		OutTradeNo:    tradeNo,
		OutRefundNo:   refundNo,
		TotalFee:      fmt.Sprintf("%d", totalFee),
		RefundFee:     fmt.Sprintf("%d", refundFee),
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

	rsp := &RefundOrderRsp{}
	if err = xml.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	if rsp.ReturnCode != Success {
		return nil, fmt.Errorf("return code %s, return msg %s", rsp.ReturnCode, rsp.ReturnMsg)
	}

	if rsp.ResultCode != Success {
		return nil, fmt.Errorf("err code %s, err code desc %s", rsp.ErrCode, rsp.ErrCodeDesc)
	}

	rspMap, err := toMap(rsp)
	if err != nil {
		return nil, err
	}

	rspSign := signature(rspMap, c.config.AppKey)
	if rspSign != rspMap["sign"] {
		return nil, fmt.Errorf("signature failed, expected %s, got %s", rspSign, rspMap["sign"])
	}

	return rsp, nil
}

func (c *Client) QueryRefund(transID, tradeNo, refundNo, refundID string) (*QueryRefundRsp, error) {
	req := queryRefundReq{
		AppID:         c.config.AppID,
		MchID:         c.config.MchID,
		NonceStr:      generateNonceStr(),
		TransactionID: transID,
		OutTradeNo:    tradeNo,
		OutRefundNo:   refundNo,
		RefundID:      refundID,
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

	rsp := &QueryRefundRsp{}
	if err = xml.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	if rsp.ReturnCode != Success {
		return nil, fmt.Errorf("return code %s, return msg %s", rsp.ReturnCode, rsp.ReturnMsg)
	}

	if rsp.ResultCode != Success {
		return nil, fmt.Errorf("err code %s, err code desc %s", rsp.ErrCode, rsp.ErrCodeDesc)
	}

	rspMap, err := toMap(rsp)
	if err != nil {
		return nil, err
	}

	rspSign := signature(rspMap, c.config.AppKey)
	if rspSign != rspMap["sign"] {
		return nil, fmt.Errorf("signature failed, expected %s, got %s", rspSign, rspMap["sign"])
	}

	return rsp, nil
}

// AsyncNotify retrieves the asynchronous notification from Weixin.
func (c *Client) AsyncNotify(req *http.Request) (*AsyncNotifyResult, error) {
	defer req.Body.Close()
	result := &AsyncNotifyResult{}
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
		return nil, fmt.Errorf("signature failed, expected %s, got %s, result %#v, rspMap %v",
			rspSign, rspMap["sign"], result, rspMap)
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

// AsyncNotifyResult is the result return from Weixin.
type AsyncNotifyResult struct {
	ReturnCode    string `xml:"return_code"`    // 返回状态码
	ReturnMsg     string `xml:"return_msg"`     // 返回信息
	AppID         string `xml:"appid"`          // 应用ID
	MchID         string `xml:"mch_id"`         // 商户号
	DeviceInfo    string `xml:"device_info"`    // 设备号
	NonceStr      string `xml:"nonce_str"`      // 随机字符串
	Sign          string `xml:"sign"`           // 签名
	ResultCode    string `xml:"result_code"`    // 业务结果
	ErrCode       string `xml:"err_code"`       // 错误代码
	ErrCodeDesc   string `xml:"err_code_des"`   // 错误代码描述
	OpenID        string `xml:"openid"`         // 用户标识
	IsSubscribe   string `xml:"is_subscribe"`   // 是否关注公众账号
	TradeType     string `xml:"trade_type"`     // 交易类型
	BankType      string `xml:"bank_type"`      // 付款银行
	TotalFee      string `xml:"total_fee"`      // 总金额
	FeeType       string `xml:"fee_type"`       // 货币种类
	CashFee       string `xml:"cash_fee"`       // 现金支付金额
	CashFeeType   string `xml:"cash_fee_type"`  // 现金支付货币类型
	CouponFee     string `xml:"coupon_fee"`     // 代金券或立减优惠金额
	CouponCount   string `xml:"coupon_count"`   // 代金券或立减优惠使用数量
	TransactionID string `xml:"transaction_id"` // 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	Attach        string `xml:"attach"`         // 商家数据包
	TimeEnd       string `xml:"time_end"`       // 支付完成时间
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
