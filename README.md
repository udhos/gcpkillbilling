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

    GCP Project:             my-main-project
    GCP Pubsub Topic:          budget-alerts
    GCP Pubsub Subscription:     billing-queue (subscription type must be 'pull')
    GCP Billing Account:     accountId
    GCP Limited Project:       my-capped-project

Create subscription 'billing-queue' under topic 'budget-alerts'.

Create a budget under project 'my-capped-project' with notification set to topic 'budget-alerts'.

Publish a test message
======================

Publish a fake budget notification under topic 'budget-alerts' for account 'accountId':

    killbill-pub my-main-project budget-alerts accountId

Consume the test message
========================

Consume notifications from subscription 'billing-queue'. 

CAUTION: All accounts found in notifications will be detached from all their projects. Those projects will stop, their services will be interrupted, their data will be lost. You have been warned.

    killbill my-main-project billing-queue


-x-
