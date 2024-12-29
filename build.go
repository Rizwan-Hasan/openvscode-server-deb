package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	download_dir string = "./downloads"
	pkgroot_dir  string = "./pkgroot"
	service_dir  string = "./service-files"
	debian_dir   string = "./debian-files"
)

type BuildArgs struct {
	architecture, version string
	clean                 bool
}

func (arg *BuildArgs) Parse() error {
	flag.StringVar(&arg.architecture, "arch", "", "Sets the target architecture (amd64, arm64)")
	flag.StringVar(&arg.version, "version", "", "Specifies the OpenVSCode-Server version to package.")
	flag.BoolFunc("clean", "Cleans up build artifacts.", func(s string) error {
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

	switch arg.architecture {
	case "amd64":
		d.vscode = filepath.Join(download_dir, fmt.Sprintf("openvscode-server-v%v-linux-%v.tar.gz", arg.version, "x64"))
	default:
		d.vscode = filepath.Join(download_dir, fmt.Sprintf("openvscode-server-v%v-linux-%v.tar.gz", arg.version, arg.architecture))
	}

	d.license = filepath.Join(download_dir, "LICENSE.txt")
}

func (d *Downloads) CreateDownloadDir() error {
	download_dir, _ := filepath.Abs(download_dir)
	var utils UtilityFunctions = &Utils{}
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
	var utils UtilityFunctions = &Utils{}
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

func updateDebianFiles(dir, version, arch string) error {
	var err error

	sed_arch_cmd := CommandWithArgs{
		"sed", "-i", fmt.Sprintf(`s\ARCHITECTURE\%v\g`, arch), filepath.Join(dir, "control"),
	}

	err = runCommand(sed_arch_cmd)
	if err != nil {
		return err
	}
	fmt.Println("Updated deb architecture")

	sed_version_cmd := CommandWithArgs{
		"sed", "-i", fmt.Sprintf(`s\VERSION\%v\g`, version), filepath.Join(dir, "control"),
	}

	err = runCommand(sed_version_cmd)
	if err != nil {
		return err
	}
	fmt.Println("Updated deb version")

	return nil
}

func fixPermissionOfPkgRoot() error {
	var err error

	chmod_cmd := CommandWithArgs{
		"chmod", "755", pkgroot_dir,
	}

	err = runCommand(chmod_cmd)
	if err != nil {
		return err
	}
	fmt.Printf("Fixed permission of %q\n", pkgroot_dir)

	return nil

}

func cleaner() error {
	var err error
	var utils UtilityFunctions = &Utils{}

	if utils.IsPathExists(download_dir) {
		fmt.Printf("Cleaning %q\n", download_dir)
		err = os.RemoveAll(download_dir)
		if err != nil {
			return err
		}
	}

	if utils.IsPathExists(pkgroot_dir) {
		fmt.Printf("Cleaning %q\n", pkgroot_dir)
		err = os.RemoveAll(pkgroot_dir)
		if err != nil {
			return err
		}
	}

	fmt.Println("Cleaning done")

	return nil
}

type UtilityFunctions interface {
	IsPathExists(local_path string) bool
	CreateFolder(folder_path, permission string) error
	ExtractTarGz(tar_gz_file, target_dir string) error
	RenameFileOrFolder(old_path, new_path string, overwrite bool) error
	CopyFileOrFolder(source_path, dest_path string, create_dest_path bool) error
}

type Utils struct{}

func (u *Utils) IsPathExists(local_path string) bool {
	stat, _ := os.Stat(local_path)
	return stat != nil
}

func (u *Utils) ExtractTarGz(tar_gz_file, target_dir string) error {
	tar_gz_extract_cmd := CommandWithArgs{
		"tar", "xzf", tar_gz_file, "--directory", target_dir,
	}

	fmt.Printf("Extracting %q to %q\n", filepath.Base(tar_gz_file), target_dir)
	err := runCommand(tar_gz_extract_cmd)
	if err != nil {
		return err
	}
	fmt.Println("Extraction done.")

	return nil
}

func (u *Utils) CreateFolder(folder_path, permission string) error {
	mkdir_cmd := CommandWithArgs{
		"mkdir", "--parents", "--mode", permission, folder_path,
	}

	fmt.Printf("Creating folder %q", folder_path)
	err := runCommand(mkdir_cmd)
	if err != nil {
		return err
	}
	fmt.Println("Folder created.")

	return nil
}

func (u *Utils) RenameFileOrFolder(old_path, new_path string, overwrite bool) error {
	var err error

	fmt.Printf("Renaming folder %q to %q\n", old_path, new_path)

	if overwrite && u.IsPathExists(new_path) {
		err = os.RemoveAll(new_path)
		if err != nil {
			return err
		}
	}

	err = os.Rename(old_path, new_path)
	if err != nil {
		return err
	}
	fmt.Println("Folder renamed.")
	return nil
}

func (u *Utils) CopyFileOrFolder(source_path, dest_path string, create_dest_path bool) error {
	var err error
	if create_dest_path {
		err = u.CreateFolder(dest_path, "0755")
		if err != nil {
			return err
		}
	}

	cp_cmd := CommandWithArgs{
		"cp", "-arf", source_path, dest_path,
	}

	fmt.Printf("Copying %q to %q\n", source_path, dest_path)
	err = runCommand(cp_cmd)
	if err != nil {
		return err
	}
	fmt.Println("Copying done.")

	return nil
}

func build_package() {
	var err error
	var utils UtilityFunctions = &Utils{}

	build_args := BuildArgs{}
	err = build_args.Parse()
	if err != nil {
		panic(err)
	}

	if build_args.clean {
		err = cleaner()
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

	vscode_extracted_dir_name := strings.TrimSuffix(filepath.Base(downloads.vscode), ".tar.gz")
	err = utils.CopyFileOrFolder(
		filepath.Join(download_dir, vscode_extracted_dir_name),
		filepath.Join(pkgroot_dir, "opt"),
		true)
	if err != nil {
		panic(err)
	}

	err = utils.RenameFileOrFolder(
		filepath.Join(pkgroot_dir, "opt", vscode_extracted_dir_name),
		filepath.Join(pkgroot_dir, "opt", "openvscode-server"),
		true)
	if err != nil {
		panic(err)
	}

	err = utils.CopyFileOrFolder(
		downloads.license,
		filepath.Join(pkgroot_dir, "usr/share/licenses"),
		true)
	if err != nil {
		panic(err)
	}

	err = utils.RenameFileOrFolder(
		filepath.Join(pkgroot_dir, "usr/share/licenses", filepath.Base(downloads.license)),
		filepath.Join(pkgroot_dir, "usr/share/licenses", "openvscode-server"),
		true)
	if err != nil {
		panic(err)
	}

	err = utils.CopyFileOrFolder(debian_dir, pkgroot_dir, true)
	if err != nil {
		panic(err)
	}

	err = utils.RenameFileOrFolder(
		filepath.Join(pkgroot_dir, filepath.Base(debian_dir)),
		filepath.Join(pkgroot_dir, "DEBIAN"),
		true)
	if err != nil {
		panic(err)
	}

	err = updateDebianFiles(
		filepath.Join(pkgroot_dir, "DEBIAN"),
		build_args.version,
		build_args.architecture)
	if err != nil {
		panic(err)
	}

	err = fixPermissionOfPkgRoot()
	if err != nil {
		panic(err)
	}

	var deb_name string

	switch build_args.architecture {
	case "amd64":
		deb_name = fmt.Sprintf("openvscode-server-v%v-x64.deb", build_args.version)
	default:
		deb_name = fmt.Sprintf("openvscode-server-v%v-arm64.deb", build_args.version)
	}

	build_deb_cmd := CommandWithArgs{
		"dpkg-deb", "--build", "--verbose", pkgroot_dir, deb_name,
	}
	err = runCommand(build_deb_cmd)
	if err != nil {
		panic(err)
	}
}

func main() {
	build_package()
}
