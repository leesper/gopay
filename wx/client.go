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

	data, err := c.doHTTPRequest(req.URI(), xmlStr)
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
		return nil, fmt.Errorf("signature failed, expected %s, got %s", rspSign, rspMap["sign"])
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

	data, err := c.doHTTPRequest(req.URI(), xmlStr)
	if err != nil {
		return nil, err
	}

	rsp := &QueryOrderRsp{}
	if err = xml.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	rspMap, err := toMap(rsp)
	if err != nil {
		return nil, err
	}

	rspSign := signature(rspMap, c.config.AppKey)
	if rspSign != rspMap["sign"] {
		return nil, fmt.Errorf("signature failed, expected %s, got %s", rspSign, rspMap["sign"])
	}

	if rsp.ReturnCode != Success {
		return nil, fmt.Errorf("return code %s, return msg %s", rsp.ReturnCode, rsp.ReturnMsg)
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
	return result, nil
}

func (c *Client) doHTTPRequest(uri string, xmlStr string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader([]byte(xmlStr)))
	if err != nil {
		return nil, err
	}
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
