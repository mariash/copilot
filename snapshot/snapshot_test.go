package snapshot_test

import (
	"os"
	"time"

	"code.cloudfoundry.org/copilot/models"
	"code.cloudfoundry.org/copilot/snapshot"
	"code.cloudfoundry.org/copilot/snapshot/fakes"
	"code.cloudfoundry.org/lager/lagertest"

	networking "istio.io/api/networking/v1alpha3"
	snap "istio.io/istio/pkg/mcp/snapshot"

	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run", func() {
	var (
		ticker    chan time.Time
		s         *snapshot.Snapshot
		collector *fakes.Collector
		setter    *fakes.Setter
		builder   *snap.InMemoryBuilder
	)

	BeforeEach(func() {
		l := lagertest.TestLogger{}
		ticker = make(chan time.Time)
		collector = &fakes.Collector{}
		setter = &fakes.Setter{}
		builder = snap.NewInMemoryBuilder()

		s = snapshot.New(l, ticker, collector, setter, builder)
	})

	It("mcp snapshots sends gateways, virutalServices, destinationRules and serviceEntries", func() {
		sig := make(chan os.Signal)
		ready := make(chan struct{})

		collector.CollectReturnsOnCall(0, routesWithBackends())
		collector.CollectReturnsOnCall(1, routesWithBackends())
		collector.CollectReturnsOnCall(2, routesWithBackends()[1:])

		go s.Run(sig, ready)
		ticker <- time.Time{}

		Eventually(setter.SetSnapshotCallCount).Should(Equal(1))
		node, shot := setter.SetSnapshotArgsForCall(0)
		Expect(node).To(Equal("default"))

		vs, dr, ga, se := verifyEnvelopes(shot, "1")

		Expect(vs.Hosts).To(Equal([]string{"foo.example.com"}))
		Expect(vs.Gateways).To(Equal([]string{"cloudfoundry-ingress"}))
		Expect(vs.Http).To(ConsistOf([]*networking.HTTPRoute{
			{
				Match: []*networking.HTTPMatchRequest{
					{
						Uri: &networking.StringMatch{
							MatchType: &networking.StringMatch_Prefix{
								Prefix: "/something",
							},
						},
					},
				},
				Route: []*networking.HTTPRouteDestination{
					{
						Destination: &networking.Destination{
							Host: "foo.example.com",
							Port: &networking.PortSelector{
								Port: &networking.PortSelector_Number{
									Number: 8080,
								},
							},
							Subset: "x-capi-guid",
						},
						Weight: 50,
					},
					{
						Destination: &networking.Destination{
							Host: "foo.example.com",
							Port: &networking.PortSelector{
								Port: &networking.PortSelector_Number{
									Number: 8080,
								},
							},
							Subset: "y-capi-guid",
						},
						Weight: 50,
					},
				},
			},
			{
				Route: []*networking.HTTPRouteDestination{
					{
						Destination: &networking.Destination{
							Host: "foo.example.com",
							Port: &networking.PortSelector{
								Port: &networking.PortSelector_Number{
									Number: 8080,
								},
							},
							Subset: "a-capi-guid",
						},
						Weight: 100,
					},
				},
			},
		}))

		Expect(dr.Host).To(Equal("foo.example.com"))
		Expect(dr.Subsets).To(ConsistOf([]*networking.Subset{
			{
				Name:   "a-capi-guid",
				Labels: map[string]string{"cfapp": "a-capi-guid"},
			},
			{
				Name:   "y-capi-guid",
				Labels: map[string]string{"cfapp": "y-capi-guid"},
			},
			{
				Name:   "x-capi-guid",
				Labels: map[string]string{"cfapp": "x-capi-guid"},
			},
		}))

		Expect(ga).To(Equal(networking.Gateway{
			Servers: []*networking.Server{
				&networking.Server{
					Port: &networking.Port{
						Number:   80,
						Protocol: "http",
						Name:     "http",
					},
					Hosts: []string{"*"},
				},
			},
		}))

		Expect(se.Hosts).To(Equal([]string{"foo.example.com"}))
		Expect(se.Addresses).To(BeNil())
		Expect(se.Ports).To(Equal([]*networking.Port{
			{
				Name:     "http",
				Number:   8080,
				Protocol: "http",
			},
		}))
		Expect(se.Location).To(Equal(networking.ServiceEntry_MESH_INTERNAL))
		Expect(se.Resolution).To(Equal(networking.ServiceEntry_STATIC))
		Expect(se.Endpoints).To(ConsistOf([]*networking.ServiceEntry_Endpoint{
			{
				Address: "10.10.10.1",
				Ports: map[string]uint32{
					"http": 65003,
				},
				Labels: map[string]string{"cfapp": "a-capi-guid"},
			},
			{
				Address: "10.0.0.0",
				Ports: map[string]uint32{
					"http": 65007,
				},
				Labels: map[string]string{"cfapp": "y-capi-guid"},
			},
			{
				Address: "10.0.0.1",
				Ports: map[string]uint32{
					"http": 65005,
				},
				Labels: map[string]string{"cfapp": "x-capi-guid"},
			},
		}))

		ticker <- time.Time{}

		Consistently(setter.SetSnapshotCallCount).Should(Equal(1))

		verifyEnvelopes(shot, "1")

		ticker <- time.Time{}

		Eventually(setter.SetSnapshotCallCount).Should(Equal(2))
		_, shot = setter.SetSnapshotArgsForCall(1)
		verifyEnvelopes(shot, "2")

		sig <- os.Kill
	})

	Context("when an internal route exists", func() {
		It("mcp snapshots sends virutalServices and serviceEntries for internal routes", func() {
			sig := make(chan os.Signal)
			ready := make(chan struct{})

			collector.CollectReturnsOnCall(0, internalRoutesWithBackends())

			go s.Run(sig, ready)
			ticker <- time.Time{}

			Eventually(setter.SetSnapshotCallCount).Should(Equal(1))
			node, shot := setter.SetSnapshotArgsForCall(0)
			Expect(node).To(Equal("default"))

			virtualServices := shot.Resources(snapshot.VirtualServiceTypeURL)
			destinationRules := shot.Resources(snapshot.DestinationRuleTypeURL)
			gateways := shot.Resources(snapshot.GatewayTypeURL)
			serviceEntries := shot.Resources(snapshot.ServiceEntryTypeURL)

			Expect(gateways).To(HaveLen(1))
			Expect(gateways[0].Metadata.Name).To(Equal("cloudfoundry-ingress"))

			Expect(virtualServices).To(HaveLen(1))
			Expect(virtualServices[0].Metadata.Name).To(Equal("copilot-service-for-foo.bar.internal"))

			var vs networking.VirtualService
			err := types.UnmarshalAny(virtualServices[0].Resource, &vs)
			Expect(err).NotTo(HaveOccurred())

			Expect(vs.Hosts).To(Equal([]string{"foo.bar.internal"}))
			Expect(vs.Gateways).To(HaveLen(0))
			Expect(vs.Http).To(ConsistOf([]*networking.HTTPRoute{
				{
					Route: []*networking.HTTPRouteDestination{
						{
							Destination: &networking.Destination{
								Host: "foo.bar.internal",
								Port: &networking.PortSelector{
									Port: &networking.PortSelector_Number{
										Number: 8080,
									},
								},
								Subset: "x-capi-guid",
							},
							Weight: 100,
						},
					},
				},
			}))

			Expect(serviceEntries).To(HaveLen(1))
			Expect(serviceEntries[0].Metadata.Name).To(Equal("copilot-service-entry-for-foo.bar.internal"))

			var se networking.ServiceEntry
			err = types.UnmarshalAny(serviceEntries[0].Resource, &se)
			Expect(err).NotTo(HaveOccurred())

			Expect(se.Hosts).To(Equal([]string{"foo.bar.internal"}))
			Expect(se.Addresses).To(Equal([]string{"127.127.0.1"}))
			Expect(se.Ports).To(Equal([]*networking.Port{
				{
					Name:     "tcp",
					Number:   8080,
					Protocol: "tcp",
				}}))
			Expect(se.Location).To(Equal(networking.ServiceEntry_MESH_INTERNAL))
			Expect(se.Resolution).To(Equal(networking.ServiceEntry_STATIC))
			Expect(se.Endpoints).To(ConsistOf([]*networking.ServiceEntry_Endpoint{
				{
					Address: "10.0.0.1",
					Ports: map[string]uint32{
						"tcp": 65005,
					},
					Labels: map[string]string{"cfapp": "x-capi-guid"},
				},
			}))

			Expect(destinationRules).To(HaveLen(0))

			sig <- os.Kill
		})
	})
})

func verifyEnvelopes(shot snap.Snapshot, version string) (
	vs networking.VirtualService,
	dr networking.DestinationRule,
	ga networking.Gateway,
	se networking.ServiceEntry) {

	virtualServices := shot.Resources(snapshot.VirtualServiceTypeURL)
	destinationRules := shot.Resources(snapshot.DestinationRuleTypeURL)
	gateways := shot.Resources(snapshot.GatewayTypeURL)
	serviceEntries := shot.Resources(snapshot.ServiceEntryTypeURL)

	vsVersion := shot.Version(snapshot.VirtualServiceTypeURL)
	Expect(vsVersion).To(Equal(version))

	drVersion := shot.Version(snapshot.DestinationRuleTypeURL)
	Expect(drVersion).To(Equal(version))

	// Gateway version is always 1
	gaVersion := shot.Version(snapshot.GatewayTypeURL)
	Expect(gaVersion).To(Equal("1"))

	seVersion := shot.Version(snapshot.ServiceEntryTypeURL)
	Expect(seVersion).To(Equal(version))

	Expect(virtualServices).To(HaveLen(1))
	Expect(virtualServices[0].Metadata.Name).To(Equal("copilot-service-for-foo.example.com"))

	Expect(destinationRules).To(HaveLen(1))
	Expect(destinationRules[0].Metadata.Name).To(Equal("copilot-rule-for-foo.example.com"))

	Expect(gateways).To(HaveLen(1))
	Expect(gateways[0].Metadata.Name).To(Equal("cloudfoundry-ingress"))

	Expect(serviceEntries).To(HaveLen(1))
	Expect(serviceEntries[0].Metadata.Name).To(Equal("copilot-service-entry-for-foo.example.com"))

	err := types.UnmarshalAny(virtualServices[0].Resource, &vs)
	Expect(err).NotTo(HaveOccurred())

	err = types.UnmarshalAny(destinationRules[0].Resource, &dr)
	Expect(err).NotTo(HaveOccurred())

	err = types.UnmarshalAny(gateways[0].Resource, &ga)
	Expect(err).NotTo(HaveOccurred())

	err = types.UnmarshalAny(serviceEntries[0].Resource, &se)
	Expect(err).NotTo(HaveOccurred())

	return vs, dr, ga, se
}

func routesWithBackends() []*models.RouteWithBackends {
	return []*models.RouteWithBackends{
		{
			Hostname: "foo.example.com",
			Path:     "/something",
			Backends: models.BackendSet{
				Backends: []*models.Backend{
					{
						Address: "10.0.0.1",
						Port:    uint32(65005),
					},
				},
			},
			CapiProcessGUID: "x-capi-guid",
			RouteWeight:     int32(50),
		},
		{
			Hostname: "foo.example.com",
			Path:     "/something",
			Backends: models.BackendSet{
				Backends: []*models.Backend{
					{
						Address: "10.0.0.0",
						Port:    uint32(65007),
					},
				},
			},
			CapiProcessGUID: "y-capi-guid",
			RouteWeight:     int32(50),
		},
		{
			Hostname: "foo.example.com",
			Path:     "",
			Backends: models.BackendSet{
				Backends: []*models.Backend{
					{
						Address: "10.10.10.1",
						Port:    uint32(65003),
					},
				},
			},
			CapiProcessGUID: "a-capi-guid",
			RouteWeight:     int32(100),
		},
	}
}

func internalRoutesWithBackends() []*models.RouteWithBackends {
	return []*models.RouteWithBackends{
		{
			Hostname: "foo.bar.internal",
			VIP:      "127.127.0.1",
			Internal: true,
			Path:     "/something",
			Backends: models.BackendSet{
				Backends: []*models.Backend{
					{
						Address: "10.0.0.1",
						Port:    uint32(65005),
					},
				},
			},
			CapiProcessGUID: "x-capi-guid",
			RouteWeight:     int32(100),
		},
	}
}
