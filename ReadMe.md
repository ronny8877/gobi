# GOBI

A simple mock server that allows users to fine-grain control over the mock data.

This is my first program in Go, and I am enjoying the language. A mock server is what I needed, so as my introduction to Go, I decided to make it in Go.

# heh GOBI

## Features
- Authenticated Routes
- Basic Input Validation
- Fine-grain control over mock data
- Easy to use and configure
- Lightweight and fast
- Hot Reload
- Free and Open Source
- Written in Go

## Installation
Either download the binary from the releases or build it yourself.  :) 

## Usage
The app is simple, It needs a json file named <b>input.json</b> to be present in the same directory as the binary. The json file should have the following structure:

```json
{
    "config":{
        "prefix":"/api",
        "latency": 200,
        "Port": 3000,
        "logging": true,
        "failRate": 0,
        "auth":{
        "apiKey": "1234567890",
        "bearer": "1234567890",
        "cookie": "auth=1234567890"
        }
    },
    "api":[...]
}
```
This is how usual setup will Look But you can omit any of the keys in the config object and start with a minimal setup.

```json
{
    "config":{},
    "api":[]
}

```
To test the server, you can ping the `/health`. This is a reserved route and will always be available you cannot change this route.

###### Note the prefix do not need to be added to the health check route.

```bash
curl http://localhost:8080/health
```


#### Description

This will be the configuration for the server and is global a latency added here will be applied to all the endpoints. same as the prefix

| Key | Description |  default |
| --- | --- | --- |
| prefix | The prefix for the server as http://domain/prefix/your_endpoint  | /api
| latency | The latency to be added to the server in (milliseconds) | 0
| Port | The port on which the server will run | 8080
| logging | If true the server will log the request and response | false
| failRate | The rate at which the server will fail the request. 0 the server will never fail 1 the server will always fail and you can fine-tune by setting it to anything in between 0 and 1 | 0
| auth | Auth have three different Methods, ApiKey, Cookie and Bearer. Read Below For More details | nil


None the values above are required and can be omitted. you just need to pass an empty configkey in the json file.

###### Endpoints

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
                "display_picture":"DiceBearImage()",
                "custom_bio":"Message(Hi, My name is Ronny and i like cookies)",
                "bio":"User(bio)",
        },
     }
    ]
}
```
##### Response 
``` json
{
  "address": "%6775 Donnie Highway\nHageneshaven, CO 56100-9243",
  "display_picture": "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=seed",
  "custom_bio": "Hi, My name is Ronny and i like cookies",
  "bio":"iusto ullam aut fuga vitae numquam dolorem et voluptatem quas rerum molestiae totam perspiciatis assumenda molestiae facere et ab cumque.",
  "name": "Tara Kuvalis"
}
```
##### Nested Response
 
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
                    "display_picture":"DiceBearImage(pixel-art)",
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
###### Response
``` json 
{
  "nested": {
    "address": "%83 Osinski Dale Suite 738\nSouth Bobbieville, MD 04525",
    "display_picture": "https://api.dicebear.com/9.x/pixel-art/svg?seed=Deja Kshlerin",
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
##### Path params
You can also add path params to the endpoint. The path params should be added in the path with a colon before the name of the param. The value of the path param will be available in the response object.

```json
{
    "config":{},
 "api":[
    {
      "path": "/test/:id",
      "method": "GET",
      "response":{
          "message":"Message(Test Successful)"
      }
  },
   {
      "path": "/test/:id/name",
      "method": "GET",
      "response":{
          "message":"Message(This will be a different Route)"
      }
  },
   {
      "path": "/test/:id/name/:name",
      "method": "GET",
      "response":{
          "message":"Message(Again will be a different Route)"
      }
  },
   {
      "path": "/test/:id/:name/:age",
      "method": "GET",
      "response":{
          "message":"Message(This will not work as the Route is similar to one above it, Whichever comes first will be executed"
      }
  },
  {
    "path": "/test/:id/:name/:age/why",
    "method": "GET",
    "response":{
        "message":"Message(Will work again as a septate route"
    }
}
  ]
}
```


### Tuning Endpoints

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
                    "message":"Message(any message you want)",
                },
            },
            {
                "path":"/finegrain",
                "method":"GET",
                "latency": 200,
                "failRate": 0.5,
                "responseType":"Array(8)",
                "response":{
                    "name":"User()",
                    "address":"Address()",
                    "display_picture":"User(image)",
                    "message":"Message(Hi there)",
                    "payment_info":{
                        "credit_card":"Finance(creditCard)",
                        "card_type":"Finance(cardType)",
                        "exp":"Finance(cardExpirationDate)",
                        "crypto":{
                            "btc":"Finance(btcAddress)",
                            "eth":"Finance(ethAddress)"
                        }
                      }
                    }
            }
        ]
    }
``` 
###### Response

``` json
Minimal
{
  "address": "%814 Stanford Locks\nEast Karson, IN 67216",
  "display_picture": "https://randomuser.me/api/portraits/med/men/12.jpg",
  "message": "any message you want",
  "name": "Art Lowe III"
}

Fine grain
[
  {
    "address": "%03 Shields Ville\nSouth Feltonside, GA 98219-5142",
    "display_picture": "https://randomuser.me/api/portraits/med/men/8.jpg",
    "message": "Hi there",
    "name": "Doyle Tremblay V",
    "payment_info": {
      "card_type": "Visa Retired",
      "credit_card": "8610166633991501",
      "crypto": {
        "btc": "bc1XXou24Yr18x1dPCu4Bu77R5G42anq811",
        "eth": "0xsZXJ9D0Z3b82cr6M2dsc14fKpA1sW7FTUMYknoRs"
      },
      "exp": "03/15"
    }
  },
  {
    "address": "%12 Callie Crest\nConstantinbury, MA 50612",
    "display_picture": "https://randomuser.me/api/portraits/med/men/8.jpg",
    "message": "Hi there",
    "name": "Sonya Mitchell I",
    "payment_info": {
      "card_type": "MasterCard",
      "credit_card": "0889521375720673",
      "crypto": {
        "btc": "3f2qCJh1GfYG1r16h14a56dC6r9PXb5m",
        "eth": "0xGbB5jtL7PZNP3lpX41EjW0zTL30noPa1875ea7mZ"
      },
      "exp": "28/24"
    }
  },
... 8 more
]

```

###### Explanation of the above json
- The first endpoint is a simple endpoint with no extra configuration
- The second endpoint is a fine-grain endpoint with a latency of 200ms and a fail rate of 0.5. The response type is an array of 8 objects.

| Key | Description |  required |
| --- | --- | --- |
| path | The path of the endpoint  this will be the path you want to access it'll be prefixed with the provided prefix | yes
| method | The method of the endpoint  | yes
| latency | The latency to be added to this path only and will be stacked on top of global latency | no
| failRate | The rate at which this request will fail. 0 the server will never fail 1 the server will always fail and you can fine-tune by setting it to anything in between 0 and 1 | no
| responseType | For now it just Support Array(num) If you want to request An array of the given response you can provide this the num specifies the length of the array | no
| response | The response object You want to receive | yes


#### Validation 
The server also support basic Body and query validation which you can use to validate the request. If any of the specified Keys are missing the server will throw an error and not return any response

``` json

{
  "config":{},
  "api":[
      {
          "path": "/validation",
          "validate":{
            "query" : ["name","email","password"]
          },
          "method": "GET",
          "response":{
              "message":"Message(No user Found)"
          }
      },
      {
          "path": "/validation",
          "validate":{
            "body" : ["name","email","password"]
          },
          "method": "POST",
          "response":{
              "message":"Message(User Created Successfully)"
          }
      }
  ]
}
```
The POST request will need to have the body with the keys name, email, and password. If any of the keys are missing the server will return a 400 error.

The GET request will need to have the query with the keys name, email, and password. If any of the keys are missing the server will return a 400 error.

Both query params can have n number of keys greater than the provided keys but the provided keys should be present in the request body or query depending on the validation type.


#### Auth
The server support Auth and you can provide the auth key in the config object. The server supports three types of Auth.
- ApiKey
- Bearer
- Cookie 

###### Right now the Auth do not offer much customization but if needed we will add more options in the future. Right now this was the most basic setup we could think of.

``` json
{
  "config":{
        "auth":{
            "apiKey": "1234567890",
            "bearer": "1234567890",
            "cookie": "auth=1234567890"
        }
    },
  "api":[
      {
          "path": "/me",
          "method": "GET",
          "auth":{
            "protected":true,
            "protectedBy":"bearer"
          },
          "response":{
              "message":"Message(No user Found)"
          }
      },
      {
          "path": "/create",
          "auth":{
            "protected":true,
            "protectedBy":"apiKey"
          },
          "method": "POST",
          "response":{
              "message":"Message(User Created Successfully)"
          }
      }
  ]
}
```

The above json will protect the `/me` route with the bearer token and the `/create` route with the apiKey. If the provided token is not correct the server will return a 401 error.

You can request the credentials from the server by sending a GET request to the `/auth` route. The server will return the credentials in the response.

```bash 
curl http://localhost:8080/auth
```
If the credentails are set this will return the credentials in the response.

```json
{
  "apiKey": "1234567890",
  "bearerToken": "1234567890",
  "cookie": "auth=1234567890"
}
```

And also in the response headers.

```yaml
set-cookie	auth=1234567890
x-api-key	1234567890
```

### Api key
The api key can be either provided in the header or in the query. The key should be `
x-api-key` in the header or `apiKey` in the query.

```bash
curl http://localhost:8080/create -X POST -H "x-api-key: 1234567890"
``` 

### Bearer
The bearer token can be provided in the header. The key should be `Authorization` in the header.

```bash
curl http://localhost:8080/me -H "Authorization : Bearer 1234567890"
```

### Cookie
The cookie can be provided in the header. The key should be `Cookie` in the header.

we plan to add more options to cookie in future but for now the value is `auth="your token"`

```bash
curl http://localhost:8080/me -H "Cookie : auth=1234567890; other=other"
```


## Mock Data Methods
The mock data is powered by [Faker](https://github.com/jaswdr/faker) and tries to provide access to most of the methods provided by the library in a simple way.

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
|Uuid() | Returns a random UUID |
| Array() | Returns a mock array of data also accepts two inputs `len` for array of len and `type`  for the type of data to be returned ie `Array(len=2000, type=User())` |
| Lorem() | Returns mock lorem ipsum text |
| Internet() | Returns mock internet-related data such as IP address, URL, etc. |
| Company() | Returns mock company-related data such as company name, BS, etc. |
| Color() | Returns mock color data such as color name, hex code, etc. |
| DiceBearImage() | Returns a mock image URL from DiceBear. Accepts one parameter Collection to generate a image from that collection DiceBearImage(pixel-art) |
| Message() | Returns a custom message. Accepts one parameter that will be returned as the message. `Message(Hello World)` ps any string that is not a defined function will be returned as it is. This method is just for the sake of uniformity you can just do `message: "any custom message"` |
| Placehold() | Returns a placehol image URL. Accepts multiple parameters `Placehold(width=1024,height=768,text=CustomText,font=Roboto,color=ff0000,bgColor=000000)` all of them are optional  ` |
| LoremPicsum() | Returns a Lorem Picsum Image URL. Accepts multiple parameters `LoremPicsum(width=400,height=300,blur=2,grayscale=true)` all of them are optional  ` |

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
          "status_code_message":"Internet(statusCodeMessage)",
          "user_agent":"Internet(userAgent)",
          "sha256":"Internet(sha256)",
          "md5":"Internet(md5)",
          "sha512":"Internet(sha512)",
          "sqlId":"Internet(sqlId)"
        },
        "user":{
          "user":"User()",
          "email":"User(email)",
          "firstName":"User(firstName)",
          "lastName":"User(lastName)",
          "female_firstName":"User(fFirstName)",
          "female_title" :"User(fTitle)",
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
          "image": "User(image)",
          "user_id":"Uuid()"
        },
        "finance":{
          "amount":"Finance()",
          "credit_card":"Finance(creditCard)",
          "card_type":"Finance(cardType)",
          "exp":"Finance(cardExpirationDate)",
          "iban":"Finance(iban)",
          "currency":"Finance(currency)",
          "currencyCode":"Finance(currencyCode)",
          "currencyAndCode":"Finance(currencyAndCode)",
          "amountWithCurrency":"Finance(amountWithCurrency)",
          "btcAddress":"Finance(btcAddress)",
          "ethAddress":"Finance(ethAddress)"
        },
        "Lorem":{
          "word":"Lorem(word)",
          "words":"Lorem(words)",
          "sentence":"Lorem(sentence)",
          "paragraph":"Lorem(paragraph)",
          "paragraphs":"Lorem(paragraphs)",
          "Lorem":"Lorem()"
        },
        "misc":{
          "boolean":"Bool()",
          "float":"Float()",
          "integer":"Int()",
          "intSmall":"Int(min=1,max=100000)",
          "message":"Message(Any custom message you want)",
          "json":"Json()",
          "array":"Array()",
          "array with len":"Array(len=20)",
          "array with len and type":"Array(len=20,type=User(email))",
          "uuid":"Uuid()"
        },
        "language":{
          "language":"Language()",
          "languageabbr":"Language(abbr)",
          "programming":"Language(programming)"
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
          "DiceBearImage": {
            "adventurerNeutral": "DiceBearImage(adventurer-neutral)",
            "avataaars": "DiceBearImage(avataaars)",
            "avataaarsNeutral": "DiceBearImage(avataaars-neutral)",
            "bigEars": "DiceBearImage(big-ears)",
            "bigEarsNeutral": "DiceBearImage(big-ears-neutral)",
            "bigSmile": "DiceBearImage(big-smile)",
            "bottts": "DiceBearImage(bottts)",
            "botttsNeutral": "DiceBearImage(bottts-neutral)",
            "croodles": "DiceBearImage(croodles)",
            "croodlesNeutral": "DiceBearImage(croodles-neutral)",
            "dylan": "DiceBearImage(dylan)",
            "funEmoji": "DiceBearImage(fun-emoji)",
            "identicon": "DiceBearImage(identicon)",
            "initials": "DiceBearImage(initials)",
            "glass": "DiceBearImage(glass)",
            "lorelei": "DiceBearImage(lorelei)",
            "loreleiNeutral": "DiceBearImage(lorelei-neutral)",
            "micah": "DiceBearImage(micah)",
            "miniavs": "DiceBearImage(miniavs)",
            "notionists": "DiceBearImage(notionists)",
            "notionistsNeutral": "DiceBearImage(notionists-neutral)",
            "openPeeps": "DiceBearImage(open-peeps)",
            "personas": "DiceBearImage(personas)",
            "pixelArt": "DiceBearImage(pixel-art)",
            "pixelArtNeutral": "DiceBearImage(pixel-art-neutral)",
            "rings": "DiceBearImage(rings)",
            "shapes": "DiceBearImage(shapes)",
            "thumbs": "DiceBearImage(thumbs)"
          },
          "Placehold": {
            "default": "Placehold()",
            "customWidth": "Placehold(width=1024)",
            "customHeight": "Placehold(height=768)",
            "customText": "Placehold(text=CustomText)",
            "customFont": "Placehold(font=Roboto)",
            "customColor": "Placehold(color=ff0000)",
            "customBgColor": "Placehold(bgColor=000000)",
            "customWidthHeight": "Placehold(width=1024,height=768)",
            "customWidthText": "Placehold(width=1024,text=CustomText)",
            "customWidthFont": "Placehold(width=1024,font=Roboto)",
            "customWidthColor": "Placehold(width=1024,color=ff0000)",
            "customWidthBgColor": "Placehold(width=1024,bgColor=000000)",
            "customHeightText": "Placehold(height=768,text=CustomText)",
            "customHeightFont": "Placehold(height=768,font=Roboto)",
            "customHeightColor": "Placehold(height=768,color=ff0000)",
            "customHeightBgColor": "Placehold(height=768,bgColor=000000)",
            "customTextFont": "Placehold(text=CustomText,font=Roboto)",
            "customTextColor": "Placehold(text=CustomText,color=ff0000)",
            "customTextBgColor": "Placehold(text=CustomText,bgColor=000000)",
            "customFontColor": "Placehold(font=Roboto,color=ff0000)",
            "customFontBgColor": "Placehold(font=Roboto,bgColor=000000)",
            "customColorBgColor": "Placehold(color=ff0000,bgColor=000000)",
            "allCustom": "Placehold(width=1024,height=768,text=CustomText,font=Roboto,color=ff0000,bgColor=000000)"
          },
      "LoremPicsum": {
          "default": "LoremPicsum()",
          "customWidth": "LoremPicsum(width=400)",
          "customHeight": "LoremPicsum(height=300)",
          "customBlur": "LoremPicsum(blur=2)",
          "customGrayscale": "LoremPicsum(grayscale=true)",
          "customWidthHeight": "LoremPicsum(width=400,height=300)",
          "customWidthBlur": "LoremPicsum(width=400,blur=2)",
          "customWidthGrayscale": "LoremPicsum(width=400,grayscale=true)",
          "customHeightBlur": "LoremPicsum(height=300,blur=2)",
          "customHeightGrayscale": "LoremPicsum(height=300,grayscale=true)",
          "customBlurGrayscale": "LoremPicsum(blur=2,grayscale=true)",
          "customWidthHeightBlur": "LoremPicsum(width=400,height=300,blur=2)",
          "customWidthHeightGrayscale": "LoremPicsum(width=400,height=300,grayscale=true)",
          "customWidthBlurGrayscale": "LoremPicsum(width=400,blur=2,grayscale=true)",
          "customHeightBlurGrayscale": "LoremPicsum(height=300,blur=2,grayscale=true)",
          "allCustom": "LoremPicsum(width=400,height=300,blur=2,grayscale=true)"
      }
 }
```

Response
``` json 
{
  "DiceBearImage": {
    "adventurerNeutral": "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=Ali Ratke",
    "avataaars": "https://api.dicebear.com/9.x/avataaars/svg?seed=Madison Bernier",
    "avataaarsNeutral": "https://api.dicebear.com/9.x/avataaars-neutral/svg?seed=Ardella Abshire DDS",
    "bigEars": "https://api.dicebear.com/9.x/big-ears/svg?seed=Savannah Hoeger Jr.",
    "bigEarsNeutral": "https://api.dicebear.com/9.x/big-ears-neutral/svg?seed=Bertram Koss",
    "bigSmile": "https://api.dicebear.com/9.x/big-smile/svg?seed=Ms. Dannie Smith",
    "bottts": "https://api.dicebear.com/9.x/bottts/svg?seed=Amparo Dach",
    "botttsNeutral": "https://api.dicebear.com/9.x/bottts-neutral/svg?seed=Jade Jacobi I",
    "croodles": "https://api.dicebear.com/9.x/croodles/svg?seed=Mr. Darien Boehm Sr.",
    "croodlesNeutral": "https://api.dicebear.com/9.x/croodles-neutral/svg?seed=Trey Stokes",
    "dylan": "https://api.dicebear.com/9.x/dylan/svg?seed=Mr. Chelsey Beer",
    "funEmoji": "https://api.dicebear.com/9.x/fun-emoji/svg?seed=Mr. Gideon Ortiz Jr.",
    "glass": "https://api.dicebear.com/9.x/glass/svg?seed=William Mante",
    "identicon": "https://api.dicebear.com/9.x/identicon/svg?seed=Ernestine Orn",
    "initials": "https://api.dicebear.com/9.x/initials/svg?seed=Ms. Domenica Mitchell II",
    "lorelei": "https://api.dicebear.com/9.x/lorelei/svg?seed=Nico Gutmann Jr.",
    "loreleiNeutral": "https://api.dicebear.com/9.x/lorelei-neutral/svg?seed=Armando Brown",
    "micah": "https://api.dicebear.com/9.x/micah/svg?seed=Mr. Ethel Jast MD",
    "miniavs": "https://api.dicebear.com/9.x/miniavs/svg?seed=Mr. Bobby Schiller",
    "notionists": "https://api.dicebear.com/9.x/notionists/svg?seed=Ms. Marion Cormier",
    "notionistsNeutral": "https://api.dicebear.com/9.x/notionists-neutral/svg?seed=Sarina Johnston IV",
    "openPeeps": "https://api.dicebear.com/9.x/open-peeps/svg?seed=Letitia Klocko",
    "personas": "https://api.dicebear.com/9.x/personas/svg?seed=Garrick Olson",
    "pixelArt": "https://api.dicebear.com/9.x/pixel-art/svg?seed=Enid Kulas",
    "pixelArtNeutral": "https://api.dicebear.com/9.x/pixel-art-neutral/svg?seed=Mr. Tony Glover I",
    "rings": "https://api.dicebear.com/9.x/rings/svg?seed=Mr. Jamie Ernser PhD",
    "shapes": "https://api.dicebear.com/9.x/shapes/svg?seed=Emile O'Conner",
    "thumbs": "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=Noel Grant III"
  },
  "Lorem": {
    "Lorem": "et distinctio nihil qui ullam quo quae in excepturi perferendis modi neque veniam sunt cum sint odio fuga nobis libero.",
    "paragraph": "vero beatae minus omnis quos explicabo est quod consequatur illum totam nostrum et dolores qui ipsum libero minus enim incidunt qui dolor odit perferendis et laborum dolor architecto sunt dignissimos voluptatem ipsum reiciendis quis possimus ut in nostrum iste nemo beatae consequuntur et totam est vitae a cum veritatis sit voluptas veniam nobis perferendis mollitia adipisci laudantium et itaque quia magnam necessitatibus aperiam molestias nobis sit doloribus et sint accusantium quo asperiores inventore dolorum optio nostrum ut non velit laborum accusamus dolorem sit nihil nisi qui rem.",
    "paragraphs": ["Two very long paragharphs that My IDE suffrs so removed them","Can easily be the size of a blog"],
    "sentence": "tempora assumenda quisquam quo dolores error sapiente fugiat aliquid sapiente vel vero illo illo natus neque non perferendis aut maiores.",
    "word": "vitae",
    "words": "quisquam non iusto quia quia praesentium rerum quia deserunt culpa totam unde recusandae nostrum atque rerum molestias quia itaque et."
  },
  "LoremPicsum": {
    "allCustom": "https://picsum.photos/400/300?blur=2&grayscale&id=883",
    "customBlur": "https://picsum.photos/200/200?blur=2&id=272",
    "customBlurGrayscale": "https://picsum.photos/200/200?blur=2&grayscale&id=410",
    "customGrayscale": "https://picsum.photos/200/200?blur=1&grayscale&id=634",
    "customHeight": "https://picsum.photos/200/300?blur=1&id=141",
    "customHeightBlur": "https://picsum.photos/200/300?blur=2&id=601",
    "customHeightBlurGrayscale": "https://picsum.photos/200/300?blur=2&grayscale&id=941",
    "customHeightGrayscale": "https://picsum.photos/200/300?blur=1&grayscale&id=814",
    "customWidth": "https://picsum.photos/400/200?blur=1&id=878",
    "customWidthBlur": "https://picsum.photos/400/200?blur=2&id=774",
    "customWidthBlurGrayscale": "https://picsum.photos/400/200?blur=2&grayscale&id=873",
    "customWidthGrayscale": "https://picsum.photos/400/200?blur=1&grayscale&id=655",
    "customWidthHeight": "https://picsum.photos/400/300?blur=1&id=794",
    "customWidthHeightBlur": "https://picsum.photos/400/300?blur=2&id=197",
    "customWidthHeightGrayscale": "https://picsum.photos/400/300?blur=1&grayscale&id=229",
    "default": "https://picsum.photos/200/300"
  },
  "Placehold": {
    "allCustom": "https://placehold.co/1024x768/ff0000/000000?text=CustomText&font=Roboto",
    "customBgColor": "https://placehold.co/600x400/000000/000000?text=Hello World&font=Montserrat",
    "customColor": "https://placehold.co/600x400/ff0000/FFF?text=Hello World&font=Montserrat",
    "customColorBgColor": "https://placehold.co/600x400/ff0000/000000?text=Hello World&font=Montserrat",
    "customFont": "https://placehold.co/600x400/000000/FFF?text=Hello World&font=Roboto",
    "customFontBgColor": "https://placehold.co/600x400/000000/000000?text=Hello World&font=Roboto",
    "customFontColor": "https://placehold.co/600x400/ff0000/FFF?text=Hello World&font=Roboto",
    "customHeight": "https://placehold.co/600x768/000000/FFF?text=Hello World&font=Montserrat",
    "customHeightBgColor": "https://placehold.co/600x768/000000/000000?text=Hello World&font=Montserrat",
    "customHeightColor": "https://placehold.co/600x768/ff0000/FFF?text=Hello World&font=Montserrat",
    "customHeightFont": "https://placehold.co/600x768/000000/FFF?text=Hello World&font=Roboto",
    "customHeightText": "https://placehold.co/600x768/000000/FFF?text=CustomText&font=Montserrat",
    "customText": "https://placehold.co/600x400/000000/FFF?text=CustomText&font=Montserrat",
    "customTextBgColor": "https://placehold.co/600x400/000000/000000?text=CustomText&font=Montserrat",
    "customTextColor": "https://placehold.co/600x400/ff0000/FFF?text=CustomText&font=Montserrat",
    "customTextFont": "https://placehold.co/600x400/000000/FFF?text=CustomText&font=Roboto",
    "customWidth": "https://placehold.co/1024x400/000000/FFF?text=Hello World&font=Montserrat",
    "customWidthBgColor": "https://placehold.co/1024x400/000000/000000?text=Hello World&font=Montserrat",
    "customWidthColor": "https://placehold.co/1024x400/ff0000/FFF?text=Hello World&font=Montserrat",
    "customWidthFont": "https://placehold.co/1024x400/000000/FFF?text=Hello World&font=Roboto",
    "customWidthHeight": "https://placehold.co/1024x768/000000/FFF?text=Hello World&font=Montserrat",
    "customWidthText": "https://placehold.co/1024x400/000000/FFF?text=CustomText&font=Montserrat",
    "default": "https://placehold.co/600x400"
  },
  "Time": {
    "ansi": "Tue Nov 28 19:31:27 1989",
    "day": 6,
    "iso": "1991-03-25T13:16:00+000",
    "monthName": "March",
    "time": "Tue, 10 Sep 2024 21:58:53 +0530",
    "timezone": "America/Santarem",
    "unix": 1725985733,
    "unixNano": 1725985733036314000
  },
  "address": {
    "address": "%969 Christ Vista Apt. 669\nEast Rozella, PA 04700-0536",
    "city": "South Laurence",
    "country": "Belize",
    "countryAbbr": "FIN",
    "country_code": "JM",
    "latitude": 62.211554,
    "longitude": 73.703326,
    "secondary": "%115 Vandervort Expressway\nBlandaland, MI 77383-4643",
    "state": "New York",
    "state_code": "LA",
    "street": "Rosalee Mews",
    "zip": "20659"
  },
  "app": {
    "appName": "WebDesk",
    "appVersion": "v7.0.6",
    "platform": "Android"
  },
  "color": {
    "color": "LightSeaGreen",
    "css": "rgb(45,207,248)",
    "hex": "#1B20F2",
    "rgb": "54,86,19",
    "rgba": [
      "207",
      "35",
      "11"
    ],
    "safe": "maroon"
  },
  "company": {
    "bs": "empower viral users",
    "catch_phrase": "Persevering local success",
    "company": "Kuhn and Sons",
    "ein": 95,
    "jobTitle": "Project Manager",
    "mail": "ruth.douglas@nicolas_inc.ipz.com",
    "suffix": "and Sons"
  },
  "finance": {
    "amount": 1428960,
    "amountWithCurrency": "23849 USD",
    "btcAddress": "36Rd57WDuV75u3qJt4yvpx4QxkWC",
    "card_type": "Discover Card",
    "credit_card": "4119788493747977",
    "currency": "Jamaican Dollar",
    "currencyAndCode": "Sri Lanka Rupee (LKR)",
    "currencyCode": "XDR",
    "ethAddress": "0xF4VJM7zm29nxYvspy69oAy8Hp2hl8F2dQBp58c7q",
    "exp": "15/16",
    "iban": "CZ8989836920895392395333"
  },
  "internet": {
    "domain": "yrm.com",
    "free_email": "charlotte.nolan@yahoo.com",
    "httpMethod": "TRACE",
    "ip": "5.239.104.19",
    "ipv6": "3446:1151:8277:5363:6356:2883:3238:6773",
    "mac": "0D:D4:8B:AF:A5:2B",
    "md5": "92c9713a49820b2c97fbd10165872df4",
    "safe_email": "jarret.halvorson@example.org",
    "sha256": "863a6b632034d24c2d404f944dde121158683787b89cd065fd00d5fa76f68853",
    "sha512": "f0bcfd81abefcdf9ae5e5de58d1a868317503ea76422309bc212d1ef25a1e67789d0bfa752a7e2abd4510f4f3e4f60cdaf6202a42883fb97bb7110ab3600785e",
    "slug": "bj-mugya",
    "sqlId": 74042,
    "status_code": 508,
    "status_code_message": "Continue",
    "tld": "com",
    "url": "https://www.wfc.info/zpgi-hdiz",
    "user_agent": "Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Safari/605.1"
  },
  "language": {
    "language": "Romanian",
    "languageabbr": "zu",
    "programming": "Java"
  },
  "misc": {
    "array": [
      "minus",
      "hic",
      "nulla",
      "accusamus",
      "vitae"
    ],
    "array with len": [
      "Lurline Steuber Jr.",
      "Mr. Leopoldo Waelchi Jr.",
      "Ms. Karelle Abernathy II",
      "Imogene Windler",
      "Wendell Rolfson",
      "Corine Doyle",
      "Ms. Leonor Gutkowski",
      "Manuela Hickle",
      "Isac Strosin",
      "Mr. Rylan Senger MD",
      "Junius McGlynn",
      "Ms. Destiny Gislason I",
      "Marcel Conn",
      "Hunter Brekke",
      "Ms. Tess Considine III",
      "Samir Schinner",
      "Marge Koss",
      "Gilberto Howell",
      "Ms. Viola Sporer",
      "Sydni Wehner"
    ],
    "array with len and type": [
      "hayes.nasir@yahoo.com",
      "trantow.adrianna@gmail.com",
      "dorcas.bergstrom@ukg.org",
      "fritsch.jayde@yahoo.com",
      "stoltenberg.margarett@mao.biz",
      "joey.kuvalis@xyh.info",
      "champlin@hotmail.com",
      "cronin@izg.com",
      "ritchie@gmail.com",
      "botsford.marcelo@hotmail.com",
      "ida@yahoo.com",
      "oda.pagac@yahoo.com",
      "borer.al@gmail.com",
      "elyssa@bec.com",
      "bergstrom@gmail.com",
      "davin.leuschke@hotmail.com",
      "smith.mohammad@wiv.net",
      "collier.perry@shu.com",
      "broderick.jakubowski@qek.net",
      "huels.magali@yej.biz"
    ],
    "boolean": true,
    "float": 483.20001220703125,
    "intSmall": 3,
    "integer": 200,
    "json": {
      "enim": [
        "incidunt",
        "quis",
        "voluptatem",
        "fugit"
      ],
      "facere": [
        "quam",
        "asperiores",
        "et"
      ],
      "iusto": [
        "harum",
        "minima",
        "repudiandae",
        "sequi",
        "atque",
        "quibusdam"
      ],
      "nihil": 844023.1875,
      "qui": [
        "distinctio",
        "deleniti",
        "quia",
        "ut",
        "dolorem",
        "sit"
      ],
      "sapiente": 714779.1875
    },
    "message": "Any custom message you want",
    "uuid": "22168e9a-891b-40ff-8d14-d16c92481d4e"
  },
  "user": {
    "age": 59,
    "bio": "placeat quibusdam nesciunt odit facilis deserunt non magni blanditiis autem odio voluptatem sint nostrum voluptatem vel officia incidunt vitae vero.",
    "birthday": "2029-02-07T21:58:53.0352641+05:30",
    "email": "lulu@yahoo.com",
    "female_firstName": "Cydney",
    "female_title": "Ms.",
    "firstName": "Willow",
    "gamer_tag": "VagaBond",
    "gender": "Female",
    "image": "https://randomuser.me/api/portraits/med/men/9.jpg",
    "lastName": "Harris",
    "maleFirstName": "Jermey",
    "maleTitle": "Mr.",
    "password": "ouwlbmqpt",
    "phone": "840-310-9781 x335",
    "ssn": "212696526",
    "title": "Ms.",
    "user": "Mr. Lazaro Renner",
    "userName": "alena",
    "user_id": "66645705-adf5-4582-ad93-f93ef378d9b8"
  },
  "vehicle": {
    "vehicle": "Minivan",
    "vehicleBrand": "Chevrolet",
    "vehiclePlate": "ВX7279KР",
    "vehicleTransmission": "Semi-auto",
    "vehicleType": "Bio Gas"
  }
}
```

### Things to be implemented
- [ ] More Mock Data Methods (If needed)
- [ ] Better Error Handling and Logging
- [ ] Better Documentation
- [ ] Code Cleanup
- [ ] Adding test
- [ ] Presistent Data

### Contributing
Feel free to do anything It's my first project to shitty code is expected.  : )
If ya find any bugs or have any suggestions feel free to open an issue or a PR. and please explain the issue or the PR in detail. It'll help me understand better.


