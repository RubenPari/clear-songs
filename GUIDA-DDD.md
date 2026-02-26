🏗️ Go DDD API Project Structure Guide
Questa guida descrive l'organizzazione del progetto seguendo i principi del Domain-Driven Design (DDD) e della Clean Architecture adattati all'ecosistema Go.

📁 Struttura delle Cartelle
Plaintext
.
├── cmd/
│ └── server/
│ └── main.go # Entry point: Inizializzazione e Dependency Injection
├── internal/
│ ├── domain/ # Cuore del business (Logica pura, nessuna dipendenza esterna)
│ │ ├── product/
│ │ │ ├── entity.go # Definizioni delle strutture (es. Product)
│ │ │ ├── repository.go # INTERFACCE per il salvataggio dati
│ │ │ └── service.go # Logica di dominio complessa (es. calcolo sconti)
│ │ └── shared/ # Value Objects condivisi o errori comuni
│ ├── application/ # Casi d'uso (Orchestratore tra API e Dominio)
│ │ └── product/
│ │ ├── service.go # App Service: coordina repository e dominio
│ │ └── dto.go # Data Transfer Objects (Input/Output per l'esterno)
│ ├── infrastructure/ # Dettagli implementativi (Framework e Driver)
│ │ ├── persistence/ # Implementazione REALE dei repository (SQL, NoSQL)
│ │ │ └── postgres/
│ │ ├── transport/ # Livello di comunicazione (HTTP/gRPC)
│ │ │ └── http/
│ │ │ ├── handlers.go
│ │ │ └── router.go
│ │ └── config/ # Gestione variabili d'ambiente e setup
└── go.mod
층 (Layers) e Responsabilità

1. Domain Layer (internal/domain)
   È il livello più importante. Contiene la "verità" del business.

Regola d'oro: Non può importare nulla dagli altri layer (application o infrastructure).

Contenuto: \* Entities: Strutture dati con ID unico.

Repository Interfaces: Contribuiscono al disaccoppiamento (Dependency Inversion).

2. Application Layer (internal/application)
   Agisce come un vigile urbano. Riceve comandi e interroga il dominio.

Responsabilità: Validazione dei dati di input, gestione delle transazioni, invio di email/notifiche dopo un'azione di successo.

DTOs: Definisce come i dati appaiono all'esterno (nascondendo campi sensibili del database).

3. Infrastructure Layer (internal/infrastructure)
   Contiene tutto ciò che è considerato un "dettaglio tecnico".

Persistence: Qui scrivi le query SQL. Se domani cambi DB, tocchi solo questa cartella.

Transport: Definisce le rotte HTTP (es. Gin, Echo) e converte le richieste JSON in chiamate ai servizi applicativi.

🛠️ Esempio Pratico: Dependency Inversion
Per mantenere il codice pulito, usiamo le interfacce nel dominio e le implementiamo nell'infrastruttura.

1. Domain (internal/domain/product/repository.go)

Go
package product

// Definiamo cosa vogliamo fare, non come.
type Repository interface {
GetByID(id string) (*Product, error)
Save(p *Product) error
} 2. Infrastructure (internal/infrastructure/persistence/postgres/product_repo.go)

Go
package postgres

import "internal/domain/product"

type PostgresRepo struct {
db \*sql.DB
}

// Implementazione reale
func (r *PostgresRepo) Save(p *product.Product) error {
return r.db.Exec("INSERT INTO products ...", p.Name)
}
🚀 Ciclo di una richiesta
Client invia una richiesta POST /products.

Infrastructure (HTTP Handler) riceve il JSON, lo mappa in un DTO.

Application (Service) riceve il DTO, chiama il Domain (Entity) per creare il prodotto.

Application chiama l'interfaccia Repository per salvare.

Infrastructure (Postgres) esegue fisicamente la query.

✅ Vantaggi
Testabilità: Puoi fare Unit Test del dominio senza database.

Manutenibilità: Il codice è diviso per contesti logici, non per tipologia di file.

Evoluzione: Puoi cambiare database o framework web senza riscrivere la logica di business.
