package dashboard

// DashboardUpdate TODO: Return a DashboardUpdate struct with the result of checking the internals of the pod
type DashboardUpdate struct {
	// TODO: make a list of all SBOMs from the pod to add to DB / check with grype
	// TODO:
	Denied bool
}

// BadPodDashUpdate TODO: expand this to have field values expressing that the pod could not be examined
func BadPodDashUpdate() DashboardUpdate {
	return DashboardUpdate{Denied: true}
}
