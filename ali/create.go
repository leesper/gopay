package ali

// CreateOrderReq requests alipay.trade.create.
type CreateOrderReq struct{}

// URI returns the uri.
func (req CreateOrderReq) URI() string {
	return "alipay.trade.create"
}

// CreateOrderRsp responses alipay.trade.create.
type CreateOrderRsp struct{}
