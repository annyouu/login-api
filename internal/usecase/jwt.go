package usecase

type JWTGenerator interface {
	Generate(username string) (string, error)
}