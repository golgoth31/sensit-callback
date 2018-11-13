provider "aws" {
    shared_credentials_file = "~/.aws/credentials"
    profile                 = "perso"
    region = "eu-west-3"
}

resource "aws_sqs_queue" "sensit" {
    name                      = "sensit"
    delay_seconds             = 0
    max_message_size          = 2048
    message_retention_seconds = 345600
    receive_wait_time_seconds = 0
    tags {
        Environment = "production"
    }
}

// Currently not possible in eu-west-3
// resource "aws_sqs_queue" "sensit_fifo" {
//     name                      = "sensit.fifo"
//     delay_seconds             = 0
//     max_message_size          = 2048
//     message_retention_seconds = 345600
//     receive_wait_time_seconds = 0
//     fifo_queue                  = true
//     content_based_deduplication = false
//     tags {
//         Environment = "production"
//     }
// }w