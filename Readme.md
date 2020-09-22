# LATA (Local Access Transport Area) parser 

This app was created to prepare database for ruby gem `city_phone_npa_npx`

Parser is using information provided by (https://localcallingguide.com)[https://localcallingguide.com] in XML format.

As output we have list NPA and NPX assigned to city in YAML format.

# Running

```
go run parser.go
```

# File structure

`data` - folder contain data from localcallingguide.com
`db` - folder contain prepared USA and Canada states
`output` - will contain prepared `yml` files after running application
`parser.go` - main function