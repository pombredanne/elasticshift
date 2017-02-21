package server

/*func TestClient(t *testing.T) {

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

		req := &api.CreateClientReq{}
		req.Name = test.name
		req.RedirectUris = test.redirectUrls
		req.Public = test.public
		req.TrustedPeers = test.trustedPeers

		res, err := s.Create(ctx, req)
		if test.expectErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			if res != nil {
				ids = append(ids, res.Id)
			}
		}
	}

	for _, id := range ids {
		s.Delete(ctx, &api.DeleteClientReq{ClientId: id})
	}
}*/
