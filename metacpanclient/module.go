package metacpanclient

import (
	"strings"
)

type Module struct {
	FileInfo
	hasUA
	pkg  *Package
	perm *Permission
}

func (m *Module) _type() Type {
	return TypeModule
}

func (m *Module) MetaCPANURL() string {
	const baseURL = MetaCPANURL + "/pod/release/"
	sb := strings.Builder{}
	sb.WriteString(baseURL)
	sb.WriteString(m.Author)
	sb.WriteRune('/')
	sb.WriteString(m.Release)
	sb.WriteRune('/')
	sb.WriteString(m.Path)
	return sb.String()
}

func (m *Module) Package() (*Package, error) {
	if m.pkg != nil {
		return m.pkg, nil
	}
	if m.mc == nil {
		return nil, ErrNilClient
	}
	var err error
	m.pkg, err = m.mc.Package(m.Documentation)
	return m.pkg, err
}

func (m *Module) Permission() (*Permission, error) {
	if m.perm != nil {
		return m.perm, nil
	}
	if m.mc == nil {
		return nil, ErrNilClient
	}
	var err error
	m.perm, err = m.mc.Permission(m.Documentation)
	return m.perm, err
}
