package main

import (
    "fmt"
    "time"
)

func main() {
    actions := make(chan *Action)

    go ConsoleActionHandler(actions)

    ob := NewOrderBook(actions)

    ob.AddOrder(&Order{isBuy:false,id:1,price:50,amount:50})
    ob.AddOrder(&Order{isBuy:false,id:2,price:45,amount:25})
    ob.AddOrder(&Order{isBuy:false,id:3,price:45,amount:25})

    // Should trigger three fills, two partial at 45 and one at 50
    ob.AddOrder(&Order{isBuy:true,id:4,price:55,amount:75})

    ob.CancelOrder(1)

    // Should go into the book
    ob.AddOrder(&Order{isBuy:true,id:5,price:55,amount:10})

    time.Sleep(time.Second * 1)
    fmt.Println("Done")
}
