package Models

type CuisineCustomer struct {
	Cusine  string
	CustId  int
	OrderId int
}

type CuisineRestaurant struct {
	Cuisine string
	RestId  int
	OrderId int
}

type OrderRestaurant struct {
	OrderId int
	RestId  int
}
