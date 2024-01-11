package application

import (
	"context"

	"github.com/diwise/integration-bigbelly/internal/pkg/domain"
)

type App interface {
	GetAssets(ctx context.Context) ([]domain.Asset, error)
}

type app struct {
	bigBellyApiUrl string
}

type BigBellyResponse struct {
	// fyll i det som saknas... :)
}

func New(bigBellyApiUrl string) App {
	return &app{
		bigBellyApiUrl: bigBellyApiUrl,
	}
}

func (a *app) GetAssets(ctx context.Context) ([]domain.Asset, error) {

	// skapa en ny struct för att kunna hantera svaret från servern
	// skapa en http client och anropa big belly API
	// plocka ut alla assets ur svaret från servern och returnera

	return nil, nil
}