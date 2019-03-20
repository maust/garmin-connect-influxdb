# Garmin Connect data converter to InfluxDB
This can convert your garmin connect data and store it in InfluxDB. So you can easily use e.g. Grafana to build some dashboards or do whatever you want.

## Requirements
* Go installed to build / run it
* InfluxDB running on localhost:8086
* optional: Grafana

## Steps to do
1. Request your data from Garmin Connect: https://www.garmin.com/en-GB/account/datamanagement/exportdata/
2. Clone this repo
3. Extract the export (zip) into the subfolder export
4. Run it: 
    $ go run main.go
