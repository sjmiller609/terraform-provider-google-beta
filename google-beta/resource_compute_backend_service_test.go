package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeBackendService_basic(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	extraCheckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_basic(serviceName, checkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
			{
				Config: testAccComputeBackendService_basicModified(
					serviceName, checkName, extraCheckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withBackend(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	igName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	itName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withBackend(
					serviceName, igName, itName, checkName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.lipsum", &svc),
				),
			},
			{
				Config: testAccComputeBackendService_withBackend(
					serviceName, igName, itName, checkName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.lipsum", &svc),
				),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if svc.TimeoutSec != 20 {
		t.Errorf("Expected TimeoutSec == 20, got %d", svc.TimeoutSec)
	}
	if svc.Protocol != "HTTP" {
		t.Errorf("Expected Protocol to be HTTP, got %q", svc.Protocol)
	}
	if len(svc.Backends) != 1 {
		t.Errorf("Expected 1 backend, got %d", len(svc.Backends))
	}
}

func TestAccComputeBackendService_withBackendAndIAP(t *testing.T) {
	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	igName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	itName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withBackendAndIAP(
					serviceName, igName, itName, checkName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExistsWithIAP("google_compute_backend_service.lipsum", &svc),
					resource.TestCheckResourceAttr("google_compute_backend_service.lipsum", "iap.0.oauth2_client_secret", "test"),
				),
			},
			{
				ResourceName:            "google_compute_backend_service.lipsum",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"iap.0.oauth2_client_secret"},
			},
			{
				Config: testAccComputeBackendService_withBackend(
					serviceName, igName, itName, checkName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExistsWithoutIAP(
						"google_compute_backend_service.lipsum", &svc),
				),
			},
		},
	})

	if svc.TimeoutSec != 10 {
		t.Errorf("Expected TimeoutSec == 10, got %d", svc.TimeoutSec)
	}
	if svc.Protocol != "HTTP" {
		t.Errorf("Expected Protocol to be HTTP, got %q", svc.Protocol)
	}
	if len(svc.Backends) != 1 {
		t.Errorf("Expected 1 backend, got %d", len(svc.Backends))
	}

}

func TestAccComputeBackendService_updatePreservesOptionalParameters(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withSessionAffinity(
					serviceName, checkName, "initial-description", "GENERATED_COOKIE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
			{
				Config: testAccComputeBackendService_withSessionAffinity(
					serviceName, checkName, "updated-description", "GENERATED_COOKIE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
		},
	})

	if svc.SessionAffinity != "GENERATED_COOKIE" {
		t.Errorf("Expected SessionAffinity == \"GENERATED_COOKIE\", got %s", svc.SessionAffinity)
	}
}

func TestAccComputeBackendService_withConnectionDraining(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withConnectionDraining(serviceName, checkName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if svc.ConnectionDraining.DrainingTimeoutSec != 10 {
		t.Errorf("Expected ConnectionDraining.DrainingTimeoutSec == 10, got %d", svc.ConnectionDraining.DrainingTimeoutSec)
	}
}

func TestAccComputeBackendService_withConnectionDrainingAndUpdate(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withConnectionDraining(serviceName, checkName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
			{
				Config: testAccComputeBackendService_basic(serviceName, checkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
		},
	})

	if svc.ConnectionDraining.DrainingTimeoutSec != 300 {
		t.Errorf("Expected ConnectionDraining.DrainingTimeoutSec == 300, got %d", svc.ConnectionDraining.DrainingTimeoutSec)
	}
}

func TestAccComputeBackendService_withHttpsHealthCheck(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withHttpsHealthCheck(serviceName, checkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withCdnPolicy(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withCdnPolicy(serviceName, checkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_withSecurityPolicy(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	polName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withSecurityPolicy(serviceName, checkName, polName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
					resource.TestMatchResourceAttr("google_compute_backend_service.foobar", "security_policy", regexp.MustCompile(polName)),
				),
			},
		},
	})
}

func testAccCheckComputeBackendServiceExists(n string, svc *compute.BackendService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.BackendServices.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Backend service %s not found", rs.Primary.ID)
		}

		*svc = *found

		return nil
	}
}

func testAccCheckComputeBackendServiceExistsWithIAP(n string, svc *compute.BackendService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.BackendServices.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Backend service %s not found", rs.Primary.ID)
		}

		if found.Iap == nil || found.Iap.Enabled == false {
			return fmt.Errorf("IAP not found or not enabled. Saw %v", found.Iap)
		}

		*svc = *found

		return nil
	}
}

func testAccCheckComputeBackendServiceExistsWithoutIAP(n string, svc *compute.BackendService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.BackendServices.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Backend service %s not found", rs.Primary.ID)
		}

		if found.Iap != nil && found.Iap.Enabled == true {
			return fmt.Errorf("IAP enabled when it should be disabled")
		}

		*svc = *found

		return nil
	}
}
func TestAccComputeBackendService_withCDNEnabled(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withCDNEnabled(
					serviceName, checkName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
		},
	})

	if svc.EnableCDN != true {
		t.Errorf("Expected EnableCDN == true, got %t", svc.EnableCDN)
	}
}

func TestAccComputeBackendService_withSessionAffinity(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withSessionAffinity(
					serviceName, checkName, "description", "CLIENT_IP"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
			{
				Config: testAccComputeBackendService_withSessionAffinity(
					serviceName, checkName, "description", "GENERATED_COOKIE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
		},
	})

	if svc.SessionAffinity != "GENERATED_COOKIE" {
		t.Errorf("Expected SessionAffinity == \"GENERATED_COOKIE\", got %s", svc.SessionAffinity)
	}
}

func TestAccComputeBackendService_withAffinityCookieTtlSec(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var svc compute.BackendService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withAffinityCookieTtlSec(
					serviceName, checkName, "description", "GENERATED_COOKIE", 300),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.foobar", &svc),
				),
			},
		},
	})

	if svc.SessionAffinity != "GENERATED_COOKIE" {
		t.Errorf("Expected SessionAffinity == \"GENERATED_COOKIE\", got %s", svc.SessionAffinity)
	}

	if svc.AffinityCookieTtlSec != 300 {
		t.Errorf("Expected AffinityCookieTtlSec == 300, got %v", svc.AffinityCookieTtlSec)
	}
}

func TestAccComputeBackendService_withMaxConnections(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	igName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	itName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	var svc compute.BackendService
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withMaxConnections(
					serviceName, igName, itName, checkName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.lipsum", &svc),
				),
			},
			{
				Config: testAccComputeBackendService_withMaxConnections(
					serviceName, igName, itName, checkName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.lipsum", &svc),
				),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if svc.Backends[0].MaxConnections != 20 {
		t.Errorf("Expected MaxConnections == 20, got %d", svc.Backends[0].MaxConnections)
	}
}

func TestAccComputeBackendService_withMaxConnectionsPerInstance(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	igName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	itName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	var svc compute.BackendService
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withMaxConnectionsPerInstance(
					serviceName, igName, itName, checkName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.lipsum", &svc),
				),
			},
			{
				Config: testAccComputeBackendService_withMaxConnectionsPerInstance(
					serviceName, igName, itName, checkName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBackendServiceExists(
						"google_compute_backend_service.lipsum", &svc),
				),
			},
			{
				ResourceName:      "google_compute_backend_service.lipsum",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	if svc.Backends[0].MaxConnectionsPerInstance != 20 {
		t.Errorf("Expected MaxConnectionsPerInstance == 20, got %d", svc.Backends[0].MaxConnectionsPerInstance)
	}
}

func TestAccComputeBackendService_withCustomHeaders(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_withCustomHeaders(serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendService_basic(serviceName, checkName),
			},
			{
				ResourceName:      "google_compute_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendService_internalLoadBalancing(t *testing.T) {
	t.Parallel()

	fr := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	proxy := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	urlmap := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendService_internalLoadBalancing(fr, proxy, backend, hc, urlmap),
			},
			{
				ResourceName:      "google_compute_backend_service.backend_service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeBackendService_basic(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = ["${google_compute_http_health_check.zero.self_link}"]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, checkName)
}

func testAccComputeBackendService_withCDNEnabled(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = ["${google_compute_http_health_check.zero.self_link}"]
  enable_cdn    = true
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, checkName)
}

func testAccComputeBackendService_basicModified(serviceName, checkOne, checkTwo string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
    name = "%s"
    health_checks = ["${google_compute_http_health_check.one.self_link}"]
}

resource "google_compute_http_health_check" "zero" {
    name = "%s"
    request_path = "/"
    check_interval_sec = 1
    timeout_sec = 1
}

resource "google_compute_http_health_check" "one" {
    name = "%s"
    request_path = "/one"
    check_interval_sec = 30
    timeout_sec = 30
}
`, serviceName, checkOne, checkTwo)
}

func testAccComputeBackendService_withBackend(
	serviceName, igName, itName, checkName string, timeout int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = %v

  backend {
    group = "${google_compute_instance_group_manager.foobar.instance_group}"
  }

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_instance_group_manager" "foobar" {
  name               = "%s"
  version {
    instance_template  = "${google_compute_instance_template.foobar.self_link}"
    name               = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = "${data.google_compute_image.my_image.self_link}"
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_http_health_check" "default" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, timeout, igName, itName, checkName)
}

func testAccComputeBackendService_withBackendAndIAP(
	serviceName, igName, itName, checkName string, timeout int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = %v

  backend {
    group = "${google_compute_instance_group_manager.foobar.instance_group}"
	}

  iap {
    oauth2_client_id = "test"
    oauth2_client_secret = "test"
  }

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_instance_group_manager" "foobar" {
  name               = "%s"
  version {
    instance_template  = "${google_compute_instance_template.foobar.self_link}"
    name               = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = "${data.google_compute_image.my_image.self_link}"
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_http_health_check" "default" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, timeout, igName, itName, checkName)
}

func testAccComputeBackendService_withSessionAffinity(serviceName, checkName, description, affinityName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name             = "%s"
  description      = "%s"
  health_checks    = ["${google_compute_http_health_check.zero.self_link}"]
  session_affinity = "%s"
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, description, affinityName, checkName)
}

func testAccComputeBackendService_withAffinityCookieTtlSec(serviceName, checkName, description, affinityName string, affinityCookieTtlSec int64) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name                    = "%s"
  description             = "%s"
  health_checks           = ["${google_compute_http_health_check.zero.self_link}"]
  session_affinity        = "%s"
  affinity_cookie_ttl_sec = %v
}
 resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, description, affinityName, affinityCookieTtlSec, checkName)
}

func testAccComputeBackendService_withConnectionDraining(serviceName, checkName string, drainingTimeout int64) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = ["${google_compute_http_health_check.zero.self_link}"]
  connection_draining_timeout_sec = %v
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, drainingTimeout, checkName)
}

func testAccComputeBackendService_withHttpsHealthCheck(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = ["${google_compute_https_health_check.zero.self_link}"]
  protocol      = "HTTPS"
}

resource "google_compute_https_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, checkName)
}

func testAccComputeBackendService_withCdnPolicy(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = ["${google_compute_http_health_check.zero.self_link}"]

  cdn_policy {
    cache_key_policy {
      include_protocol       = true
      include_host           = true
      include_query_string   = true
      query_string_whitelist = ["foo", "bar"]
    }
  }
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, checkName)
}

func testAccComputeBackendService_withSecurityPolicy(serviceName, checkName, polName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = ["${google_compute_http_health_check.zero.self_link}"]
  security_policy = "${google_compute_security_policy.policy.self_link}"
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_security_policy" "policy" {
	name        = "%s"
	description = "basic security policy"
}
`, serviceName, checkName, polName)
}

func testAccComputeBackendService_withMaxConnections(
	serviceName, igName, itName, checkName string, maxConnections int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "TCP"

  backend {
    group = "${google_compute_instance_group_manager.foobar.instance_group}"
    max_connections = %v
  }

  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_instance_group_manager" "foobar" {
  name               = "%s"
  version {
    instance_template  = "${google_compute_instance_template.foobar.self_link}"
    name               = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = "${data.google_compute_image.my_image.self_link}"
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_health_check" "default" {
  name               = "%s"
  tcp_health_check {
      port = "110"
  }
}
`, serviceName, maxConnections, igName, itName, checkName)
}

func testAccComputeBackendService_withMaxConnectionsPerInstance(
	serviceName, igName, itName, checkName string, maxConnectionsPerInstance int64) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_backend_service" "lipsum" {
  name        = "%s"
  description = "Hello World 1234"
  port_name   = "http"
  protocol    = "TCP"

  backend {
    group = "${google_compute_instance_group_manager.foobar.instance_group}"
    max_connections_per_instance = %v
  }

  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_instance_group_manager" "foobar" {
  name               = "%s"
  version {
    instance_template  = "${google_compute_instance_template.foobar.self_link}"
    name               = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "%s"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = "${data.google_compute_image.my_image.self_link}"
    auto_delete  = true
    boot         = true
  }
}

resource "google_compute_health_check" "default" {
  name               = "%s"
  tcp_health_check {
      port = "110"
  }
}
`, serviceName, maxConnectionsPerInstance, igName, itName, checkName)
}

func testAccComputeBackendService_withCustomHeaders(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  health_checks = ["${google_compute_http_health_check.zero.self_link}"]

  custom_request_headers =  ["Client-Region: {client_region}", "Client-Rtt: {client_rtt_msec}"]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, serviceName, checkName)
}

func testAccComputeBackendService_internalLoadBalancing(fr, proxy, backend, hc, urlmap string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  name                  = "%s"
  target                = "${google_compute_target_http_proxy.default.self_link}"
  port_range            = "80"
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"
  ip_address            = "0.0.0.0"
}

resource "google_compute_target_http_proxy" "default" {
  name        = "%s"
  description = "a description"
  url_map     = "${google_compute_url_map.default.self_link}"
}

resource "google_compute_backend_service" "backend_service" {
  name                  = "%s"
  port_name             = "http"
  protocol              = "HTTP"
  timeout_sec           = 10
  load_balancing_scheme = "INTERNAL_SELF_MANAGED"

  backend {
    group = "${google_compute_instance_group_manager.foobar.instance_group}"
    balancing_mode = "RATE"
    capacity_scaler = 0.4
    max_rate_per_instance = 50
  }

  health_checks = ["${google_compute_health_check.default.self_link}"]
}

resource "google_compute_health_check" "default" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}

resource "google_compute_url_map" "default" {
  name            = "%s"
  description     = "a description"
  default_service = "${google_compute_backend_service.backend_service.self_link}"

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = "${google_compute_backend_service.backend_service.self_link}"

    path_rule {
      paths   = ["/*"]
      service = "${google_compute_backend_service.backend_service.self_link}"
    }
  }
}

data "google_compute_image" "debian_image" {
  family   = "debian-9"
  project  = "debian-cloud"
}

resource "google_compute_instance_group_manager" "foobar" {
  name               = "igm-internal"
  version {
    instance_template  = "${google_compute_instance_template.foobar.self_link}"
    name               = "primary"
  }
  base_instance_name = "foobar"
  zone               = "us-central1-f"
  target_size        = 1
}

resource "google_compute_instance_template" "foobar" {
  name         = "instance-template-internal"
  machine_type = "n1-standard-1"

  network_interface {
    network = "default"
  }

  disk {
    source_image = "${data.google_compute_image.debian_image.self_link}"
    auto_delete  = true
    boot         = true
  }
}`, fr, proxy, backend, hc, urlmap)
}
