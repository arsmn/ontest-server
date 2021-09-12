package user

type ServiceProvider interface {
	UserService() Service
}

type Service interface {
	Register() error
}
