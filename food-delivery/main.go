package main

import (
	"fmt"
	"log"
	"net/http"

	"food-delivery/internal/handler"
	"food-delivery/internal/repository"
	"food-delivery/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// --- Слой репозиториев (in-memory) ---
	catalogRepo := repository.NewInMemoryCatalogRepository()
	cartRepo    := repository.NewInMemoryCartRepository()
	orderRepo   := repository.NewInMemoryOrderRepository()
	userRepo    := repository.NewInMemoryUserRepository()

	// --- Слой use cases (внедряем репозитории явно) ---
	catalogUC := usecase.NewCatalogUseCase(catalogRepo)
	cartUC    := usecase.NewCartUseCase(cartRepo, catalogRepo)
	orderUC   := usecase.NewOrderUseCase(orderRepo, userRepo)

	// --- Слой обработчиков (внедряем use cases явно) ---
	catalogH := handler.NewCatalogHandler(catalogUC)
	cartH    := handler.NewCartHandler(cartUC)
	orderH   := handler.NewOrderHandler(orderUC)
	userH    := handler.NewUserHandler(userRepo)

	// --- Маршрутизатор ---
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Компонент управления каталогом
	r.Get("/restaurants", catalogH.ListRestaurants)
	r.Get("/restaurants/{restaurant_id}/menu", catalogH.GetMenu)

	// Компонент корзины
	r.Post("/cart", cartH.UpdateCart)

	// Компонент оформления и отслеживания заказа
	r.Post("/orders", orderH.PlaceOrder)
	r.Get("/orders/{order_id}", orderH.TrackOrder)

	// Пользователи (для регистрации и просмотра админом)
	r.Post("/users/register", userH.Register)
	r.Get("/users", userH.GetAll)
	r.Patch("/users/{user_id}/block-status", userH.UpdateBlockStatus)

	// Раздача HTML-страниц из папки frontend (должна быть последней!)
	fileServer := http.FileServer(http.Dir("./frontend"))
	r.Handle("/*", fileServer)

	addr := ":8080"
	fmt.Println("Server listening on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
