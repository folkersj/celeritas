package celeritas

import (
    "encoding/json"
    "encoding/xml"
    "errors"
    "fmt"
    "io"
    "net/http"
    "path"
    "path/filepath"
)

func (c *Celeritas) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
    maxBytes := 1024 * 1024 // one MiB
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

    dec := json.NewDecoder(r.Body)
    err := dec.Decode(data)
    if err != nil {
        return err
    }

    err = dec.Decode(&struct{}{})
    if err != io.EOF {
        return errors.New("body must only have a single json value")
    }

    return nil
}

func (c *Celeritas) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
    out, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return err
    }

    if len(headers) > 0 {
        for key, value := range headers[0] {
            w.Header()[key] = value
        }
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    _, err = w.Write(out)
    return err
}

func (c *Celeritas) WriteXNL(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
    out, err := xml.MarshalIndent(data, "", "\t")
    if err != nil {
        return err
    }

    if len(headers) > 0 {
        for key, value := range headers[0] {
            w.Header()[key] = value
        }
    }

    w.Header().Set("Content-Type", "application/xml")
    w.WriteHeader(status)

    _, err = w.Write(out)
    return err
}

func (c *Celeritas) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
    fp := path.Join(pathToFile, fileName)

    fileToServe := filepath.Clean(fp)
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; file=\"%s\"", fileName))
    http.ServeFile(w, r, fileToServe)
    return nil
}

func (c *Celeritas) Error404(w http.ResponseWriter, r *http.Request) {
    c.ErrorStatus(w, http.StatusNotFound)
}

func (c *Celeritas) Error500(w http.ResponseWriter, r *http.Request) {
    c.ErrorStatus(w, http.StatusInternalServerError)
}

func (c *Celeritas) ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
    c.ErrorStatus(w, http.StatusUnauthorized)
}

func (c *Celeritas) ErrorForbidden(w http.ResponseWriter, r *http.Request) {
    c.ErrorStatus(w, http.StatusForbidden)
}

func (c *Celeritas) ErrorStatus(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}
