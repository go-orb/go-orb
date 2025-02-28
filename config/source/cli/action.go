package cli

// SelectedAction is a helper struct to hold the service name and the action when an cli app chooses it.
type SelectedAction struct {
	ServiceName string
	Action      string
}

// Set updates its values.
func (s *SelectedAction) Set(sName string, action string) {
	s.ServiceName = sName
	s.Action = action
}

// String returns a human readable string.
func (s *SelectedAction) String() string {
	return s.ServiceName + ": " + s.Action
}
