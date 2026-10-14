package repository

import (
	"errors"
	"food-delivery/internal/domain"
	"strings"
)

// InMemoryCatalogRepository — заглушка каталога (in-memory)
type InMemoryCatalogRepository struct {
	restaurants map[string]domain.Restaurant
	menus       map[string][]domain.MenuItem
	zones       map[string][]string
}

func NewInMemoryCatalogRepository() *InMemoryCatalogRepository {
	restaurants := map[string]domain.Restaurant{
		"r1": {ID: "r1", Name: "Суши Токио", Cuisine: "японская", Rating: 4.8, DeliveryTime: 30},
		"r2": {ID: "r2", Name: "Пицца Рома", Cuisine: "итальянская", Rating: 4.5, DeliveryTime: 45},
		"r3": {ID: "r3", Name: "Бургер Хаус", Cuisine: "американская", Rating: 4.3, DeliveryTime: 25},
		"r4": {ID: "r4", Name: "Тандыр", Cuisine: "узбекская", Rating: 4.7, DeliveryTime: 40},
		"r5": {ID: "r5", Name: "Том Ям", Cuisine: "тайская", Rating: 4.6, DeliveryTime: 35},
		"r6": {ID: "r6", Name: "Шаурма 24/7", Cuisine: "восточная", Rating: 4.2, DeliveryTime: 20},
		"r7": {ID: "r7", Name: "Грузинский дворик", Cuisine: "грузинская", Rating: 4.9, DeliveryTime: 50},
		"r8": {ID: "r8", Name: "Мексикана", Cuisine: "мексиканская", Rating: 4.4, DeliveryTime: 35},
	}

	menus := map[string][]domain.MenuItem{
		"r1": {
			{ID: "r1m1", Name: "Ролл Калифорния", Price: 350, Description: "Краб, авокадо, огурец, кунжут"},
			{ID: "r1m2", Name: "Суши Лосось", Price: 120, Description: "Свежий лосось, рис, нори"},
			{ID: "r1m3", Name: "Ролл Дракон", Price: 490, Description: "Угорь, авокадо, икра тобико"},
			{ID: "r1m4", Name: "Сашими Тунец", Price: 280, Description: "Тунец премиум, соевый соус, васаби"},
			{ID: "r1m5", Name: "Мисо суп", Price: 150, Description: "Тофу, водоросли вакамэ, зелёный лук"},
			{ID: "r1m6", Name: "Гункан Икра", Price: 160, Description: "Красная икра, рис, нори"},
			{ID: "r1m7", Name: "Ролл Феникс", Price: 420, Description: "Креветка темпура, сыр, огурец, острый соус"},
			{ID: "r1m8", Name: "Том Ям с морепродуктами", Price: 380, Description: "Кальмар, креветки, грибы эноки, кокосовое молоко"},
		},
		"r2": {
			{ID: "r2m1", Name: "Маргарита", Price: 450, Description: "Томатный соус, моцарелла, базилик"},
			{ID: "r2m2", Name: "Пепперони", Price: 520, Description: "Острая колбаса, моцарелла, томат"},
			{ID: "r2m3", Name: "Четыре сыра", Price: 580, Description: "Моцарелла, пармезан, горгонзола, рикотта"},
			{ID: "r2m4", Name: "Карбонара", Price: 490, Description: "Бекон, яйцо, пармезан, паста"},
			{ID: "r2m5", Name: "Тирамису", Price: 290, Description: "Маскарпоне, савоярди, эспрессо, какао"},
			{ID: "r2m6", Name: "Брускетта", Price: 220, Description: "Томаты черри, базилик, оливковое масло, чиабатта"},
			{ID: "r2m7", Name: "Лазанья Болоньезе", Price: 540, Description: "Говяжий фарш, соус бешамель, пармезан"},
			{ID: "r2m8", Name: "Паста Путтанеска", Price: 420, Description: "Оливки, каперсы, анчоусы, томаты"},
		},
		"r3": {
			{ID: "r3m1", Name: "Классик Бургер", Price: 380, Description: "Говядина 150г, салат, томат, лук, соус"},
			{ID: "r3m2", Name: "Чизбургер Двойной", Price: 490, Description: "Двойная котлета, чеддер, маринованный огурец"},
			{ID: "r3m3", Name: "Барбекю Бургер", Price: 450, Description: "Говядина, бекон, соус BBQ, луковые кольца"},
			{ID: "r3m4", Name: "Куриный Бургер", Price: 360, Description: "Хрустящая курица, коул сло, соус ранч"},
			{ID: "r3m5", Name: "Картофель Фри", Price: 150, Description: "Хрустящий картофель, соль, соус на выбор"},
			{ID: "r3m6", Name: "Луковые кольца", Price: 180, Description: "В панировке, соус тартар"},
			{ID: "r3m7", Name: "Чикен Наггетсы", Price: 240, Description: "8 шт., соусы кетчуп и горчица"},
			{ID: "r3m8", Name: "Молочный шейк", Price: 220, Description: "Клубника / ваниль / шоколад, 400мл"},
		},
		"r4": {
			{ID: "r4m1", Name: "Плов Фергана", Price: 420, Description: "Баранина, рис девзира, морковь, зира"},
			{ID: "r4m2", Name: "Шашлык Баранина", Price: 580, Description: "Маринованная баранина на углях, лук, зелень"},
			{ID: "r4m3", Name: "Манты", Price: 320, Description: "6 шт., говядина с луком, каймак"},
			{ID: "r4m4", Name: "Лагман", Price: 360, Description: "Тянутая лапша, говядина, овощи, зира"},
			{ID: "r4m5", Name: "Самса Говядина", Price: 80, Description: "Слоёное тесто, говядина, лук, специи"},
			{ID: "r4m6", Name: "Шурпа", Price: 340, Description: "Наваристый суп, баранина, овощи"},
			{ID: "r4m7", Name: "Дамлама", Price: 390, Description: "Баранина с овощами на пару"},
			{ID: "r4m8", Name: "Нан", Price: 60, Description: "Свежая лепёшка из тандыра"},
		},
		"r5": {
			{ID: "r5m1", Name: "Том Ям Кунг", Price: 420, Description: "Креветки, грибы, галангал, лемонграсс, кокос"},
			{ID: "r5m2", Name: "Пад Тай", Price: 380, Description: "Рисовая лапша, курица, ростки сои, арахис"},
			{ID: "r5m3", Name: "Зелёное карри", Price: 410, Description: "Курица, баклажан, кокосовое молоко, тайский базилик"},
			{ID: "r5m4", Name: "Жареный рис с морепродуктами", Price: 390, Description: "Кальмар, креветки, яйцо, соус"},
			{ID: "r5m5", Name: "Спринг роллы", Price: 280, Description: "4 шт., овощи, рисовая бумага, соус нам чим"},
			{ID: "r5m6", Name: "Манго со клейким рисом", Price: 240, Description: "Спелое манго, клейкий рис, кокосовый крем"},
		},
		"r6": {
			{ID: "r6m1", Name: "Шаурма Классик", Price: 250, Description: "Курица, овощи, соус чесночный, лаваш"},
			{ID: "r6m2", Name: "Шаурма Острая", Price: 270, Description: "Курица, острый соус, перец, лаваш"},
			{ID: "r6m3", Name: "Шаурма Двойная", Price: 350, Description: "Двойная порция мяса, все соусы"},
			{ID: "r6m4", Name: "Фалафель", Price: 220, Description: "Нут, зелень, тхина, лаваш"},
			{ID: "r6m5", Name: "Хумус с питой", Price: 180, Description: "Домашний хумус, оливковое масло, питы"},
			{ID: "r6m6", Name: "Картофель по-деревенски", Price: 140, Description: "С розмарином и чесноком"},
		},
		"r7": {
			{ID: "r7m1", Name: "Хинкали Мясные", Price: 360, Description: "6 шт., говядина и свинина, зелень"},
			{ID: "r7m2", Name: "Хачапури по-аджарски", Price: 420, Description: "Лодочка с сыром, яйцом и маслом"},
			{ID: "r7m3", Name: "Шашлык из телятины", Price: 620, Description: "Нежная телятина, маринад из вина и трав"},
			{ID: "r7m4", Name: "Чахохбили", Price: 390, Description: "Курица тушёная с томатами и зеленью"},
			{ID: "r7m5", Name: "Лобиани", Price: 280, Description: "Пирог с фасолью и луком"},
			{ID: "r7m6", Name: "Баклажаны с орехами", Price: 260, Description: "Жареные баклажаны, ореховая паста, чеснок"},
			{ID: "r7m7", Name: "Мцвади", Price: 580, Description: "Свиной шашлык на шампуре, соус ткемали"},
		},
		"r8": {
			{ID: "r8m1", Name: "Буррито с говядиной", Price: 420, Description: "Говядина, рис, фасоль, сальса, сметана"},
			{ID: "r8m2", Name: "Тако Карнитас", Price: 180, Description: "Тушёная свинина, лук, кинза, лайм, 1 шт."},
			{ID: "r8m3", Name: "Начос с сыром", Price: 320, Description: "Кукурузные чипсы, сырный соус, халапеньо, гуакамоле"},
			{ID: "r8m4", Name: "Кесадилья Курица", Price: 380, Description: "Курица, перец, лук, чеддер, сметана"},
			{ID: "r8m5", Name: "Гуакамоле с чипсами", Price: 260, Description: "Авокадо, томат, кинза, лайм, перец"},
			{ID: "r8m6", Name: "Фахитас", Price: 450, Description: "Говядина, перцы, лук, тортилья, соусы"},
			{ID: "r8m7", Name: "Чили кон карне", Price: 370, Description: "Говядина, фасоль, томаты, острый перец"},
		},
	}

	// Зоны доставки: ключевые слова адреса → список ID ресторанов
	zones := map[string][]string{
		// Центр
		"невский":      {"r1", "r2", "r7"},
		"литейный":     {"r1", "r2", "r7"},
		"садовая":      {"r2", "r3", "r7"},
		"центр":        {"r1", "r2", "r3", "r7"},
		"адмиралтейск": {"r2", "r7", "r8"},
		// Васильевский остров
		"васильев":     {"r2", "r3", "r5"},
		"васильевск":   {"r2", "r3", "r5"},
		"линия":        {"r3", "r5", "r6"},
		// Петроградская
		"петроград":    {"r1", "r4", "r6", "r7"},
		"кронверкск":   {"r1", "r6", "r7"},
		"большой пр":   {"r4", "r6", "r7"},
		// Выборгская / Калининский
		"выборгск":     {"r4", "r5", "r6"},
		"калинин":      {"r3", "r4", "r6"},
		"пискарёвск":   {"r4", "r6"},
		"север":        {"r5", "r6", "r8"},
		// Московский
		"московск":     {"r3", "r4", "r8"},
		"купчино":      {"r3", "r4", "r8"},
		"фрунзенск":    {"r3", "r8"},
		"юг":           {"r3", "r4", "r6"},
		// Приморский
		"приморск":     {"r1", "r5", "r8"},
		"комендантск":  {"r5", "r8"},
		"чёрная речка": {"r1", "r5"},
		// Красногвардейский / Невский р-н
		"красногвард":  {"r2", "r5", "r8"},
		"рыбацкое":     {"r4", "r6", "r8"},
		"обухово":      {"r4", "r6"},
		// Кировский / Красносельский
		"кировск":      {"r3", "r6", "r8"},
		"красносельск": {"r3", "r6"},
		"автово":       {"r3", "r6", "r8"},
		// Пушкин / Колпино
		"пушкин":       {"r4", "r7"},
		"колпино":      {"r4", "r8"},
		// Улицы
		"ленина":       {"r1", "r4", "r7"},
		"ленинград":    {"r1", "r2", "r3"},
	}

	return &InMemoryCatalogRepository{
		restaurants: restaurants,
		menus:       menus,
		zones:       zones,
	}
}

// GetRestaurantsByAddress возвращает рестораны по зоне адреса
func (r *InMemoryCatalogRepository) GetRestaurantsByAddress(address string) ([]domain.Restaurant, error) {
	if address == "" {
		return nil, errors.New("address is required")
	}

	lower := strings.ToLower(address)
	seen := map[string]bool{}
	result := []domain.Restaurant{}

	for keyword, ids := range r.zones {
		if strings.Contains(lower, keyword) {
			for _, id := range ids {
				if !seen[id] {
					seen[id] = true
					result = append(result, r.restaurants[id])
				}
			}
		}
	}

	// Пустой список — район не обслуживается
	return result, nil
}

// GetMenuByRestaurantID возвращает меню ресторана по его ID
func (r *InMemoryCatalogRepository) GetMenuByRestaurantID(restaurantID string) ([]domain.MenuItem, error) {
	menu, ok := r.menus[restaurantID]
	if !ok {
		return nil, errors.New("restaurant not found")
	}
	return menu, nil
}
