#!/bin/bash

go get golang.org/x/oauth2/google
go get cloud.google.com/go/pubsub
go get cloud.google.com/go/pubsub/apiv1
go get google.golang.org/api/cloudbilling/v1

build() {
	local i="$1"

	gofmt -s -w "$i"
	go tool fix "$i"
	go tool vet "$i"

	hash golint 2>/dev/null && golint "$i"

	go test "$i"
	go install "$i"
}

build ./killbill
build ./killbill-pub
