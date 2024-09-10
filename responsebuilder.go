package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jaswdr/faker/v2"
)

func ResponseBuilder(rawData map[string]interface{}) map[string]interface{} {
	response := make(map[string]interface{})
	//LOL FUNNY FIX
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed + rand.Int63())
	fake := faker.NewWithSeed(source)
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
			case *args == "sha256":
				return fake.Hash().SHA256()
			case *args == "md5":
				return fake.Hash().MD5()
			case *args == "sha512":
				return fake.Hash().SHA512()

			case *args == "userAgent":
				var userAgent = [...]string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
					"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/117.0",
					"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:119.0) Gecko/20100101 Firefox/119.0",
					"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/537.36",
					"Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1",
					"Mozilla/5.0 (Android 13; Mobile; rv:119.0) Gecko/119.0 Firefox/119.0",
					"Mozilla/5.0 (Android 13; Tablet; rv:117.0) Gecko/117.0 Firefox/117.0",
					"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/114.0.0.0 Safari/537.36",
					"Mozilla/5.0 (Linux; Android 10; SM-G960U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36",
					"Mozilla/5.0 (Linux; Android 12; SM-T820) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"}
				return userAgent[rand.Intn(len(userAgent))]

			case *args == "sqlId":
				// return an integer between 1 and 100000
				return fake.Int64Between(1, 100000)
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
			case *args == "username":
				return fake.Internet().User()
			case *args == "title":
				return fake.Person().Title()
			case *args == "gender":
				return fake.Person().Gender()
			case *args == "ssn":
				return fake.Person().SSN()
			case *args == "bio":
				return fake.Lorem().Sentence(20)
			case *args == "gamerTag":
				return fake.Gamer().Tag()
			case *args == "birthday":
				rand.Seed(time.Now().UnixNano())

				// TODO ChAT GPT CODE DEBUG LATER
				years := rand.Intn(21) - 10
				months := rand.Intn(25) - 12
				days := rand.Intn(61) - 30

				// Add the random values to the current date
				randomDate := time.Now().AddDate(years, months, days)

				// Return the random date
				return randomDate
			case *args == "image":
				return fmt.Sprintf("https://randomuser.me/api/portraits/med/men/%d.jpg", fake.RandomDigit())
			default:
				return fake.Person().Name()
			}
		},
		"Age": func(args *string) interface{} { return fake.Int8Between(10, 100) },
		"Finance": func(args *string) interface{} {
			switch {
			case *args == "":
				return fake.Int64Between(1000, 10000000)
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
			case *args == "amountWithCurrency":
				return fmt.Sprintf("%d %s", fake.Int64Between(1000, 100000), fake.Currency().Code())
			case *args == "btcAddress":
				return fake.Crypto().BitcoinAddress()
			case *args == "ethAddress":
				return fake.Crypto().EtheriumAddress()
			default:
				return fmt.Sprintf("%s Not a valid Finance Argument", *args)
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
				return fake.Int64Between(1, 1000000)
			}
			if val, err := strconv.Atoi(*args); err == nil {
				return rand.Intn(val)
			}
			return fake.Int64Between(1, 1000000)
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
			var alen int
			var dataType string
			var result []interface{}
			alen = 5
			dataType = "User()"

			if *args == "" {
				for i := 0; i < alen; i++ {
					result = append(result, fake.Lorem().Word())
				}
				return result
			}
			pargs, err := parseArguments(*args)

			if err != nil {
				return err.Error()
			}
			if val, ok := pargs["len"]; ok {
				alen, _ = strconv.Atoi(val)
			}
			if val, ok := pargs["type"]; ok {
				dataType = val
			}
			for i := 0; i < alen; i++ {
				if dataType == "Json()" {
					result = append(result, fake.Map())
					continue
				}
				if dataType == "Array()" {
					result = append(result, fake.Lorem().Word())
					continue
				}
				response := ResponseBuilder(map[string]interface{}{"type": dataType})
				if value, ok := response["type"].(string); ok {
					result = append(result, value)
				} else if value, ok := response["type"].(map[string]interface{}); ok {
					result = append(result, value)
				} else if value, ok := response["type"].([]interface{}); ok {
					result = append(result, value)
				} else {
					result = append(result, "unexpected_type")
				}
			}
			return result
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
		"DiceBearImage": func(args *string) interface{} {
			seed := fake.Person().Name()
			switch {
			case args == nil:
				return formatAvatarURL(seed, "adventurer-neutral")
			case *args == "adventurer-neutral":
				return formatAvatarURL(seed, "adventurer-neutral")
			case *args == "avataaars":
				return formatAvatarURL(seed, "avataaars")
			case *args == "avataaars-neutral":
				return formatAvatarURL(seed, "avataaars-neutral")
			case *args == "big-ears":
				return formatAvatarURL(seed, "big-ears")
			case *args == "big-ears-neutral":
				return formatAvatarURL(seed, "big-ears-neutral")
			case *args == "big-smile":
				return formatAvatarURL(seed, "big-smile")
			case *args == "bottts":
				return formatAvatarURL(seed, "bottts")
			case *args == "bottts-neutral":
				return formatAvatarURL(seed, "bottts-neutral")
			case *args == "croodles":
				return formatAvatarURL(seed, "croodles")
			case *args == "croodles-neutral":
				return formatAvatarURL(seed, "croodles-neutral")
			case *args == "dylan":
				return formatAvatarURL(seed, "dylan")
			case *args == "fun-emoji":
				return formatAvatarURL(seed, "fun-emoji")
			case *args == "identicon":
				return formatAvatarURL(seed, "identicon")
			case *args == "initials":
				return formatAvatarURL(seed, "initials")
			case *args == "glass":
				return formatAvatarURL(seed, "glass")
			case *args == "lorelei":
				return formatAvatarURL(seed, "lorelei")
			case *args == "lorelei-neutral":
				return formatAvatarURL(seed, "lorelei-neutral")
			case *args == "micah":
				return formatAvatarURL(seed, "micah")
			case *args == "miniavs":
				return formatAvatarURL(seed, "miniavs")
			case *args == "notionists":
				return formatAvatarURL(seed, "notionists")
			case *args == "notionists-neutral":
				return formatAvatarURL(seed, "notionists-neutral")
			case *args == "open-peeps":
				return formatAvatarURL(seed, "open-peeps")
			case *args == "personas":
				return formatAvatarURL(seed, "personas")
			case *args == "pixel-art":
				return formatAvatarURL(seed, "pixel-art")
			case *args == "pixel-art-neutral":
				return formatAvatarURL(seed, "pixel-art-neutral")
			case *args == "rings":
				return formatAvatarURL(seed, "rings")
			case *args == "shapes":
				return formatAvatarURL(seed, "shapes")
			case *args == "thumbs":
			default:
				return formatAvatarURL(seed, "adventurer-neutral")
			}
			return formatAvatarURL(seed, "adventurer-neutral")
		},
		"Placehold": func(args *string) interface{} {
			var width, height int
			var text, font, color, bgColor string
			if *args == "" {
				return "https://placehold.co/600x400"
			} else {
				width = 600
				height = 400
				text = "Hello World"
				font = "Montserrat"
				color = "000000"
				bgColor = "FFF"
				parsedArgs, err := parseArguments(*args)
				if err != nil {
					return err.Error()
				}
				if val, ok := parsedArgs["width"]; ok {
					width, _ = strconv.Atoi(val)
				}
				if val, ok := parsedArgs["height"]; ok {
					height, _ = strconv.Atoi(val)
				}
				if val, ok := parsedArgs["text"]; ok {
					text = val
				}
				if val, ok := parsedArgs["font"]; ok {
					font = val
				}
				if val, ok := parsedArgs["color"]; ok {
					color = val
				}
				if val, ok := parsedArgs["bgColor"]; ok {
					bgColor = val
				}
			}
			return fmt.Sprintf("https://placehold.co/%dx%d/%s/%s?text=%s&font=%s", width, height, color, bgColor, text, font)

		},
		"LoremPicsum": func(args *string) interface{} {
			var width, height, blur, id int
			var grayscale bool

			if *args == "" {
				return "https://picsum.photos/200/300"
			} else {
				width = 200
				height = 200
				blur = 1
				grayscale = false
				id = rand.Intn(1000)
				parsedArgs, err := parseArguments(*args)
				if err != nil {
					return err.Error()
				}
				if val, ok := parsedArgs["width"]; ok {
					width, _ = strconv.Atoi(val)
				}
				if val, ok := parsedArgs["height"]; ok {
					height, _ = strconv.Atoi(val)
				}
				if val, ok := parsedArgs["blur"]; ok {
					blur, _ = strconv.Atoi(val)
				}
				if val, ok := parsedArgs["grayscale"]; ok {
					grayscale, _ = strconv.ParseBool(val)
				}
			}
			if grayscale {
				return fmt.Sprintf("https://picsum.photos/%d/%d?blur=%d&grayscale&id=%d", width, height, blur, id)
			}
			return fmt.Sprintf("https://picsum.photos/%d/%d?blur=%d&id=%d", width, height, blur, id)

		},
		"Message": func(args *string) interface{} {
			if args == nil {
				return "Success"
			}
			return *args
		},
		"Json": func(args *string) interface{} {
			return fake.Map()
		},
		"Uuid": func(args *string) interface{} {
			return fake.UUID().V4()
		},
	}

	// Usage
	response = processMap(rawData, funcMap)
	return response
}
