package ix

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xkortex/vprint"
	"log"
	"os"
)

var (
	Version = "unset"
)

func parseSliceArgs(args []string) (settings *MultiSlice, err error) {
	if len(args) == 0 {
		return
	}
	return
	//lineSlicer := &SliceIndex{}
	//fieldSlicer := &SliceIndex{}
	//if len(args) == 1 {
	//	// no space between commas
	//	sliceParts := strings.Split(args[0], ",")
	//	if len(sliceParts) == 1 {
	//		rowSlices := strings.Split(sliceParts[0], ":")
	//		if len(rowSlices[0]) != 0 {
	//			start, err := strconv.Atoi(rowSlices[0])
	//			if err != nil { return nil, err }
	//			lineSlicer.Start = start
	//			lineSlicer.HasStart = true
	//		}
	//
	//		if len(rowSlices) == 2 {}
	//
	//	}
	//}
}

func PrintVersionAndQuit() {
	fmt.Println(Version)
	os.Exit(0)
}

func RunIx(args []string) {
	// todo: stdin/file routing
	if !HasStdinPipe() {
		log.Fatal("No stdin found")
	}
	arg := ""
	if len(args) > 0 {
		arg = args[0]
	}
	multiSlice, err := ParseMultiSlice(arg)
	if err != nil {
		log.Fatal(err)
	}
	vprint.Printf("LineSlice: %v \n", multiSlice.LineSlicer)
	RunIxStdin(multiSlice)
}

// RootCmd represents the root command
var RootCmd = &cobra.Command{
	Use:   "ix",
	Short: "Utility for splitting and slicing file output",
	Long: `Use python-like slice notation instead of head/tail  
`,
	Run: func(cmd *cobra.Command, args []string) {
		doVersion, _ := cmd.Flags().GetBool("version")
		if doVersion {
			PrintVersionAndQuit()
		}

		vprint.Println("args: ", args)

		RunIx(args)
		//_ = cmd.Help()
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}

func init() {

	RootCmd.PersistentFlags().BoolP("silent", "S", false, "Suppress errors")
	RootCmd.PersistentFlags().StringP("input", "i", "", "Input file to read (optional)")
	RootCmd.PersistentFlags().StringP("sep", "s", "", "Separator between horizontal fields")
	RootCmd.PersistentFlags().StringP("recordsep", "R", "\n", "Separator between vertical records/lines")

	// Runtime
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose tracing (in progress)")
	RootCmd.PersistentFlags().BoolP("version", "V", false, "Print version and quit")

}
