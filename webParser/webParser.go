package main

//go get github.com/PuerkitoBio/goquery
//go get github.com/denisenkom/go-mssqldb
//go get github.com/streadway/amqp
//go get github.com/go-redis/redis/
import (
	"fmt"
	"strconv"
	"strings"
	"webParser/rabbitMQ"
	"webParser/sql"
)

func main() {
	//SQL
	var db = sql.SqlOpen()
	defer db.Close()

	//GetSqlContent
	names, prices, trPrices, exchangeType, exchangeValue, err := sql.GetSqlContent(db)
	if err != nil {
		fmt.Println("(sqltest) Error getting content: " + err.Error())
	}
	fmt.Println(strings.Repeat("-", 100))
	// Now read the contents
	for i := range names {
		fmt.Println("Product " + strconv.Itoa(i) + ": " + names[i] + ", ExchangeType: " + exchangeType[i] + ", ExchangeValue: " + strconv.FormatFloat(exchangeValue[i], 'f', 2, 64) + ", Price: " + strconv.FormatFloat(prices[i], 'f', 2, 64)+ ", TrPrice: " + strconv.FormatFloat(trPrices[i], 'f', 2, 64))
	}

	//RabbitMQ Consumer
	rabbitMQ.ConsumeRabbitMQ(db)
}
