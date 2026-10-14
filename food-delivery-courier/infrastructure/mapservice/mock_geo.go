package mapservice

import (
	"math"
	"food-delivery-courier/internal/domain"
)

type MockGeoService struct{}

func NewMockGeoService() *MockGeoService {
	return &MockGeoService{}
}

func (s *MockGeoService) GetDeliveryRoute(courier, restaurant, delivery domain.Location) (domain.Route, error) {
	d1 := s.CalculateDistance(courier, restaurant)
	d2 := s.CalculateDistance(restaurant, delivery)

	return domain.Route{
		TotalDistanceKm:   d1 + d2,
		TotalEstimatedMin: int((d1+d2)*2) + 5, // ~2 мин/км + 5 мин на сборку
		Legs: []domain.RouteLeg{
			{Type: "to_restaurant", From: courier, To: restaurant,
				DistanceKm: d1, EstimatedMinutes: int(d1 * 2)},
			{Type: "to_client", From: restaurant, To: delivery,
				DistanceKm: d2, EstimatedMinutes: int(d2*2) + 5},
		},
	}, nil
}

func (s *MockGeoService) CalculateDistance(from, to domain.Location) float64 {
	const R = 6371.0
	lat1 := from.Lat * math.Pi / 180
	lat2 := to.Lat * math.Pi / 180
	dlat := (to.Lat - from.Lat) * math.Pi / 180
	dlon := (to.Lon - from.Lon) * math.Pi / 180
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
