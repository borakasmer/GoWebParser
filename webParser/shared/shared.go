package shared

type Configuration struct {
	AMQPURL string
	SQLURL string
}

var Config = Configuration{
	AMQPURL: "amqp://test:test@localhost:5672/",
	SQLURL: "sqlserver://**userName**:*****@192.168.**.**:1433?database=Deno&connection+timeout=30",
}

type AddProduct struct {
	Price float64
	TrPrice float64
	Name string
	ExchangeType int
	ExchangeName string
	ExchangeValue float64
}

type exchangeType struct {
	Dolar   string
	Euro    string
	Sterlin string
	Altin   string
}

func newExchangeType() *exchangeType {
	return &exchangeType{
		Dolar:   "DOLAR",
		Euro:    "EURO",
		Sterlin: "STERLİN",
		Altin:   "GRAM ALTIN",
	}
}

var Exchanges = newExchangeType()