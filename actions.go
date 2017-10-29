package bexchange

import "fmt"

type ActionType string
const (
    AT_BUY = "BUY"
    AT_SELL = "SELL"
    AT_CANCEL = "CANCEL"
    AT_CANCELLED = "CANCELLED"
    AT_PARTIAL_FILLED = "PARTIAL_FILLED"
    AT_FILLED = "FILLED"
    AT_DONE = "DONE"
)

type Action struct {
    actionType ActionType `json:"actionType"`
    orderId uint64 `json:"orderId"`
    fromOrderId uint64 `json:"fromOrderId"`
    amount uint32 `json:"amount"`
    price uint32 `json:"price"`
}

func (a *Action) String() string {
    return fmt.Sprintf("\nAction{actionType:%v,orderId:%v,fromOrderId:%v,amount:%v,price:%v}",
        a.actionType, a.orderId, a.fromOrderId, a.amount, a.price)
}

func NewBuyAction(o *Order) *Action {
    return &Action{actionType: AT_BUY, orderId: o.id, amount: o.amount,
        price: o.price}
}

func NewSellAction(o *Order) *Action {
    return &Action{actionType: AT_SELL, orderId: o.id, amount: o.amount,
        price: o.price}
}

func NewCancelAction(id uint64) *Action {
    return &Action{actionType: AT_CANCEL, orderId: id}
}

func NewCancelledAction(id uint64) *Action {
    return &Action{actionType: AT_CANCELLED, orderId: id}
}

func NewPartialFilledAction(o *Order, fromOrder *Order) *Action {
    return &Action{actionType: AT_PARTIAL_FILLED, orderId: o.id, fromOrderId: fromOrder.id,
        amount: fromOrder.amount, price: fromOrder.price}
}

func NewFilledAction(o *Order, fromOrder *Order) *Action {
    return &Action{actionType: AT_FILLED, orderId: o.id, fromOrderId: fromOrder.id,
        amount: o.amount, price: fromOrder.price}
}

func NewDoneAction() *Action {
    return &Action{actionType: AT_DONE}
}

func ConsoleActionHandler(actions <-chan *Action, done chan<- bool) {
    for {
        a := <-actions
        switch a.actionType {
        case AT_BUY, AT_SELL:
            fmt.Printf("%s - Order: %v, Amount: %v, Price: %v\n",
                a.actionType, a.orderId, a.amount, a.price)
        case AT_CANCEL, AT_CANCELLED:
            fmt.Printf("%s - Order: %v\n", a.actionType, a.orderId)
        case AT_PARTIAL_FILLED, AT_FILLED:
            fmt.Printf("%s - Order: %v, Filled %v@%v, From: %v\n",
                a.actionType, a.orderId, a.amount, a.price, a.fromOrderId)
        case AT_DONE:
            fmt.Printf("%s\n", a.actionType)
            done <- true
            return
        default:
            panic("Unknown action type.")
        }
    }
}

func NoopActionHandler(actions <-chan *Action) {
    for { <-actions }
}
