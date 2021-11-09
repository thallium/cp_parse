package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/thallium/cp_parse/util"
)

type data struct {
	Body string
	Url  string
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}

	var d data

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	parsedURL, err := url.Parse(d.Url)
	if err != nil {
		fmt.Errorf("URL cannot be parsed.")
		return
	}

	website := util.GetWebsiteName(parsedURL.Hostname())
	util.BodyByExtension = d.Body
	cmd, _, _ := rootCmd.Find([]string{website})
	if cmd == nil {
		fmt.Errorf("This website isn't supported!")
		return
	}
	// cmd.SetArgs([]string{d.Url})
	cmd.Run(nil, []string{d.Url})
	// cmd.Execute()
	util.BodyByExtension = ""
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start a server for parsing data from the broswer extension",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", handler)
		fmt.Println("Server started!")
		http.ListenAndServe(":8090", nil)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cfCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cfCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
