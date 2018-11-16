package snapshot_test

import (
	"code.cloudfoundry.org/copilot/certs"
	"code.cloudfoundry.org/copilot/certs/fakes"
	"code.cloudfoundry.org/copilot/models"
	"code.cloudfoundry.org/copilot/snapshot"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	networking "istio.io/api/networking/v1alpha3"
)

var _ = Describe("Config", func() {
	var (
		config      *snapshot.Config
		fakeLocator *fakes.Locator
	)

	BeforeEach(func() {
		fakeLocator = &fakes.Locator{}
		config = snapshot.NewConfig(fakeLocator, lagertest.NewTestLogger("config"))
	})

	Describe("CreateGatewayEnvelopes", func() {
		It("creates gateway envelopes", func() {
			gateways := config.CreateGatewayEnvelopes()
			var ga networking.Gateway

			Expect(gateways).To(HaveLen(1))
			Expect(gateways[0].Metadata.Name).To(Equal("cloudfoundry-ingress"))

			err := types.UnmarshalAny(gateways[0].Resource, &ga)
			Expect(err).NotTo(HaveOccurred())

			Expect(ga).To(Equal(networking.Gateway{
				Servers: []*networking.Server{
					{
						Port: &networking.Port{
							Number:   80,
							Protocol: "http",
							Name:     "http",
						},
						Hosts: []string{"*"},
					},
				},
			}))
		})

		Context("When locator returns cert pair paths", func() {
			It("creates gateway envelopes with http and https servers", func() {
				certPairs := []certs.CertPairPaths{
					{
						Hosts:    []string{"example.com"},
						CertPath: "/some/path/not/important.crt",
						KeyPath:  "/some/path/not/important.key",
					},
				}
				fakeLocator.LocateReturns(certPairs, nil)

				gateways := config.CreateGatewayEnvelopes()
				var ga networking.Gateway

				Expect(gateways).To(HaveLen(1))
				Expect(gateways[0].Metadata.Name).To(Equal("cloudfoundry-ingress"))

				err := types.UnmarshalAny(gateways[0].Resource, &ga)
				Expect(err).NotTo(HaveOccurred())

				Expect(ga).To(Equal(networking.Gateway{
					Servers: []*networking.Server{
						{
							Port: &networking.Port{
								Number:   80,
								Protocol: "http",
								Name:     "http",
							},
							Hosts: []string{"*"},
						},
						{
							Port: &networking.Port{
								Number:   443,
								Protocol: "https",
								Name:     "https",
							},
							Tls: &networking.Server_TLSOptions{
								Mode:              networking.Server_TLSOptions_SIMPLE,
								ServerCertificate: "/some/path/not/important.crt",
								PrivateKey:        "/some/path/not/important.key",
							},
							Hosts: []string{"example.com"},
						},
					},
				}))
			})
		})

		PContext("internal routes", func() {
			It("should not create any gateways", func() {
				gateways := config.CreateGatewayEnvelopes()
				Expect(gateways).To(HaveLen(0))
			})
		})
	})

	Describe("CreateVirtualServiceEnvelopes", func() {
		It("creates virtualService envelopes", func() {
			virtualServices := config.CreateVirtualServiceEnvelopes(routesWithBackends(), "1")
			var vs networking.VirtualService

			Expect(virtualServices).To(HaveLen(1))
			Expect(virtualServices[0].Metadata.Name).To(Equal("copilot-service-for-foo.example.com"))

			err := types.UnmarshalAny(virtualServices[0].Resource, &vs)
			Expect(err).NotTo(HaveOccurred())

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
		})

		Context("internal routes", func() {
			It("creates virtualService envelopes that are super special", func() {
				virtualServices := config.CreateVirtualServiceEnvelopes(internalRoutesWithBackends(), "1")
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
			})
		})
	})

	Describe("CreateDestinationRuleEnvelopes", func() {
		It("creates destinationRule envelopes", func() {
			destinationRules := config.CreateDestinationRuleEnvelopes(routesWithBackends(), "1")
			var dr networking.DestinationRule

			Expect(destinationRules).To(HaveLen(1))
			Expect(destinationRules[0].Metadata.Name).To(Equal("copilot-rule-for-foo.example.com"))

			err := types.UnmarshalAny(destinationRules[0].Resource, &dr)
			Expect(err).NotTo(HaveOccurred())

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
		})

		Context("internal routes", func() {
			It("should not create any destination rules", func() {
				destinationRules := config.CreateDestinationRuleEnvelopes(internalRoutesWithBackends(), "1")
				Expect(destinationRules).To(HaveLen(0))
			})
		})
	})

	Describe("CreateServiceEntryEnvelopes", func() {
		It("creates serviceEntry envelopes", func() {
			serviceEntries := config.CreateServiceEntryEnvelopes(routesWithBackends(), "1")
			var se networking.ServiceEntry

			Expect(serviceEntries).To(HaveLen(1))
			Expect(serviceEntries[0].Metadata.Name).To(Equal("copilot-service-entry-for-foo.example.com"))

			err := types.UnmarshalAny(serviceEntries[0].Resource, &se)
			Expect(err).NotTo(HaveOccurred())

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
		})

		Context("internal routes", func() {
			It("creates service entries that are super special", func() {
				serviceEntries := config.CreateServiceEntryEnvelopes(internalRoutesWithBackends(), "1")
				var se networking.ServiceEntry

				Expect(serviceEntries).To(HaveLen(1))
				Expect(serviceEntries[0].Metadata.Name).To(Equal("copilot-service-entry-for-foo.bar.internal"))

				err := types.UnmarshalAny(serviceEntries[0].Resource, &se)
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

			})
		})
	})
})

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
