package isolated

import (
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("curl command", func() {
	It("returns the expected request", func() {
		session := helpers.CF("curl", "/v3/banana")
		Eventually(session).Should(Say(`"title":"CF-NotFound"`))
		Eventually(session).Should(Exit(0))
	})

	When("using -v", func() {
		It("returns the expected request with verbose output", func() {
			session := helpers.CF("curl", "-v", "/v3/banana")
			Eventually(session).Should(Say("GET /v3/banana HTTP/1.1"))
			Eventually(session).Should(Say(`"title": "CF-NotFound"`))
			Eventually(session).Should(Exit(0))
		})
	})
})
