/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : map.go

* Purpose :

* Creation Date : 04-12-2017

* Last Modified : Thu 13 Apr 2017 02:12:53 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type googleData struct {
	Key string
}

const googlemap = `<!DOCTYPE html>
<html>
  <head>
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no">
    <meta charset="utf-8">
    <title>Marker Clustering</title>
    <style>
      /* Always set the map height explicitly to define the size of the div
       * element that contains the map. */
      #map {
        height: 100%;
      }
      /* Optional: Makes the sample page fill the window. */
      html, body {
        height: 100%;
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <div id="map"></div>
<script
  src="https://code.jquery.com/jquery-3.2.1.min.js"
  integrity="sha256-hwg4gsxgFZhOsEEamdOYGBf13FyQuiTwlAQgxVSNgt4="
  crossorigin="anonymous"></script>
    <script>

      function initMap() {

        var map = new google.maps.Map(document.getElementById('map'), {
          zoom: 2,
          center: {lat: -28.024, lng: 140.887}
        });

        // Create an array of alphabetical characters used to label the markers.
        var labels = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

        // Add some markers to the map.
        // Note: The code uses the JavaScript Array.prototype.map() method to
        // create an array of markers based on a given "locations" array.
        // The map() method here has nothing to do with the Google Maps API.

	    var locations;
		function update() {
          $.getJSON('/json')
            .done(function(data) {
              locations = data;
              var markers = locations.map(function(location, i) {
                return new google.maps.Marker({
                  position: location,
                  label: labels[i % labels.length]
                });
              });
              var options = {
                imagePath: 'https://raw.githubusercontent.com/googlemaps/js-marker-clusterer/gh-pages/images/m'
              };
              // Add a marker clusterer to manage the markers.
              var markerCluster = new MarkerClusterer(map, markers, options);
			  setTimeout(update, 10000);
            });
		}
		update();
      }
    </script>
    <script src="https://developers.google.com/maps/documentation/javascript/examples/markerclusterer/markerclusterer.js">
    </script>
    <script async defer
    src="https://maps.googleapis.com/maps/api/js?key={{.Key}}&callback=initMap">
    </script>
  </body>
</html>
`

var googleApiKey *string = flag.String("google-map-api-key", "AIzaSyAc0R5epXJUgKcLAIlie8GTOt7lZwjiqas", "googlemap api key")
var googleMapInterface *string = flag.String("google-map-listen", ":7676", "google map listen interface")

func toJson(i interface{}) string {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	return string(b)
}

func join(a []string, sep string) template.JS {
	return template.JS(strings.Join(a, sep))
}

func runHttp() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.New("index")
		data := &googleData{
			Key: *googleApiKey,
		}
		t.Parse(googlemap)
		t.Execute(w, data)
	})
	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		var list []string
		mapData.RLock()
		for _, v := range mapData.m {
			line := fmt.Sprintf(`{"lat": %v, "lng": %v}`, v.Latitude, v.Longitude)
			list = append(list, line)
		}
		mapData.RUnlock()
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "[%s]", strings.Join(list, ","))
	})
	log.Fatal(http.ListenAndServe(*googleMapInterface, gziphandler.GzipHandler(mux)))
}
