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

// AsyncNotificationResult is the result return from Alipay.
type AsyncNotificationResult struct {
	NotifyTime        string `json:"notify_time"`         // 通知时间
	NotifyType        string `json:"notify_type"`         // 通知类型
	NotifyID          string `json:"notify_id"`           // 通知校验ID
	AppID             string `json:"app_id"`              // 开发者的app_id
	Charset           string `json:"charset"`             // 编码格式
	Version           string `json:"version"`             // 接口版本
	SignType          string `json:"sign_type"`           // 签名类型
	Sign              string `json:"sign"`                // 签名
	TradeNo           string `json:"trade_no"`            // 支付宝交易号
	OutTradeNo        string `json:"out_trade_no"`        // 商户订单号
	OutBizNo          string `json:"out_biz_no"`          // 商户业务号
	BuyerID           string `json:"buyer_id"`            // 买家支付宝用户号
	BuyerLogonID      string `json:"buyer_logon_id"`      // 买家支付宝账号
	SellerID          string `json:"seller_id"`           // 卖家支付宝用户号
	SellerEmail       string `json:"seller_email"`        // 卖家支付宝账号
	TradeStatus       string `json:"trade_status"`        // 交易状态
	TotalAmount       string `json:"total_amount"`        // 订单金额
	ReceiptAmount     string `json:"receipt_amount"`      // 实收金额
	InvoiceAmount     string `json:"invoice_amount"`      // 开票金额
	BuyerPayAmount    string `json:"buyer_pay_amount"`    // 付款金额
	PointAmount       string `json:"point_amount"`        // 集分宝金额
	RefundFee         string `json:"refund_fee"`          // 总退款金额
	Subject           string `json:"subject"`             // 总退款金额
	Body              string `json:"body"`                // 商品描述
	GmtCreate         string `json:"gmt_create"`          // 交易创建时间
	GmtPayment        string `json:"gmt_payment"`         // 交易付款时间
	GmtRefund         string `json:"gmt_refund"`          // 交易退款时间
	GmtClose          string `json:"gmt_close"`           // 交易结束时间
	FundBillList      string `json:"fund_bill_list"`      // 支付金额信息
	PassbackParams    string `json:"passback_params"`     // 回传参数
	VoucherDetailList string `json:"voucher_detail_list"` // 优惠券信息
}

// AsyncNotification retrieves the asynchronous notification from Weixin.
func (c *Client) AsyncNotification(req *http.Request) (*AsyncNotificationResult, error) {
	if req == nil {
		return nil, errors.New("http request nil")
	}
	req.ParseForm()

	result := &AsyncNotificationResult{}
	result.NotifyTime = req.PostFormValue("notify_time")
	result.NotifyType = req.PostFormValue("notify_type")
	result.NotifyID = req.PostFormValue("notify_id")
	result.AppID = req.PostFormValue("app_id")
	result.Charset = req.PostFormValue("charset")
	result.Version = req.PostFormValue("version")
	result.SignType = req.PostFormValue("sign_type")
	result.Sign = req.PostFormValue("sign")
	result.TradeNo = req.PostFormValue("trade_no")
	result.OutTradeNo = req.PostFormValue("out_trade_no")
	result.OutBizNo = req.PostFormValue("out_biz_no")
	result.BuyerID = req.PostFormValue("buyer_id")
	result.BuyerLogonID = req.PostFormValue("buyer_logon_id")
	result.SellerID = req.PostFormValue("seller_id")
	result.SellerEmail = req.PostFormValue("seller_email")
	result.TradeStatus = req.PostFormValue("trade_status")
	result.TotalAmount = req.PostFormValue("total_amount")
	result.ReceiptAmount = req.PostFormValue("receipt_amount")
	result.InvoiceAmount = req.PostFormValue("invoice_amount")
	result.BuyerPayAmount = req.PostFormValue("buyer_pay_amount")
	result.PointAmount = req.PostFormValue("point_amount")
	result.RefundFee = req.PostFormValue("refund_fee")
	result.Subject = req.PostFormValue("subject")
	result.Body = req.PostFormValue("body")
	result.GmtCreate = req.PostFormValue("gmt_create")
	result.GmtPayment = req.PostFormValue("gmt_payment")
	result.GmtRefund = req.PostFormValue("gmt_refund")
	result.GmtClose = req.PostFormValue("gmt_close")
	result.FundBillList = req.PostFormValue("fund_bill_list")
	result.PassbackParams = req.PostFormValue("passback_params")
	result.VoucherDetailList = req.PostFormValue("voucher_detail_list")

	if result.NotifyID == "" {
		return nil, errors.New("invalid notify ID")
	}

	fmt.Println("ASYNC RESULT", result, result.Sign)

	ok := verify(req.PostForm, c.config.AliPublicKey, c.config.SignType)
	if ok {
		return result, nil
	}
	return nil, errors.New("verify signature failed")
}
