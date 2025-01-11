# Notification Service.
This service takes its messages from RabbitMQ and logs them to the console.

## Steps 
1. Do a docker compose build
2. Do a docker compose up
3. Open a browser and go to [RabbitMQ UI](http://localhost:15672/)
4. Credentials: username: `guest` password: `guest`
5. Go to Queues and Streams and select the **notification** queue.
6. Send a message and it would appear on the logs of the docker compose service `cli`

### Message Sample
```
{
    "id": 1,
    "user_id": "qwerty-121",
    "message": "airtime purchase successful",
    "read": false
}
```