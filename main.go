package main

func main() {
    actions := make(chan *Action)
    done := make(chan bool)
    ob := NewOrderBook(actions)

    go ConsoleActionHandler(actions, done)

    // Should all go into the book
    ob.AddOrder(&Order{isBuy:false,id:1,price:50,amount:50})
    ob.AddOrder(&Order{isBuy:false,id:2,price:45,amount:25})
    ob.AddOrder(&Order{isBuy:false,id:3,price:45,amount:25})
    // Should trigger three fills, two partial at 45 and one at 50
    ob.AddOrder(&Order{isBuy:true,id:4,price:55,amount:75})
    // Should cancel immediately
    ob.CancelOrder(1)
    // Should all go into the book
    ob.AddOrder(&Order{isBuy:true,id:5,price:55,amount:20})
    ob.AddOrder(&Order{isBuy:true,id:6,price:50,amount:15})
    // Should trigger two fills, one partial at 55 and one at 50
    ob.AddOrder(&Order{isBuy:false,id:7,price:45,amount:25})
    ob.Done()

    <-done
}
