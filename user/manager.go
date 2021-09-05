package user

type Manager interface {
}

type ManagementProvider interface {
	UserManager() Manager
}
