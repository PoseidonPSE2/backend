package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type data struct {
	ID        string `json:"id"`
	Ml        string `json:"ml,omitempty"`
	WaterType string `json:"waterType,omitempty"`
}

type databaseEntry struct {
	Ml        string
	WaterType string
}

var database_temp = make(map[string]databaseEntry)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var request_data data
	err := json.NewDecoder(r.Body).Decode(&request_data)
	if err != nil {
		log.Printf("Fehler beim Dekodieren der Anfrage: %v", err)
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	log.Printf("Anfrage erhalten für ID: %s", request_data.ID)

	if entry, ok := database_temp[request_data.ID]; ok {
		var dataResponse = data{
			ID:        request_data.ID,
			Ml:        entry.Ml,
			WaterType: entry.WaterType,
		}

		response, err := json.Marshal(dataResponse)
		if err != nil {
			log.Printf("Fehler beim Marshaling der Antwort: %v", err)
			http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
			return
		}

		log.Printf("Antwort gesendet: %s", response)
		fmt.Fprintf(w, "%s\n", response)
	} else {
		log.Printf("ID nicht gefunden: %s", request_data.ID)
		http.Error(w, "ID nicht gefunden", http.StatusNotFound)
	}
}

func addData(w http.ResponseWriter, r *http.Request) {
	var data data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Fehler beim Dekodieren der Anfrage: %v", err)
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	database_temp[data.ID] = databaseEntry{Ml: data.Ml, WaterType: data.WaterType}
	response := fmt.Sprintf("ID %s mit Füllstand %s ml und Wasserart %s wurde hinzugefügt\n", data.ID, data.Ml, data.WaterType)
	log.Printf("Daten hinzugefügt: %s", response)
	fmt.Fprint(w, response)
}

func addDataManually(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ml := r.URL.Query().Get("ml")
	waterType := r.URL.Query().Get("waterType")

	if id == "" || ml == "" || waterType == "" {
		http.Error(w, "Bitte geben Sie id, ml und Wasserart an", http.StatusBadRequest)
		return
	}

	database_temp[id] = databaseEntry{Ml: ml, WaterType: waterType}
	response := fmt.Sprintf("ID %s mit Füllstand %s ml und Wasserart %s wurde hinzugefügt\n", id, ml, waterType)
	log.Printf("Daten hinzugefügt: %s", response)
	fmt.Fprint(w, response)
}

func addInitialData() {
	database_temp["13:8E:BD:0C"] = databaseEntry{Ml: "500", WaterType: "still"}
	database_temp["13:E0:0B:35"] = databaseEntry{Ml: "100", WaterType: "sprudel"}
}
func getAllEntries(w http.ResponseWriter, r *http.Request) {

	var allEntries []data

	for id, entry := range database_temp {
		allEntries = append(allEntries, data{
			ID:        id,
			Ml:        entry.Ml,
			WaterType: entry.WaterType,
		})
	}

	response, err := json.Marshal(allEntries)
	if err != nil {
		log.Printf("Fehler beim Marshaling der Antwort: %v", err)
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	log.Printf("Antwort gesendet: %s", response)
	fmt.Fprintf(w, "%s\n", response)
}

// This was main, but needed to renamed it, because there is already a main in endpoint_server.go
func temp() {
	addInitialData()
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/add", addData)
	http.HandleFunc("/addManually", addDataManually)
	http.HandleFunc("/getAllEntries", getAllEntries)
	log.Println("Server läuft und hört auf Port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server konnte nicht gestartet werden: %v", err)
	}
}
