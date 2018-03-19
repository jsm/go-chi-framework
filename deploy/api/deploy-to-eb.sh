#! /bin/bash
 
set -ex

deploy-to-region () {
    AWS_ACCOUNT_ID=0000000000000
    API_ECR_REPO=$AWS_ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/backend-api
    WORKER_ECR_REPO=$AWS_ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com/backend-worker
    EB_BUCKET=$AWS_ACCOUNT_ID-eb-$REGION-versions

    # Deploy API image to Docker Hub
    docker tag api:$SHA1 $API_ECR_REPO:$VERSION
    docker tag -f api:$SHA1 $API_ECR_REPO:latest
    docker push $API_ECR_REPO:$VERSION
    docker push $API_ECR_REPO:latest

    # Deploy Worker image to Docker Hub
    docker tag worker:$SHA1 $WORKER_ECR_REPO:$VERSION
    docker tag -f worker:$SHA1 $WORKER_ECR_REPO:latest
    docker push $WORKER_ECR_REPO:$VERSION
    docker push $WORKER_ECR_REPO:latest

    # Create new Elastic Beanstalk version
    sed -e "s/<TAG>/$VERSION/" -e "s/<REGION>/$REGION/" < eb/Dockerrun.aws.json.template > eb/Dockerrun.aws.json
    (cd eb && zip -r ../$VERSION.zip .)
    aws s3 cp $VERSION.zip s3://$EB_BUCKET/$APPLICATION_NAME/$VERSION.zip
    aws --region $REGION elasticbeanstalk create-application-version --application-name $APPLICATION_NAME \
        --version-label $VERSION --source-bundle S3Bucket=$EB_BUCKET,S3Key=$APPLICATION_NAME/$VERSION.zip

    sleep 10

    # If it's a hotfix, skip the dev deploy
    if [ "$2" != "hotfix" ]; then
        aws --region $REGION elasticbeanstalk update-environment --application-name $APPLICATION_NAME --environment-name $APPLICATION_NAME-dev --version-label $VERSION
    fi
}

APPLICATION_NAME=api
VERSION=$(cat ../../VERSION)
SHA1=$1

# Handle Hotfix
if [ "$2" == "hotfix" ]; then
    VERSION="$VERSION-hotfix-$(git rev-list HEAD ^$VERSION --count)"
fi

# Deploy to us-east-1
REGION=us-east-1 deploy-to-region
