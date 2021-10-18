package auth

type FakeConfig struct {
	Hosts map[string]*HostConfig
}

func (f *FakeConfig) GetToken(hostname string) string {
	return f.Hosts[hostname].Token
}

func (f *FakeConfig) GetUser(hostname string) string {
	return f.Hosts[hostname].User
}

func (f *FakeConfig) HasHosts() bool {
	return len(f.Hosts) > 0
}

func (f *FakeConfig) GetHosts() []string {
	var hosts []string
	for k := range f.Hosts {
		hosts = append(hosts, k)
	}
	return hosts
}
