package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"sync"
	"wb_l0/domain"
	"wb_l0/service"
)

var NC, _ = nats.Connect(nats.DefaultURL)

func Subscribe(o_service *service.OrderService) {

	NC.Subscribe("order", func(m *nats.Msg) {
		var new_order domain.Order
		err := json.Unmarshal(m.Data, &new_order)
		if err != nil {
			fmt.Printf("ERROR: can't unmarshal json with new order in subscriber: %v\n", err)
			return
		}
		o_service.AddOrder(&new_order)
	})
	defer NC.Close()

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}

