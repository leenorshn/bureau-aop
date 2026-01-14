package main

import (
	"fmt"
	"os"

	"bureau/internal/auth"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate-password-hash.go <password>")
		os.Exit(1)
	}

	password := os.Args[1]

	// Valider le mot de passe
	if err := auth.ValidatePassword(password); err != nil {
		fmt.Printf("❌ Erreur de validation: %v\n", err)
		os.Exit(1)
	}

	// Générer le hash
	hash, err := auth.HashPassword(password)
	if err != nil {
		fmt.Printf("❌ Erreur lors du hashage: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Hash généré avec succès:\n%s\n", hash)
	fmt.Printf("\nPour mettre à jour dans MongoDB:\n")
	fmt.Printf("db.admins.updateOne(\n")
	fmt.Printf("  { email: \"admin@mlm.com\" },\n")
	fmt.Printf("  { $set: { passwordHash: \"%s\" } }\n", hash)
	fmt.Printf(")\n")
}



