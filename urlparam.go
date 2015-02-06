package neo

// Representing named url parameters.
// For example if we have url like: ``/some/:name/other/:lastname``,
// then named parameters are ``name`` and ``lastname``.
type UrlParam map[string]string

func (u UrlParam) Get(key string) string {
	return u[key]
}

func (u UrlParam) Exist(key string) bool {
	_, ok := u[key]
	return ok
}
