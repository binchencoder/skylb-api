package skylb

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/naming"

	"binchencoder.com/skylb-api/client/option"
	pb "binchencoder.com/skylb-api/proto"
	jt "binchencoder.com/skylb-api/testing"
	vexpb "binchencoder.com/gateway-proto/data"
)

func TestServiceLocator(t *testing.T) {
	Convey("Skylb service locator", t, func() {
		keeper := jt.SkyLbKeeperMock{}
		keeper.On("Start",
			vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE,
			"shared-test-client-service",
			true)
		locator := NewServiceLocator(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE, map[string]string{})
		locator.keeper = &keeper

		Convey("when resolve a service", func() {
			spec1 := pb.ServiceSpec{
				Namespace:   "ns",
				ServiceName: "svc1",
				PortName:    "ptn",
			}

			handler := jt.LoadBalanceHandlerMock{}
			handler.On("ServiceSpec").Return(&spec1)
			handler.On("BeforeResolve", &spec1,
				mock.MatchedBy(func(r naming.Resolver) bool {
					if r == nil {
						return false
					}
					_, ok := r.(*skylbResolver)
					return ok
				}),
				mock.MatchedBy(func(opts *option.ResolveOptions) bool {
					return opts != nil
				}))
			handler.On("AfterResolve", &spec1,
				vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE,
				"shared-test-client-service",
				mock.MatchedBy(func(keeper option.SkyLbKeeper) bool {
					if keeper == nil {
						return false
					}
					_, ok := keeper.(*jt.SkyLbKeeperMock)
					return ok
				}),
				mock.MatchedBy(func(tracer opentracing.Tracer) bool {
					return tracer != nil
				}),
				false)

			updates := make(chan []*naming.Update)
			keeper.On("RegisterService", &spec1).Return(updates)

			opts := func(*option.ResolveOptions) {}
			locator.Resolve(&handler, opts)

			So(locator.lbHandlers, ShouldContainKey, spec1.String())
			So(len(locator.specs), ShouldEqual, 1)
			So(len(locator.opts), ShouldEqual, 1)
			So(locator.specs[0].String(), ShouldResemble, spec1.String())

			Convey("when resolve another service", func() {
				spec2 := pb.ServiceSpec{
					Namespace:   "ns",
					ServiceName: "svc2",
					PortName:    "ptn",
				}

				h := jt.LoadBalanceHandlerMock{}
				h.On("ServiceSpec").Return(&spec2)
				h.On("BeforeResolve", &spec2,
					mock.MatchedBy(func(r naming.Resolver) bool {
						if r == nil {
							return false
						}
						_, ok := r.(*skylbResolver)
						return ok
					}),
					mock.MatchedBy(func(opts *option.ResolveOptions) bool {
						return opts != nil
					}))
				h.On("AfterResolve", &spec2,
					vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE,
					"shared-test-client-service",
					mock.MatchedBy(func(k option.SkyLbKeeper) bool {
						if k == nil {
							return false
						}
						_, ok := k.(*jt.SkyLbKeeperMock)
						return ok
					}),
					mock.MatchedBy(func(tracer opentracing.Tracer) bool {
						return tracer != nil
					}),
					false)

				updates := make(chan []*naming.Update)
				keeper.On("RegisterService", &spec2).Return(updates)

				locator.Resolve(&h, opts)

				So(locator.lbHandlers, ShouldContainKey, spec2.String())
				So(len(locator.specs), ShouldEqual, 2)
				So(len(locator.opts), ShouldEqual, 2)
				So(locator.specs[1].String(), ShouldResemble, spec2.String())

				Convey("next, when start the locator", func() {
					keeper.On("WaitUntilReady")
					locator.Start()
				})
			})
		})
	})
}

func TestParseEndpoint(t *testing.T) {
	ie := ParseEndpoint("192.168.10.41:4100")
	Convey("Parse endpoint ip:port", t, func() {
		So(ie.Host, ShouldEqual, "192.168.10.41")
		So(ie.Port, ShouldEqual, 4100)
	})

	ie = ParseEndpoint("avatarservice:4100")
	Convey("Parse endpoint host:port", t, func() {
		So(ie.Host, ShouldEqual, "avatarservice")
		So(ie.Port, ShouldEqual, 4100)
	})

	// Though this case is rared used in reality, still verify ParseEndpoint works.
	ie = ParseEndpoint(":4100")
	Convey("Parse endpoint :port", t, func() {
		So(ie.Host, ShouldEqual, "")
		So(ie.Port, ShouldEqual, 4100)
	})

	Convey("Parse bad port", t, func() {
		So(func() {
			ParseEndpoint("host:port")
		}, ShouldPanic)
	})

	Convey("Parse bad format", t, func() {
		So(func() {
			ParseEndpoint("host")
		}, ShouldPanic)
	})

	Convey("Parse blank", t, func() {
		So(func() {
			ParseEndpoint("")
		}, ShouldPanic)
	})

	var ep string
	Convey("Parse nil", t, func() {
		So(func() {
			ParseEndpoint(ep)
		}, ShouldPanic)
	})
}
