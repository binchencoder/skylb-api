package tests

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/binchencoder/letsgo/trace"
	pb "github.com/binchencoder/skylb-api/cmd/stress/proto"
)

func TestContextMetadataPassing(t *testing.T) {
	tid := trace.GenerateTraceId()
	url := "http://localhost:6666/stress/v1/say-hello/dory"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("NewRequest: ", err)
		return
	}
	req.Header.Set("x-source", "web")
	req.Header.Set("x-Uid", "marlin")
	req.Header.Set("X-Cid", "disney")
	req.Header.Set("X-Request-Id", tid)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	Convey("Request should have no error", t, func() {
		So(err, ShouldBeNil)
	})

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	Convey("Read response body should have no error", t, func() {
		So(err, ShouldBeNil)
	})

	record := pb.SayHelloResp{}
	err = jsonpb.UnmarshalString(string(bodyBytes), &record)

	Convey("Metadata verification", t, func() {
		Convey("Status code should be 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})

		Convey("Response decoding should succeed", func() {
			So(err, ShouldBeNil)
		})

		Convey("Verify URL parameter 'name'", func() {
			So(record.Greeting, ShouldEqual, "Hello dory")
		})

		Convey("Verify context metadata through grpc", func() {
			So(record.Peer, ShouldEqual, "marlin")

			// Janus Gateway service ID is expected.
			So(record.ServiceId, ShouldEqual, 12)
		})

		Convey("Verify SkyLB auto to-and-from metadata of context", func() {
			// Should be the same as x-request-id.
			So(record.Tid, ShouldEqual, tid)
		})
	})
}
