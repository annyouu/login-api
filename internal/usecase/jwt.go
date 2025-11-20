package usecase

type JWTGeneratorInterface interface {
	Generate(username string) (string, error)
}