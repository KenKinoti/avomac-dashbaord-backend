package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var router *chi.Mux
var db *sql.DB

// Company details struct
type Company struct {
	ID             int
	FakeName       string `json:"fake_company_name"`
	Description    string `json:"description"`
	Tagline        string `json:"tagline"`
	CompanyEmail   string `json:"company_email"`
	BusinessNumber string `json:"business_number"`
	Restricted     string `json:"restricted"`
}

type Delivery struct {
	ID 			int
	Supplier 	string `json:"supplier"`
	Source 		string `json:"source"`
	Amount  	string `json:"amount"`
	Price 		string `json:"price"`
	Comments 	string `json:"comments"`
	CreatedAt	string `json:"created_at"`
}

type Supplier struct {
	ID 			int
	Company 	string `json:"company"`
	VehicleType	string `json:"vehicle_type"`
	Phone	  	string `json:"phone"`
	DriverFName	string `json:"driver_f_name"`
	DriverLName	string `json:"driver_l_name"`
	Address		string `json:"address"`
	County		string `json:"county"`
	Country		string `json:"country"`
}

type SupplierName struct {
	ID 			int
	Company 	string `json:"company"`
}

func init() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	var err error

	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV")
	}

	dbUser := os.Getenv("dbUser")
	dbName := os.Getenv("dbName")
	dbPass := os.Getenv("dbPass")
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")

	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPass, dbHost, dbPort, dbName)
	
	db, err = sql.Open("mysql", dbSource)
	catch(err)	
}

func routers() *chi.Mux {

	router.Get("/", ping)
	router.Get("/deliveries", AllDeliveries)
	router.Get("/suppliers", AllSuppliers)

	router.Get("/suppliers-name", SuppliersNames)

	router.Post("/create-delivery", CreateDelivery)
	router.Post("/create-supplier", CreateSupplier)
	
	router.Get("/company/{id}", DetailCompany)

	return router
}

//-------------- API ENDPOINT ------------------//
func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondwithJSON(w, http.StatusOK, map[string]string{"message": "Welcome to AVOMAC API status check"})
}

//-------------- API ENDPOINT ------------------//

// AllCompanys get all the companys
func AllDeliveries(w http.ResponseWriter, r *http.Request) {

	errors := []error{}
	payload := []Delivery{}

	rows, err := db.Query("Select * From delivery")
	catch(err)

	defer rows.Close()

	for rows.Next() {
		data := Delivery{}

		er := rows.Scan(
			&data.ID,
			&data.Supplier,
			&data.Source,
			&data.Amount,
			&data.Price,
			&data.Comments,
			&data.CreatedAt,
		)

		if er != nil {
			errors = append(errors, er)
		}
		payload = append(payload, data)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	respondwithJSON(w, http.StatusOK, payload)
}

func AllSuppliers(w http.ResponseWriter, r *http.Request) {
	errors := []error{}
	payload := []Supplier{}

	rows, err := db.Query("Select * From supplier")
	catch(err)

	defer rows.Close()

	for rows.Next() {
		data := Supplier{}

		er := rows.Scan(
			&data.ID,
			&data.Company,
			&data.VehicleType,
			&data.Phone,
			&data.DriverFName,
			&data.DriverLName,
			&data.Address,
			&data.County,
			&data.Country,
		)

		if er != nil {
			errors = append(errors, er)
		}
		payload = append(payload, data)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	respondwithJSON(w, http.StatusOK, payload)
}

func SuppliersNames(w http.ResponseWriter, r *http.Request) {
	errors := []error{}
	payload := []SupplierName{}

	rows, err := db.Query("Select id, company From supplier")
	catch(err)

	defer rows.Close()

	for rows.Next() {
		data := SupplierName{}

		er := rows.Scan(
			&data.ID,
			&data.Company,
		)

		if er != nil {
			errors = append(errors, er)
		}
		payload = append(payload, data)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	respondwithJSON(w, http.StatusOK, payload)
}


func CreateDelivery(w http.ResponseWriter, r *http.Request) {

	deliveryStruct := Delivery{}
	json.NewDecoder(r.Body).Decode(&deliveryStruct)

	query, err := db.Prepare("Insert delivery SET supplier=?, source=?, amount=?, price=?, comments=?, created_at=?")
	catch(err)

	// Define values for the insert
	_, er := query.Exec(deliveryStruct.Supplier, deliveryStruct.Source, deliveryStruct.Amount, deliveryStruct.Price, deliveryStruct.Comments, deliveryStruct.CreatedAt)
	catch(er)

	defer query.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Printf("Delivery doesn't exist... Creating ")

		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Delivery successfully created", "status": "true"})
	}
}

func CreateSupplier(w http.ResponseWriter, r *http.Request) {

	supplyStruct := Supplier{}
	json.NewDecoder(r.Body).Decode(&supplyStruct)

	query, err := db.Prepare("Insert supplier SET company=?, vehicle_type=?, phone=?, driver_f_name=?, driver_l_name=?, address=?, county=?, country=?")
	catch(err)


	// Define values for the insert
	_, er := query.Exec(supplyStruct.Company, supplyStruct.VehicleType, supplyStruct.Phone, supplyStruct.DriverFName, supplyStruct.DriverLName, supplyStruct.Address, supplyStruct.County, supplyStruct.Country)
	catch(er)

	defer query.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Printf("Supplier doesn't exist... Creating ")

		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Supplier successfully created", "status": "true"})
	}
}


// DetailCompany get specific company details
func DetailCompany(w http.ResponseWriter, r *http.Request) {
	payload := Company{}
	// inda := IndustryString{}

	id := chi.URLParam(r, "id")
	// fmt.Println(id)

	row := db.QueryRow("Select * From faux_id_fake_companies where `business_number`=?", id)

	err := row.Scan(
		&payload.ID,
		&payload.FakeName,
		&payload.Description,
		&payload.Tagline,
		&payload.CompanyEmail,
		&payload.BusinessNumber,
		&payload.Restricted,
	)

	// payload.Industry = strings.Split(inda.Industry, ",")

	if err != nil {
		respondWithError(w, http.StatusNotFound, "no rows in result set")
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	respondwithJSON(w, http.StatusOK, payload)
}


func main() {
	routers()
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading port configuration from ENV")
	}
	
	port := os.Getenv("port")

	if port == "" {
        port = "8090"
	}

	http.ListenAndServe(":"+port, Logger())

}
