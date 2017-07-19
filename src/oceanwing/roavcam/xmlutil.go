package roavcam

import (
	"encoding/xml"
	"log"
)

// just for test.
var myXMLData = `
	<LIST>
		<ALLFile>
			<File>
				<NAME>2017-07-17-01.mp4</NAME>
				<FPATH>01PATH</FPATH>
				<SIZE>1000</SIZE>
				<TIMECODE>12345</TIMECODE>
				<TIME>18:00:03</TIME>
				<ATTR>32</ATTR>
			</File>
		</ALLFile>
		<ALLFile>
			<File>
				<NAME>2018-07-17-01.mp4</NAME>
				<FPATH>02PATH</FPATH>
				<SIZE>2000</SIZE>
				<TIMECODE>22345</TIMECODE>
				<TIME>19:00:03</TIME>
				<ATTR>32</ATTR>
			</File>
		</ALLFile>
	</LIST>
`

// RoavXML hh.
type RoavXML struct {
	XMLName xml.Name `xml:"LIST"`
	AllFile []Allf   `xml:"ALLFile"`
}

// Allf hh.
type Allf struct {
	Name     string `xml:"File>NAME"`
	Fpath    string `xml:"File>FPATH"`
	Size     string `xml:"File>SIZE"`
	TimeCode string `xml:"File>TIMECODE"`
	Time     string `xml:"File>TIME"`
	Attr     string `xml:"File>ATTR"`
}

// myXml22 just for test.
func myXml22() {
	v := &RoavXML{}
	err := xml.Unmarshal([]byte(myXMLData), v)
	if err != nil {
		log.Printf("unmarsal failt: %s\n", err)
		return
	}
	log.Printf("allfile count: %d\n", len(v.AllFile))
	for i, x := range v.AllFile {
		log.Printf("Num.%d file info\n", i)
		log.Printf("Name: %s\n", x.Name)
		log.Printf("Fpath: %s\n", x.Fpath)
		log.Printf("Size: %s\n", x.Size)
		log.Printf("Timecode: %s\n", x.TimeCode)
		log.Printf("Time: %s\n", x.Time)
		log.Printf("Attr: %s\n", x.Attr)
	}
}
