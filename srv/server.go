package srv

type Server interface {
	Listen(port int) error
}

type DummyServer struct {
	realm Realm

	accept chan Client
}

func (s *DummyServer) Listen(port int) error {
	s.accept = make(chan Client)
	defer close(s.accept)
	for {
		conn := <-s.accept
		if err := s.realm.Accept(conn); err != nil {
			return err
		}
	}
}

func (s *DummyServer) Connect(c Client) {
	s.accept <- c
}

type Realm interface {
	Accept(Client) error
}

type testRealm struct {
	area  Area
	chars map[string]*Unit
}

func NewRealm(area Area) *testRealm {
	return &testRealm{
		area:  area,
		chars: make(map[string]*Unit),
	}
}

func (r *testRealm) Accept(c Client) error {
	// read client token
	token, err := c.ReadToken()
	if err != nil {
		c.Drop(err.Error())
		return err
	}

	// todo: authenticate token

	// todo: actually load character
	player, exists := r.chars[token.Character]
	if !exists {
		player = &Unit{
			name: token.Character,
		}
	}

	// join area
	r.area.Join(player)

	// todo: send observe entity event
	c.Observe(r.area, player)

	return nil
}
