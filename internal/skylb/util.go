package skylb

import (
	"fmt"

	"google.golang.org/grpc/naming"

	pb "github.com/binchencoder/skylb-api/proto"
)

func calcServiceKey(spec *pb.ServiceSpec) string {
	return fmt.Sprintf("%s.%s:%s", spec.Namespace, spec.ServiceName, spec.PortName)
}

func toNamingOp(op pb.Operation) naming.Operation {
	switch op {
	case pb.Operation_Add:
		return naming.Add
	case pb.Operation_Delete:
		return naming.Delete
	}

	return naming.Add
}

func opToString(op naming.Operation) string {
	switch op {
	case naming.Add:
		return "ADD"
	case naming.Delete:
		return "DELETE"
	}
	return ""
}
