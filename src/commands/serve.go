/*
    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package commands

import (
    "strings"
    "io/ioutil"
    "net/http"
    "net/url"

    "github.com/spf13/cobra"
    "github.com/gorilla/mux"

    clamor "github.com/clamor-vms/clamor-go-core"

    "clamor/controllers"
    "clamor/core"
)

func writeProxyResponse(w http.ResponseWriter, resp *http.Response, err error) {
    if err != nil {
        // err
        http.Error(w, err.Error(), 500)
        return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        // err
        http.Error(w, err.Error(), 500)
        return
    }

    for k, _ := range resp.Header {
        w.Header().Set(k, resp.Header.Get(k))
    }

    w.WriteHeader(resp.StatusCode)
    w.Write([]byte(string(body)))
}

func proxyRequest(serviceDomain string, serviceURLPrefix string, w http.ResponseWriter, r *http.Request) {
    newUrl, err := url.Parse(r.URL.RequestURI())
    newUrl.Scheme = "http"
    newUrl.Host = serviceDomain

    trimmedPath := strings.TrimPrefix(r.URL.Path, "/") // apparently starting with / is optional. So lets just trim left that first
    trimmedPath = strings.TrimPrefix(trimmedPath, serviceURLPrefix) // so if it exists we'll be safe.
    newUrl.Path = trimmedPath

    req, err := http.NewRequest(r.Method, newUrl.String(), nil)
    req.Header = r.Header
    req.Body = r.Body

    resp, err := http.DefaultClient.Do(req)

    writeProxyResponse(w, resp, err)
}

func redirectAuth(w http.ResponseWriter, r *http.Request) {
    proxyRequest("auth-service", "auth", w, r)
}

func redirectCampaign(w http.ResponseWriter, r *http.Request) {
    proxyRequest("campaign-service", "campaign", w, r)
}

func redirectTask(w http.ResponseWriter, r *http.Request) {
    proxyRequest("task-service", "task", w, r)
}

func redirectUser(w http.ResponseWriter, r *http.Request) {
    proxyRequest("user-service", "user", w, r)
}

//TODO: For testing only while in dev. I don't think this will be part of the accessable api.
func redirectVoter(w http.ResponseWriter, r *http.Request) {
    proxyRequest("voter-service", "voter", w, r)
}

//TODO: I think I need to think this through a bit more? should this even be part of this?
func serveStaticWebApp(w http.ResponseWriter, r *http.Request) {
    proxyRequest("web-client", "", w, r)
}

func LowerCaseURI(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        r.URL.Path = strings.ToLower(r.URL.Path)
        h.ServeHTTP(w, r)
    }

    return http.HandlerFunc(fn)
}

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "runs the rest api",
    Long:  `runs the rest api`,
    Run: func(cmd *cobra.Command, args []string) {
        r := mux.NewRouter()

        //setup controllers
        aboutController := clamor.NewControllerProcessor(controllers.NewAboutController())

        //setup all proxies to other services
        r.HandleFunc("/auth/{rest:.*}", redirectAuth)
        r.HandleFunc("/campaign/{rest:.*}", redirectCampaign)
        r.HandleFunc("/task/{rest:.*}", redirectTask)
        r.HandleFunc("/user/{rest:.*}", redirectUser)
        
        //TODO: For testing only while in dev. I don't think this will be part of the accessable api.
        r.HandleFunc("/voter/{rest:.*}", redirectVoter)
        

        //TODO: This should return the static web client?
        r.HandleFunc("/about", aboutController.Logic)
        
        //server up app
        if err := http.ListenAndServe(":" + core.PORT_NUMBER, clamor.PanicHandler(LowerCaseURI(r))); err != nil {
            panic(err)
        }
    },
}

//Entry
func init() {
    RootCmd.AddCommand(serveCmd)
}
