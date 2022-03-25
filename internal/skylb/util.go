package skylb

import (
	"fmt"

	pb "github.com/binchencoder/skylb-apiv2/proto"
)

func calcServiceKey(spec *pb.ServiceSpec) string {
	return fmt.Sprintf("%s.%s:%s", spec.Namespace, spec.ServiceName, spec.PortName)
}
