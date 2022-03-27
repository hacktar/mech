package paramount

import (
   "encoding/xml"
   "github.com/89z/format"
   "net/http"
   "net/url"
   "github.com/89z/mech"
   "strconv"
   "strings"
)

const (
   aid = 2198311517
   sid = "dJ5BDC"
)

var LogLevel format.LogLevel

type Media struct {
   Body struct {
      Seq  struct {
         Video []Video `xml:"video"`
      } `xml:"seq"`
   } `xml:"body"`
}

func NewMedia(guid string) (*Media, error) {
   buf := []byte("https://link.theplatform.com/s/")
   buf = append(buf, sid...)
   buf = append(buf, "/media/guid/"...)
   buf = strconv.AppendInt(buf, aid, 10)
   buf = append(buf, '/')
   buf = append(buf, guid...)
   req, err := http.NewRequest("GET", string(buf), nil)
   if err != nil {
      return nil, err
   }
   // We need "MPEG4", otherwise you get a "EXT-X-KEY" with "skd" scheme:
   req.URL.RawQuery = url.Values{
      "format": {"SMIL"},
      "formats": {"MPEG4,M3U"},
   }.Encode()
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   med := new(Media)
   if err := xml.NewDecoder(res.Body).Decode(med); err != nil {
      return nil, err
   }
   return med, nil
}

func (m Media) Video() *Video {
   for _, vid := range m.Body.Seq.Video {
      return &vid
   }
   return nil
}

type Video struct {
   Title string `xml:"title,attr"`
   Src string `xml:"src,attr"`
   Param []struct {
      Name string `xml:"name,attr"`
      Value string `xml:"value,attr"`
   } `xml:"param"`
}

func (v Video) Base() string {
   var buf strings.Builder
   buf.WriteString(v.Title)
   buf.WriteByte('-')
   buf.WriteString(v.SeasonNumber())
   buf.WriteByte('-')
   buf.WriteString(v.EpisodeNumber())
   return mech.Clean(buf.String())
}

func (v Video) EpisodeNumber() string {
   for _, par := range v.Param {
      if par.Name == "EpisodeNumber" {
         return par.Value
      }
   }
   return ""
}

func (v Video) SeasonNumber() string {
   for _, par := range v.Param {
      if par.Name == "SeasonNumber" {
         return par.Value
      }
   }
   return ""
}
