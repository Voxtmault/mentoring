package internal

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
    ID       int  `json:"user_id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Role     string `json:"role"`
}

type Supplier struct {
    ID         int  `json:"supplier_id"`
    Name       string `json:"name"`
    ContactInfo string `json:"contact_info"`
}

type Stock struct {
    ID            int   `json:"stock_id"`
    FruitName     string  `json:"fruit_name"`
    Quantity      int     `json:"quantity"`
    PricePerUnit  float64 `json:"price_per_unit"`
}

type Transaction struct {
    ID            int   `json:"transaction_id"`
    UserID        int   `json:"user_id"`  
    TransactionDate string `json:"transaction_date"`
    TotalAmount    float64 `json:"total_amount"`
}

type TransactionDetail struct {
    ID               int   `json:"id"`
    TransactionID    int   `json:"transaction_detail_id"`  
    StockID          int   `json:"stock_id"`        
    Quantity         int     `json:"quantity"`
    PriceAtTimeOfSale float64 `json:"price_at_time_of_sale"`
}

type PurchaseStock struct {
    ID         int  `json:"purchase_id"`
    SupplierID int  `json:"supplier_id"`
    StockID   int  `json:"stock_id"`
    Date       string `json:"date"`
    Details    string `json:"details"`
}
