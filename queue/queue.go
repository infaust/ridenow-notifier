package queue

import "github.com/streadway/amqp"

type QueueConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// queue amqp.Queue
}

func NewQueueConsumer(dataSourceName string) (*QueueConsumer, error) {
	conn, err := amqp.Dial(dataSourceName)
	if err != nil {
		return nil, err
	}
	// defer conn.Close() // ?
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// defer ch.Close() // ?
	err = ch.ExchangeDeclare(
		"ridenow_matcher", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}
	return &QueueConsumer{conn, ch}, nil
}

func (q *QueueConsumer) Subscribe(routings ...string) (<-chan amqp.Delivery, error) {
	queue, err := q.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	for _, routing := range routings {
		err = q.channel.QueueBind(
			queue.Name,        // queue name
			routing,           // routing key
			"ridenow_matcher", // exchange
			false,
			nil)
		if err != nil {
			return nil, err
		}
	}

	msgs, err := q.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil

}
