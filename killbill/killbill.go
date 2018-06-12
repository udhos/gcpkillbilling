package main

import (
	"context"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudbilling/v1"
)

func killbill(billingAccountID string) error {

	log.Printf("killbill: DRY=%v account=%s", dry, billingAccountID)

	ctx := context.Background()
	hc, errDc := google.DefaultClient(ctx, cloudbilling.CloudPlatformScope)
	if errDc != nil {
		log.Printf("killbill: DRY=%v account=%s could not get default client: %v", dry, billingAccountID, errDc)
		return errDc
	}
	client, errNew := cloudbilling.New(hc)
	if errNew != nil {
		log.Printf("killbill: DRY=%v account=%s could not create billing client: %v", dry, billingAccountID, errNew)
		return errNew
	}

	var projInfoList []*cloudbilling.ProjectBillingInfo

	name := "billingAccounts/" + billingAccountID
	call := client.BillingAccounts.Projects.List(name)
	if errPages := call.Pages(ctx, func(page *cloudbilling.ListProjectBillingInfoResponse) error {
		for _, v := range page.ProjectBillingInfo {
			log.Printf("killbill: DRY=%v account=%s found: project=%s", dry, billingAccountID, v.ProjectId)
			projInfoList = append(projInfoList, v)
		}
		return nil // NOTE: returning a non-nil error stops pagination.
	}); errPages != nil {
		log.Printf("killbill: DRY=%v account=%s could not page billing project info: %v", dry, billingAccountID, errPages)
		return errPages
	}

	log.Printf("killbill: DRY=%v account=%s found %d projects", dry, billingAccountID, len(projInfoList))

	var lastErr error
	for _, i := range projInfoList {
		if errKill := killprojbill(ctx, client, i); errKill != nil {
			log.Printf("killbill: DRY=%v account=%s project=%s error: %v", dry, billingAccountID, i.ProjectId, errKill)
			lastErr = errKill
		}
	}

	return lastErr
}

func killprojbill(ctx context.Context, client *cloudbilling.APIService, info *cloudbilling.ProjectBillingInfo) error {

	log.Printf("killprojbill: DRY=%v project=%s before: account=[%s]", dry, info.ProjectId, info.BillingAccountName)

	name := "projects/" + info.ProjectId

	if !dry {
		info.BillingAccountName = "" // unlink project from billing account
	}

	resp, errUpdate := client.Projects.UpdateBillingInfo(name, info).Context(ctx).Do()
	if errUpdate != nil {
		log.Printf("killprojbill: DRY=%v project=%s account=%s update error: %v", dry, info.ProjectId, info.BillingAccountName, errUpdate)
		return errUpdate
	}

	log.Printf("killprojbill: DRY=%v project=%s after: account=[%s]", dry, resp.ProjectId, resp.BillingAccountName)

	return nil
}
