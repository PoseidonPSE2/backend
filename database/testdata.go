package database

import (
	"log"
	"math/rand"

	"gorm.io/gorm"
)

func CreateTestData(db *gorm.DB) *gorm.DB {
	log.Print("Test data creation started")

	db = CreateUsers(db)
	db = CreateBottles(db)
	db = CreateRefillStations(db)
	db = CreateRefillStationReviews(db)
	db = CreateRefillStationProblems(db)
	db = CreateWaterTransactions(db)
	db = CreateLikes(db)

	log.Print("Test data creation finished")

	return db
}

func CreateUsers(db *gorm.DB) *gorm.DB {
	users := []User{
		{
			FirstName: "Jonas",
			LastName:  "Blum",
		},
		{
			FirstName: "Heiko",
			LastName:  "Michel",
		},
		{
			FirstName: "David",
			LastName:  "Miller",
		},
		{
			FirstName: "Yusuf Can",
			LastName:  "Özdemirkan",
		},
		{
			FirstName: "Alejandro",
			LastName:  "Restrepo Klinge",
		},
	}

	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	log.Print("Created users successfully")

	return db
}

func CreateBottles(db *gorm.DB) *gorm.DB {
	bottle0 := ImageToBase64("./images/bottles/bottle0.png")
	bottle1 := ImageToBase64("./images/bottles/bottle1.jpg")
	bottle2 := ImageToBase64("./images/bottles/bottle2.jpg")
	bottle3 := ImageToBase64("./images/bottles/bottle3.jpg")

	bottles := []Bottle{
		{
			UserID:      4,
			NFCID:       "04:72:52:1A:94:11:90",
			FillVolume:  100,
			WaterType:   "tap",
			Title:       "Letzte App-Wahl",
			BottleImage: &bottle0,
		},
		{
			UserID:      5,
			NFCID:       "13:E0:0B:35",
			FillVolume:  500,
			WaterType:   "tap",
			Title:       "Daily Bottle",
			BottleImage: &bottle1,
		},
		{
			UserID:      5,
			NFCID:       "13:8E:BD:0C",
			FillVolume:  250,
			WaterType:   "mineral",
			Title:       "Fancy Bottle",
			BottleImage: &bottle2,
		},
		{
			UserID:      3,
			NFCID:       "",
			FillVolume:  1000,
			WaterType:   "mineral",
			Title:       "Sports Bottle",
			BottleImage: &bottle3,
		},
	}

	if err := db.Create(&bottles).Error; err != nil {
		log.Fatalf("failed to create bottle: %v", err)
	}

	log.Print("Created bottles successfully")

	return db
}

func CreateRefillStations(db *gorm.DB) *gorm.DB {
	stadtpark_picture := ImageToBase64("./images/refill_stations/stadtpark.jpg")
	rewe_picture := ImageToBase64("./images/refill_stations/rewe.jpg")
	gartenschau_picture := ImageToBase64("./images/refill_stations/gartenschau.jpeg")
	wochenmarkt_picture := ImageToBase64("./images/refill_stations/wochenmarkt.jpg")

	refillStations := []*RefillStation{
		{
			Name:               "Stadtpark KL",
			Description:        "Irgendwas mit einem Schnitzeljagd?",
			Latitude:           49.437551349217266,
			Longitude:          7.761465081072085,
			Address:            "Trippstadter Str. 2, 67663 Kaiserslautern",
			WaterSource:        "Spitzrainbrunnen",
			OpeningTimes:       "Mon - Son / 00:00 AM - 12:59 PM",
			Type:               "smart",
			OfferedWaterTypes:  "both",
			RefillStationImage: &stadtpark_picture,
		},
		{
			Name:               "Rewe Station",
			Description:        "Wasserhahn REWE",
			Latitude:           49.44490159879211,
			Longitude:          7.767478778642334,
			Address:            "Fruchthallstraße 29, 67655 Kaiserslautern",
			WaterSource:        "Spitzrainbrunnen",
			OpeningTimes:       "Mon - Sam / 7:00 AM - 10:00 PM",
			Type:               "manual",
			OfferedWaterTypes:  "tap",
			RefillStationImage: &rewe_picture,
		},
		{
			Name:               "Gartenschau KL",
			Description:        "Die Gartenschau Kaiserslautern ist ein atemberaubendes jährliches Ereignis, das die Schönheit der Natur und die Freude am Gartenbau feiert.",
			Latitude:           49.44694255672088,
			Longitude:          7.751210803377673,
			Address:            "Lauterstraße 51, 67659 Kaiserslautern",
			WaterSource:        "St. Georgsbrunnen",
			OpeningTimes:       "Mon - Son / 00:00 AM - 12:59 PM",
			Type:               "smart",
			OfferedWaterTypes:  "tap",
			RefillStationImage: &gartenschau_picture,
		},
		{
			Name:               "Wochenmarkt",
			Description:        "Frische Lebensmittel von Obst über Käse, Gemüse und Wurstwaren, bis hin zu Fisch und Backwaren, sowie Blumen und Pflanzen",
			Latitude:           49.44025,
			Longitude:          7.75878,
			Address:            "Königstraße 68, 67655 Kaiserslautern",
			WaterSource:        "Stadtwerke",
			OpeningTimes:       "Do / 07:00 AM - 1:59 PM",
			Active:             NullBool{Valid: true, Bool: false},
			Type:               "smart",
			OfferedWaterTypes:  "tap",
			RefillStationImage: &wochenmarkt_picture,
		},
	}

	if err := db.Create(&refillStations).Error; err != nil {
		log.Fatalf("failed to create refill station: %v", err)
	}

	log.Print("Created refill stations successfully")

	return db
}

func CreateRefillStationReviews(db *gorm.DB) *gorm.DB {
	reviews := []RefillStationReview{
		{
			StationID:     1,
			UserID:        1,
			Cleanness:     4,
			Accessibility: 5,
			WaterQuality:  3,
		},
		{
			StationID:     1,
			UserID:        2,
			Cleanness:     5,
			Accessibility: 4,
			WaterQuality:  5,
		},
		{
			StationID:     2,
			UserID:        1,
			Cleanness:     3,
			Accessibility: 3,
			WaterQuality:  4,
		},
		{
			StationID:     2,
			UserID:        5,
			Cleanness:     2,
			Accessibility: 2,
			WaterQuality:  5,
		},
		{
			StationID:     3,
			UserID:        1,
			Cleanness:     3,
			Accessibility: 3,
			WaterQuality:  4,
		},
		{
			StationID:     4,
			UserID:        5,
			Cleanness:     5,
			Accessibility: 4,
			WaterQuality:  5,
		},
		{
			StationID:     4,
			UserID:        3,
			Cleanness:     1,
			Accessibility: 1,
			WaterQuality:  1,
		},
	}

	if err := db.Create(&reviews).Error; err != nil {
		log.Fatalf("failed to create refill station reviews: %v", err)
	}

	log.Print("Created refill station reviews successfully")

	return db
}

func CreateRefillStationProblems(db *gorm.DB) *gorm.DB {
	problem1 := ImageToBase64("./images/problems/dripping.jpg")
	problem2 := ImageToBase64("./images/problems/broken.jpg")
	problem3 := ImageToBase64("./images/problems/dirty.jpg")
	problems := []RefillStationProblem{
		{
			StationID:                 1,
			Title:                     "Undichte Wasserhähne",
			Description:               "Der Wasserhahn an der Nachfüllstation tropft kontinuierlich.",
			Status:                    "OPEN",
			RefillStationProblemImage: &problem1,
		},
		{
			StationID:                 2,
			Title:                     "Beschädigter Spender",
			Description:               "Der Wasserspender an der Nachfüllstation ist beschädigt und gibt kein Wasser ordnungsgemäß ab.",
			Status:                    "INPROGRESS",
			RefillStationProblemImage: &problem2,
		},
		{
			StationID:                 3,
			Title:                     "Wasserkontamination",
			Description:               "Benutzer meldeten Probleme mit Wasserkontamination an dieser Nachfüllstation.",
			Status:                    "SOLVED",
			RefillStationProblemImage: &problem3,
		},
	}

	if err := db.Create(&problems).Error; err != nil {
		log.Fatalf("failed to create refill station problems: %v", err)
	}

	log.Print("Created refill station problems successfully")

	return db
}

func CreateWaterTransactions(db *gorm.DB) *gorm.DB {
	transactions := []WaterTransaction{}

	for i := 0; i < 50; i++ {
		randomStationID := uint(rand.Intn(3) + 1) // Random station ID between 1 and 3
		randomBottleID := uint(rand.Intn(2) + 1)  // Random Bottle ID between 1 and 2
		randomVolume := rand.Intn(1500) + 500     // Random volume between 500 and 2000
		randomWaterType := []string{"TAP", "MINERAL"}[rand.Intn(2)]
		userID := uint(1)

		transaction := WaterTransaction{
			StationID: randomStationID,
			BottleID:  &randomBottleID,
			UserID:    &userID,
			Volume:    randomVolume,
			WaterType: randomWaterType,
		}

		transactions = append(transactions, transaction)
	}

	if err := db.Create(&transactions).Error; err != nil {
		log.Fatalf("failed to create water transactions: %v", err)
	}

	log.Print("Created water transactions successfully")

	return db
}

func CreateLikes(db *gorm.DB) *gorm.DB {
	likes := []Like{
		{
			StationID: 1,
			UserID:    1,
		},
		{
			StationID: 2,
			UserID:    2,
		},
		{
			StationID: 1,
			UserID:    2,
		},
		{
			StationID: 2,
			UserID:    3,
		},
		{
			StationID: 1,
			UserID:    3,
		},
		{
			StationID: 2,
			UserID:    5,
		},
		{
			StationID: 1,
			UserID:    5,
		},
		{
			StationID: 2,
			UserID:    4,
		},
	}

	if err := db.Create(&likes).Error; err != nil {
		log.Fatalf("failed to create likes: %v", err)
	}

	log.Print("Created likes successfully")

	return db
}
