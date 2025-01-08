package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"

    "golang.org/x/crypto/ssh/terminal"
    "syscall"
)

var masterPassword string

func main() {
    if err := LoadStore(); err != nil {
        log.Fatalf("Error loading password store: %v", err)
    }

    fmt.Print("Enter master password: ")
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatalf("\nError reading password: %v", err)
    }
    fmt.Println()
    masterPassword = string(bytePassword)

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("password-keeper> ")
        input, err := reader.ReadString('\n')
        if err != nil {
            log.Println("Error reading input:", err)
            continue
        }

        input = strings.TrimSpace(input)
        if input == "" {
            continue
        }

        args := strings.Split(input, " ")
        command := strings.ToLower(args[0])

        switch command {
        case "add":
            handleAdd(args)
        case "get":
            handleGet(args)
        case "delete":
            handleDelete(args)
        case "list":
            handleList()
        case "help":
            printHelp()
        case "exit", "quit":
            fmt.Println("Exiting...")
            return
        default:
            fmt.Println("Unknown command. Type 'help' for available commands.")
        }
    }
}

func handleAdd(args []string) {
    if len(args) < 2 {
        fmt.Println("Usage: add <service>")
        return
    }

    service := args[1]
    fmt.Printf("Enter password for %s: ", service)
    bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        fmt.Println("\nError reading password:", err)
        return
    }
    fmt.Println()

    encryptedPassword, err := encrypt(string(bytePassword), masterPassword)
    if err != nil {
        fmt.Println("Error encrypting password:", err)
        return
    }

    entry := PasswordEntry{
        Service:  service,
        Password: encryptedPassword,
    }

    if err := AddEntry(entry); err != nil {
        fmt.Println("Error adding entry:", err)
        return
    }

    fmt.Println("Password added successfully.")
}

func handleGet(args []string) {
    if len(args) < 2 {
        fmt.Println("Usage: get <service>")
        return
    }

    service := args[1]
    entry, err := GetEntry(service)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    decryptedPassword, err := decrypt(entry.Password, masterPassword)
    if err != nil {
        fmt.Println("Error decrypting password:", err)
        return
    }

    fmt.Printf("Password for %s: %s\n", service, decryptedPassword)
}

func handleDelete(args []string) {
    if len(args) < 2 {
        fmt.Println("Usage: delete <service>")
        return
    }

    service := args[1]
    if err := DeleteEntry(service); err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Entry deleted successfully.")
}

func handleList() {
    entries := ListEntries()
    if len(entries) == 0 {
        fmt.Println("No entries found.")
        return
    }

    fmt.Println("Services:")
    for _, e := range entries {
        fmt.Println("- " + e.Service)
    }
}

func printHelp() {
    fmt.Println("Available commands:")
    fmt.Println("  add <service>    - Add a new password for a service")
    fmt.Println("  get <service>    - Retrieve the password for a service")
    fmt.Println("  delete <service> - Delete the password entry for a service")
    fmt.Println("  list             - List all services")
    fmt.Println("  help             - Show this help message")
    fmt.Println("  exit, quit       - Exit the application")
}
