package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const download_dir = "./downloads"

type BuildArgs struct {
	architecture, version string
	clean                 bool
}

func (arg *BuildArgs) Parse() error {
	flag.StringVar(&arg.architecture, "arch", "", "Target architecture: amd64, arm64")
	flag.StringVar(&arg.version, "version", "", "Build version")
	flag.BoolFunc("clean", "Clean existing build", func(s string) error {
		switch s {
		case "true":
			arg.clean = true
		case "false":
			arg.clean = false
		default:
			return fmt.Errorf("error: Invalid value %q provided to the parameter \"clean\"", s)
		}
		return nil
	})
	flag.Parse()

	if arg.clean {
		return nil
	}

	switch arg.architecture {
	case "amd64", "arm64":
		// Do nothing
	case "":
		return fmt.Errorf("error: \"arch\" flag is not set")
	default:
		return fmt.Errorf("error: Invalid value %q provided to the parameter \"arch\"", arg.architecture)
	}

	switch arg.version {
	case "":
		return fmt.Errorf("error: \"version\" flag is not set")
	default:
		// Do nothing
	}

	return nil

}

type Downloads struct {
	vscode, license string
}

func (d *Downloads) Parse(arg BuildArgs) {
	download_dir, _ := filepath.Abs(download_dir)
	d.vscode = filepath.Join(download_dir, fmt.Sprintf("openvscode-server-v%v-linux-%v.tar.gz", arg.version, arg.architecture))
	d.license = filepath.Join(download_dir, "LICENSE.txt")
}

func (d *Downloads) CreateDownloadDir() error {
	download_dir, _ := filepath.Abs(download_dir)
	utils := Utils{}
	if !utils.IsPathExists(download_dir) {
		err := os.Mkdir(download_dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

type CommandWithArgs []string

func runCommand(cmd_with_args CommandWithArgs) error {
	var cmd *exec.Cmd
	var err error

	switch len(cmd_with_args) {
	case 1:
		cmd = exec.Command(cmd_with_args[0])
	default:
		cmd = exec.Command(cmd_with_args[0], cmd_with_args[1:]...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	return err
}

func downloadAll(arg BuildArgs, downloads Downloads) error {
	var err error
	var utils = Utils{}
	if arg.architecture == "amd64" {
		arg.architecture = "x64"
	}

	err = downloads.CreateDownloadDir()
	if err != nil {
		return err
	}

	if !utils.IsPathExists(downloads.vscode) {

		vscode_download_cmd := CommandWithArgs{
			"curl", "-L", "--output", downloads.vscode,
			fmt.Sprintf("https://github.com/gitpod-io/openvscode-server/releases/download/openvscode-server-v%v/openvscode-server-v%v-linux-%v.tar.gz", arg.version, arg.version, arg.architecture),
		}

		err = runCommand(vscode_download_cmd)
		if err != nil {
			return err
		}
	}

	if !utils.IsPathExists(downloads.license) {
		license_download_cmd := CommandWithArgs{
			"curl", "-L", "--output", downloads.license,
			fmt.Sprintf("https://raw.githubusercontent.com/gitpod-io/openvscode-server/refs/tags/openvscode-server-v%v/LICENSE.txt", arg.version),
		}

		err = runCommand(license_download_cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func clean() error {
	var err error
	var utils = Utils{}

	if utils.IsPathExists(download_dir) {
		err = os.RemoveAll(download_dir)
		if err != nil {
			return err
		}
	}

	return nil
}

type Utils struct{}

func (u Utils) IsPathExists(local_path string) bool {
	stat, _ := os.Stat(local_path)
	return stat != nil
}

func (u Utils) ExtractTarGz(tar_gz_file, target_dir string) error {
	tar_gz_extract_cmd := CommandWithArgs{
		"tar", "xzf", tar_gz_file, "--directory", target_dir,
	}

	fmt.Printf("Extracting %q\n", filepath.Base(tar_gz_file))
	err := runCommand(tar_gz_extract_cmd)
	if err != nil {
		return err
	}
	fmt.Println("Extraction done.")

	return nil
}

func (u Utils) CopyFileOrFolder(source_path, dest_path string, create_dest_path_tree bool) error {
	cp_cmd := CommandWithArgs{
		"cp", "-avrf", source_path, dest_path,
	}

	fmt.Println("Copying file folder...")
	err := runCommand(cp_cmd)
	if err != nil {
		return err
	}
	fmt.Println("Copying done.")

	return nil
}

func build_package() {
	var err error
	utils := Utils{}

	build_args := BuildArgs{}
	err = build_args.Parse()
	if err != nil {
		panic(err)
	}

	if build_args.clean {
		err = clean()
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	downloads := Downloads{}
	downloads.Parse(build_args)
	err = downloadAll(build_args, downloads)
	if err != nil {
		panic(err)
	}

	err = utils.ExtractTarGz(downloads.vscode, download_dir)
	if err != nil {
		panic(err)
	}
}

func main() {
	build_package()
}
