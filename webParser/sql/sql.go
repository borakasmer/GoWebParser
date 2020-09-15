package sql

import (
	"context"
	"database/sql"
	"log"
	"time"
	shared2 "webParser/shared"
)

func SqlOpen() *sql.DB {
	db, err := sql.Open("sqlserver", shared2.Config.SQLURL)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetSqlContent(db *sql.DB) ([]string, []float64, []float64, []string, []float64, error) {
	var (
		Name          []string
		Price         []float64
		TrPrice       []float64
		ExchangeType  []string
		ExchangeValue []float64
		ctx           context.Context
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select Name,Price,TrPrice,ExchangeType,ExchangeValue from [dbo].[Products2]")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var _name string
		var _price float64
		var _trPrice float64
		var _exchangeType string
		var _exchangeValue float64
		err := rows.Scan(&_name, &_price, &_trPrice, &_exchangeType, &_exchangeValue)
		if err != nil {
			return Name, Price, TrPrice, ExchangeType, ExchangeValue, err
		} else {
			Name = append(Name, _name)
			Price = append(Price, _price)
			TrPrice = append(TrPrice, _trPrice)
			ExchangeType = append(ExchangeType, _exchangeType)
			ExchangeValue = append(ExchangeValue, _exchangeValue)
		}
	}
	return Name, Price, TrPrice, ExchangeType, ExchangeValue, nil
}
func InsertSqlContent(db *sql.DB, product *shared2.AddProduct) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO Products2(Name,Price,TrPrice,IsActive,ExchangeType,ExchangeValue) VALUES (@p1, @p2,@p3,@p4,@p5,@p6); select ID = convert(bigint, SCOPE_IDENTITY())")
	if err != nil {
		handleError(err, "Could not insert SqlDB")
		return 0, err
	}
	var ctx context.Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	defer stmt.Close()
	rows := stmt.QueryRowContext(ctx, product.Name, product.Price, product.TrPrice, 1, product.ExchangeName, product.ExchangeValue)
	if rows.Err() != nil {
		return 0, err
	}
	var _id int64
	rows.Scan(&_id)
	return _id, nil
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}
