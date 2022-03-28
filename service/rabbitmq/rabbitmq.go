package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

/**
* @creator: xuwuruoshui
* @date: 2022-03-19 16:51:34
* @content: rabbitmq订阅模式
 */

var Conn *amqp.Connection

func initMQ() {
	conn, err := amqp.Dial("amqp://root:root@192.168.0.110:5672/")
	if err != nil {
		fmt.Println("连接失败: ", err)
		return
	}
	Conn = conn
}

// 发布消息
func Publish(exchange, queueName, body string) (err error) {
	initMQ()
	defer Conn.Close()
	
	// 1.创建Channel
	channel, err := Conn.Channel()
	if err != nil {
		fmt.Println("创建channel失败: ", err)
		return
	}
	defer channel.Close()

	// 2.创建队列
	// durable持久化
	que, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		fmt.Println("创建queue失败: ", err)
		return
	}

	// 3.发送消息
	// DeliveryMode: amqp.Persistent 持久化
	channel.Publish(exchange, que.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	return
}

type CallBack func(msg string) error

// 接受消息
func Consumer(exchange, queueName string, callback CallBack) (err error) {
	initMQ()
	defer Conn.Close()
	
	// 1.创建Channel
	channel, err := Conn.Channel()
	if err != nil {
		fmt.Println("创建channel失败: ", err)
		return
	}
	defer channel.Close()

	// 2.创建队列
	que, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		fmt.Println("创建queue失败: ", err)
		return
	}

	// autoAck: 自动/手动ack, 手动调用ack,待业务执行完后ack
	msgs, err := channel.Consume(que.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println("消费失败: ", err)
		return
	}
	fmt.Println("等待消息来临............")
	// 读取消息
	for {
		select {
		case d := <-msgs:
			err := callback(string(d.Body))
			if err!=nil{
				break
			}
			d.Ack(false)
		}
	}
	return
}
