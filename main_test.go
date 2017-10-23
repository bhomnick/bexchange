package main

import (
    "testing"
    "reflect"
    "math"
    "math/rand"
    "time"
    "fmt"
)

func TestBehavior(t *testing.T) {
    actions := make(chan *Action)
    done := make(chan bool)
    ob := NewOrderBook(actions)

    log := make([]*Action, 0)
    go func() {
        for {
            action := <-actions
            log = append(log, action)
            if action.actionType == AT_DONE {
                done <- true
                return
            }
        }
    }()

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

    expected := []*Action{
        &Action{AT_SELL,1,0,50,50},
        &Action{AT_SELL,2,0,25,45},
        &Action{AT_SELL,3,0,25,45},
        &Action{AT_BUY,4,0,75,55},
        &Action{AT_PARTIAL_FILLED,4,2,25,45},
        &Action{AT_PARTIAL_FILLED,4,3,25,45},
        &Action{AT_FILLED,4,1,25,50},
        &Action{AT_CANCEL,1,0,0,0},
        &Action{AT_CANCELLED,1,0,0,0},
        &Action{AT_BUY,5,0,20,55},
        &Action{AT_BUY,6,0,15,50},
        &Action{AT_SELL,7,0,25,45},
        &Action{AT_PARTIAL_FILLED,7,5,20,55},
        &Action{AT_FILLED,7,6,5,50},
        &Action{AT_DONE,0,0,0,0},
    }
    if !reflect.DeepEqual(log, expected) {
        t.Error("\n\nExpected:\n\n", expected, "\n\nGot:\n\n", log, "\n\n")
    }
}

func buildOrders(n int, priceMean, priceStd float64, maxAmount int32) []*Order {
    orders := make([]*Order, 0)
    var price uint32
    for i := 0 ; i < n ; i++ {
        price = uint32(math.Abs(rand.NormFloat64()*priceStd+priceMean))
        orders = append(orders, &Order{
            id: uint64(i)+1,
            isBuy: float64(price) >= priceMean,
            price: price,
            amount: uint32(rand.Int31n(maxAmount)),
        })
    }
    return orders
}

func doPerfTest(n int, priceMean, priceStd float64, maxAmount int32) {
    orders := buildOrders(n, priceMean, priceStd, maxAmount)
    actions := make(chan *Action)
    done := make(chan bool)
    ob := NewOrderBook(actions)
    actionCount := 0

    go func() {
        for {
            action := <-actions
            actionCount++
            if action.actionType == AT_DONE {
                done <- true
                return
            }
        }
    }()

    start := time.Now()
    for _, order := range orders {
        ob.AddOrder(order)
    }
    ob.Done()
    <-done
    elapsed := time.Since(start)

    fmt.Printf("Handled %v actions in %v at %v actions/second.\n",
        actionCount, elapsed, int(float64(actionCount)/elapsed.Seconds()))
}

func TestPerf(t *testing.T) {
    doPerfTest(10000, 5000, 10, 50)
    doPerfTest(10000, 5000, 1000, 5000)
    doPerfTest(100000, 5000, 10, 50)
    doPerfTest(100000, 5000, 1000, 5000)
    doPerfTest(1000000, 5000, 10, 50)
    doPerfTest(1000000, 5000, 1000, 5000)
}
