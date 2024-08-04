package aws

import (
	awsSdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetCurrentSession(region string) *session.Session {
	if cached := getServiceFromCache(sessionCacheKey, region); cached != nil {
		return cached.(*session.Session)
	}

	s := session.Must(
		session.NewSession(
			awsSdk.NewConfig().WithRegion(region),
		),
	)

	addServiceToCache(sessionCacheKey, region, s)

	return s
}
