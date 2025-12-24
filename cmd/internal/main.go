package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/nanoteck137/kricketune/apis"
	"github.com/nanoteck137/pyrin/spark"
	"github.com/nanoteck137/pyrin/spark/typescript"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "internal",
}

var genCmd = &cobra.Command{
	Use: "gen",
	Run: func(cmd *cobra.Command, args []string) {
		router := spark.Router{}
		apis.RegisterHandlers(nil, &router)

		nameFilter := spark.NameFilter{}

		serverDef, err := spark.CreateServerDef(&router, nameFilter)
		if err != nil {
			slog.Error("failed to create server def", "err", err)
			os.Exit(-1)
		}

		err = serverDef.SaveToFile("misc/pyrin.json")
		if err != nil {
			slog.Error("failed save server def", "err", err)
			os.Exit(-1)
		}

		slog.Info("Wrote 'misc/pyrin.json'")

		resolver, err := spark.CreateResolverFromServerDef(&serverDef)
		if err != nil {
			slog.Error("failed to create resolver", "err", err)
			os.Exit(-1)
		}

		{
			gen := typescript.TypescriptGenerator{}

			err = gen.Generate(&serverDef, resolver, "web/src/lib/api")
			if err != nil {
				slog.Error("failed to generate typescript client", "err", err)
				os.Exit(-1)
			}
		}
	},
}

// var genClientCmd = &cobra.Command{
// 	Use: "gen-client",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		dwebbleVersion := "main"
//
// 		url := ""
// 		if dwebbleVersion == "main" {
// 			url = "https://raw.githubusercontent.com/Nanoteck137/dwebble/refs/heads/main/misc/pyrin.json"
// 		} else {
// 			url = fmt.Sprintf("https://raw.githubusercontent.com/Nanoteck137/dwebble/refs/tags/v%s/misc/pyrin.json", dwebbleVersion)
// 		}
//
// 		resp, err := http.Get(url)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer resp.Body.Close()
//
// 		s := spec.Server{}
//
// 		d := json.NewDecoder(resp.Body)
// 		err = d.Decode(&s)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		// TODO(patrik): Save the spec to disk
// 		err = gen.GenerateGolang(&s, "client/api")
// 		if err != nil {
// 			log.Fatal("Failed to generate golang code", err)
// 		}
// 	},
// }

func init() {
	rootCmd.AddCommand(genCmd)
	// rootCmd.AddCommand(genClientCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
