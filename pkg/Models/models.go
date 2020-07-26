package Models

type Order struct {
	Id          int
	Discount    float64
	Amount      float64
	PaymentMode string
	Rating      int
	Duration    int
	Cuisine     string
	Time        int
	CustId      int
	CustName    string
	RestId      int
	RestName    string
	State       string
}

type Customer struct {
	Id    int
	Name  string
	State string
}

type Restaurant struct {
	Id    int
	Name  string
	State string
}
