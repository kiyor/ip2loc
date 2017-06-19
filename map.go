/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : map.go

* Purpose :

* Creation Date : 04-12-2017

* Last Modified : Fri 14 Apr 2017 06:20:15 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"flag"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type googleData struct {
	Key      string
	Hostname string
}

const googlemap = `<!DOCTYPE html>
<html>
  <head>
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no">
    <meta charset="utf-8">
    <title>{{.Hostname}} | ip2loc map</title>
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

        var haightAshbury = {lat: 37.769, lng: -122.446};

        var map = new google.maps.Map(document.getElementById('map'), {
          zoom: 2,
          center: haightAshbury,
          mapTypeId: 'terrain'
        });

        // Create an array of alphabetical characters used to label the markers.
        var labels = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

        // Add some markers to the map.
        // Note: The code uses the JavaScript Array.prototype.map() method to
        // create an array of markers based on a given "locations" array.
        // The map() method here has nothing to do with the Google Maps API.

	    var locations;
		var markers = [];
		var markerCluster;
        var options = {
          imagePath: 'https://raw.githubusercontent.com/googlemaps/js-marker-clusterer/gh-pages/images/m'
        };

        // Adds a marker to the map and push to the array.
        function addMarker(location) {
          var marker = new google.maps.Marker({
            position: location,
            map: map
          });
          markers.push(marker);
        }

        // Sets the map on all markers in the array.
        function setMapOnAll(map) {
          for (var i = 0; i < markers.length; i++) {
            markers[i].setMap(map);
          }
        }
  
        // Removes the markers from the map, but keeps them in the array.
        function clearMarkers() {
          setMapOnAll(null);
        }
  
        // Shows any markers currently in the array.
        function showMarkers() {
          setMapOnAll(map);
        }
  
        // Deletes all markers in the array by removing references to them.
        function deleteMarkers() {
          clearMarkers();
          markers = [];
        }

		function update() {
          $.getJSON('/json')
            .done(function(data) {
              locations = data;
              deleteMarkers();
              $.each( locations, function( key, value ) {
                addMarker(value);
              });
              // Add a marker clusterer to manage the markers.
              if ( markers.length > 0 ) {
                markerCluster = new MarkerClusterer(map, markers, options);
              } else {
                markerCluster = new MarkerClusterer(map, [], options);
              }
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

var (
	mapData            = newMapData()
	googleApiKey       = flag.String("google-map-api-key", "AIzaSyAc0R5epXJUgKcLAIlie8GTOt7lZwjiqas", "googlemap api key")
	googleMapInterface = flag.String("google-map-listen", ":7676", "google map listen interface")
	flagExpire         = flag.Duration("expire", 5*time.Minute, "expire time")
	hostname, _        = os.Hostname()
)

type MapData struct {
	m map[mapKey]*ipLoc
	*sync.RWMutex
}

type mapKey struct {
	t time.Time
	r string
}

func newMapData() MapData {
	return MapData{
		m:       make(map[mapKey]*ipLoc),
		RWMutex: new(sync.RWMutex),
	}
}

func cronHttp() {
	ticker := time.NewTicker(*flagExpire / 10)
	go func() {
		for range ticker.C {
			cleanMap()
		}
	}()
}

func cleanMap() {
	mapData.Lock()
	defer mapData.Unlock()
	s1 := len(mapData.m)
	t1 := time.Now()
	for k := range mapData.m {
		if time.Now().Sub(k.t) > *flagExpire {
			delete(mapData.m, k)
		}
	}
	s2 := len(mapData.m)
	log.Printf("map cleaned %d -> %d (%d in %v)\n", s1, s2, s2-s1, time.Since(t1))
}

func runHttp() {
	go cronHttp()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.New("index")
		data := &googleData{
			Key:      *googleApiKey,
			Hostname: hostname,
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
