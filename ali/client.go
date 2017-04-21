package ali

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// misc constants
const (
	RSA  = "RSA"
	RSA2 = "RSA2"
)

// PayParam is the interface of all AliPay APIs.
type PayParam interface {
	URI() string
	ExtraParams() map[string]string
	BizContent() string
}

// Config contains all configuration info.
type Config struct {
	APIGateway    string
	AppID         string
	NotifyURL     string
	SandBox       bool
	SignType      string
	AliPublicKey  []byte
	AppPublicKey  []byte
	AppPrivateKey []byte
}

// Client handles all transactions.
type Client struct {
	config    Config
	tlsClient http.Client
}

// NewClient returns a *Client for Alipay.
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

// CreateTrade creates order from Alipay.
func (c *Client) CreateTrade(p CreateTradeParam) (*CreateTradeRsp, error) {
	p.NotifyURL = c.config.NotifyURL
	data, err := c.doHTTPRequest(p)
	if err != nil {
		return nil, err
	}
	rsp := &CreateTradeRsp{}
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	if rsp.TradeCreateResponse.Code != "10000" {
		return nil, fmt.Errorf("code %s msg %s err %s err msg %s",
			rsp.TradeCreateResponse.Code, rsp.TradeCreateResponse.Msg,
			rsp.TradeCreateResponse.SubCode, rsp.TradeCreateResponse.SubMsg)
	}

	responseStr := marshalJSON(rsp.TradeCreateResponse)
	var ok bool
	if c.config.SignType == RSA {
		ok = verifyPKCS1v15([]byte(responseStr), []byte(rsp.Sign), c.config.AliPublicKey, crypto.SHA1)
	} else if c.config.SignType == RSA2 {
		ok = verifyPKCS1v15([]byte(responseStr), []byte(rsp.Sign), c.config.AliPublicKey, crypto.SHA256)
	}

	if !ok {
		return nil, errors.New("verify signature failed")
	}

	return rsp, nil
}

// QueryTrade queries order from Alipay.
func (c *Client) QueryTrade(p PayParam) (*QueryTradeRsp, error) {
	data, err := c.doHTTPRequest(p)
	if err != nil {
		return nil, err
	}
	rsp := &QueryTradeRsp{}
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	if rsp.TradeQueryResponse.Code != "10000" {
		return nil, fmt.Errorf("code %s msg %s err %s err msg %s",
			rsp.TradeQueryResponse.Code, rsp.TradeQueryResponse.Msg,
			rsp.TradeQueryResponse.SubCode, rsp.TradeQueryResponse.SubMsg)
	}

	responseStr := marshalJSON(rsp.TradeQueryResponse)
	var ok bool
	if c.config.SignType == RSA {
		ok = verifyPKCS1v15([]byte(responseStr), []byte(rsp.Sign), c.config.AliPublicKey, crypto.SHA1)
	} else if c.config.SignType == RSA2 {
		ok = verifyPKCS1v15([]byte(responseStr), []byte(rsp.Sign), c.config.AliPublicKey, crypto.SHA256)
	}

	if !ok {
		return nil, errors.New("verify signature failed")
	}

	return rsp, nil
}

func (c *Client) doHTTPRequest(param PayParam) ([]byte, error) {
	reader := strings.NewReader(urlValues(c, param).Encode())
	req, err := http.NewRequest(http.MethodPost, c.config.APIGateway, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

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

// AsyncNotifyResult is the result return from Alipay.
type AsyncNotifyResult struct {
	NotifyTime       string `json:"notify_time"`
	NotifyType       string `json:"notify_type"`
	NotifyID         string `json:"notify_id"`
	SignType         string `json:"sign_type"`
	Sign             string `json:"sign"`
	OutTradeNo       string `json:"out_trade_no"`
	Subject          string `json:"subject"`
	PaymentType      string `json:"payment_type"`
	TradeNo          string `json:"trade_no"`
	TradeStatus      string `json:"trade_status"`
	GmtCreate        string `json:"gmt_create"`
	GmtPayment       string `json:"gmt_payment"`
	GmtClose         string `json:"gmt_close"`
	SellerEmail      string `json:"seller_email"`
	BuyerEmail       string `json:"buyer_email"`
	SellerID         string `json:"seller_id"`
	BuyerID          string `json:"buyer_id"`
	Price            string `json:"price"`
	TotalFee         string `json:"total_fee"`
	Quantity         string `json:"quantity"`
	Body             string `json:"body"`
	Discount         string `json:"discount"`
	IsTotalFeeAdjust string `json:"is_total_fee_adjust"`
	UseCoupon        string `json:"use_coupon"`
	RefundStatus     string `json:"refund_status"`
	GmtRefund        string `json:"gmt_refund"`
}

// AsyncNotification retrieves the asynchronous notification from Weixin.
func (c *Client) AsyncNotification(req *http.Request) (*AsyncNotifyResult, error) {
	if req == nil {
		return nil, errors.New("http request nil")
	}
	req.ParseForm()

	result := &AsyncNotifyResult{}
	result.NotifyTime = req.PostFormValue("notify_time")
	result.NotifyType = req.PostFormValue("notify_type")
	result.NotifyID = req.PostFormValue("notify_id")
	result.SignType = req.PostFormValue("sign_type")
	result.Sign = req.PostFormValue("sign")
	result.OutTradeNo = req.PostFormValue("out_trade_no")
	result.Subject = req.PostFormValue("subject")
	result.PaymentType = req.PostFormValue("payment_type")
	result.TradeNo = req.PostFormValue("trade_no")
	result.TradeStatus = req.PostFormValue("trade_status")
	result.GmtCreate = req.PostFormValue("gmt_create")
	result.GmtPayment = req.PostFormValue("gmt_payment")
	result.GmtClose = req.PostFormValue("gmt_close")
	result.SellerEmail = req.PostFormValue("seller_email")
	result.BuyerEmail = req.PostFormValue("buyer_email")
	result.SellerID = req.PostFormValue("seller_id")
	result.BuyerID = req.PostFormValue("buyer_id")
	result.Price = req.PostFormValue("price")
	result.TotalFee = req.PostFormValue("total_fee")
	result.Quantity = req.PostFormValue("quantity")
	result.Body = req.PostFormValue("body")
	result.Discount = req.PostFormValue("discount")
	result.IsTotalFeeAdjust = req.PostFormValue("is_total_fee_adjust")
	result.UseCoupon = req.PostFormValue("use_coupon")
	result.RefundStatus = req.PostFormValue("refund_status")
	result.GmtRefund = req.PostFormValue("gmt_refund")

	if result.NotifyID == "" {
		return nil, errors.New("invalid notify ID")
	}

	fmt.Printf("ASYNC RESULT %#v %s\n", result, result.Sign)

	ok := verify(req.PostForm, c.config.AliPublicKey, c.config.SignType)

	if ok {
		return result, nil
	}
	return nil, errors.New("verify signature failed")
}
