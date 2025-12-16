# Optimisations de Performance pour Cloud Run

## Problèmes Identifiés et Solutions

### 1. ✅ Configuration MongoDB Non Optimisée

**Problème:** Le client MongoDB utilisait les paramètres par défaut, ce qui peut causer des problèmes de performance sur Cloud Run.

**Solution:** Configuration optimisée du pool de connexions MongoDB:
- `MaxPoolSize: 50` - Maximum de connexions dans le pool
- `MinPoolSize: 5` - Minimum de connexions maintenues (évite les cold starts)
- `MaxConnIdleTime: 30s` - Ferme les connexions inactives après 30s
- `ConnectTimeout: 10s` - Timeout pour la connexion initiale
- `ServerSelectionTimeout: 5s` - Timeout pour la sélection du serveur
- `SocketTimeout: 30s` - Timeout pour les opérations socket
- `HeartbeatInterval: 10s` - Intervalle de heartbeat pour le monitoring

**Fichier modifié:** `server.go`

### 2. ✅ Configuration Cloud Run Sous-Optimale

**Problème:** 
- `minScale: 0` causait des cold starts à chaque requête
- CPU/Memory insuffisants (1 CPU, 1Gi)

**Solution:**
- `minScale: 1` - Maintient au moins 1 instance active (évite les cold starts)
- `cpu: 2` - Augmentation à 2 CPUs pour meilleure performance
- `memory: 2Gi` - Augmentation à 2Gi de mémoire
- `startup-cpu-boost: true` - Boost CPU pendant le démarrage

**Fichier modifié:** `cloud-run.yaml`

### 3. ✅ Serveur HTTP Sans Timeouts

**Problème:** Le serveur HTTP n'avait pas de timeouts configurés, ce qui pouvait causer des connexions qui traînent.

**Solution:** Configuration des timeouts HTTP:
- `ReadTimeout: 15s` - Durée maximale pour lire la requête complète
- `WriteTimeout: 15s` - Durée maximale avant timeout des écritures
- `IdleTimeout: 60s` - Temps maximum d'attente pour la prochaine requête
- `ReadHeaderTimeout: 5s` - Temps alloué pour lire les en-têtes

**Fichier modifié:** `server.go`

## Recommandations Supplémentaires

### Monitoring
- Surveiller les métriques Cloud Run: latence, CPU, mémoire, nombre d'instances
- Surveiller les connexions MongoDB dans MongoDB Atlas
- Activer les logs structurés pour identifier les requêtes lentes

### Optimisations Futures Possibles

1. **Caching:**
   - Implémenter un cache Redis pour les requêtes fréquentes
   - Utiliser le cache GraphQL existant plus efficacement

2. **Indexes MongoDB:**
   - Vérifier que tous les indexes nécessaires sont créés
   - Analyser les requêtes lentes avec `explain()`

3. **Connection Pooling:**
   - Surveiller l'utilisation du pool de connexions
   - Ajuster `MaxPoolSize` et `MinPoolSize` selon la charge

4. **GraphQL Query Optimization:**
   - Implémenter DataLoader pour éviter les N+1 queries
   - Optimiser les requêtes GraphQL complexes

5. **Database Queries:**
   - Utiliser `projection` pour limiter les champs retournés
   - Implémenter la pagination efficace pour les grandes listes

## Tests de Performance

Pour tester les améliorations:

1. **Avant déploiement:**
   ```bash
   # Test local avec load testing
   ab -n 1000 -c 10 http://localhost:8080/query
   ```

2. **Après déploiement:**
   - Surveiller les métriques Cloud Run
   - Comparer les temps de réponse avant/après
   - Vérifier la réduction des cold starts

## Notes Importantes

- Les timeouts MongoDB sont maintenant gérés au niveau du client
- Les timeouts HTTP protègent contre les connexions qui traînent
- `minScale: 1` augmente légèrement les coûts mais améliore significativement la performance
- Surveiller les coûts Cloud Run après ces changements












