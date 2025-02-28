package cli

type SelectedAction struct {
	ServiceName string
	Action      string
}

func (s *SelectedAction) Set(sName string, action string) {
	s.ServiceName = sName
	s.Action = action
}

func (s *SelectedAction) String() string {
	return s.ServiceName + ": " + s.Action
}
