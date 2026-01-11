package domain

type ProductSize string

const (
	SizeSmall  ProductSize = "Small"
	SizeMedium ProductSize = "Medium"
	SizeLarge  ProductSize = "Large"
)

type Product struct {
	ID                int64       `json:"id"`
	Name              string      `json:"name"`
	Type              string      `json:"type"`
	Flavor            string      `json:"flavor"`
	Size              ProductSize `json:"size"`
	Price             float64     `json:"price"`
	Quantity          int         `json:"quantity"`
	ManufacturingDate string      `json:"manufacturing_date"` // YYYY-MM-DD
}
