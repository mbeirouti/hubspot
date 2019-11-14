# HubSpot Challenge

This was pretty fun.

The one thing that stopped me from submitting the correct solution on time was that I realized too late that in the case of multiple possible dates, the earliest one should be submitted. Realized that about 20 minutes after the time ran out and then fixed the problem.

## Running the Code

## Locally

Put the extracted hubspot directory into your `GOPATH` (usually something like `~/Go/src`). Navigate to that directory and run the command:

```bash
go run ./main.go
```

You should see the output as a result of a successful `POST` request:
```
You did it.. Woot!
```

## Docker

Running the code in docker is just as simple, navigate into the extracted "Hubspot" directory then build the Dockerfile using the following (may require login):

```bash
sudo docker build .
```

Note the resulting Image ID, then run the container using the following replacing `<IMAGE_ID>` below with the Image ID you noted previously.

```bash
docker run -p 8080:8080 -d <IMAGE_ID>
```

You should see the output as a result of a successful `POST` request:
```
You did it.. Woot!
```
