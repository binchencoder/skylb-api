package skylb

import (
	"fmt"

	"google.golang.org/grpc/attributes"

	pb "github.com/binchencoder/skylb-apiv2/proto"
)

type AddressOpKey struct{}

func calcServiceKey(spec *pb.ServiceSpec) string {
	return fmt.Sprintf("%s.%s:%s", spec.Namespace, spec.ServiceName, spec.PortName)
}

func toAttributes(op pb.Operation) *attributes.Attributes {
	switch op {
	case pb.Operation_Add:
		return attributes.New(AddressOpKey{}, pb.Operation_Add)
	case pb.Operation_Delete:
		return attributes.New(AddressOpKey{}, pb.Operation_Delete)
	}

	return attributes.New(AddressOpKey{}, pb.Operation_Add)
}

func opToString(op pb.Operation) string {
	switch op {
	case pb.Operation_Add:
		return "ADD"
	case pb.Operation_Delete:
		return "DELETE"
	}
	return ""
}
