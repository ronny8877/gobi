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

	"github.com/jaswdr/faker/v2"
)

type API struct {
	Method   string                 `json:"method"`
	Path     string                 `json:"path"`
	Latency  *int                   `json:"latency"`
	FailRate *float64               `json:"failRate"`
	Validate *Validate              `json:"validate"`
	Response map[string]interface{} `json:"response"`
}

type Validate struct {
	Query  map[string]interface{} `json:"query"`
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
	FailRate    *float32 `json:"failRate"`
	Timeout     *int32   `json:"timeout"`
	ApiKey      *string  `json:"apiKey"`
}

var defaultLatency = 0
var defaultFailRate = float32(0.4)
var defaultTimeout = int32(5000)
var healthCheck = "/health"
var defaultPort = 8080

var defaultConfig = AppConfig{
	Prefix:      "/api",
	HealthCheck: &healthCheck, // Optional field
	Port:        8080,
	Latency:     &defaultLatency,  // Optional field
	FailRate:    &defaultFailRate, // Optional field
	Timeout:     &defaultTimeout,  // Optional field
	ApiKey:      nil,              // Optional field
}

type App struct {
	Config AppConfig `json:"config"`
	APIs   []API     `json:"api"`
}

func main() {
	file, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatal(err)
	}
	//reading the config from the file
	var app App
	err = json.Unmarshal(file, &app)
	if err != nil {
		fmt.Println("Aight Something went wrong")
		log.Fatal(err)

	}

	if int(app.Config.Port) == 0 {
		app.Config.Port = defaultPort
	}

	if app.Config.Prefix == "" {
		app.Config.Prefix = defaultConfig.Prefix
	}

	//starting a server with the port
	http.HandleFunc(fmt.Sprintf("%s/", app.Config.Prefix), func(w http.ResponseWriter, r *http.Request) {
		var path = r.URL.Path
		//check if the path is in the APIs
		var found = false

		for _, api := range app.APIs {
			if strings.HasPrefix(path, app.Config.Prefix) && api.Path == strings.TrimPrefix(path, app.Config.Prefix) && api.Method == r.Method {
				found = true
				//check if the latency is provided
				if api.Latency != nil || app.Config.Latency != nil || *app.Config.Latency != 0 {
					//sleep for the latency
					fmt.Printf("Sleeping for %d ms\n", *api.Latency)
					time.Sleep(time.Duration(*api.Latency) * time.Millisecond)

				}
				//check if the fail rate is provided
				if api.FailRate != nil {
					//generate a random number between 0 and 1
					//if the number is less than the fail rate, return an error
					random := rand.Float64()
					if random < *api.FailRate {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}

				}
				//check if the validate is provided
				if api.Validate != nil {
					//check if the query parameters are provided
					if api.Validate.Query != nil {
						//check if the query parameters are correct
					}
					//check if the body is provided
					if api.Validate.Body != nil {
						//check if the body is correct
					}
					//check if the params are provided
					if api.Validate.Params != nil {
						//check if the params are correct
					}
				}
				//return the response
				response := ResponseBuilder(api.Response)
				json.NewEncoder(w).Encode(response)

				fmt.Println("Response: ", response)
				//Send the response and stop the loop

				break

			}
		}
		if !found {
			http.Error(w, "Not Found", http.StatusNotFound)
		}

	})
	fmt.Print(app.APIs[0].Method)
	fmt.Printf("Server is running on http://localhost:%d%s/\n", app.Config.Port, app.Config.Prefix)
	http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)
	// Wait for interrupt signal to gracefully shutdown the application
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down...")
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
		//TODO : Implement the Date function
		// "time": func(args *string) interface{} {
		// 	switch {
		// 	case args == nil:
		// 		return fake.Time().DayOfWeek()
		// 	case *args == "long":
		// 		return fake.Date().Long()
		// 	case *args == "short":
		// 		return fake.Date().Short()
		// 	case *args == "birthday":
		// 		return fake.Date().Birthday()
		// 	case *args == "unix":
		// 		return fake.Date().Unix()
		// 	case *args == "unixNano":
		// 		return fake.Date().UnixNano()
		// 	default:
		// 		return fake.Date().Short()
		// 	}
		// },
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
			var funcName = strings.Split(v, "(")[0]
			var args = strings.Split(v, "(")[1]
			args = strings.TrimRight(args, ")")
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
