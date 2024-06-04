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
	db = CreateConsumerTests(db)
	db = CreateConsumerTestsQuestions(db)
	db = CreateConsumerTestsAnswers(db)
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

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("failed to create user: %v", err)
		}
	}

	log.Print("Created users successfully")

	return db
}

func CreateBottles(db *gorm.DB) *gorm.DB {
	bottles := []Bottle{
		{
			UserID:     1,
			NFCID:      "13:E0:0B:35",
			FillVolume: 500,
			WaterType:  "tap",
			Title:      "First Bottle",
		},
		{
			UserID:     1,
			NFCID:      "",
			FillVolume: 1000,
			WaterType:  "mineral",
			Title:      "Second Bottle",
		},
	}

	for _, bottle := range bottles {
		if err := db.Create(&bottle).Error; err != nil {
			log.Fatalf("failed to create bottle: %v", err)
		}
	}

	log.Print("Created bottles successfully")

	return db
}

func CreateConsumerTests(db *gorm.DB) *gorm.DB {
	// Test are only a logical container, no parameters
	tests := []ConsumerTest{
		{},
		{},
	}

	for _, test := range tests {
		if err := db.Create(&test).Error; err != nil {
			log.Fatalf("failed to create consumer test: %v", err)
		}
	}

	log.Print("Created consumer tests successfully")

	return db
}

func CreateConsumerTestsQuestions(db *gorm.DB) *gorm.DB {
	// Test are only a logical container, no parameters
	testQuestions := []ConsumerTestQuestion{
		{
			Text:   "Is this a Question?",
			TestID: 1,
		},
		{
			Text:   "An this one?",
			TestID: 1,
		},
		{
			Text:   "How many?",
			TestID: 2,
		},
	}

	for _, testQuestion := range testQuestions {
		if err := db.Create(&testQuestion).Error; err != nil {
			log.Fatalf("failed to create consumer test question: %v", err)
		}
	}

	log.Print("Created consumer tests questions successfully")

	return db
}

func CreateConsumerTestsAnswers(db *gorm.DB) *gorm.DB {
	// Test are only a logical container, no parameters
	testAnswers := []ConsumerTestAnswer{
		{
			UserID:     1,
			QuestionID: 1,
			Answer:     5.0,
		},
		{
			UserID:     1,
			QuestionID: 2,
			Answer:     1.0,
		},
		{
			UserID:     2,
			QuestionID: 3,
			Answer:     4.20,
		},
	}

	for _, testAnswer := range testAnswers {
		if err := db.Create(&testAnswer).Error; err != nil {
			log.Fatalf("failed to create consumer test answer: %v", err)
		}
	}

	log.Print("Created consumer tests answers successfully")

	return db
}

func CreateRefillStations(db *gorm.DB) *gorm.DB {
	refillStations := []RefillStation{
		{
			Name:              "Rewe Stations",
			Description:       "Wasserhahn REWE",
			Latitude:          49.44490159879211,
			Longitude:         7.767478778642334,
			Address:           "Fruchthallstraße 29, 67655 Kaiserslautern",
			WaterSource:       "Spitzrainbrunnen",
			OpeningTimes:      "Mon - Sam / 7:00 AM - 10:00 PM",
			Type:              "Manual",
			OfferedWaterTypes: "Tap",
		},
		{
			Name:              "Stadtpark KL",
			Description:       "Irgendwas mit einem Schnitzeljagd?",
			Latitude:          49.437551349217266,
			Longitude:         7.761465081072085,
			Address:           "Trippstadter Str. 2, 67663 Kaiserslautern",
			WaterSource:       "Spitzrainbrunnen",
			OpeningTimes:      "Mon - Son / 00:00 AM - 12:59 PM",
			Type:              "Smart",
			OfferedWaterTypes: "Mineral, Tap",
		},
		{
			Name:              "Gartenschau KL",
			Description:       "Die Gartenschau Kaiserslautern ist ein atemberaubendes jährliches Ereignis, das die Schönheit der Natur und die Freude am Gartenbau feiert.",
			Latitude:          49.44694255672088,
			Longitude:         7.751210803377673,
			Address:           "Lauterstraße 51, 67659 Kaiserslautern",
			WaterSource:       "St. Georgsbrunnen",
			OpeningTimes:      "Mon - Son / 00:00 AM - 12:59 PM",
			Type:              "Smart",
			OfferedWaterTypes: "Tap",
		},
		{
			Name:              "Wochenmarkt",
			Description:       "Frische Lebensmittel von Obst über Käse, Gemüse und Wurstwaren, bis hin zu Fisch und Backwaren, sowie Blumen und Pflanzen",
			Latitude:          49.44025,
			Longitude:         7.75878,
			Address:           "Königstraße 68 67655 Kaiserslautern",
			WaterSource:       "Stadtwerke",
			OpeningTimes:      "Do / 07:00 AM - 13:59 PM",
			Active:            false,
			Type:              "Smart",
			OfferedWaterTypes: "Tap",
		},
	}

	for _, station := range refillStations {
		if err := db.Create(&station).Error; err != nil {
			log.Fatalf("failed to create refill station: %v", err)
		}
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
	}

	for _, review := range reviews {
		if err := db.Create(&review).Error; err != nil {
			log.Fatalf("failed to create refill station review: %v", err)
		}
	}

	log.Print("Created refill station reviews successfully")

	return db
}

func CreateRefillStationProblems(db *gorm.DB) *gorm.DB {
	problems := []RefillStationProblem{
		{
			StationID:   1,
			Title:       "Undichte Wasserhähne",
			Description: "Der Wasserhahn an der Nachfüllstation tropft kontinuierlich.",
			Status:      "Active",
		},
		{
			StationID:   2,
			Title:       "Beschädigter Spender",
			Description: "Der Wasserspender an der Nachfüllstation ist beschädigt und gibt kein Wasser ordnungsgemäß ab.",
			Status:      "Active",
		},
		{
			StationID:   3,
			Title:       "Wasserkontamination",
			Description: "Benutzer meldeten Probleme mit Wasserkontamination an dieser Nachfüllstation.",
			Status:      "Active",
		},
	}

	for _, problem := range problems {
		if err := db.Create(&problem).Error; err != nil {
			log.Fatalf("failed to create refill station problem: %v", err)
		}
	}

	log.Print("Created refill station problems successfully")

	return db
}

func CreateWaterTransactions(db *gorm.DB) *gorm.DB {
	transactions := []WaterTransaction{}

	for i := 0; i < 10; i++ {
		randomStationID := uint(rand.Intn(3) + 1) // Random station ID between 1 and 3
		randomBottleID := uint(rand.Intn(2) + 1)  // Random Bottle ID between 1 and 2
		randomVolume := rand.Intn(1500) + 500     // Random volume between 500 and 2000
		randomWaterType := []string{"Tap", "Mineral"}[rand.Intn(2)]
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

	for _, transaction := range transactions {
		if err := db.Create(&transaction).Error; err != nil {
			log.Fatalf("failed to create water transaction: %v", err)
		}
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

	for _, like := range likes {
		if err := db.Create(&like).Error; err != nil {
			log.Fatalf("failed to create like: %v", err)
		}
	}

	log.Print("Created likes successfully")

	return db
}
