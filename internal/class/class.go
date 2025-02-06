package internal

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	ID         int        `json:"user_id"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	Role       UserRole   `json:"role"`
	Status     UserStatus `json:"status"`
	Created_at string     `json:"created_at"`
	Updated_at string     `json:"updated_at"`
}

type Supplier struct {
	ID          int    `json:"supplier_id"`
	Name        string `json:"name"`
	ContactInfo string `json:"contact_info"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
}

type Stock struct {
	ID           int     `json:"stock_id"`
	FruitName    string  `json:"fruit_name"`
	Quantity     int     `json:"quantity"`
	PricePerUnit float64 `json:"price_per_unit"`
	Created_at   string  `json:"created_at"`
	Updated_at   string  `json:"updated_at"`
}

type Transaction struct {
	ID              int     `json:"transaction_id"`
	UUID            string  `json:"uuid"`
	UserID          int     `json:"user_id"`
	Status          int     `json:"status"`
	TransactionDate string  `json:"transaction_date"`
	TotalAmount     float64 `json:"total_amount"`
	Created_at      string  `json:"created_at"`
	Updated_at      string  `json:"updated_at"`
}

type TransactionDetail struct {
	ID                int     `json:"id"`
	TransactionID     int     `json:"transaction_detail_id"`
	StockID           int     `json:"stock_id"`
	Quantity          int     `json:"quantity"`
	PriceAtTimeOfSale float64 `json:"price_at_time_of_sale"`
}

type TransactionDetailed struct {
	ID              int                 `json:"transaction_id"`
	UUID            string              `json:"uuid"`
	Status          int                 `json:"status"`
	UserID          int                 `json:"user_id"`
	TransactionDate string              `json:"transaction_date"`
	TotalAmount     float64             `json:"total_amount"`
	Created_at      string              `json:"created_at"`
	Updated_at      string              `json:"updated_at"`
	Detail          []TransactionDetail `json:"transaction_detail"`
}

type PurchaseStock struct {
	ID         int    `json:"purchase_id"`
	SupplierID int    `json:"supplier_id"`
	StockID    int    `json:"stock_id"`
	Date       string `json:"date"`
	Details    string `json:"details"`
}

type TransactionStatus struct {
	ID         int    `json:"id"`
	Status     string `json:"status"`
	Detail     string `json:"detail"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type UserRole struct {
	ID         int    `json:"id"`
	Role       string `json:"role"`
	Detail     string `json:"detail"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type UserStatus struct {
	ID         int    `json:"id"`
	Status     string `json:"status"`
	Detail     string `json:"detail"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type Karyawan struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
