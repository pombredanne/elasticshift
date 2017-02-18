package server

import (
	"testing"

	"gitlab.com/conspico/elasticshift/pb"

	context "golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serv, err := newTestServer(ctx, t, nil)
	assert.Nil(t, err, "Failed to initialize server")
	defer serv.Store.GetSession().Close()

	s := NewClientServer(serv)

	tests := []struct {
		name         string
		id           string
		secret       string
		redirectUrls []string
		public       bool
		trustedPeers []string

		expectErr bool
	}{
		{
			name:         "foo",
			redirectUrls: []string{"http://testserver.com/foo"},
			public:       true,
			trustedPeers: []string{"esh"},
			expectErr:    false,
		},
		{
			name:         "foo",
			redirectUrls: []string{"http://testserver.com/foo"},
			public:       true,
			trustedPeers: []string{"esh"},
			expectErr:    false,
		},
	}

	var ids []string
	ids = append(ids, "")
	for _, test := range tests {

		req := &pb.CreateReq{}
		req.Client = &pb.ClientMsg{}
		req.Client.Name = test.name
		req.Client.RedirectUris = test.redirectUrls
		req.Client.Public = test.public
		req.Client.TrustedPeers = test.trustedPeers

		res, err := s.Create(ctx, req)
		if test.expectErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			if res.Client != nil {
				ids = append(ids, res.Client.Id)
			}
		}
	}

	for _, id := range ids {
		s.Delete(ctx, &pb.DeleteReq{ClientId: id})
	}
}
