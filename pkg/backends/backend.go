package backends

func DeleteBackend(backends []string, backendUrl string) ([]string, bool) {
	for pos, backend := range backends {
		if backend == backendUrl {
			backends = append(backends[:pos], backends[pos+1:]...)
			return backends, true
		}
	}
	return nil, false
}
