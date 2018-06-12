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

    GCP Main Account:        main-account
    GCP Project:               main-project
    GCP Pubsub Topic:            budget-alerts
    GCP Pubsub Subscription:       killbill-queue (subscription type must be 'pull')
    GCP Limited Account:     capped-account
    GCP Budget:                capped-budget
    GCP Limited Project:       capped-project

Create pubsub topic 'budget-alerts' under project 'main-project'.

Create pull-type subscription 'killbill-queue' under topic 'budget-alerts'.

Create a budget 'capped-budget' under account 'capped-account' with notification set to topic 'budget-alerts'.

Create a project 'capped-project' linked to account 'capped-account'.

Publish a test message
======================

Publish a fake notification under topic 'budget-alerts' for account 'capped-account':

    killbill-pub main-project budget-alerts capped-account

Consume the test message
========================

Consume notifications from subscription 'killbill-queue'.

Expected result is: All projects linked to account 'capped-account' should be detached from it.

CAUTION: All **accounts** found in notifications sent to subscription 'killbill-queue' will be detached from all their projects. All those unlinked projects will stop, their services will be interrupted, their data will be lost. You have been warned.

    killbill main-project killbill-queue


-x-

