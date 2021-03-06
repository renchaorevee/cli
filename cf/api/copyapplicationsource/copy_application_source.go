package copyapplicationsource

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/cf/configuration/coreconfig"
	"github.com/cloudfoundry/cli/cf/net"
)

//go:generate counterfeiter . CopyApplicationSourceRepository

type CopyApplicationSourceRepository interface {
	CopyApplication(sourceAppGUID, targetAppGUID string) error
}

type CloudControllerApplicationSourceRepository struct {
	config  coreconfig.Reader
	gateway net.Gateway
}

func NewCloudControllerCopyApplicationSourceRepository(config coreconfig.Reader, gateway net.Gateway) *CloudControllerApplicationSourceRepository {
	return &CloudControllerApplicationSourceRepository{
		config:  config,
		gateway: gateway,
	}
}

func (repo *CloudControllerApplicationSourceRepository) CopyApplication(sourceAppGUID, targetAppGUID string) error {
	url := fmt.Sprintf("/v2/apps/%s/copy_bits", targetAppGUID)
	body := fmt.Sprintf(`{"source_app_guid":"%s"}`, sourceAppGUID)
	return repo.gateway.CreateResource(repo.config.APIEndpoint(), url, strings.NewReader(body), new(interface{}))
}
