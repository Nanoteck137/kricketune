package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nanoteck137/kricketune/apis"
	"github.com/nanoteck137/pyrin/spec"
	"github.com/nanoteck137/pyrin/tools/gen"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "internal",
}

var genCmd = &cobra.Command{
	Use: "gen",
	Run: func(cmd *cobra.Command, args []string) {
		router := spec.Router{}

		apis.RegisterHandlers(nil, &router)

		s, err := spec.GenerateSpec(router.Routes)
		if err != nil {
			log.Fatal("Failed to generate spec", err)
		}

		d, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			log.Fatal("Failed to marshal server", err)
		}

		err = os.WriteFile("misc/pyrin.json", d, 0644)
		if err != nil {
			log.Fatal("Failed to write pyrin.json", err)
		}

		fmt.Println("Wrote 'misc/pyrin.json'")

		err = gen.GenerateTypescript(s, "web/src/lib/api")
		if err != nil {
			log.Fatal("Failed to generate golang code", err)
		}
	},
}

var genClientCmd = &cobra.Command{
	Use: "gen-client",
	Run: func(cmd *cobra.Command, args []string) {
		dwebbleVersion := "main"

		url := ""
		if dwebbleVersion == "main" {
			url = "https://raw.githubusercontent.com/Nanoteck137/dwebble/refs/heads/main/misc/pyrin.json"
		} else {
			url = fmt.Sprintf("https://raw.githubusercontent.com/Nanoteck137/dwebble/refs/tags/v%s/misc/pyrin.json", dwebbleVersion)
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		s := spec.Server{}

		d := json.NewDecoder(resp.Body)
		err = d.Decode(&s)
		if err != nil {
			log.Fatal(err)
		}

		// TODO(patrik): Save the spec to disk
		err = gen.GenerateGolang(&s, "client/api")
		if err != nil {
			log.Fatal("Failed to generate golang code", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(genClientCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
