package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"bureau/internal/auth"
	"bureau/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Charger la configuration
	cfg := config.Load()

	// Demander le nouveau mot de passe
	var newPassword string
	if len(os.Args) > 1 {
		newPassword = os.Args[1]
	} else {
		fmt.Print("Entrez le nouveau mot de passe pour admin@mlm.com: ")
		fmt.Scanln(&newPassword)
	}

	if newPassword == "" {
		fmt.Println("❌ Le mot de passe ne peut pas être vide")
		os.Exit(1)
	}

	// Valider le mot de passe
	if err := auth.ValidatePassword(newPassword); err != nil {
		fmt.Printf("❌ Erreur de validation du mot de passe: %v\n", err)
		fmt.Println("\nLe mot de passe doit respecter ces règles:")
		fmt.Println("- Minimum 8 caractères")
		fmt.Println("- Au moins une majuscule (A-Z)")
		fmt.Println("- Au moins une minuscule (a-z)")
		fmt.Println("- Au moins un chiffre (0-9)")
		fmt.Println("- Au moins un caractère spécial parmi : @$!%*?&")
		os.Exit(1)
	}

	// Générer le hash du mot de passe
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		fmt.Printf("❌ Erreur lors du hashage: %v\n", err)
		os.Exit(1)
	}

	// Se connecter à MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("❌ Erreur de connexion à MongoDB: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	// Vérifier la connexion
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Printf("❌ Erreur de ping MongoDB: %v\n", err)
		os.Exit(1)
	}

	db := client.Database(cfg.MongoDBName)
	adminsCollection := db.Collection("admins")

	// Vérifier si l'admin existe
	adminEmail := "admin@mlm.com"
	var adminDoc bson.M
	err = adminsCollection.FindOne(ctx, bson.M{"email": adminEmail}).Decode(&adminDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("❌ Aucun admin trouvé avec l'email: %s\n", adminEmail)
			fmt.Println("💡 Voulez-vous créer un nouvel admin? (y/N)")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				os.Exit(1)
			}

			// Créer un nouvel admin
			newAdmin := bson.M{
				"email":        adminEmail,
				"name":         "Admin",
				"role":         "admin",
				"passwordHash": hashedPassword,
				"createdAt":    time.Now(),
			}
			_, err = adminsCollection.InsertOne(ctx, newAdmin)
			if err != nil {
				fmt.Printf("❌ Erreur lors de la création de l'admin: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Nouvel admin créé avec succès!\n")
			fmt.Printf("📧 Email: %s\n", adminEmail)
			fmt.Printf("🔑 Mot de passe: %s\n", newPassword)
			return
		}
		fmt.Printf("❌ Erreur lors de la recherche de l'admin: %v\n", err)
		os.Exit(1)
	}

	// Mettre à jour le mot de passe
	updateResult, err := adminsCollection.UpdateOne(
		ctx,
		bson.M{"email": adminEmail},
		bson.M{"$set": bson.M{"passwordHash": hashedPassword}},
	)
	if err != nil {
		fmt.Printf("❌ Erreur lors de la mise à jour: %v\n", err)
		os.Exit(1)
	}

	if updateResult.MatchedCount == 0 {
		fmt.Printf("❌ Aucun admin trouvé avec l'email: %s\n", adminEmail)
		os.Exit(1)
	}

	fmt.Println("✅ Mot de passe réinitialisé avec succès!")
	fmt.Printf("📧 Email: %s\n", adminEmail)
	fmt.Printf("🔑 Nouveau mot de passe: %s\n", newPassword)
	fmt.Println("\n💡 Vous pouvez maintenant vous connecter avec ces identifiants.")
}
