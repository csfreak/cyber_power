package cyberpower

type mocker interface {
	calls() int
	called() bool
	called_once() bool
}

type mock struct {
	_calls int
}

func (m *mock) calls() int {
	return m._calls
}

func (m *mock) called() bool {
	return m._calls > 0
}

func (m *mock) called_once() bool {
	return m._calls == 1
}

type mock_cpmodule struct {
	mock
	parent       CyberPower
	update_error error
}

func (m *mock_cpmodule) update() error {
	m._calls++
	return m.update_error
}

func (m *mock_cpmodule) getParent() CyberPower {
	m._calls++
	return m.parent
}
