package dto

type SellerOrderDetails struct {
	OrderRefNumber  int    `json:"order_ref_number"`
	OrderStatus     string `json:"order_status"`
	CreatedAt       string `json:"created_at"`
	OrderItemId     int    `json:"order_item_id"`
	ProductId       int    `json:"product_id"`
	Name            string `json:"name"`
	ImageUrl        string `json:"image_url"`
	Price           int    `json:"price"`
	Qty             int    `json:"qty"`
	CustomerName    string `json:"customer_name"`
	CustomerEmail   string `json:"customer_email"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerAddress string `json:"customer_address"`
}
