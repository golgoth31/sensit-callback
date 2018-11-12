# sensit-callback

A simple receiver for sensit v2 callbacks from Sigfox.

## usage

To avoid opening any network port, this project use an AWS SQS Queue. This queue is connected to the sigfox backend through an API gateway and a lambda (lambda given).
Once API gateway, lambda and SQS are defined, go to the Sigfox backend to devlare the associated callback with the following body format:
{
  "device": "{device}",
  "time": "{time}",
  "data": "{data}",
  "messageId": "{seqNumber}",
  "ack": "{ack}"
}

Then create apropriate AWS token to access SQS queue from sensit-callback, place it in ~/.aws/ and create a config.yaml file by coping config.yaml.sample; you will be able to make sensit-callback run with:
go run sensit-callback.go