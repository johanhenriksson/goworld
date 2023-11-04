package game_test

import (
	"log"
	"testing"
	"time"

	"github.com/johanhenriksson/goworld/game/client"
	"github.com/johanhenriksson/goworld/game/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGame(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Game Suite")
}

var _ = Describe("networking", func() {
	var srv *server.Server
	var cli *client.Client

	BeforeEach(func() {
		var err error

		srv, err = server.NewServer()
		if err != nil {
			Fail(err.Error())
		}

		cli = client.NewClient(func(e client.Event) {
			log.Println("client event:", e)
		})
		if err = cli.Connect("localhost"); err != nil {
			Fail(err.Error())
		}
	})

	It("works?", func() {
		err := cli.SendAuthToken(1337)
		Expect(err).ToNot(HaveOccurred())

		// some time later..
		time.Sleep(2 * time.Second)
		Expect(srv.Instance.Entities).To(HaveLen(1))
	})
})
