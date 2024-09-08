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
			case *args == "passwordMd5":
				return fake.Hash().MD5()
			case *args == "passwordSha256":
				return fake.Hash().SHA256()
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
			var result []string
			for i := 0; i < 20; i++ {
				result = append(result, fake.RandomStringElement([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}))
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
			//TODO: Implement ImagePlaceholder API
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
	}

	// Usage
	response = processMap(rawData, funcMap)
	return response
}
