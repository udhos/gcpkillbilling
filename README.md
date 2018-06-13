# gcpkillbilling

gcpkillbilling removes all projects from a GCP billing account.
CAUTION DANGEROUS.

Install GCP SDK
===============

More information here: https://cloud.google.com/sdk/docs/#linux

    wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-192.0.0-linux-x86_64.tar.gz
    tar xf google-cloud-sdk-192.0.0-linux-x86_64.tar.gz
    gcloud init

Build
=====

    go get github.com/udhos/gcpkillbilling
    cd ~/go/src/github.com/udhos/gcpkillbilling
    ./build.sh

Create GCP pubsub topic and subscription
========================================

Example:

    GCP Main Billing Account:    main-account
    GCP Project:                   main-project
    GCP IAM Service Account:         main-killbill (main-killbill@main-project.iam.gserviceaccount.com)
    GCP Pubsub Topic:                budget-alerts
    GCP Pubsub Subscription:           killbill-queue (subscription type must be 'pull')
    GCP Limited Billing Account: capped-account
    GCP Budget:                    capped-budget
    GCP Limited Project:           capped-project

Create pubsub topic 'budget-alerts' under project 'main-project'.

Create pull-type subscription 'killbill-queue' under topic 'budget-alerts'.

Create a budget 'capped-budget' under account 'capped-account' with notification set to topic 'budget-alerts'.

Create a project 'capped-project' linked to account 'capped-account'.

Create an IAM service account 'main-killbill' to authorize the killbill application. Grant the service account these privileges:
- Project OWNER in the capped-project. The killbill application needs this permission in order to change the billing info for the project.
- Pubsub EDITOR (publish/consume messages). The killbill application needs this permission to publish and consume pubsub messages.
- Billing ADMIN the capped-account. Killbill application needs this permission to list and remove its attached projects.

Save the service account credentials as: $HOME/killbill_credentials.json

Publish a test message
======================

Publish a fake notification under topic 'budget-alerts' for account 'capped-account':

    export GOOGLE_APPLICATION_CREDENTIALS=$HOME/killbill_credentials.json
    killbill-pub main-project budget-alerts capped-account

Consume the test message
========================

Consume notifications from subscription 'killbill-queue'.

Expected result is: All projects linked to account 'capped-account' should be detached from it.

CAUTION: All **accounts** found in notifications sent to subscription 'killbill-queue' will be detached from all their projects. All those unlinked projects will stop, their services will be interrupted, their data will be lost. You can limit damage by granting project OWNER permisson to the IAM service account only on projects you can safely destroy.

    export GOOGLE_APPLICATION_CREDENTIALS=$HOME/killbill_credentials.json
    killbill main-project killbill-queue

gcloud cli recipes
==================

Some gcloud cli recipes.

    # publish message to topic
    gcloud pubsub topics publish projects/PROJECT/topics/TOPIC --attribute="billingAccountId=000000-111111-222222"

    # delete subscription
    gcloud pubsub subscriptions delete projects/PROJECT/subscriptions/SUBSCRIPTION

    # create subscription
    gcloud pubsub subscriptions create --topic=projects/PROJECT/topics/TOPIC SUBSCRIPTION


-x-

