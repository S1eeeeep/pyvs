package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/S1eeeeep/pyvs/utils/file"
	"github.com/S1eeeeep/pyvs/utils/python"
	"github.com/S1eeeeep/pyvs/utils/web"
	"github.com/urfave/cli"
	"golang.org/x/sys/windows/registry"
)

var version = "1.0.0"

const PATH = "PATH"

const (
	defaultOriginalpath = "https://hub.gitmirror.com/raw.githubusercontent.com/S1eeeeep/pyvs/refs/heads/main/pyindex.json"
)

type Config struct {
	PythonHome           string `json:"python_home"`
	CurrentPythonVersion string `json:"current_python_version"`
	Originalpath         string `json:"original_path"`
	store                string
	download             string
}

var config Config

type PythonVersion struct {
	Version string `json:"version"`
	Url     string `json:"url"`
}

func main() {
	app := cli.NewApp()
	app.Name = "pyvs"
	app.Usage = `Python Version Switcher (PYVS) for Windows`
	app.Version = version

	app.CommandNotFound = func(c *cli.Context, command string) {
		log.Fatal("Command Not Found")
	}
	app.Commands = commands()
	app.Before = startup
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "init",
			Usage:       "Initialize config file",
			Description: `before init you should clear PYTHON_HOME, PATH Environment variable。`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "python_home",
					Usage: "the PYTHON_HOME location",
					Value: filepath.Join(os.Getenv("ProgramFiles"), "python"),
				},
				cli.StringFlag{
					Name:  "originalpath",
					Usage: "the python download index file url.",
					Value: defaultOriginalpath,
				},
			},
			Action: func(c *cli.Context) error {
				if c.IsSet("python_home") || config.PythonHome == "" {
					config.PythonHome = c.String("python_home")
				}

				if c.IsSet("originalpath") || config.Originalpath == "" {
					config.Originalpath = c.String("originalpath")
				}

				osPath, err := getSystemEnvVariable(PATH)
				if err != nil {
					return errors.New("get Environment variable `path` failure: Please run as admin user")
				}

				if !strings.Contains(osPath, config.PythonHome) {
					err := setSystemEnvVariable(PATH, fmt.Sprintf("%s;%s", config.PythonHome, osPath))
					if err != nil {
						return errors.New("set Environment variable `PYTHON_HOME` failure: Please run as admin user")
					}
					fmt.Println("add PYTHON_HOME to `path` Environment variable")
				}

				osPath, err = getSystemEnvVariable(PATH)
				if err != nil {
					return errors.New("get Environment variable `path` failure: Please run as admin user")
				}

				if !strings.Contains(osPath, file.GetCurrentPath()) {
					err := setSystemEnvVariable(PATH, fmt.Sprintf("%s%s;", osPath, file.GetCurrentPath()))
					if err != nil {
						return errors.New("set Environment variable `CURRENT_PATH` failure: Please run as admin user")
					}
					fmt.Println("add pyvs.exe to `path` Environment variable")
				}

				return nil
			},
		},
		{
			Name:      "list",
			ShortName: "ls",
			Usage:     "List current Python installations.",
			Action: func(c *cli.Context) error {
				fmt.Println("Installed python (* marks in use):")
				v := python.GetInstalled(config.store)
				for i, version := range v {
					str := ""
					if config.CurrentPythonVersion == version {
						str = fmt.Sprintf("%s  * %d) %s", str, i+1, version)
					} else {
						str = fmt.Sprintf("%s    %d) %s", str, i+1, version)
					}
					fmt.Println(str)
				}
				if len(v) == 0 {
					fmt.Println("No installations recognized.")
				}
				return nil
			},
		},
		{
			Name:      "install",
			ShortName: "i",
			Usage:     "Install available remote python",
			Action: func(c *cli.Context) error {
				v := c.Args().Get(0)
				if v == "" {
					return errors.New("invalid version., Type \"pyvs rls\" to see what is available for install")
				}

				if python.IsVersionInstalled(config.store, v) {
					fmt.Println("Version " + v + " is already installed.")
					return nil
				}
				versions, err := getPythonVersions()
				if err != nil {
					return err
				}

				if !file.Exists(config.download) {
					os.MkdirAll(config.download, 0777)
				}
				if !file.Exists(config.store) {
					os.MkdirAll(config.store, 0777)
				}

				for _, version := range versions {
					if version.Version == v {
						dlzipfile, success := web.GetPython(config.download, v, version.Url)
						if success {
							fmt.Printf("Installing Python %s ...\n", v)

							// Extract python to the temp directory
							pythontempfile := filepath.Join(config.download, fmt.Sprintf("%s_temp", v))
							if file.Exists(pythontempfile) {
								err := os.RemoveAll(pythontempfile)
								if err != nil {
									panic(err)
								}
							}
							err := file.Unzip(dlzipfile, pythontempfile)
							if err != nil {
								return fmt.Errorf("unzip failed: %w", err)
							}

							// Copy the jdk files to the installation directory
							temPythonHome := getPythonHome(pythontempfile)
							err = os.Rename(temPythonHome, filepath.Join(config.store, v))
							if err != nil {
								return fmt.Errorf("unzip failed: %w", err)
							}

							// Remove the temp directory
							// may consider keep the temp files here
							os.RemoveAll(pythontempfile)

							fmt.Println("Installation complete. If you want to use this version, type\n\npyvs switch", v)
						} else {
							fmt.Println("Could not download Python " + v + " executable.")
						}
						return nil
					}
				}
				return errors.New("invalid version., Type \"pyvs rls\" to see what is available for install")
			},
		},
		{
			Name:  "rls",
			Usage: "Show a list of versions available for download. ",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "a",
					Usage: "list all the version",
				},
			},
			Action: func(c *cli.Context) error {
				versions, err := getPythonVersions()
				if err != nil {
					return err
				}
				for i, version := range versions {
					fmt.Printf("    %d) %s\n", i+1, version.Version)
					if !c.Bool("a") && i >= 9 {
						fmt.Println("\nuse \"pyvs rls -a\" show all the versions ")
						break
					}
				}
				if len(versions) == 0 {
					fmt.Println("No availabled python version for download.")
				}

				fmt.Printf("\nFor a complete list, visit %s\n", config.Originalpath)
				return nil
			},
		},
		{
			Name:      "switch",
			ShortName: "s",
			Usage:     "Switch to use the specified version or index number.",
			Action: func(c *cli.Context) error {
				v := c.Args().Get(0)
				if v == "" {
					return errors.New("you should input a version or index number, Type \"pyvs list\" to see what is installed")
				}

				// Check if input is a number (index)
				index, err := strconv.Atoi(v)
				if err == nil && index > 0 {
					// Input is a valid number, get the list of installed python
					installedPys := python.GetInstalled(config.store)
					if len(installedPys) == 0 {
						return errors.New("no python installations found")
					}

					if index > len(installedPys) {
						return fmt.Errorf("invalid index: %d, should be between 1 and %d", index, len(installedPys))
					}

					v = installedPys[index-1]
					fmt.Printf("Using index %d to select python %s\n", index, v)
				}

				if !python.IsVersionInstalled(config.store, v) {
					fmt.Printf("python %s is not installed. ", v)
					return nil
				}

				// Create or update the symlink
				err = os.Remove(config.PythonHome)
				if err != nil {
					return errors.New("Switch python failed, please manually remove " + config.PythonHome)
				}

				osPath, err := getSystemEnvVariable(PATH)
				if err != nil {
					return errors.New("get Environment variable `path` failure: Please run as admin user")
				}

				if !strings.Contains(osPath, config.PythonHome) {
					fmt.Println("Please run \"pyvs init\" first.")
				}
				err = os.Symlink(filepath.Join(config.store, v), config.PythonHome)
				if err != nil {
					return errors.New("Switch python failed, " + err.Error())
				}
				fmt.Println("Switch success.\nNow using python " + v)
				config.CurrentPythonVersion = v
				return nil
			},
		},
		{
			Name:      "remove",
			ShortName: "rm",
			Usage:     "Remove a specific version.",
			Action: func(c *cli.Context) error {
				v := c.Args().Get(0)
				if v == "" {
					return errors.New("you should input a version, Type \"pyvs list\" to see what is installed")
				}
				if python.IsVersionInstalled(config.store, v) {
					fmt.Printf("Remove Python %s ...\n", v)
					if config.CurrentPythonVersion == v {
						os.Remove(config.PythonHome)
					}
					dir := filepath.Join(config.store, v)
					e := os.RemoveAll(dir)
					if e != nil {
						fmt.Println("Error removing python " + v)
						fmt.Println("Manually remove " + dir + ".")
					} else {
						fmt.Printf(" done")
					}
				} else {
					fmt.Println("python " + v + " is not installed. Type \"pyvs list\" to see what is installed.")
				}
				return nil
			},
		},
	}
}

func startup(c *cli.Context) error {
	s := file.GetCurrentPath()
	config.store = filepath.Join(s, "store")
	config.download = filepath.Join(s, "download")
	config.Originalpath = defaultOriginalpath
	config.PythonHome = filepath.Join(os.Getenv("ProgramFiles"), "python")
	return nil
}

// 设置系统环境变量（需要管理员权限）
func setSystemEnvVariable(name, value string) error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	return key.SetStringValue(name, value)
}

// 获取系统环境变量（需要管理员权限）
func getSystemEnvVariable(name string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ)
	if err != nil {
		return "", err
	}
	defer key.Close()

	value, _, err := key.GetStringValue(name)
	return value, err
}

func getPythonVersions() ([]PythonVersion, error) {
	jsonContent, err := web.GetRemoteTextFile(config.Originalpath)
	if err != nil {
		return nil, err
	}
	var versions []PythonVersion
	err = json.Unmarshal([]byte(jsonContent), &versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func getPythonHome(pythonTempFile string) string {
	var pythonHome string
	fs.WalkDir(os.DirFS(pythonTempFile), ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Base(path) == "python.exe" {
			temPath := strings.Replace(path, "python.exe", "", -1)
			pythonHome = filepath.Join(pythonTempFile, temPath)
			return fs.SkipDir
		}
		return nil
	})
	return pythonHome
}

