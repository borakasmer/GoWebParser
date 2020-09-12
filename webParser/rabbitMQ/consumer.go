package rabbitMQ

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"webParser/parser"
	"webParser/redis"
	shared2 "webParser/shared"
	sql2 "webParser/sql"
)

func ConsumeRabbitMQ(db *sql.DB) {
	conn, err := amqp.Dial(shared2.Config.AMQPURL)
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("product", false, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")
	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			fmt.Println(strings.Repeat("-", 100))
			log.Printf("Received a message: %s", d.Body)

			addProduct := &shared2.AddProduct{}

			err := json.Unmarshal(d.Body, addProduct)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			//Check Redis
			var redisClient = redis.GetRedisClient()

			exchangeRedisVal := &redis.ExchangeRedisValue{}
			err2 := redisClient.GetKey(addProduct.ExchangeName, exchangeRedisVal)
			if err2 != nil && err2.Error() != "redis: nil" {
				log.Fatalf("Error: %v", err2.Error())
			}
			//-------------------------
			var exchange, exchangeFlt string

			if exchangeRedisVal.Name == "" { //If there is no value on Redis, we will parse the web page
				//Web Parser
				//var exchange = parser.ParseWeb(shared2.Exchanges.Dolar)
				exchange = parser.ParseWeb(addProduct.ExchangeName)
				exchangeFlt = strings.Replace(exchange, ",", ".", 1)

				exchangeRedisVal.Name = addProduct.ExchangeName
				exchangeRedisVal.Value = exchangeFlt
				exchangeRedisVal.CreatedDate = time.Now()
				redisClient.SetKey(addProduct.ExchangeName, exchangeRedisVal, time.Minute*1)
			} else { // Get ExchangeData From Redis
				exchange = exchangeRedisVal.Name
				exchangeFlt = exchangeRedisVal.Value
			}
			fmt.Println(strings.Repeat("-", 100))
			//fmt.Printf("Kur :%s - %s\n", shared2.Exchanges.Dolar,exchange)
			fmt.Printf("Kur :%s - %s\n", addProduct.ExchangeName, exchange)

			exchangeValue, err2 := strconv.ParseFloat(exchangeFlt, 64)
			if err2 != nil {
				return
			}
			//--------------------------------------------

			//Convert addProduct price to ₺ ==>TrPrice
			addProduct.TrPrice = addProduct.Price * exchangeValue
			addProduct.ExchangeValue = exchangeValue

			log.Printf("Price %f of %s. ExchangeType : %s", addProduct.Price, addProduct.Name, addProduct.ExchangeName)
			res, err2 := sql2.InsertSqlContent(db, addProduct)
			handleError(err2, "Could not Insert Product to Sql")
			log.Printf("Inserted Product ID : %d", res)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

			// SQL List All Product
			names, prices, trPrices, exchangeType, exchangeVal, err := sql2.GetSqlContent(db)
			if err != nil {
				fmt.Println("(sqltest) Error getting content: " + err.Error())
			}
			fmt.Println(strings.Repeat("-", 100))
			// Now read the contents
			for i := range names {
				fmt.Println("Product " + strconv.Itoa(i) + ": " + names[i] + ", ExchangeType: " + exchangeType[i] + ", ExchangeValue: " + strconv.FormatFloat(exchangeVal[i], 'f', 2, 64) + ", Price: " + strconv.FormatFloat(prices[i], 'f', 2, 64) + ", TrPrice: " + strconv.FormatFloat(trPrices[i], 'f', 2, 64))
			}
		}
	}()

	// Stop for program termination
	<-stopChan
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}
