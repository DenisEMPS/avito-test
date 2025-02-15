package types

type InfoResponse struct {
	Coins        int `db:"coins"`
	Inventory    []Inventory
	CoinsHistory CoinsHistory
}

type CoinsHistory struct {
	Received []Received `db:"received"`
	Sent     []Sent     `db:"sent"`
}

type Inventory struct {
	Type     string `db:"type"`
	Quantity int    `db:"quantity"`
}

type Received struct {
	FromUser string `db:"from_user"`
	Amount   string `db:"amount"`
}

type Sent struct {
	ToUser string `db:"to_user"`
	Amount int    `db:"amount"`
}

type SendCoinRequest struct {
	ToUser string `json:"to_user"`
	Amount int    `json:"amount"`
}

type SentCoinsRequest struct {
	ToUser string `json:"to_user"`
	Amount int    `json:"amount"`
}

type BuyRequest struct {
	Item     string  `json:"item"`
	Quantity float64 `json:"quantity"`
}
