# GOBI

A simple mock server that allows users to fine-grain control over the mock data.

This is my first program in Go, and I am enjoying the language. A mock server is what I needed, so as my introduction to Go, I decided to make it in Go.

# Why ?
Creating frontend application suck when you don't have a backend and m too poor to pay for a mocking service and I wanted to learn GO so two birds one stone.  : )

So far as a first experience with the language it's been a great better than js TBH XD

## Features
- Fine-grain control over mock data
- Easy to use and configure
- Lightweight and fast
- Written in Go
- Download and run

## Installation
Either download the binary from the releases or build it yourself.  :)

## Usage
The app is simple for now it needs a json file named <b>input.json</b> to be present in the same directory as the binary. The json file should have the following structure:

```json
{
    "config":{
        "prefix":"/api",
        "healthCheck": "/to_be_implemented",
        "latency": 200,
        "Port": 3000,
        "logging": true,
        "failRate": 0,
        "timeout": "to be implemented",
        "apiKey":  "to be implemented"
    },
}

```
This will be the configuration for the server and is global a latency added here will be applied to all the endpoints. same as the prefix

| Key | Description |  default |
| --- | --- | --- |
| prefix | The prefix for the server as http://domain/prefix/your_endpoint  | /api
| healthCheck | A route where you can ping without (TBI) | /health
| latency | The latency to be added to the server in (milliseconds) | 0
| Port | The port on which the server will run | 8080
| logging | If true the server will log the request and response | false
| failRate | The rate at which the server will fail the request. 0 the server will never fail 1 the server will always fail and you can fine-tune by setting it to anything in between 0 and 1 | 0
| timeout | The server will timeout after this interval (TBI) | 0
| apiKey | Protecting routes By providing api auth  (TBI) | nil


None the values above are required and can be omitted. you just need to pass an empty configkey in the json file.

After The config you can add the endpoints you want to mock. The structure of the endpoint is as follows:

```json
{
    "config":{},
    "api":[
        {
            "path": "/test",
            "method": "GET",
            "response":{
                "name":"User()",
                "address":"Address()",
                "display_picture":"Image(seed)",
                "message":"Success(any message you want)",
        },
     }
    ]
}
```
Response 
``` json
{
  "address": "%6775 Donnie Highway\nHageneshaven, CO 56100-9243",
  "display_picture": "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=seed",
  "message": "any message you want",
  "name": "Tara Kuvalis"
}
```

You can request objects Deeply nested 
```json
{
    "config":{},
    "endpoints":[
        {
            "path": "/test",
            "method": "GET",
            "response":{
                "nested":{
                    "name":"User()",
                    "address":"Address()",
                    "display_picture":"Image(seed)",
                    "message":"Success(any message you want)",
                    "payment_info":{
                        "credit_card":"Finance(creditCard)",
                        "card_type":"Finance(cardType)",
                        "exp":"Finance(cardExpirationDate)",
                    }
                }
        },
     }
    ]
}
```

``` json 
{
  "nested": {
    "address": "%83 Osinski Dale Suite 738\nSouth Bobbieville, MD 04525",
    "display_picture": "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=seed",
    "message": "any message you want",
    "name": "Merlin Gislason",
    "payment_info": {
      "card_type": "Visa",
      "credit_card": "2418850258956182",
      "exp": "30/18"
    }
  }
}
```

You can fine tune Each endpoint you want 
    
```json
    {
        "config":{},
        "api":[
            {
                "path": "/minimal",
                "method": "GET",
                "response":{
                    "name":"User()",
                    "address":"Address()",
                    "display_picture":"Image(seed)",
                    "message":"Success(any message you want)",
                },
            },
            {
                "path":"/finegrain",
                "method":"GET",
                "latency": 200,
                "failRate": 0.5,
                "responseType":"Array(8)",
                "validate":{
                "query":["id","name","age","email"],
                    "body": "to be implemented",
                },
                "response":{
                    "name":"User()",
                    "address":"Address()",
                    "display_picture":"Image(seed)",
                    "message":"Success(any message you want)",
                },
            }
        ]
    }
``` 
Explanation of the above json
- The first endpoint is a simple endpoint with no extra configuration
- The second endpoint is a fine-grain endpoint with a latency of 200ms and a fail rate of 0.5. The response type is an array of 8 objects. The query is validated for id, name, age, and email. The response is the same as the first endpoint.

| Key | Description |  required |
| --- | --- | --- |
| path | The path of the endpoint  this will be the path you want to access it'll be prefixed with the provided prefix | yes
| method | The method of the endpoint  | yes
| latency | The latency to be added to this path only and will be stacked on top of global latency | no
| failRate | The rate at which this request will fail. 0 the server will never fail 1 the server will always fail and you can fine-tune by setting it to anything in between 0 and 1 | no
| responseType | For now it just Support Array(num) If you want to request An array of the given response you can provide this the num specifies the length of the array | no
| validate | The query and body you want to validate For now the query can be validated Body validation is yet to be implemented | no
| response | The response object You want to receive | yes

## Mock Data Methods
The mock data is powered by [Faker](github.com/jaswdr/faker/v2) and tries to provide access to most of the methods provided by the library in a simple way.

## Available Methods
| Method | Description |
| --- | --- |
| User() | Used to mock user-related data such as email, name, phone, etc. |
| Address() | Used to mock address-related data such as street, city, country, etc. |
| Image() | Returns a mock image URL (Work In Progress) |
| Finance() | Returns finance-related data such as crypto, banks, etc. |
| Time() | Returns mock time data such as current time, date, etc. |
| Vehicle() | Returns mock vehicle data such as make, model, year, etc. |
| App() | Returns mock application-related data such as app name, version, etc. |
| Language() | Returns mock language-related data such as language name, programming language, etc. |
| Bool() | Returns a random boolean value (true or false) |
| Float() | Returns a random float value |
| Int() | Returns a random integer value |
| Json() | Returns mock JSON data |
| Array() | Returns a mock array of data |
| Lorem() | Returns mock lorem ipsum text |
| Internet() | Returns mock internet-related data such as IP address, URL, etc. |
| Company() | Returns mock company-related data such as company name, BS, etc. |
| Color() | Returns mock color data such as color name, hex code, etc. |

Most of the methods above accepts one parameter that can be used to fine-tune the data. For example, `User(name)` will return a random name and `User(mFirstName)` will return a male first name. 

# A list of everything 
``` json 
      "response":{
          "color":{ 
            "color":"Color()",
            "hex":"Color(hex)",
            "rgb":"Color(rgb)",
            "rgba":"Color(rgba)",
            "safe" :"Color(safe)",
            "css":"Color(css)"
          },
          "address" :{
            "address":"Address()",
            "city":"Address(city)",
            "state":"Address(state)",
            "country":"Address(country)",
            "zip":"Address(zip)",
            "latitude":"Address(latitude)",
            "longitude":"Address(longitude)",
            "secondary":"Address(full)",
            "street":"Address(street)",
            "country_code":"Address(countryCode)",
            "state_code":"Address(stateAbbr)",
            "countryAbbr":"Address(countryAbbr)"
          },
          "company":{
            "company":"Company()",
            "catch_phrase":"Company(catchPhrase)",
            "bs":"Company(bs)",
            "jobTitle":"Company(jobTitle)",
            "ein":"Company(ein)",
            "suffix":"Company(suffix)",
            "mail":"Company(mail)"
          },
          "internet":{
            "url":"Internet()",
            "domain":"Internet(domain)",
            "ip":"Internet(ip)",
            "ipv6":"Internet(ipv6)",
            "mac":"Internet(mac)",
            "httpMethod":"Internet(httpMethod)",
            "tld":"Internet(tld)",
            "slug":"Internet(slug)",
            "status_code":"Internet(statusCode)",
            "free_email":"Internet(freeEmail)",
            "safe_email":"Internet(safeEmail)",
            "status_code_message":"Internet(statusCodeMessage)"
          },
          "user":{
            "user":"User()",
            "email":"User(email)",
            "firstName":"User(firstName)",
            "lastName":"User(lastName)",
            "femaleFirstName":"User(fFirstName)",
            "female Title" :"User(fTitle)",
            "maleFirstName":"User(mFirstName)",
            "maleTitle":"User(mTitle)",
            "phone":"User(phone)",
            "userName":"User(userName)",
            "password":"User(password)",
            "title":"User(title)",
            "gender":"User(gender)",
            "ssn":"User(ssn)",
            "bio":"User(bio)",
            "birthday":"User(birthday)",
            "age":"Age()",
            "gamer_tag":"User(gamerTag)",
            "uuid":"User(uuid)",
            "sqlId":"User(sqlId)"
          },
          "finance":{
            "credit_card":"Finance(creditCard)",
            "card_type":"Finance(cardType)",
            "exp":"Finance(cardExpirationDate)",
            "amount":"Finance(amount)",
            "iban":"Finance(iban)",
            "currency":"Finance(currency)",
            "currencyCode":"Finance(currencyCode)",
            "currencyAndCode":"Finance(currencyAndCode)",
            "amountWithCurrency":"Finance(amountWithCurrency)",
            "transactionDescription":"Finance(transactionDescription)",
            "btcAddress":"Finance(btcAddress)",
            "ethAddress":"Finance(ethAddress)"
          },
          "Lorem":{
            "word":"Lorem(word)",
            "words":"Lorem(words)",
            "sentence":"Lorem(sentence)",
            "paragraph":"Lorem(paragraph)",
            "paragraphs":"Lorem(paragraphs)", //Returns a lot of data
            "Lorem":"Lorem()"
          },
          "misc":{
            "boolean":"Bool()",
            "float":"Float()",
            "integer":"Int()",
            "intSmall":"Int(20)",
            "json":"Json()",
            "array":"Array()"
          },
          "language":{
            "language":"Language()",
            "languageabbr":"Language(abbr)",
            "languageName":"Language(programming)"
          },
          "app":{
            "appName":"App()",
            "appVersion":"App(version)",
            "platform":"App(platform)"
          },
          "vehicle":{
            "vehicle":"Vehicle()",
            "vehicleBrand":"Vehicle(brand)",
            "vehicleType":"Vehicle(type)",
            "vehicleTransmission":"Vehicle(transmission)",
            "vehiclePlate":"Vehicle(plate)"
            },
          "Time":{
            "time":"Time()",
            "unix":"Time(unix)",
            "unixNano":"Time(unixNano)",
            "iso":"Time(iso)",
            "day":"Time(day)",
            "ansi":"Time(ansi)",
            "monthName":"Time(monthName)",
            "timezone":"Time(timezone)"
          },
          "Image" :{
            "image":"Image(dfas)"
          }
        }
```

Response
``` json 
{
  "Image": {
    "image": "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=dfas"
  },
  "Lorem": {
    "Lorem": "quia voluptate voluptatum assumenda necessitatibus veritatis accusantium consequuntur aut occaecati rerum eum natus dicta vitae ex vero similique quibusdam sed.",
    "paragraph": "quaerat dolor veniam saepe quisquam nihil id soluta consequatur beatae vitae voluptas explicabo doloremque laboriosam est similique unde temporibus accusamus architecto enim asperiores sit autem illum dolor nulla magnam omnis enim aut aliquam neque non est doloremque sit sunt est soluta qui omnis molestiae ut sed numquam consequatur explicabo et ut a et dolor rem ratione.",
    "sentence": "aut iusto qui doloremque vero perferendis libero aliquid repellendus voluptatem blanditiis iusto error fugit placeat quo adipisci incidunt illo dolores.",
    "word": "aut",
    "words": "voluptatem iste officiis ratione maiores hic et voluptates laboriosam distinctio numquam provident doloribus architecto nemo fugiat praesentium molestiae a exercitationem."
  },
  "Time": {
    "ansi": "Sun May  3 03:33:53 1981",
    "day": 4,
    "iso": "2013-07-09T06:40:23+000",
    "monthName": "December",
    "time": "Thu, 05 Sep 2024 22:00:01 +0530",
    "timezone": "Asia/Pyongyang",
    "unix": 1725553801,
    "unixNano": 1725553801159471600
  },
  "address": {
    "address": "%815 Beier Wall Apt. 837\nWelchshire, MT 42878",
    "city": "Port Liza",
    "country": "Senegal",
    "countryAbbr": "GRD",
    "country_code": "BR",
    "latitude": 34.235314,
    "longitude": 36.882405,
    "secondary": "%21 Emard Junctions Apt. 832\nAmieberg, MT 52090",
    "state": "Minnesota",
    "state_code": "DC",
    "street": "Robel Freeway",
    "zip": "06285"
  },
  "app": {
    "appName": "App Joy",
    "appVersion": "v4.2.5",
    "platform": "Mac"
  },
  "color": {
    "color": "SlateBlue",
    "css": "rgb(95,65,54)",
    "hex": "#06DD48",
    "rgb": "193,115,4",
    "rgba": [
      "181",
      "54",
      "56"
    ],
    "safe": "blue"
  },
  "company": {
    "bs": "utilize web-enabled markets",
    "catch_phrase": "Cloned eco-centric success",
    "company": "Rolfson Group",
    "ein": 77,
    "jobTitle": "Merchandise Displayer OR Window Trimmer",
    "mail": "durgan.ernesto@shields-shields.cpb.net",
    "suffix": "Ltd"
  },
  "finance": {
    "amount": 3512812,
    "amountWithCurrency": "98691 TZS",
    "btcAddress": "1zMi2V87gt7BJmP24415wz95J9yXDmJju7",
    "card_type": "MasterCard",
    "credit_card": "6529527462718439",
    "currency": "Kenyan Shilling",
    "currencyAndCode": "Bulgarian Lev (BGN)",
    "currencyCode": "SCR",
    "ethAddress": "0xOxbH4eZjDJ0QAlZ0AdtWCmaH6mfW9F311dcM9I01",
    "exp": "15/22",
    "iban": "SI75634784348549603",
    "transactionDescription": "Finance Needs an argument"
  },
  "internet": {
    "domain": "yvp.net",
    "free_email": "bessie.johns@gmail.com",
    "httpMethod": "HEAD",
    "ip": "88.178.165.59",
    "ipv6": "3723:8626:3625:5775:4471:6852:4313:4625",
    "mac": "84:46:86:75:0F:4F",
    "safe_email": "giuseppe@example.org",
    "slug": "qhshe-nwpg",
    "status_code": 403,
    "status_code_message": "Temporary Redirect",
    "tld": "net",
    "url": "http://emo.org/kaoz-abrj"
  },
  "language": {
    "language": "Sunda",
    "languageabbr": "tg",
    "programming": "ASP / ASP.NET"
  },
  "misc": {
    "array": [
      "Aniya Rolfson DVM",
      "Kennith Kertzmann",
      "Carlotta Schmeler"
    ],
    "boolean": true,
    "float": 682.0999755859375,
    "intSmall": 0,
    "integer": 4,
    "json": {
      "aliquid": 840732.3125,
      "autem": 4486365,
      "nam": [
        "incidunt",
        "laborum",
        "eum",
        "et",
        "ut",
        "sed",
        "qui",
        "voluptatem"
      ],
      "quos": {
        "quisquam": 620453.125
      }
    }
  },
  "user": {
    "age": 71,
    "bio": "et corporis repellendus autem vero et incidunt voluptatum placeat dolorem recusandae ex nulla quasi quia mollitia voluptatem quibusdam qui ipsa omnis perferendis officia autem rerum ex labore mollitia repudiandae exercitationem aliquam consequatur dolores illum tempora et dolores veniam rem quis odio rerum voluptatem et deserunt consectetur dolor omnis recusandae dicta in numquam magnam at non perspiciatis officia iste sed maxime porro veniam reiciendis suscipit temporibus magnam deserunt voluptatem voluptatem maiores est cum voluptas vitae vel et earum excepturi excepturi laudantium impedit dolore labore dolores autem sit a sed ducimus quia et sint perspiciatis eligendi ullam blanditiis.",
    "birthday": "2034-01-31T22:00:01.1594715+05:30",
    "email": "coleman.ledner@gmail.com",
    "female_firstName": "Mireya",
    "female_title": "Ms.",
    "firstName": "Misty",
    "gamer_tag": "TheZodiac",
    "gender": "Male",
    "lastName": "Kessler",
    "maleFirstName": "Selmer",
    "maleTitle": "Mr.",
    "password": "exugfi|bln|r{",
    "phone": "920-325-2779 x4587",
    "sqlId": 77701,
    "ssn": "295281269",
    "title": "Ms.",
    "user": "Janelle Wunsch",
    "userName": "pasquale.upton",
    "uuid": "c7722f94-59d8-459e-96a9-93ae6b844866"
  },
  "vehicle": {
    "vehicle": "Hatchback",
    "vehicleBrand": "Ford",
    "vehiclePlate": "АC4543MА",
    "vehicleTransmission": "CVT",
    "vehicleType": "Ethanol"
  }
}
```

### Things to be implemented
- [ ] Timeout
- [ ] Health Check
- [ ] API Key Auth
- [ ] Better Query Validation
- [ ] Better Image Mocking
- [ ] Body Validation
- [ ] More Mock Data Methods (If needed)
- [ ] More Fine Grained Control
- [ ] Better Error Handling and Logging
- [ ] Better Documentation
- [ ] Code Cleanup
- [ ] Adding test
- [ ] Improve Project Structure
### Contributing
Feel free to do anything It's my first project to shitty code is expected.  : )
If ya find any bugs or have any suggestions feel free to open an issue or a PR. and please explain the issue or the PR in detail. It'll help me understand better.


