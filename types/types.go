package types

type User struct {
	Name string	`json:"name"`
	Age  int	`json:"age"`
}

type App struct{
	Name 	string	 `json:"name"`
	Born 	int 	 `json:"born"`
	Price 	float32  `json:"price"`
}

