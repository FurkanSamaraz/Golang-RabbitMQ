package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

func main() {
	// Create a new RabbitMQ connection.
	connRabbitMQ, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	// Create a new Fiber instance.
	app := fiber.New()

	// Add middleware.
	app.Use(
		logger.New(), // add simple logger
	)

	// Add route.
	app.Get("/send", func(c *fiber.Ctx) error {
		// Checking, if query is empty.
		if c.Query("msg") == "" {
			log.Println("Missing 'msg' query parameter")
		}

		// Let's start by opening a channel to our RabbitMQ instance
		// over the connection we have already established
		ch, err := connRabbitMQ.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()

		// With this channel open, we can then start to interact.
		// With the instance and declare Queues that we can publish and subscribe to.
		kuyruk, err := ch.QueueDeclare(
			"kuyruk1", // kuyruğumuzun ismi
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,
		)
		// Handle any errors if we were unable to create the queue.
		if err != nil {
			return err
		}

		// Attempt to publish a message to the queue.
		err = ch.Publish(
			"",          // exchange
			kuyruk.Name, // Gönderilecek kuyruk ismi. Bu şekilde önceki oluşturduğumuz kuyruğun ismini alabiliriz
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType: "text/plain", //mesajımızın tipi
				Body:        []byte("mesajımız"),
			},
		)
		if err != nil {
			return err
		}

		return nil
	})

	// Start Fiber API server.
	log.Fatal(app.Listen(":3000"))
}
