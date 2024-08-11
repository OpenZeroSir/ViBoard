package Lib

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// the complete array of all items
type KML struct {
	XMLName xml.Name `xml:"kml"`
	//XMLNS    string   `xml:"xmlns,attr" json:"xmlns,attr"`
	Document Document `xml:"Document"`
}
type KML2 struct {
	XMLName xml.Name `xml:"kml"`
	//XMLNS    string   `xml:"xmlns,attr" json:"xmlns,attr"`
	Document Document2 `xml:"Document"`
}

// the Document struct, containing Schema
// and Folders
type Document struct {
	//XMLName xml.Name `xml:"Document" json:"Document"`
	Id     string      `xml:"id,attr"`
	Schema Schema      `xml:"Schema"`
	Folder Folder      `xml:"Folder"`
	Styles []MainStyle `xml:"Style"`
}
type Document2 struct {
	//XMLName xml.Name `xml:"Document" json:"Document"`
	Id         string      `xml:"id,attr"`
	Schema     Schema      `xml:"Schema"`
	Placemarks []Placemark `xml:"Placemark"`
	Styles     []MainStyle `xml:"Style"`
}

type MainStyle struct {
	Id         string     `xml:"id,attr,omitempty"`
	LineStyle  LineStyle  `xml:"LineStyle,omitempty"`
	PolyStyle  PolyStyle  `xml:"PolyStyle,omitempty"`
	PointStyle PointStyle `xml:"PointStyle,omitempty"`
}

// Schema describes the types of the extended attributes
type Schema struct {
	Name        string        `xml:"name,attr"`
	ID          string        `xml:"id,attr"`
	SimpleField []SimpleField `xml:"SimpleField"`
}

type SimpleField struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// the Folder struct holds all the Placemarks (drawing items)
type Folder struct {
	Name       string      `xml:"name"`
	Placemarks []Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name          string        `xml:"name"`
	Style         Style         `xml:"Style"`
	MultiGeometry MultiGeometry `xml:"MultiGeometry"`
	Point         Point         `xml:"Point"`
	// ExtendedData  ExtendedData  `xml:"ExtendedData"`
	Polygon    Polygon    `xml:"Polygon"`
	StyleUrl   string     `xml:"styleUrl"`
	LineString LineString `xml:"LineString"`
}

// Styles are the colors pre assigned
// The following are all STYLES
type Style struct {
	LineStyle  LineStyle  `xml:"LineStyle"`
	PolyStyle  PolyStyle  `xml:"PolyStyle"`
	PointStyle PointStyle `xml:"PointStyle"`
}

type PointStyle struct {
	Color string `xml:"color"`
	Fill  string `xml:"fill"`
}

type LineStyle struct {
	Color string `xml:"color"`
	Fill  string `xml:"fill"`
}

type PolyStyle struct {
	Color string `xml:"color"`
	Fill  string `xml:"fill"`
}

// GEOMS
// MultiGeometry is the container
type MultiGeometry struct {
	LineString LineString `xml:"LineString"`
	Polygon    Polygon    `xml:"Polygon"`
	LinearRing LinearRing `xml:"LinearRing"`
}

type Point struct {
	StringCoords string `xml:"coordinates"`
}

type LineString struct {
	StringCoords string `xml:"coordinates"`
}

type Polygon struct {
	OuterBoundaryIs OuterBoundary `xml:"outerBoundaryIs"`
}

type OuterBoundary struct {
	LinearRing LinearRing `xml:"LinearRing"`
}

type LinearRing struct {
	StringCoords string `xml:"coordinates"`
}

// EXTENDED ATTRIBUTES
// added to kmls that weren't in original schema
// type ExtendedData struct {
// 	SchemaData SchemaData `xml:"SchemaData"`
// }

// type SchemaData struct {
// 	Schema     string       `xml:"schemaUrl,attr"`
// 	SimpleData []SimpleData `xml:"SimpleData"`
// }

// type SimpleData struct {
// 	Key   string `xml:"name,attr"`
// 	Value string `xml:",chardata"`
// }

// unravel the xml into a single KML struct
// all the code basically sets up the structs
// the rest is just using xml to parse
// func KMLDecode(f *bytes.Buffer, kml *KML) {

// 	// xml.Unmarshal(f, kml)
// 	d := xml.NewDecoder(f)

// 	// place all the xml where it belongs
// 	d.Decode(kml)

// 	// change coords from string >> [][]float64
// 	// for i, geom := range kml.Document.Folder.Placemarks {
// 	// 	if i > 5000 {
// 	// 		fmt.Print(i, geom)
// 	// 	}
// 	// 	// if point
// 	// 	if len(geom.Point.StringCoords) > 0 {
// 	// 		str := geom.Point.StringCoords

// 	// 		payload := coordStringDecode(str)
// 	// 	}

// 	// 	// if linestring
// 	// 	if len(geom.MultiGeometry.LineString.StringCoords) > 0 {
// 	// 		str := geom.MultiGeometry.LineString.StringCoords

// 	// 		payload := coordStringDecode(str)

// 	// 		kml.Document.Folder.Placemarks[i].MultiGeometry.LineString.Coordinates = payload
// 	// 	}

// 	// 	// if polygon
// 	// 	if len(geom.MultiGeometry.Polygon.OuterBoundaryIs.LinearRing.StringCoords) > 0 {
// 	// 		str := geom.MultiGeometry.Polygon.OuterBoundaryIs.LinearRing.StringCoords

// 	// 		payload := coordStringDecode(str)

// 	// 		kml.Document.Folder.Placemarks[i].MultiGeometry.Polygon.OuterBoundaryIs.LinearRing.Coordinates = payload
// 	// 	}
// 	// }
// }

// coordStringDecode performs the parsing and number conversion
// func coordStringDecode(str string) [][]float64 {
// 	temp := strings.Split(str, " ")

// 	var payload [][]float64

// 	for _, coord := range temp {

// 		xyz := strings.Split(coord, ",")

// 		x, _ := strconv.ParseFloat(string(xyz[0]), 64)
// 		y, _ := strconv.ParseFloat(string(xyz[1]), 64)

// 		floatcoord := []float64{x, y}

// 		// test if a third coord (elevation) is present!!
// 		if len(xyz) > 2 {
// 			z, _ := strconv.ParseFloat(string(xyz[2]), 64)
// 			floatcoord = append(floatcoord, z)
// 		}

// 		payload = append(payload, floatcoord)
// 	}

// 	return payload
// }

// kml文件中可能有内容需要替换
func KMLReplaceString(txt string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(txt, "\t", ""), "\r", ""), "\n", " "))
}

func KMLCoordStringToFloat(txt string) [][]float64 {
	str := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(txt, "\t", ""), "\r", ""), "\n", " ")
	arr := strings.Split(str, " ")
	var points [][]float64
	for _, v := range arr {
		if len(v) > 3 && strings.Contains(v, ",") {
			tmp := strings.Split(v, ",")
			if len(tmp) >= 2 {
				lon, err := strconv.ParseFloat(tmp[0], 64)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				lat, err := strconv.ParseFloat(tmp[1], 64)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				var point []float64
				point = append(point, lon)
				point = append(point, lat)
				points = append(points, point)
			}
		}
	}
	return points
}
