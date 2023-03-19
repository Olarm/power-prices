package main

import(
    "encoding/json"
    "net/http"
    "io/ioutil"
    "fmt"
    "time"
)


type Hour struct {
    Eur       float64 `json:"EUR_per_kWh"`
    TimeStop  time.Time `json:"time_end"`
    Exr       float64 `json:"EXR"`
    TimeStart time.Time `json:"time_start"`
    Nok       float64 `json:"NOK_per_kWh"`
}


func getDay(day, month, year, zone string) ([]Hour, error) {
    url := fmt.Sprintf("https://www.hvakosterstrommen.no/api/v1/prices/%s/%s-%s_%s.json", year, month, day, zone)
    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("Error getting data: %+v", err)
        return nil, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body) // response body is []byte
    if err != nil {
        fmt.Printf("error reading byte array: %+v\n", err)
    }
    var prices []Hour
    if err := json.Unmarshal(body, &prices); err != nil {   // Parse []byte to go struct pointer
        fmt.Println("Can not unmarshal JSON")
    }
    return prices, nil
}

func getToday(zone string) ([]Hour, error) {
    now := time.Now()
    y := fmt.Sprintf("%04d", now.Year())
    m := fmt.Sprintf("%02d", now.Month())
    d := fmt.Sprintf("%02d", now.Day())
    prices, _ := getDay(d, m, y, zone)
    return prices, nil
}

func getTodayArray(zone string) ([24]float64, error) {
    prices, _ := getToday(zone)
    var pricesNok [24]float64
    for i, hour := range prices {
        pricesNok[i] = hour.Nok
    }
    return pricesNok, nil
}

func minMax(prices [24]float64) (float64, float64) {
    min := prices[0]
    max := prices[0]
    for _, price := range prices[1:24] {
        if price < min {
            min = price
        } else if price > max {
            max = price
        }
    }
    return min, max
}

func main() {
    pricesNok, _ := getTodayArray("NO5")
    min, max := minMax(pricesNok)
    diff := max - min
    p25 := diff / 100 * 25 + min
    p75 := diff / 100 * 75 + min

    fmt.Println(PrettyPrint(pricesNok))
    fmt.Printf("max: %f, min: %f\n", max, min)
    fmt.Printf("p25: %f, p75: %f\n", p25, p75)

    //now := time.Now()
    //hour := now.Hour()


}


func PrettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}
