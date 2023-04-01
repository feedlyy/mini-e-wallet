# HOW TO USE

### Pre-requisites

Before running the server,
you need to install or running several server
on your local machine.

##### You need to run or install:
- zookeper
- kafka (binding with your zookeper)
- postgres

##### Kafka Setup
- Create docker network name kafka
```azure
docker network create kafka
```

- Run zookeeper
```azure
docker run --net=kafka -d --name=zookeeper -e ZOOKEEPER_CLIENT_PORT=2181 confluentinc/cp-zookeeper:latest
```

- Run Kafka
```azure
docker run --net=kafka -d -p 9092:9092 --name=kafka -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092 -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 confluentinc/cp-kafka:latest
```
after finish run zookeeper and kafka, do not
forget to set DNS with env KAFKA_ADVERTISED_LISTENERS
that we set when running kafka.

For linux user, you can edit file with this command
on terminal
```azure
sudo vim /etc/hosts
```

and add to be like this
```azure
127.0.0.1 localhost
127.0.0.1 kafka
```

##### Postgres
You can use postgres from direct install or you could
using docker, if you using docker, get the image:
```azure
docker pull postgres
```

After successfully getting postgres image, run this
block of command on your terminal for create container
and running the postgres image.
```azure
docker run --name my-postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD= -e POSTGRES_DB=julo -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 -d postgres:latest
```
This will run postgres on default port
with user and create db name julo.

You just finish all the setup, and the kafka
will be running on your local machine.

### API
Make sure you installed all pre-requisites to
use the api.

After everythings already setup, you need to
go open cmd folder and consumer folder inside
of the cmd folder. Open both of it on 2
separate tab terminal, and running the app
by syntax:
```azure
go run .
```

The reason why we need to run 2 separate server
because main.go for producer server (or main server)
and main.go inside consumer folder for consumer kafka.

After run the application for the first time,
the tables will auto migrate, the tables:
```azure
- accounts
- tokens
- transactions
- wallets
```


