package client

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
	remoteDomain = RemoteDomain
	localDomain = LocalDomain
	isSecure = IsSecure
}
