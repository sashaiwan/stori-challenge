import json
import boto3
import base64
import os

s3 = boto3.client('s3')
lambda_client = boto3.client('lambda')

def lambda_handler(event, context):
    try:
        if 'body' in event:
            body = json.loads(event['body'])
        else:
            body = event
            
        file_content_b64 = body['file_content']
        email = body['email']
    except Exception as e:
        return {
            'statusCode': 400,
            'body': json.dumps({'error': 'Invalid request format'})
        }
    
    file_content = base64.b64decode(file_content_b64)
    
    bucket = os.environ['S3_BUCKET']
    key = f"uploads/{email.replace('@', '_')}_{context.aws_request_id}.csv"
    
    s3.put_object(Bucket=bucket, Key=key, Body=file_content)
    
    payload = {
        'detail': {
            'bucket': bucket,
            'key': key,
            'email': email
        }
    }

    lambda_client.invoke(
        FunctionName=os.environ['EC2_CALLER_FUNCTION_NAME'],
        InvocationType='Event',
        Payload=json.dumps(payload)
    )
    
    return {
        'statusCode': 200,
        'headers': {'Content-Type': 'application/json'},
        'body': json.dumps({
            'message': 'File received and processing started',
            's3_key': key
        })
    }
