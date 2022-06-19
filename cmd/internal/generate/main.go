package generate

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jakseer/any2struct"
	"github.com/spf13/cobra"
)

var (
	// inputType input data type
	inputType string

	// outputTags output data tags
	outputTags string

	// inputPath input file path. read data from stdin if inputPath is omitted
	inputPath string
)

var CmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "generate go struct",
	Long:  "generate go struct with various tags",
	Run:   run,
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.any2struct.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	CmdGenerate.Flags().StringVarP(&inputType, "data-type", "d", "", "input data type, support: yaml,json,sql")
	CmdGenerate.Flags().StringVarP(&outputTags, "tags", "t", "", "generated struct tags, support: json,gorm,yaml")
	CmdGenerate.Flags().StringVarP(&inputPath, "input", "i", "", "input file path")

	_ = CmdGenerate.MarkFlagRequired("data-type")
}

func run(cmd *cobra.Command, args []string) {
	// check decode type
	if !strInArray(inputType, []string{
		any2struct.DecodeTypeSQL,
		any2struct.DecodeTypeJSON,
		any2struct.DecodeTypeYaml,
	}) {
		log.Fatal(errors.New("invalid param data-type"))
	}

	var inputData string
	if inputPath != "" {
		fp, err := os.Open(inputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()

		fileData, err := io.ReadAll(fp)
		if err != nil {
			log.Fatal(err)
		}

		inputData = string(fileData)
	} else {
		// get from stdin
		if len(args) == 0 {
			log.Fatal(errors.New("no input data"))
		}
		inputData = args[0]
	}

	// check and parse encode tags
	var tags []string
	for _, tag := range strings.Split(outputTags, ",") {
		if tag == "" {
			continue
		}
		if !strInArray(tag, []string{
			any2struct.EncodeTypeJSON,
			any2struct.EncodeTypeGorm,
			any2struct.EncodeTypeYaml,
		}) {
			log.Fatal(errors.New("invalid param tags"))
		}
		tags = append(tags, tag)
	}

	out, err := any2struct.NewConvertor().Convert(inputData, inputType, tags)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}

func strInArray(str string, arr []string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}
