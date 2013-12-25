package sb

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "html/template"
    "io/ioutil"
    "net"
    "net/http"
    "path/filepath"
    "reflect"
    "strings"
)

type JSONMessage interface {
    /* All messages must have these fields:
       Success bool
    */
}

//json response helper
type JSONResponse struct {
    Success bool
    Message string
}

func (r JSONResponse) String() string {
    b, err := json.Marshal(r)
    if err != nil {
        return ""
    }
    return string(b)
}

func FailureResponse(message string) string {
    return JSONResponse{Success: false, Message: message}.String()
}

func SuccessResponse(message string) string {
    return JSONResponse{Success: true, Message: message}.String()
}

/* HTTP request utilities */

func BuildURL(protocol string, address string, route string,
    params map[string]string) string {
    query := ""
    if params != nil {
        var _params []string
        for key, value := range params {
            _params = append(_params, fmt.Sprintf("%s=%s", key, value))
        }
        query = "?" + strings.Join(_params, "&")
    }

    if protocol != "" && !strings.HasSuffix(protocol, "://") {
        protocol += "://"
    }
    return fmt.Sprintf("%s%s%s%s", protocol, address, route, query)
}

/* JSON Helper */

func unpackSuccessField(message JSONMessage) (success reflect.Value, err error) {
    /* Uses reflection to retrieve the "Success bool" field that all
       messages must have */
    rm := reflect.ValueOf(message)
    if rm.Kind() == reflect.Ptr {
        rm = reflect.Indirect(rm)
    }
    if !rm.IsValid() {
        err = errors.New("Invalid message")
    } else {
        success = rm.FieldByName("Success")
        if !success.IsValid() {
            err = errors.New("Message must contain a \"Success bool\" field")
        } else if !success.CanSet() {
            err = errors.New("Can't set Success")
        } else if success.Kind() != reflect.Bool {
            err = errors.New("Message's Success field must be a bool")
        }
    }
    return
}

func setSuccess(message JSONMessage, val bool) (err error) {
    success, err := unpackSuccessField(message)
    if err == nil {
        success.SetBool(val)
    }
    return
}

func getSuccess(message JSONMessage) (val bool, err error) {
    val = false
    success, err := unpackSuccessField(message)
    if err == nil {
        val = success.Bool()
    }
    return
}

func SendJSON(w http.ResponseWriter, message JSONMessage) error {
    setSuccess(message, true)
    out, err := json.Marshal(message)
    if err == nil {
        _, err := w.Write(out)
        if err != nil {
            return err
        }
    }
    return err
}

func GetJSON(address string, route string, message JSONMessage) error {
    return GetJSONParams(address, route, message, nil)
}

func GetJSONParams(address string, route string, message JSONMessage,
    params map[string]string) error {
    setSuccess(message, false)

    url := _make_url(address, route, params)
    fmt.Printf("Making request to %v\n", url)

    // request peers
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // read response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    _print_response(body)

    err = json.Unmarshal(body, message)
    if err != nil {
        return err
    }

    success, err := getSuccess(message)
    if err != nil {
        return err
    }
    if !success {
        return fmt.Errorf("Failed to parse json response: %s\n", body)
    }

    return nil
}

func PostJSON(address string, route string, message JSONMessage,
    outgoing interface{}) error {
    return PostJSONParams(address, route, message, outgoing, nil)
}

func PostJSONParams(address string, route string, message JSONMessage,
    outgoing interface{}, params map[string]string) error {
    setSuccess(message, false)

    // marshal our data
    data, err := json.Marshal(outgoing)
    if err != nil {
        return err
    }
    reader := bytes.NewReader(data)

    // form the complete url
    url := _make_url(address, route, params)
    fmt.Printf("Making request to %v\n", url)

    // make request
    resp, err := http.Post(url, "application/json", reader)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // read response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    _print_response(body)

    // unmarshal response
    err = json.Unmarshal(body, message)
    if err != nil {
        return err
    }

    success, err := getSuccess(message)
    if err != nil {
        return err
    }
    if !success {
        return fmt.Errorf("Failed to parse json response: %s\n", body)
    }

    return nil
}

func _print_response(body []byte) {
    var PRINT_MAX int = 1024
    if len(body) < PRINT_MAX {
        fmt.Printf("JSONResponse: %s\n", body)
    } else {
        fmt.Printf("JSONResponse: %s (...truncated)\n", body[:PRINT_MAX])
    }
}

func _make_url(address string, route string, params map[string]string) string {
    protocol := "http"
    if strings.HasPrefix(address, "http://") ||
        strings.HasPrefix(address, "https://") {
        protocol = ""
    }
    return BuildURL(protocol, address, route, params)
}

/* HTTP Server Utilities */

func ListenAndServeBackground(address string, mux *http.ServeMux,
    errchan chan<- error) {
    err := http.ListenAndServe(address, mux)
    errchan <- err
}

func HttpError(w http.ResponseWriter, status int, default_message string, messages []string) {
    message := default_message
    if len(messages) != 0 {
        message = strings.Join(messages, "<br>")
    }
    http.Error(w, message, status)
}

func Error400(w http.ResponseWriter, messages ...string) {
    HttpError(w, http.StatusBadRequest, "Bad request", messages)
}

func Error404(w http.ResponseWriter, messages ...string) {
    HttpError(w, http.StatusNotFound, "Not found", messages)
}

func Error405(w http.ResponseWriter, messages ...string) {
    HttpError(w, http.StatusMethodNotAllowed, "Method not allowed", messages)
}

func Error501(w http.ResponseWriter, messages ...string) {
    HttpError(w, http.StatusNotImplemented, "Not implemented", messages)
}

func Error500(w http.ResponseWriter, messages ...string) {
    HttpError(w, http.StatusInternalServerError, "Internal server error", messages)
}

/* Template helpers */

func LoadTemplate(html_file string) (*template.Template, error) {
    const template_dir = "./static"
    if !strings.HasPrefix(html_file, template_dir) {
        html_file = filepath.Join(template_dir, html_file)
    }

    t, err := template.ParseFiles(html_file, "./static/common.html")
    return t, err
}

func ShowTemplate(w http.ResponseWriter, html_file string, p interface{}) {
    t, err := LoadTemplate(html_file)
    if err != nil {
        Error500(w, err.Error())
    }
    err = t.Execute(w, p)
    if err != nil {
        Error500(w, err.Error())
    }
}

/* Networking */

func LocalIP() (net.IP, error) {
    tt, err := net.Interfaces()
    if err != nil {
        return nil, err
    }
    for _, t := range tt {
        aa, err := t.Addrs()
        if err != nil {
            return nil, err
        }
        for _, a := range aa {
            ipnet, ok := a.(*net.IPNet)
            if !ok {
                continue
            }
            v4 := ipnet.IP.To4()
            if v4 == nil || v4[0] == 127 { // loopback address
                continue
            }
            return v4, nil
        }
    }
    return nil, errors.New("cannot find local IP address")
}

func LocalIPString() (string, error) {
    _ip, err := LocalIP()
    var ip string = ""
    if err == nil {
        ip = _ip.String()
    }
    return ip, err
}
