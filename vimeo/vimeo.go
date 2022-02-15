package vimeo

import (
   "encoding/json"
   "github.com/89z/format"
   "net/http"
   "strconv"
   "strings"
   "time"
)

var LogLevel format.LogLevel

type Clip struct {
   ID, UnlistedHash int64
}

func NewClip(address string) (*Clip, error) {
   var (
      clipPage Clip
      err error
   )
   fields := strings.FieldsFunc(address, func(r rune) bool {
      return r < '0' || r > '9'
   })
   for key, val := range fields {
      switch key {
      case 0:
         clipPage.ID, err = strconv.ParseInt(val, 10, 64)
      case 1:
         clipPage.UnlistedHash, err = strconv.ParseInt(val, 10, 64)
      }
      if err != nil {
         return nil, err
      }
   }
   return &clipPage, nil
}

type JsonWeb struct {
   Token string
}

func NewJsonWeb() (*JsonWeb, error) {
   req, err := http.NewRequest("GET", "https://vimeo.com/_rv/jwt", nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("X-Requested-With", "XMLHttpRequest")
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   web := new(JsonWeb)
   if err := json.NewDecoder(res.Body).Decode(web); err != nil {
      return nil, err
   }
   return web, nil
}

func (w JsonWeb) Video(clip *Clip) (*Video, error) {
   buf := []byte("https://api.vimeo.com/videos/")
   buf = strconv.AppendInt(buf, clip.ID, 10)
   if clip.UnlistedHash >= 1 {
      buf = append(buf, ':')
      buf = strconv.AppendInt(buf, clip.UnlistedHash, 10)
   }
   req, err := http.NewRequest("GET", string(buf), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("Authorization", "JWT " + w.Token)
   req.URL.RawQuery = "fields=duration,download"
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errorString(res.Status)
   }
   vid := new(Video)
   if err := json.NewDecoder(res.Body).Decode(vid); err != nil {
      return nil, err
   }
   return vid, nil
}

type Download struct {
   Video_File_ID int64
   Quality string
   Width int64
   Height int64
   Size_Short string
   Link string
}

type Video struct {
   Download []Download
   Duration int64
   Release_Time string
}

func (d Download) Format(link bool) string {
   buf := []byte("ID:")
   buf = strconv.AppendInt(buf, d.Video_File_ID, 10)
   buf = append(buf, " Quality:"...)
   buf = append(buf, d.Quality...)
   buf = append(buf, " Width:"...)
   buf = strconv.AppendInt(buf, d.Width, 10)
   buf = append(buf, " Height:"...)
   buf = strconv.AppendInt(buf, d.Height, 10)
   buf = append(buf, " Size:"...)
   buf = append(buf, d.Size_Short...)
   if link {
      buf = append(buf, " Link:"...)
      buf = append(buf, d.Link...)
   }
   return string(buf)
}

func (v Video) Format(link bool) string {
   buf := []byte("Duration: ")
   buf = append(buf, v.Time().String()...)
   buf = append(buf, "\nDownloads:"...)
   for _, down := range v.Download {
      buf = append(buf, '\n')
      buf = append(buf, down.Format(link)...)
   }
   return string(buf)
}

func (v Video) Time() time.Duration {
   return time.Duration(v.Duration) * time.Second
}

type errorString string

func (e errorString) Error() string {
   return string(e)
}
