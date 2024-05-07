import json

import pika


class RabbitMQConnection:
    def __init__(
        self,
        host="127.0.0.1",
        port=5672,
        virtual_host="/",
        username="guest",
        password="guest",
    ):
        self.credentials = pika.PlainCredentials(username, password)
        self.parameters = pika.ConnectionParameters(
            host, port, virtual_host, self.credentials, heartbeat=0
        )

    def open_connection(self):
        connection = pika.BlockingConnection(self.parameters)
        return connection

    @staticmethod
    def close_connection(connection):
        if connection and connection.is_open:
            connection.close()

    @staticmethod
    def send_message_to_queue(
        connection: pika.BlockingConnection, services: [str], conf: dict, user_name: str
    ):
        channel = connection.channel()
        message = json.dumps(
            {"services": services, "conf": conf, "user_name": user_name}
        )
        channel.queue_declare(queue="train_queue")
        channel.basic_publish(exchange="", routing_key="train_queue", body=message)

    def consume_messages_from_queue(self, callback):
        connection = self.open_connection()
        channel = connection.channel()
        try:
            channel.queue_declare(queue="train_queue")
            channel.basic_consume(
                queue="train_queue", on_message_callback=callback, auto_ack=True
            )
            channel.start_consuming()
        finally:
            self.close_connection(connection)
