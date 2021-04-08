package pkg

type (
	PKG interface {
		Version() string
		Writer()
	}
)