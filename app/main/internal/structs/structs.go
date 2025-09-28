package structs

type Product struct {
	ID          int
	Name        string
	Price       float64
	Description string
	Category    string
	Stock       int
}
type Categories struct {
	ID   int
	Name string
}
type MainpageData struct {
	Categories []Categories
	Products   []Product
	Title      string
}
