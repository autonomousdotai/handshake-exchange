package bean

type GdaxPlaceOrderRequest struct {
	Size      string
	Price     string
	Side      string
	ProductId string
}

type GdaxOrderResponse struct {
	Id            string
	Prize         string
	Size          string
	ProductId     string
	Side          string
	Stp           string
	Type          string
	CreatedAt     string
	FillFees      string
	FilledSize    string
	ExecutedValue string
	Status        string
	Settled       bool
}

type GdaxWithdrawRequest struct {
	Address  string
	Amount   string
	Currency string
}

type GdaxWithdrawResponse struct {
	Id       string
	Amount   string
	Currency string
}

func (request GdaxPlaceOrderRequest) GetRequestBody() map[string]interface{} {
	return map[string]interface{}{
		"size":       request.Size,
		"price":      request.Price,
		"side":       request.Side,
		"product_id": request.ProductId,
	}
}

func (request GdaxWithdrawRequest) GetRequestBody() map[string]interface{} {
	return map[string]interface{}{
		"crypto_address": request.Address,
		"amount":         request.Amount,
		"currency":       request.Currency,
	}
}
