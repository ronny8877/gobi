package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jaswdr/faker/v2"
)

type API struct {
	Method       string                 `json:"method"`
	Path         string                 `json:"path"`
	Latency      *int                   `json:"latency"`
	FailRate     *float64               `json:"failRate"`
	Validate     *Validate              `json:"validate"`
	ResponseType *string                `json:"responseType"`
	Response     map[string]interface{} `json:"response"`
}

type Validate struct {
	Query  []string               `json:"query"`
	Body   map[string]interface{} `json:"body"`
	Params map[string]interface{} `json:"params"`
}

type AppInput struct {
	Config AppConfig `json:"config"`
	APIs   []API     `json:"api"`
}

type AppConfig struct {
	Prefix      string   `json:"prefix"`
	HealthCheck *string  `json:"healthCheck"`
	Port        int      `json:"port"`
	Latency     *int     `json:"latency"`
	Logging     *bool    `json:"logging"`
	FailRate    *float32 `json:"failRate"`
	Timeout     *int32   `json:"timeout"`
	ApiKey      *string  `json:"apiKey"`
}

var defaultLatency = 0
var defaultFailRate = float32(0.4)
var defaultTimeout = int32(5000)
var healthCheck = "/health"
var defaultPort = 8080
var defaultLogging = false

var defaultConfig = AppConfig{
	Prefix:      "/api",
	HealthCheck: &healthCheck, // Optional field
	Port:        8080,
	Logging:     &defaultLogging,  // Optional field
	Latency:     &defaultLatency,  // Optional field
	FailRate:    &defaultFailRate, // Optional field
	Timeout:     &defaultTimeout,  // Optional field
	ApiKey:      nil,              // Optional field
}

type App struct {
	Config AppConfig `json:"config"`
	APIs   []API     `json:"api"`
}

func loadConfig(app *AppInput) {
	// Read the file
	file, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// Check if the file is empty
	if len(file) == 0 {
		fmt.Println(string(file))
		log.Fatal("Error: input.json is empty")
	}

	// Create a new instance of AppInput
	//This is a workaround as i had issues where if the key was removed from input the value was not getting updated
	var newApp AppInput

	// Unmarshal the JSON into the new instance
	err = json.Unmarshal(file, &newApp)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		log.Fatal(err)
	}

	// Copy the values from the new instance to the existing app instance
	*app = newApp

	// Set default values if necessary
	if int(app.Config.Port) == 0 {
		app.Config.Port = defaultPort
	}

	if app.Config.Prefix == "" {
		app.Config.Prefix = defaultConfig.Prefix
	}

	if app.Config.HealthCheck == nil {
		app.Config.HealthCheck = defaultConfig.HealthCheck
	}

	if app.Config.Latency == nil {
		app.Config.Latency = defaultConfig.Latency
	}

	fmt.Println("Configuration Loaded: ")
}

func watchConfigFile(filename string, app *AppInput) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Config file modified, reloading...")
					time.Sleep(100 * time.Millisecond) // Add a small delay before reloading as the it was crashing without it
					loadConfig(app)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal("Error adding file to watcher:", err)
	}
	<-done
}

// TODO Implement the timeout function
// TODO Implement auth	 using API key cookie or bearer token
func main() {
	var app AppInput
	loadConfig(&app)
	go watchConfigFile("input.json", &app)
	//starting a server with the port
	go startServer(&app)
	fmt.Printf("Server is running on http://localhost:%d%s/\n", app.Config.Port, app.Config.Prefix)
	http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)
	// Wait for interrupt signal to gracefully shutdown the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down...")
}

// func validateQueryParams(query map[string]interface{}) bool {
// 	return true
// }

func startServer(app *AppInput) {
	http.HandleFunc(fmt.Sprintf("%s/", app.Config.Prefix), func(w http.ResponseWriter, r *http.Request) {
		// if logging is enabled, log the request
		if app.Config.Logging != nil && *app.Config.Logging {
			fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		}

		//If global latency is provided, sleep for that time
		if app.Config.Latency != nil || *app.Config.Latency != 0 {
			fmt.Printf("Adding Latency %d ms\n", *app.Config.Latency)
			time.Sleep(time.Duration(*app.Config.Latency) * time.Millisecond)
		}
		//If global fail rate is provided, return an error if the random number is less than the fail rate
		if app.Config.FailRate != nil {
			random := rand.Float64()
			if random < float64(*app.Config.FailRate) {
				http.Error(w, `{"error":"Internal server Error"}`, http.StatusInternalServerError)
				return
			}
		}

		var path = r.URL.Path
		//check if the path is in the APIs
		var found = false

		for _, api := range app.APIs {
			if strings.HasPrefix(path, app.Config.Prefix) && api.Path == strings.TrimPrefix(path, app.Config.Prefix) && api.Method == r.Method {
				found = true
				//check if the latency is provided
				if api.Latency != nil && *api.Latency != 0 {
					//sleep for the latency
					fmt.Printf("Adding Latency %d ms\n", *api.Latency)
					time.Sleep(time.Duration(*api.Latency) * time.Millisecond)

				}
				//check if the fail rate is provided
				if api.FailRate != nil && *api.FailRate != 0 {
					//generate a random number between 0 and 1
					//if the number is less than the fail rate, return an error
					random := rand.Float64()
					if random < *api.FailRate {
						http.Error(w, `{"error":"Internal server Error"}`, http.StatusInternalServerError)
						return
					}
				}

				//check if the validate is provided
				if api.Validate != nil {
					//check if the query parameters are provided
					if api.Validate.Query != nil {
						fmt.Println("Query Params: ", api.Validate.Query)
						fmt.Println("Query Params: ", r.URL.Query())
						_, err := validateQueryParam(api.Validate.Query, r.URL.Query())
						if err != nil {
							http.Error(w, `{"error": "Invalid Query Params"}`, http.StatusBadRequest)
							return
						}
						if app.Config.Logging != nil && *app.Config.Logging {
							fmt.Println("Valid Query Params")
						}

					}
					//check if the body is provided
					if api.Validate.Body != nil {
						//check if the body is correct
					}
					//check if the params are provided
					if api.Validate.Params != nil {
						//check if the params are correct.
					}
				}
				//return the response
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Powered-By", "Gobi")
				// fmt.Println("Response Type: ", string(*api.ResponseType))
				if api.ResponseType != nil {
					respType, arg := parseValueBwBrackets(*api.ResponseType)
					response := []interface{}{}
					if respType == "Array" {
						arrLenInt, _ := strconv.Atoi(arg)
						for i := 0; i < arrLenInt; i++ {
							response = append(response, ResponseBuilder(api.Response))
						}
						json.NewEncoder(w).Encode(response)
						break
					}
				} else {
					response := ResponseBuilder(api.Response)
					json.NewEncoder(w).Encode(response)
					break
				}
				// fmt.Println("Response: ", response)
				//Send the response and stop the loop

			}
		}
		if !found {
			http.Error(w, `{"error": "API not found"}`, http.StatusNotFound)
		}

	})
}

func parseValueBwBrackets(value string) (string, string) {
	var name = strings.Split(value, "(")[0]
	var args = strings.Split(value, "(")[1]
	args = strings.TrimRight(args, ")")
	return name, args
}

func validateQueryParam(query []string, rQuery map[string][]string) (bool, error) {
	//  if both dose not have the same length then return false
	if len(query) != len(rQuery) {
		return false, fmt.Errorf("mismatch in query params")
	}
	for _, q := range query {
		if _, ok := rQuery[q]; !ok {
			return false, fmt.Errorf("mismatch in query params")
		}
	}
	return true, nil
}
func ResponseBuilder(rawData map[string]interface{}) map[string]interface{} {
	response := make(map[string]interface{})
	fake := faker.New()
	funcMap := map[string]func(*string) interface{}{
		"Address": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.Address().Address()
			case *args == "city":
				return fake.Address().City()
			case *args == "country":
				return fake.Address().Country()
			case *args == "state":
				return fake.Address().State()
			case *args == "countryCode":
				return fake.Address().CountryCode()
			case *args == "zip":
				return fake.Address().PostCode()
			case *args == "countryAbbr":
				return fake.Address().CountryAbbr()
			case *args == "street":
				return fake.Address().StreetName()
			case *args == "stateAbbr":
				return fake.Address().StateAbbr()
			case *args == "secondary":
				return fake.Address().SecondaryAddress()
			case *args == "latitude":
				return fake.Address().Latitude()
			case *args == "longitude":
				return fake.Address().Longitude()

			default:
				return fake.Address().Address()
			}
		},
		"Company": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.Company().Name()
			case *args == "suffix":
				return fake.Company().Suffix()
			case *args == "jobTitle":
				return fake.Company().JobTitle()
			case *args == "bs":
				return fake.Company().BS()
			case *args == "catchPhrase":
				return fake.Company().CatchPhrase()
			case *args == "ein":
				return fake.Company().EIN()
			case *args == "mail":
				return fake.Internet().CompanyEmail()
			default:
				return fake.Company().Name()
			}
		},
		"Internet": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.Internet().URL()
			case *args == "ip":
				return fake.Internet().Ipv4()
			case *args == "ipv6":
				return fake.Internet().Ipv6()
			case *args == "mac":
				return fake.Internet().MacAddress()
			case *args == "httpMethod":
				return fake.Internet().HTTPMethod()
			case *args == "domain":
				return fake.Internet().Domain()
			case *args == "tld":
				return fake.Internet().TLD()
			case *args == "slug":
				return fake.Internet().Slug()
			case *args == "statusCode":
				return fake.Internet().StatusCode()
			case *args == "freeEmail":
				return fake.Internet().FreeEmail()
			case *args == "safeEmail":
				return fake.Internet().SafeEmail()
			case *args == "statusCodeMessage":
				return fake.Internet().StatusCodeMessage()
			default:
				return fake.Internet().URL()
			}
		},
		"User": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.Person().Name()
			case *args == "email":
				return fake.Internet().Email()
			case *args == "firstName":
				return fake.Person().FirstName()
			case *args == "lastName":
				return fake.Person().LastName()
			case *args == "mFirstName":
				return fake.Person().FirstNameMale()
			case *args == "fFirstName":
				return fake.Person().FirstNameFemale()
			case *args == "fTitle":
				return fake.Person().TitleFemale()
			case *args == "mTitle":
				return fake.Person().TitleMale()
			case *args == "phone":
				return fake.Phone().Number()
			case *args == "password":
				return fake.Internet().Password()
			case *args == "userName":
				return fake.Internet().User()
			case *args == "title":
				return fake.Person().Title()
			case *args == "gender":
				return fake.Person().Gender()
			case *args == "ssn":
				return fake.Person().SSN()
			case *args == "bio":
				return fake.Lorem().Paragraph(1)
			case *args == "gamerTag":
				return fake.Gamer().Tag()
			case *args == "uuid":
				return fake.UUID().V4()
			case *args == "sqlId":
				// return an integer between 1 and 100000
				return fake.Int64Between(1, 100000)

			case *args == "birthday":
				rand.Seed(time.Now().UnixNano())

				// Generate random values within a reasonable range
				years := rand.Intn(21) - 10  // Random number between -10 and 10
				months := rand.Intn(25) - 12 // Random number between -12 and 12
				days := rand.Intn(61) - 30   // Random number between -30 and 30

				// Add the random values to the current date
				randomDate := time.Now().AddDate(years, months, days)

				// Return the random date
				return randomDate
			case *args == "image":
				return fmt.Sprintf("https://randomuser.me/api/portraits/med/%d.jpg", fake.RandomDigit())
			default:
				return fake.Person().Name()
			}
		},
		"Age": func(args *string) interface{} { return rand.Intn(100) },
		"Finance": func(args *string) interface{} {
			switch {
			case args == nil:
				return "Finance Needs an argument"
			case *args == "creditCard":
				return fake.Payment().CreditCardNumber()
			case *args == "cardType":
				return fake.Payment().CreditCardType()
			case *args == "cardExpirationDate":
				return fake.Payment().CreditCardExpirationDateString()
			case *args == "iban":
				return fake.Payment().Iban()
			case *args == "currency":
				return fake.Currency().Currency()
			case *args == "currencyCode":
				return fake.Currency().Code()
			case *args == "currencyAndCode":
				curr, code := fake.Currency().CurrencyAndCode()
				return fmt.Sprintf("%s (%s)", curr, code)
			case *args == "amount":
				return fake.Int64Between(1000, 10000000)
			case *args == "amountWithCurrency":
				return fmt.Sprintf("%d %s", fake.Int64Between(1000, 100000), fake.Currency().Code())
			case *args == "btcAddress":
				return fake.Crypto().BitcoinAddress()
			case *args == "ethAddress":
				return fake.Crypto().EtheriumAddress()
			default:
				return "Finance Needs an argument"
			}

		},
		"Lorem": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.Lorem().Text(20)
			case *args == "word":
				return fake.Lorem().Word()
			case *args == "sentence":
				return fake.Lorem().Sentence(20)
			case *args == "paragraph":
				return fake.Lorem().Paragraph(1)
			case *args == "paragraphs":
				return fake.Lorem().Paragraphs(2)
			default:
				return fake.Lorem().Sentence(20)
			}

		},
		"Color": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.Color().ColorName()
			case *args == "hex":
				return fake.Color().Hex()
			case *args == "rgb":
				return fake.Color().RGB()
			case *args == "css":
				return fake.Color().CSS()
			case *args == "rgba":
				return fake.Color().RGBAsArray()
			case *args == "safe":
				return fake.Color().SafeColorName()
			default:
				return fake.Color().ColorName()
			}
		},
		"Int": func(args *string) interface{} {
			// If no arguments are provided, return a random integer between 0 and 100
			if args == nil {
				return rand.Intn(10000)
			}
			//if the argument is passed then convert it to int
			if val, err := strconv.Atoi(*args); err == nil {
				return rand.Intn(val)
			}
			return fake.RandomDigit()
		},
		"Float": func(args *string) interface{} {
			// If no arguments are provided, return a random float between 0 and 100
			if args == nil {
				return fake.Float64(3, 1, 1000)
			}
			//if the argument is passed then convert it to float
			if val, err := strconv.Atoi(*args); err == nil {
				return fake.Float64(3, 1, val)
			}
			return fake.Float64(3, 1, 1000)
		},
		"Bool": func(args *string) interface{} {
			return fake.BoolWithChance(50)
		},
		"Array": func(args *string) interface{} {
			//TODO : implement the SO  the array length can be passed
			return []interface{}{fake.Person().Name(), fake.Person().Name(), fake.Person().Name()}
		},
		"Language": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.Language().Language()
			case *args == "abbr":
				return fake.Language().LanguageAbbr()
			case *args == "programming":
				return fake.Language().ProgrammingLanguage()
			default:
				return fake.Language().Language()
			}
		},
		"App": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.App().Name()
			case *args == "version":
				return fake.App().Version()
			case *args == "platform":
				// return a random platform from iOS, Android, web, Mac, Windows, Linux
				platforms := []string{"iOS", "Android", "web", "Mac", "Windows", "Linux"}
				return platforms[rand.Intn(len(platforms))]
			default:
				return fake.App().Name()
			}
		},
		"Vehicle": func(args *string) interface{} {
			switch {
			case args == nil:
				return fake.Car().Category()
			case *args == "brand":
				return fake.Car().Maker()
			case *args == "transmission":
				return fake.Car().TransmissionGear()
			case *args == "plate":
				return fake.Car().Plate()
			case *args == "model":
				return fake.Car().Model()
			case *args == "type":
				return fake.Car().FuelType()
			default:
				return fake.Car().Category()
			}
		},
		"Time": func(args *string) interface{} {
			switch {
			case args == nil:
				return time.Now().Format(time.RFC1123Z)
			case *args == "unix":
				return time.Now().Unix()
			case *args == "unixNano":
				return time.Now().UnixNano()
			case *args == "iso":
				return fake.Time().ISO8601(time.Now())
			case *args == "amPm":
				return fake.Time().AmPm()
			case *args == "month":
				return fake.Time().Month()
			case *args == "day":
				return fake.Time().DayOfWeek()
			case *args == "ansi":
				return fake.Time().ANSIC(time.Now())
			case *args == "monthName":
				return fake.Time().MonthName()
			case *args == "timezone":
				return fake.Time().Timezone()
			default:
				return time.Now().Format(time.RFC1123Z)
			}
		},
		"Image": func(args *string) interface{} {
			//TODO : Implement so type of dicebear can be passed along with seed
			return fmt.Sprintf("https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=%s", *args)
			//TODO : Implement so Images form other sources can be passed
		},
		"Success": func(args *string) interface{} {
			if args == nil {
				return "Success"
			}
			return *args
		},
		"Json": func(args *string) interface{} {
			return fake.Map()
		},
	}

	// Usage
	response = processMap(rawData, funcMap)
	return response
}

func processMap(data map[string]interface{}, funcMap map[string]func(*string) interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case string:
			//Getting the value In bw the brackets also removing the brackets
			funcName, args := parseValueBwBrackets(v)
			if fn, ok := funcMap[funcName]; ok {
				result[key] = fn(&args)
			} else {
				result[key] = v
			}
		case map[string]interface{}:
			result[key] = processMap(v, funcMap)
		default:
			result[key] = v
		}
	}
	return result
}
