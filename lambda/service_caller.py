import json
import urllib3
import os

http = urllib3.PoolManager()

def lambda_handler(event, context):
    detail = event['detail']
    
    response = http.request(
        'POST',
        os.environ['EC2_ENDPOINT'],
        body=json.dumps(detail),
        headers={'Content-Type': 'application/json'}
    )
    
    return {
        'statusCode': response.status,
        'body': response.data.decode('utf-8')
    }
