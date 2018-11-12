import boto3  
import json

def sensit(event, context):
    sqs = boto3.resource('sqs')
    queue = sqs.get_queue_by_name(QueueName='sensit.fifo')
    response = queue.send_message(MessageBody=json.dumps(event), MessageGroupId='sensit', MessageDeduplicationId=event['messageId'])
