package dto

type CreateCartRequest struct {
	ProductID uint `json:"product_id"`
	Qty       uint `json:"qty"`
}
