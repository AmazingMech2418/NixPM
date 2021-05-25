package main

import (
    "encoding/json"
    "fmt"
    "os/exec"
    "log"
    "os"
    "strings"
    "io/ioutil"
)


func search(query string) {
    alreadyInstalled := toOriginal(getInstalled())

    args := []string{"search", "--json"}
    args = append(args, strings.Split(query, " ")...)

    cmd := exec.Command("nix", args...)

    out, err := cmd.Output()

    errStr := fmt.Sprintf("%s", err)


    if err != nil && errStr != "exit status 1" {
        log.Fatal(err)
    }

    var dat map[string]map[string]string
	if err := json.Unmarshal([]byte(out), &dat); err != nil {
		log.Fatal(err)
	}

    fmt.Println("Results:")
    for pkg := range dat {
        fmt.Println("\t" + pkg)

        thisPkg := dat[pkg]

        fmt.Printf("\t\tName: %s\n", thisPkg["pkgName"])
        fmt.Printf("\t\tVersion: %s\n", thisPkg["version"])
        fmt.Printf("\t\tDescription: %s\n", thisPkg["description"])
        if contains(alreadyInstalled, pkg) {
            fmt.Println("\t\tInstalled")
        }
        fmt.Print("\n")
    }
}

func printData(pkgID string) {
    if pkgID == "" {
        return
    }

    alreadyInstalled := toOriginal(getInstalled())

    cmd := exec.Command("nix", "search", "--json", pkgID)

    out, err := cmd.Output()

    errStr := fmt.Sprintf("%s", err)


    if err != nil && errStr != "exit status 1" {
        log.Fatal(err)
    }

    var dat map[string]map[string]string
	if err := json.Unmarshal([]byte(out), &dat); err != nil {
		log.Fatal(err)
	}

    fmt.Printf("\t%s\n", pkgID)
    thisPkg := dat[pkgID]
    fmt.Printf("\t\tName: %s\n", thisPkg["pkgName"])
    fmt.Printf("\t\tVersion: %s\n", thisPkg["version"])
    fmt.Printf("\t\tDescription: %s\n", thisPkg["description"])
    if contains(alreadyInstalled, pkgID) {
        fmt.Println("\t\tInstalled")
    }
    fmt.Print("\n")
}

func getInstalled() []string {
    dat, err := ioutil.ReadFile("replit.nix")
    if err != nil {
        log.Fatal(err)
    }

    data := string(dat)

    indx := strings.Index(data, "deps")

    var res string = ""
    var inside bool = false
    for pos, chr := range data {
        if chr == '[' && pos > indx {
            inside = true
        } else if chr == ']' && inside {
            break
        } else if inside {
            res += fmt.Sprintf("%c", chr)
        }
    }

    var splitted []string = strings.Split(strings.Trim(res, " \t\n"), "\n")

    
    for i := range splitted {
        splitted[i] = strings.Trim(splitted[i], " \t")
    }

    return splitted
}

func toOriginal(in []string) []string {
    var dat []string = in

    for i := range in {
        dat[i] = strings.Replace(dat[i], "pkgs", "nixpkgs", 1)
    }

    return dat
}

func install(pkg string) {
    splitted := getInstalled()

    splitted = append(splitted, strings.Replace(pkg, "nixpkgs", "pkgs", 1))

    out := "{ pkgs }: {\n" + 
           "\tdeps = [\n\t\t" +
           strings.Join(splitted, "\n\t\t") + "\n" +
           "\t];\n" +
           "}"

    ioutil.WriteFile("replit.nix", []byte(out), 0644)

    remove("")
}

func initNix() {
    out1 := "{ pkgs }: {\n" + 
           "\tdeps = [\n\t\t\n" +
           "\t];\n" +
           "}"

    ioutil.WriteFile("replit.nix", []byte(out1), 0644)

    out2 := "run = \"nix-shell /opt/nixproxy.nix --argstr repldir $PWD --command 'bash main.sh'\""

    ioutil.WriteFile(".replit", []byte(out2), 0644)
}

func remove(pkg string) {
    old := getInstalled()

    splitted := []string{}

    for _, item := range old {
        if item != strings.Replace(pkg, "nixpkgs", "pkgs", 1) {
            splitted = append(splitted, item)
        }
    }


    out := "{ pkgs }: {\n" + 
           "\tdeps = [\n\t\t" +
           strings.Join(splitted, "\n\t\t") + "\n" +
           "\t];\n" +
           "}"

    ioutil.WriteFile("replit.nix", []byte(out), 0644)

    if pkg != "" {
        remove("")
    }
}


func contains(arr []string, val string) bool {
    for _, item := range arr {
        if item == val {
            return true
        }
    }
    return false
}


func main() {
    args := os.Args[1:]

    if contains(args, "-h") || contains(args, "--help") || len(args) == 0 {
        fmt.Println("NixPM - Nix Package Manager for Repl.it\n")
        fmt.Println("Commands:")
        fmt.Println("-h / --help    - display this message")
        fmt.Println("search <query> - search for a package")
        fmt.Println("install <pkg>  - install a package")
        fmt.Println("remove <pkg>   - uninstall a package")
        fmt.Println("installed      - list installed packages")
        fmt.Println("init           - initialize Nix environment")
    } else if args[0] == "search" {
        search(strings.Join(args[1:], " "))
    } else if args[0] == "install" {
        pkg := args[1]
        install(pkg)
    } else if args[0] == "remove" {
        pkg := args[1]
        remove(pkg)
    } else if args[0] ==  "installed" {
        installed := toOriginal(getInstalled())
        for _, pkg := range installed {
            printData(pkg)
        }
    } else if args[0] ==  "init" {
        initNix()
    } else {
        fmt.Println("Command " + args[0] + " not found.")
        fmt.Println("Use -h or --help to learn how to use this.")
    }
    
}
