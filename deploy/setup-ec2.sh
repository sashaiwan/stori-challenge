#!/bin/bash

echo "Starting Stori Challenge deploy..."

sudo yum update -y

echo "Installing Dependencies"

echo "Installing Docker..."
sudo yum install -y docker
sudo service docker start
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user

echo "Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

if ! command -v aws &> /dev/null; then
    echo "Installing AWS CLI..."
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    sudo ./aws/install
    rm -rf awscliv2.zip aws/
fi

cd /home/ec2-user
APP_DIR=$(pwd)

if [ ! -f "docker-compose.yml" ]; then
    echo "Error: docker-compose.yml not found!"
    exit 1
fi


echo "Retrieving configuration from AWS Secrets Manager..."
SECRET_NAME="stori-challenge" 
REGION="us-east-2" 

aws secretsmanager get-secret-value \
    --secret-id $SECRET_NAME \
    --region $REGION \
    --query SecretString \
    --output text > .env

if [ $? -ne 0 ]; then
    echo "Failed to retrieve secret from Secrets Manager"
    echo "Make sure the EC2 instance has needed permissions"
    exit 1
fi

echo ".env file created from Secrets Manager"

echo "Starting services..."
sudo /usr/local/bin/docker-compose pull
sudo /usr/local/bin/docker-compose up -d

echo "Waiting for services to start..."
sleep 10
sudo /usr/local/bin/docker-compose ps

echo "Deployment complete!"
