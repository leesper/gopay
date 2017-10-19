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
	"net/url"
	"reflect"
	"strings"

	"github.com/leesper/holmes"
)

// misc constants
const (
	RSA     = "RSA"
	RSA2    = "RSA2"
	Success = "10000"
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

	if rsp.TradeCreateResponse.Code != Success {
		return nil, fmt.Errorf("code %s msg %s err %s err msg %s",
			rsp.TradeCreateResponse.Code, rsp.TradeCreateResponse.Msg,
			rsp.TradeCreateResponse.SubCode, rsp.TradeCreateResponse.SubMsg)
	}

	responseStr := marshalJSON(rsp.TradeCreateResponse)
	var ok bool
	if c.config.SignType == RSA {
		ok = verifyPKCS1v15([]byte(responseStr), []byte(rsp.Sign), c.config.AppPublicKey, crypto.SHA1)
	} else if c.config.SignType == RSA2 {
		ok = verifyPKCS1v15([]byte(responseStr), []byte(rsp.Sign), c.config.AppPublicKey, crypto.SHA256)
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

	if rsp.TradeQueryResponse.Code != Success {
		return nil, fmt.Errorf("code %s msg %s", rsp.TradeQueryResponse.Code, rsp.TradeQueryResponse.Msg)
	}

	// values, err := toValues(rsp.TradeQueryResponse)
	// if err != nil {
	// 	return nil, err
	// }
	// values.Add("sign", rsp.Sign)
	// ok := verify(values, c.config.AliPublicKey, c.config.SignType)
	// if !ok {
	// 	return nil, errors.New("verify signature failed")
	// }

	return rsp, nil
}

// RefundTrade refunds trade .
func (c *Client) RefundTrade(p PayParam) (*RefundTradeRsp, error) {
	data, err := c.doHTTPRequest(p)
	if err != nil {
		return nil, err
	}
	rsp := &RefundTradeRsp{}
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	if rsp.TradeRefundResponse.Code != "10000" {
		return nil, fmt.Errorf("code %s msg %s", rsp.TradeRefundResponse.Code, rsp.TradeRefundResponse.Msg)
	}

	return rsp, nil
}

// QueryRefund queries the result of refund.
func (c *Client) QueryRefund(p PayParam) (*QueryRefundRsp, error) {
	data, err := c.doHTTPRequest(p)
	if err != nil {
		return nil, err
	}
	rsp := &QueryRefundRsp{}
	if err = json.NewDecoder(bytes.NewReader(data)).Decode(rsp); err != nil {
		return nil, err
	}

	if rsp.RefundQueryResponse.Code != "10000" {
		return nil, fmt.Errorf("code %s msg %s", rsp.RefundQueryResponse.Code, rsp.RefundQueryResponse.Msg)
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
	AppID          string `json:"app_id"`
	AuthAPPID      string `json:"auth_app_id"`
	BuyerID        string `json:"buyer_id"`
	BuyerLogonID   string `json:"buyer_logon_id"`
	BuyerPayAmount string `json:"buyer_pay_amount"`
	Charset        string `json:"charset"`
	FundBillList   string `json:"fund_bill_list"`
	GmtCreate      string `json:"gmt_create"`
	GmtPayment     string `json:"gmt_payment"`
	InvoiceAmount  string `json:"invoice_amount"`
	NotifyID       string `json:"notify_id"`
	NotifyTime     string `json:"notify_time"`
	NotifyType     string `json:"notify_type"`
	OutTradeNo     string `json:"out_trade_no"`
	PointAmount    string `json:"point_amount"`
	ReceiptAmount  string `json:"receipt_amount"`
	SellerEmail    string `json:"seller_email"`
	SellerID       string `json:"seller_id"`
	Sign           string `json:"sign"`
	SignType       string `json:"sign_type"`
	Subject        string `json:"subject"`
	TotalAmount    string `json:"total_amount"`
	TradeNo        string `json:"trade_no"`
	TradeStatus    string `json:"trade_status"`
	Version        string `json:"version"`
}

func newAsyncNotifyResult(values url.Values) *AsyncNotifyResult {
	result := &AsyncNotifyResult{}
	typ := reflect.TypeOf(result)
	val := reflect.ValueOf(result)
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		tag := sf.Tag.Get("json")
		val.Field(i).SetString(values.Get(tag))
	}
	return result
}

// AsyncNotify retrieves the asynchronous notification from Weixin.
func (c *Client) AsyncNotify(req *http.Request) (*AsyncNotifyResult, error) {
	if req == nil {
		return nil, errors.New("http request nil")
	}

	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	holmes.Debugln("values", values)

	result := newAsyncNotifyResult(values)

	if result.NotifyID == "" {
		return nil, errors.New("invalid notify ID")
	}

	ok := verify(values, c.config.AliPublicKey, values.Get("sign_type"))

	if ok {
		return result, nil
	}
	return nil, errors.New("verify signature failed")
}
