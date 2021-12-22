package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/majst01/ssm"
)

const (
	CreatingOrder     ssm.StateType = "CreatingOrder"
	OrderFailed       ssm.StateType = "OrderFailed"
	OrderPlaced       ssm.StateType = "OrderPlaced"
	ChargingCard      ssm.StateType = "ChargingCard"
	TransactionFailed ssm.StateType = "TransactionFailed"
	OrderShipped      ssm.StateType = "OrderShipped"

	CreateOrder     ssm.EventType = "CreateOrder"
	FailOrder       ssm.EventType = "FailOrder"
	PlaceOrder      ssm.EventType = "PlaceOrder"
	ChargeCard      ssm.EventType = "ChargeCard"
	FailTransaction ssm.EventType = "FailTransaction"
	ShipOrder       ssm.EventType = "ShipOrder"
)

type OrderCreationContext struct {
	items []string
	err   error
}

func (c *OrderCreationContext) String() string {
	return fmt.Sprintf("OrderCreationContext [ items: %s, err: %v ]",
		strings.Join(c.items, ","), c.err)
}

type OrderShipmentContext struct {
	cardNumber string
	address    string
	err        error
}

func (c *OrderShipmentContext) String() string {
	return fmt.Sprintf("OrderShipmentContext [ cardNumber: %s, address: %s, err: %v ]",
		c.cardNumber, c.address, c.err)
}

type CreatingOrderAction struct{}

func (a *CreatingOrderAction) Execute(eventCtx ssm.EventContext) ssm.EventType {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Validating, order:", order)
	if len(order.items) == 0 {
		order.err = errors.New("insufficient number of items in order")
		return FailOrder
	}
	return PlaceOrder
}

type OrderFailedAction struct{}

func (a *OrderFailedAction) Execute(eventCtx ssm.EventContext) ssm.EventType {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Order failed, err:", order.err)
	return ssm.NoOp
}

type OrderPlacedAction struct{}

func (a *OrderPlacedAction) Execute(eventCtx ssm.EventContext) ssm.EventType {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Order placed, items:", order.items)
	return ssm.NoOp
}

type ChargingCardAction struct{}

func (a *ChargingCardAction) Execute(eventCtx ssm.EventContext) ssm.EventType {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Validating card, shipment:", shipment)
	if shipment.cardNumber == "" {
		shipment.err = errors.New("card number is invalid")
		return FailTransaction
	}
	return ShipOrder
}

type TransactionFailedAction struct{}

func (a *TransactionFailedAction) Execute(eventCtx ssm.EventContext) ssm.EventType {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Transaction failed, err:", shipment.err)
	return ssm.NoOp
}

type OrderShippedAction struct{}

func (a *OrderShippedAction) Execute(eventCtx ssm.EventContext) ssm.EventType {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Order shipped, address:", shipment.address)
	return ssm.NoOp
}

func newOrderFSM() *ssm.StateMachine {
	return &ssm.StateMachine{
		States: ssm.States{
			ssm.Default: ssm.State{
				Events: ssm.Events{
					CreateOrder: CreatingOrder,
				},
			},
			CreatingOrder: ssm.State{
				Action: &CreatingOrderAction{},
				Events: ssm.Events{
					FailOrder:  OrderFailed,
					PlaceOrder: OrderPlaced,
				},
			},
			OrderFailed: ssm.State{
				Action: &OrderFailedAction{},
				Events: ssm.Events{
					CreateOrder: CreatingOrder,
				},
			},
			OrderPlaced: ssm.State{
				Action: &OrderPlacedAction{},
				Events: ssm.Events{
					ChargeCard: ChargingCard,
				},
			},
			ChargingCard: ssm.State{
				Action: &ChargingCardAction{},
				Events: ssm.Events{
					FailTransaction: TransactionFailed,
					ShipOrder:       OrderShipped,
				},
			},
			TransactionFailed: ssm.State{
				Action: &TransactionFailedAction{},
				Events: ssm.Events{
					ChargeCard: ChargingCard,
				},
			},
			OrderShipped: ssm.State{
				Action: &OrderShippedAction{},
			},
		},
	}
}
