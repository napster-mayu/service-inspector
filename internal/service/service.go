package service

type baseService struct {
}

type BaseSevice interface {
	Welcome() string
}

func NewService() BaseSevice {
	return &baseService{}
}

func (s *baseService) Welcome() string {
	return "welcome"
}
