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

#### Endpoints

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
##### Nested Data
 
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


#### Endpoints Controls

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

<i> Right now the Auth do not offer much customization but if needed we will add more options in the future. Right now this was the most basic setup we could think of. </i>

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


## Ref 
A lot of endpoints can share similar Response Schema Or you might need an array of Custom data. You can easily solve this problem by using the Ref key. The Ref key will allow you to reference a response object from reference.

This allows you to declare Most of the common Schema in one Place then just reference it where-ever needed.

<i> Both `ref.ObjectName` and `ObjectName` are identified as a valid argument</i>

``` json
{
  "config":{ },
  "ref":{
    "UserData":{
      "email":"User(email)",
      "username":"User(username)",
      "password":"User(password)",
      "uuid":"Uuid()",
      "nesting":{
        "testing":"Message(This is just to test)",
      }
    }
  },
  "api":[
    {
      "path":"/collections",
      "method":"GET",
      "response":{
        "test":"Array(len=2,type=Ref(ref.UserData))",
        "test2":"Ref(UserData)",
      }
    }
   
  ]
}
```
#### Response

```json 

{
  "test": [
    {
      "email": "mcdermott@lyx.net",
      "nesting": {
        "testing": "This is just to test"
      },
      "password": "lqxnmq",
      "username": "vito.lowe",
      "uuid": "539752a5-c99e-4187-8108-884184464fe7"
    },
    {
      "email": "madonna.hettinger@lil.com",
      "nesting": {
        "testing": "This is just to test"
      },
      "password": "rqxkn|",
      "username": "alanis.paucek",
      "uuid": "9f5ce36b-962a-4399-bfee-7f50836c1efa"
    }
  ],
  "test2": {
    "email": "runolfsdottir.monica@yahoo.com",
    "nesting": {
      "testing": "This is just to test"
    },
    "password": "~sydshcjexbroh",
    "username": "wunsch",
    "uuid": "4ec20bb2-9381-438c-be8f-87ea43c73e8b"
  }
}
```


## Mock Data Methods
The mock data is powered by [Faker](https://github.com/jaswdr/faker) and tries to provide access to most of the methods provided by the library in a simple way.

## Available Methods
| Method | Description |
| --- | --- |
| User() | Used to mock user-related data such as email, name, phone, etc. |
| Address() | Used to mock address-related data such as street, city, country, etc. |
| Finance() | Returns finance-related data such as crypto, banks, etc. |
| Time() | Returns mock time data such as current time, date, etc. |
| Vehicle() | Returns mock vehicle data such as make, model, year, etc. |
| App() | Returns mock application-related data such as app name, version, etc. |
| Language() | Returns mock language-related data such as language name, programming language, etc. |
| Bool() | Returns a random boolean value (true or false) randomly, If you need one or Other You can do `Bool(true)` |
| Float() | Returns a random float value, Accepts Two optional Parameters `min` `max` `Float(min=0,max=1)` will only return value bw 0 and 1 |
| Int() | Returns a random integer value Can Accepts two args `min` `max` `Int(min=10,max=20)`  |
| Json() | Returns mock JSON data |
|Uuid() | Returns a random UUID |
| Array() | Returns a mock array of data also accepts two inputs `len` for array of len and `type`  for the type of data to be returned ie `Array(len=2000, type=User())` or `Array(len=2000, type=Ref(ref.UserData))`  |
| Lorem() | Returns mock lorem ipsum text. Accepts two argument `len` for the length and `type` for the type Like `word` `sentence` `sentences` `paragraph` `paragraphs` Each returning more data than before.  |
| Internet() | Returns mock internet-related data such as IP address, URL, etc. |
| Company() | Returns mock company-related data such as company name, BS, etc. |
| Color() | Returns mock color data such as color name, hex code, etc. |
| DiceBearImage() | Returns a mock image URL from DiceBear. Accepts one parameter Collection to generate a image from that collection DiceBearImage(pixel-art) |
| Message() | Returns a custom message. Accepts one parameter that will be returned as the message. `Message(Hello World)` ps any string that is not a defined function will be returned as it is. This method is just for the sake of uniformity you can just do `message: "any custom message"` |
| Placehold() | Returns a placehol image URL. Accepts multiple parameters `Placehold(width=1024,height=768,text=CustomText,font=Roboto,color=ff0000,bgColor=000000)` all of them are optional  ` |
| LoremPicsum() | Returns a Lorem Picsum Image URL. Accepts multiple parameters `LoremPicsum(width=400,height=300,blur=2,grayscale=true)` all of them are optional  ` |

<i> Note: If a function only accepts one parameter it can be unnamed But if It accepts More than 1 parameter they have to be nammed  ie `User(username)` as user only accepts one param. `Array(len=20)` as Array accepts two param this is done so It's easy to read and Understand the JSON Otherwise It would be `Int(20,40)` which is hard to understand in a glance instead `Int(min=20,max=40)` Is a bit more to type but easy to understand. and it was easy to code XD  </i>

Most of the methods above accepts one parameter that can be used to fine-tune the data. For example, `User(name)` will return a random name and `User(mFirstName)` will return a male first name. 

# A list of everything 
``` json 
{
  "config": {
    "prefix": "/api",
    "latency": 200,
    "Port": 3000,
    "logging": true,
    "failRate": 0,
    "timeout": 1000,
    "auth": {
      "apiKey": "1234567890",
      "bearer": "1234567890",
      "cookie": "auth=1234567890"
    }
  },
  "ref": {
    "reusableUser": {
      "user": "User()",
      "email": "User(email)",
      "firstName": "User(firstName)",
      "lastName": "User(lastName)"
    }
  },
  "api": [
    {
      "path": "/all/data",
      "latency": 200,
      "method": "GET",
      "failRate": 0,
      "auth": {
        "protected": true,
        "protectedBy": "apiKey"
      },
      "response": {
        "color": {
          "color": "Color()",
          "hex": "Color(hex)",
          "rgb": "Color(rgb)",
          "rgba": "Color(rgba)",
          "safe": "Color(safe)",
          "css": "Color(css)"
        },
        "address": {
          "address": "Address()",
          "city": "Address(city)",
          "state": "Address(state)",
          "country": "Address(country)",
          "zip": "Address(zip)",
          "latitude": "Address(latitude)",
          "longitude": "Address(longitude)",
          "secondary": "Address(full)",
          "street": "Address(street)",
          "country_code": "Address(countryCode)",
          "state_code": "Address(stateAbbr)",
          "countryAbbr": "Address(countryAbbr)"
        },
        "company": {
          "company": "Company()",
          "catch_phrase": "Company(catchPhrase)",
          "bs": "Company(bs)",
          "jobTitle": "Company(jobTitle)",
          "ein": "Company(ein)",
          "suffix": "Company(suffix)",
          "mail": "Company(mail)"
        },
        "internet": {
          "url": "Internet()",
          "domain": "Internet(domain)",
          "ip": "Internet(ip)",
          "ipv6": "Internet(ipv6)",
          "mac": "Internet(mac)",
          "httpMethod": "Internet(httpMethod)",
          "tld": "Internet(tld)",
          "slug": "Internet(slug)",
          "status_code": "Internet(statusCode)",
          "free_email": "Internet(freeEmail)",
          "safe_email": "Internet(safeEmail)",
          "status_code_message": "Internet(statusCodeMessage)",
          "user_agent": "Internet(userAgent)",
          "sha256": "Internet(sha256)",
          "md5": "Internet(md5)",
          "sha512": "Internet(sha512)",
          "sqlId": "Internet(sqlId)"
        },
        "user": {
          "user": "User()",
          "email": "User(email)",
          "firstName": "User(firstName)",
          "lastName": "User(lastName)",
          "female_firstName": "User(fFirstName)",
          "female_title": "User(fTitle)",
          "maleFirstName": "User(mFirstName)",
          "maleTitle": "User(mTitle)",
          "phone": "User(phone)",
          "userName": "User(userName)",
          "password": "User(password)",
          "title": "User(title)",
          "gender": "User(gender)",
          "ssn": "User(ssn)",
          "bio": "User(bio)",
          "birthday": "User(birthday)",
          "age": "Int(min=20,max=100)",
          "gamer_tag": "User(gamerTag)",
          "image": "User(image)",
          "user_id": "Uuid()"
        },
        "finance": {
          "amount": "Finance()",
          "credit_card": "Finance(creditCard)",
          "card_type": "Finance(cardType)",
          "exp": "Finance(cardExpirationDate)",
          "iban": "Finance(iban)",
          "currency": "Finance(currency)",
          "currencyCode": "Finance(currencyCode)",
          "currencyAndCode": "Finance(currencyAndCode)",
          "amountWithCurrency": "Finance(amountWithCurrency)",
          "btcAddress": "Finance(btcAddress)",
          "ethAddress": "Finance(ethAddress)"
        },
        "Lorem": {
          "word": "Lorem(type=word,len=1)",
          "words": "Lorem(type=words,len=2)",
          "sentence": "Lorem(type=sentence,len=2)",
          "sentences": "Lorem(type=sentences,len=2)",
          "paragraph": "Lorem(type=paragraph,len=2)",
          "paragraphs": "Lorem(type=paragraphs,len=1)",

          "Lorem": "Lorem()"
        },
        "misc": {
          "random bool": "Bool()",
          "always true": "Bool(true)",
          "random float": "Float()",
          "float": "Float(min=0,max=1)",
          "random integer": "Int()",
          "int": "Int(min=20,max=100)",
          "message": "Message(Any custom message you want)",
          "json": "Json()",
          "array": "Array()",
          "array with len": "Array(len=20)",
          "array with len and type": "Array(len=20,type=User(email))",
          "uuid": "Uuid()",
          "ref": "Ref(ref.reusableUser)",
          "array with custom data": "Array(len=2,type=Ref(ref.reusableUser))"
        },
        "language": {
          "language": "Language()",
          "languageabbr": "Language(abbr)",
          "programming": "Language(programming)"
        },
        "app": {
          "appName": "App()",
          "appVersion": "App(version)",
          "platform": "App(platform)"
        },
        "vehicle": {
          "vehicle": "Vehicle()",
          "vehicleBrand": "Vehicle(brand)",
          "vehicleType": "Vehicle(type)",
          "vehicleTransmission": "Vehicle(transmission)",
          "vehiclePlate": "Vehicle(plate)"
        },
        "Time": {
          "time": "Time()",
          "unix": "Time(unix)",
          "unixNano": "Time(unixNano)",
          "iso": "Time(iso)",
          "day": "Time(day)",
          "ansi": "Time(ansi)",
          "monthName": "Time(monthName)",
          "timezone": "Time(timezone)",
          "year": "Time(year)",
          "month": "Time(month)",
          "months": "Time(months)",
          "days": "Time(days)"
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
    }
  ]
}


```

Response
``` json 
{
  "DiceBearImage": {
    "adventurerNeutral": "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=Jameson Wisoky",
    "avataaars": "https://api.dicebear.com/9.x/avataaars/svg?seed=Ollie Pfannerstill",
    "avataaarsNeutral": "https://api.dicebear.com/9.x/avataaars-neutral/svg?seed=Ms. Violette Gislason",
    "bigEars": "https://api.dicebear.com/9.x/big-ears/svg?seed=Art Moore",
    "bigEarsNeutral": "https://api.dicebear.com/9.x/big-ears-neutral/svg?seed=Ms. Kayli Adams DDS",
    "bigSmile": "https://api.dicebear.com/9.x/big-smile/svg?seed=Eleazar Bartoletti",
    "bottts": "https://api.dicebear.com/9.x/bottts/svg?seed=Ms. Sister Murray",
    "botttsNeutral": "https://api.dicebear.com/9.x/bottts-neutral/svg?seed=Jerrold Fritsch",
    "croodles": "https://api.dicebear.com/9.x/croodles/svg?seed=Mr. Moises Schmeler DDS",
    "croodlesNeutral": "https://api.dicebear.com/9.x/croodles-neutral/svg?seed=Obie Mosciski",
    "dylan": "https://api.dicebear.com/9.x/dylan/svg?seed=Alford Baumbach",
    "funEmoji": "https://api.dicebear.com/9.x/fun-emoji/svg?seed=Ms. Pasquale Collier",
    "glass": "https://api.dicebear.com/9.x/glass/svg?seed=Ms. Grace Frami IV",
    "identicon": "https://api.dicebear.com/9.x/identicon/svg?seed=Edd Heathcote",
    "initials": "https://api.dicebear.com/9.x/initials/svg?seed=Aron Haley I",
    "lorelei": "https://api.dicebear.com/9.x/lorelei/svg?seed=Angelica Funk",
    "loreleiNeutral": "https://api.dicebear.com/9.x/lorelei-neutral/svg?seed=Ms. Aylin Wisozk PhD",
    "micah": "https://api.dicebear.com/9.x/micah/svg?seed=Skylar Sauer",
    "miniavs": "https://api.dicebear.com/9.x/miniavs/svg?seed=Ms. Assunta Braun I",
    "notionists": "https://api.dicebear.com/9.x/notionists/svg?seed=Dario Cassin",
    "notionistsNeutral": "https://api.dicebear.com/9.x/notionists-neutral/svg?seed=Yasmine Grimes",
    "openPeeps": "https://api.dicebear.com/9.x/open-peeps/svg?seed=Antonina Boyer",
    "personas": "https://api.dicebear.com/9.x/personas/svg?seed=Annamae Murphy",
    "pixelArt": "https://api.dicebear.com/9.x/pixel-art/svg?seed=Nikki Kunde",
    "pixelArtNeutral": "https://api.dicebear.com/9.x/pixel-art-neutral/svg?seed=Elfrieda Howell",
    "rings": "https://api.dicebear.com/9.x/rings/svg?seed=Ms. Lilla Mills DVM",
    "shapes": "https://api.dicebear.com/9.x/shapes/svg?seed=Margarita Volkman",
    "thumbs": "https://api.dicebear.com/9.x/adventurer-neutral/svg?seed=Ms. Felipa Jacobi"
  },
  "Lorem": {
    "Lorem": "et.",
    "paragraph": "sed ut voluptatibus perferendis aliquam maxime nulla voluptas ipsum autem enim adipisci dicta delectus cumque animi eius eum sit cumque reiciendis a quo esse delectus harum dolor quis vel et ullam est a neque omnis praesentium iusto aut numquam dignissimos aut autem quasi minus assumenda aut cum illo nisi nostrum sit autem neque. ut deleniti ex doloribus ullam est et molestiae eos quasi quia provident rerum dignissimos est qui natus quisquam illo ut vero est voluptas sequi eos consequatur velit non dolorem sunt et natus dolor esse rerum et accusamus quis dolor ab distinctio ut expedita quia nostrum rerum fuga aspernatur placeat adipisci in illum et eos et odio nisi.",
    "paragraphs": [
      "Over 500 words here removed them for the sake of readability",
    ],
    "sentence": "labore maiores.",
    "sentences": [
      "earum quam deserunt saepe quod dolorum aperiam vel nihil voluptas et velit occaecati ut qui ipsa dolorem omnis consequuntur mollitia nobis totam recusandae animi qui quam autem temporibus distinctio minima id non eos eos mollitia ut eum dolorem aut nobis sit voluptatem quam qui animi temporibus corrupti quos cupiditate qui voluptatem rerum qui assumenda labore consequuntur temporibus sint ducimus ut enim eos ex aut perspiciatis nobis inventore laboriosam dolor expedita sit facilis vel totam omnis minus sit itaque necessitatibus.",
      "molestiae dolore nemo recusandae accusantium suscipit qui dolor ea suscipit saepe accusamus est ut perspiciatis deleniti itaque dolor vitae voluptatem perspiciatis necessitatibus sed distinctio voluptatem pariatur dolores et vel culpa praesentium expedita saepe quia doloribus fuga aspernatur autem sed amet illo quasi voluptas aut fugiat possimus soluta fugit dolorem qui et quibusdam consequatur temporibus sed soluta alias nam non."
    ],
    "word": "necessitatibus",
    "words": [
      "aut",
      "iure"
    ]
  },
  "LoremPicsum": {
    "allCustom": "https://picsum.photos/400/300?blur=2&grayscale&id=186",
    "customBlur": "https://picsum.photos/200/200?blur=2&id=910",
    "customBlurGrayscale": "https://picsum.photos/200/200?blur=2&grayscale&id=718",
    "customGrayscale": "https://picsum.photos/200/200?blur=1&grayscale&id=37",
    "customHeight": "https://picsum.photos/200/300?blur=1&id=404",
    "customHeightBlur": "https://picsum.photos/200/300?blur=2&id=545",
    "customHeightBlurGrayscale": "https://picsum.photos/200/300?blur=2&grayscale&id=201",
    "customHeightGrayscale": "https://picsum.photos/200/300?blur=1&grayscale&id=712",
    "customWidth": "https://picsum.photos/400/200?blur=1&id=206",
    "customWidthBlur": "https://picsum.photos/400/200?blur=2&id=99",
    "customWidthBlurGrayscale": "https://picsum.photos/400/200?blur=2&grayscale&id=205",
    "customWidthGrayscale": "https://picsum.photos/400/200?blur=1&grayscale&id=0",
    "customWidthHeight": "https://picsum.photos/400/300?blur=1&id=353",
    "customWidthHeightBlur": "https://picsum.photos/400/300?blur=2&id=65",
    "customWidthHeightGrayscale": "https://picsum.photos/400/300?blur=1&grayscale&id=526",
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
    "ansi": "Sun Nov  1 02:53:49 1981",
    "day": 3,
    "days": [
      "Sunday",
      "Monday",
      "Tuesday",
      "Wednesday",
      "Thursday",
      "Friday",
      "Saturday"
    ],
    "iso": "1980-11-12T21:45:51+000",
    "month": 5,
    "monthName": "January",
    "months": [
      "January",
      "February",
      "March",
      "April",
      "May",
      "June",
      "July",
      "August",
      "September",
      "October",
      "November",
      "December"
    ],
    "time": "Thu, 12 Sep 2024 02:49:01 +0530",
    "timezone": "America/Yakutat",
    "unix": 1726089541,
    "unixNano": 1726089541416112000,
    "year": 1992
  },
  "address": {
    "address": "%1437 Angelica Highway\nNorth Marguerite, VT 89836",
    "city": "East Jeffry",
    "country": "Bahrain",
    "countryAbbr": "RUS",
    "country_code": "IQ",
    "latitude": 22.037837,
    "longitude": 13.150227,
    "secondary": "%91 Major Flats Apt. 416\nHomenickhaven, LA 67248-4736",
    "state": "Florida",
    "state_code": "VA",
    "street": "Schaefer Cape",
    "zip": "81810-4769"
  },
  "app": {
    "appName": "Essential Web",
    "appVersion": "v5.3.0",
    "platform": "Windows"
  },
  "color": {
    "color": "GreenYellow",
    "css": "rgb(156,82,158)",
    "hex": "#0BA0BB",
    "rgb": "186,216,107",
    "rgba": [
      "143",
      "39",
      "2"
    ],
    "safe": "teal"
  },
  "company": {
    "bs": "incentivize mission-critical portals",
    "catch_phrase": "Open-architected upward-trending opensystem",
    "company": "Schinner Ltd",
    "ein": 12,
    "jobTitle": "Immigration Inspector OR Customs Inspector",
    "mail": "denesik.christina@mann-mann.tfy.com",
    "suffix": "Inc"
  },
  "finance": {
    "amount": 5412186,
    "amountWithCurrency": "63071 PKR",
    "btcAddress": "bc15dqZEG5k8sE85N7A2pZQTq1",
    "card_type": "Visa",
    "credit_card": "0674127062363508",
    "currency": "Euro",
    "currencyAndCode": "US Dollar (USD)",
    "currencyCode": "EUR",
    "ethAddress": "0xw0MDQGLYJU1c2AySbm4YC4sETwl5ZWz4gq7033jL",
    "exp": "00/17",
    "iban": "CR24764567723760078084"
  },
  "internet": {
    "domain": "zyj.biz",
    "free_email": "sydney@hotmail.com",
    "httpMethod": "DELETE",
    "ip": "151.152.176.105",
    "ipv6": "7348:3336:3327:5732:5847:2617:8443:4145",
    "mac": "4C:D3:F8:75:D0:12",
    "md5": "a81ba73d43118ec091f856197a2d0e8d",
    "safe_email": "haley@example.org",
    "sha256": "f22e89c251e50eec7c9e184f5584b99fdf48fdbde7317acd6388ce05ea691a43",
    "sha512": "708a0ac37986b927ab8de7c99044f5e06f0751765b974f2fc6b4ce896c6ee06a526e03d44304e0f8c2bfd53caedf856b1596edee3d55495f799d988b153a4bfa",
    "slug": "ney-nnnh",
    "sqlId": 96563,
    "status_code": 413,
    "status_code_message": "Locked (WebDAV)",
    "tld": "com",
    "url": "https://www.oyl.org/wmag-upyw",
    "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/117.0"
  },
  "language": {
    "language": "Magahi",
    "languageabbr": "hi",
    "programming": "Ruby on Rails"
  },
  "misc": {
    "always true": true,
    "array": [
      "tenetur",
      "magni",
      "omnis",
      "quis",
      "nihil"
    ],
    "array with custom data": [
      {
        "email": "janae@yahoo.com",
        "firstName": "Reginald",
        "lastName": "Toy",
        "user": "Ashton Turner DVM"
      },
      {
        "email": "kamron@yahoo.com",
        "firstName": "Timmothy",
        "lastName": "Berge",
        "user": "Francisca Paucek DDS"
      }
    ],
    "array with len": [
      "Allan Stark",
      "Mariela Hayes",
      "Liana Rodriguez",
      "Colt Little",
      "Eda Barton",
      "Mr. Edwin Beier PhD",
      "Brenden Klein PhD",
      "Mr. Edgar Ondricka DDS",
      "Adella Flatley",
      "Ms. Luz Donnelly III",
      "Mr. Kale Effertz",
      "Jerrod Lueilwitz",
      "Yadira VonRueden",
      "Nikita Predovic MD",
      "Ms. Greta Morissette",
      "Cedrick Gorczany",
      "Jorge Schaefer",
      "Edd Erdman III",
      "Mr. Jaeden Kuphal",
      "Jamison Lockman"
    ],
    "array with len and type": [
      "hildegard.parisian@hotmail.com",
      "tatyana@fiv.org",
      "mayer.caden@hah.com",
      "danial@amo.org",
      "ryan@yahoo.com",
      "enos@ymz.org",
      "joe.hickle@gmail.com",
      "darron@yahoo.com",
      "o_keefe@yahoo.com",
      "pacocha@gmail.com",
      "jonas@hotmail.com",
      "candice@hotmail.com",
      "laury@egy.info",
      "ondricka@gmail.com",
      "jessie@cdz.biz",
      "lilla@gmail.com",
      "joshuah@yahoo.com",
      "grimes@yahoo.com",
      "kozey@gmail.com",
      "lia.parker@hotmail.com"
    ],
    "float": 0.10000000149011612,
    "int": 91,
    "json": {
      "illum": "%33 Gladyce Lodge Suite 048\nLake Nestor, WI 25169",
      "iste": "%39 Kozey Views Suite 396\nWest Jazminbury, NM 68509-8371",
      "quod": {
        "aut": 6799611
      },
      "sunt": 9186739
    },
    "message": "Any custom message you want",
    "random bool": false,
    "random float": 141.1999969482422,
    "random integer": 73,
    "ref": {
      "email": "considine.bulah@wia.com",
      "firstName": "Geovany",
      "lastName": "Upton",
      "user": "King Rohan MD"
    },
    "uuid": "c9d1f77f-b8c0-46f3-a07b-d2749ccacf2f"
  },
  "user": {
    "bio": "asperiores omnis veniam sunt porro rerum velit consequatur et ut nesciunt itaque eveniet fugit libero amet facere distinctio quia illum.",
    "birthday": "2021-02-10T02:49:01.4166387+05:30",
    "email": "moore.justine@hotmail.com",
    "female_firstName": "Vallie",
    "female_title": "Ms.",
    "firstName": "Edd",
    "gamer_tag": "Grave",
    "gender": "Female",
    "image": "https://randomuser.me/api/portraits/med/men/0.jpg",
    "lastName": "Renner",
    "maleFirstName": "Edgardo",
    "maleTitle": "Mr.",
    "password": "u}hbcfnsmv",
    "phone": "(450) 550-2879",
    "ssn": "983266484",
    "title": "Ms.",
    "user": "Carolyne Willms",
    "userName": "Mr. Larue Considine V",
    "user_id": "972503c5-91af-4ab1-a6b6-03ef8c8e5c5b"
  },
  "vehicle": {
    "vehicle": "Pickup",
    "vehicleBrand": "Maserati",
    "vehiclePlate": "ВO5892OМ",
    "vehicleTransmission": "Tiptronic",
    "vehicleType": "Hybrid"
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


