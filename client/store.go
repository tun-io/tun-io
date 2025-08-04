package client

import "strings"

var (
	remoteDomain = "tunio.test"
	localDomain  = "localhost:8000"
	isSecure     = false
)

func SetupStore(
	RemoteDomain string,
	LocalDomain string,
	IsSecure bool,
) {

	if !strings.HasPrefix(LocalDomain, "http://") && !strings.HasPrefix(LocalDomain, "https://") {
		localDomain = "http://" + LocalDomain
	} else {
		localDomain = LocalDomain
	}

	if !strings.HasPrefix(RemoteDomain, "http://") && !strings.HasPrefix(RemoteDomain, "https://") {
		remoteDomain = "http://" + RemoteDomain
	} else {
		remoteDomain = RemoteDomain
	}

	isSecure = IsSecure
}
