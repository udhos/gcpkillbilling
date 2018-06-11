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

    GCP Project:             my-billing-project
    GCP Pubsub Topic:        budget-alerts
    GCP Pubsub Subscription: billing-queue (subscription type must be 'pull')
    GCP Billing Account:     accountId

Create a budget under project 'my-billing-project' with notification set to topic 'budget-alerts'.

Publish a test message
======================

    killbill-pub my-billing-project budget-alerts accountId

Consume the test message
========================

    killbill my-billing-project billing-queue


-x-
