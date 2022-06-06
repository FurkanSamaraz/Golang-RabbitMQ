package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	//Önce RabbitMQ sunucumuzla bağlantı kuralım.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	//Mesajımızı göndermek için bir kanal oluşturalım.
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer ch.Close()

	//Mesajımızı göndereceğimiz kuyruğu tanımlayalım.
	kuyruk, err := ch.QueueDeclare(
		"kuyruk1", // kuyruğumuzun ismi
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalln(err)
	}

	mesajımız := "baba yaga!"

	//Mesajımızı paylaşmak için yapmamız gerekenler
	err = ch.Publish(
		"",          // exchange
		kuyruk.Name, // Gönderilecek kuyruk ismi. Bu şekilde önceki oluşturduğumuz kuyruğun ismini alabiliriz
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain", //mesajımızın tipi
			Body:        []byte(mesajımız),
		})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Mesajımız başarılı bir şekilde gönderildi!")

}
