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

func (request GdaxPlaceOrderRequest) GetRequestBody() map[string]interface{} {
	return map[string]interface{}{
		"size":       request.Size,
		"price":      request.Price,
		"side":       request.Side,
		"product_id": request.ProductId,
	}
}
