package ali

import (
	"fmt"
	"net/url"
	"testing"
)

var (
	AliPayPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnxj/9qwVfgoUh/y2W89L6BkRA
FljhNhgPdyPuBV64bfQNN1PjbCzkIM6qRdKBoLPXmKKMiFYnkd6rAoprih3/PrQE
B/VsW8OoM8fxn67UDYuyBTqA23MML9q1+ilIZwBC2AQ2UBVOrFXfFl75p6/B5Ksi
NG9zpgmLCUYuLkxpLQIDAQAB
-----END PUBLIC KEY-----
`
)

func TestVerify(t *testing.T) {
	values := url.Values{}
	values.Add("invoice_amount", "0.02")
	values.Add("buyer_pay_amount", "0.02")
	values.Add("gmt_payment", "2017-10-18 17:57:41")
	values.Add("out_trade_no", "59e7167cea1ecb56134e789b")
	values.Add("auth_app_id", "2017070707671420")
	values.Add("buyer_id", "2088902709749474")
	values.Add("notify_id", "ffa14de9b3c6036d7fa6b90339d62e1jmm")
	values.Add("receipt_amount", "0.02")
	values.Add("sign_type", "RSA")
	values.Add("notify_time", "2017-10-18 17:57:42")
	values.Add("charset", "utf-8")
	values.Add("notify_type", "trade_status_sync")
	values.Add("trade_status", "TRADE_SUCCESS")
	values.Add("seller_id", "2088721352602446")
	values.Add("version", "1.0")
	values.Add("app_id", "2017070707671420")
	values.Add("total_amount", "0.02")
	values.Add("trade_no", "2017101821001004470216470005")
	values.Add("gmt_create", "2017-10-18 17:57:41")
	values.Add("seller_email", "guizhouquzu@qq.com")
	values.Add("subject", "中天会展城国际会议中心B座【趣猪总部】交租")
	values.Add("sign", "Tx4ugSsM5XE9HXcOQGivGN7g0+RiIKFqf3kSS1YbtxPZMyfenmjnGKpcmCrAgCGj439377QR24/bk4wrSR3sDxGRH6DiMQRtTO5FO/2Hy6mxavxqDgdagnYt6SbN2oK0AIvapWgJz769Q2VJNgWym0FI4db75kKIPOQuqWpYdVI=")
	values.Add("fund_bill_list", `[{"amount":"0.02","fundChannel":"ALIPAYACCOUNT"}]`)
	values.Add("buyer_logon_id", "131****6107")
	values.Add("point_amount", "0.00")
	ok := verify(values, []byte(AliPayPublicKey), values.Get("sign_type"))
	if !ok {
		t.Error("verification error")
	}
	result := newAsyncNotifyResult(values)
	if result.AppID != values.Get("app_id") {
		t.Errorf("returned: %s, expected: %s", result.AppID, values.Get("app_id"))
	}
	fmt.Printf("%#v", result)
}
