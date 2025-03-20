# Go-Orb: Vom Nullpunkt zum Helden in der Go-Entwicklung

Hey! Go-Orb ist ein umfassendes, leistungsstarkes Framework zum Erstellen von Monolithen und verteilten Systemen in Go. Wir haben es als Nachfolger von go-micro konzipiert, mit jeder Menge Verbesserungen bei Architektur, Performance und Entwicklererfahrung.

> **Beta-Release kommt bald!** Wir entwickeln Go-Orb aktiv und es wird bald als Beta-Version verfügbar sein. Bleib dran für Updates!

## Warum Go-Orb wählen?

### Produktionsreife Architektur

Go-Orb gibt dir eine felsenfeste Grundlage für den Aufbau verteilter Systeme mit:

- **Fast keine Reflection**: Bessere Typsicherheit und schnellere Performance durch Verzicht auf Laufzeit-Reflection
- **Wire-basierte Dependency Injection**: Kompilierzeit-Sicherheit ohne fiese Globals oder Laufzeit-Überraschungen
- **Steckbare Architektur**: Tausche Komponenten aus, ohne deinen Anwendungscode zu verändern

### Sauberes, schnittstellenbasiertes Design

Go-Orb hält die Dinge übersichtlich mit einer klaren Trennung der Belange:

- **Nur Kernschnittstellen**: [go-orb/go-orb](https://github.com/go-orb/go-orb) besteht nur aus Interfaces und minimalem Verbindungscode
- **Plugins machen die Arbeit**: Die eigentlichen Implementierungen findest du in [go-orb/plugins](https://github.com/go-orb/plugins)
- **Einfaches Plugin-System**: Füge einfach einen leeren Import hinzu (`import _ "github.com/go-orb/plugins/..."`) und schon kann's losgehen
- **Mix und Match**: Wähle genau die Plugins, die zu deinen Anforderungen passen
- **Erweiterbar**: Erstelle ganz einfach deine eigenen Plugins, die Kernschnittstellen implementieren

### Vereinfachte Entwicklung verteilter Systeme

Konzentriere dich auf deine Geschäftslogik, während Go-Orb alle Kopfschmerzen mit verteilten Systemen für dich übernimmt:

- **Service-Erkennung**: Automatische Dienstregistrierung und Namensauflösung
- **Lastverteilung**: Intelligente Anfragenverteilung über Service-Instanzen
- **Fehlertoleranz**: Eingebaute Wiederholungsversuche und Circuit-Breaking, damit alles glatt läuft
- **Nachrichtenkodierung**: Dynamische, inhaltsbasierte Kodierung und Dekodierung

### Super-schnelle In-Memory-Kommunikation

Go-Orb hat einen High-Performance In-Memory-Adapter, der dir Folgendes ermöglicht:

- **Direkte Handler-Aufrufe**: Blitzschnelle prozessinterne Kommunikation ohne Serialisierungs-Overhead
- **Gleiche API wie Netzwerkaufrufe**: Nutze dieselbe Client-Schnittstelle, egal ob du lokale oder entfernte Dienste aufrufst
- **Perfekt für Monolithen**: Starte mit allem in einem Prozess und splitte Services auf, wenn nötig
- **Nahtloses Testen**: Teste deine Dienste isoliert und ohne Netzwerkabhängigkeiten

### Unterstützung für moderne Protokolle

Kommuniziere genau wie du es brauchst mit Unterstützung für:

- **gRPC**: Hochleistungs-RPC mit bidirektionalem Streaming
- **HTTP/HTTPS**: RESTful APIs mit voller Unterstützung für HTTP/1.1 und HTTP/2
- **DRPC**: Ultra-schnelle RPC-Alternative mit reduziertem Overhead
- **HTTP/3**: Next-Gen HTTP mit QUIC
- **Ereignisgesteuerte Kommunikation**: Asynchrones Messaging für entkoppelte Architekturen

Schau dir unsere [Benchmarks](https://github.com/go-orb/go-orb/wiki/RPC-Benchmarks) für mehr Details an.

### Entwicklerfreundliche Erfahrung

Schnell loslegen und produktiv bleiben:

- **Intuitive APIs**: Saubere, konsistente Schnittstellen, die leicht zu verstehen sind
- **Flexible Konfiguration**: Konfiguriere über Dateien, Umgebungsvariablen oder Code
- **Umfassende Dokumentation**: Detaillierte Anleitungen und Beispiele für den Einstieg
- **Starke Testunterstützung**: Mit Testbarkeit im Hinterkopf entwickelt

### Qualitätsfokussierte Entwicklung

Bei Go-Orb dreht sich alles um Codequalität und Zuverlässigkeit:

- **Umfassende statische Analyse**: Die gesamte Codebasis wird mit golangci-lint nach strengen Regeln überprüft
- **Umfangreiche Testsuite**: Hohe Testabdeckung über alle Komponenten und Plugins
- **CI-Durchsetzung**: Qualitätschecks laufen automatisch für alle Pull Requests
- **Keine Kompromisse**: Strenge Linting- und Testanforderungen sorgen für gleichbleibende Qualität
- **Produktionsreif**: Unsere strengen Entwicklungspraktiken bedeuten, dass du Go-Orb auch nach der Beta-Version vertrauen kannst

## Was macht Go-Orb besonders?

### Klein anfangen, groß skalieren

Go-Orb ist darauf ausgelegt, mit deiner Anwendung zu wachsen:

- **Starte mit einem Monolithen**: Beginne die Entwicklung mit allen Diensten in einem Prozess über den In-Memory-Adapter
- **Umstieg auf Microservices**: Extrahiere Dienste schrittweise, ohne deine Geschäftslogik zu ändern
- **Hybride Architektur**: Führe leistungskritische Komponenten prozessintern aus, während du andere verteilst
- **Progressive Skalierung**: Füge mehr Instanzen bestimmter Dienste hinzu, wenn deine Last steigt

### Mehrere Einstiegspunkte

Im Gegensatz zu herkömmlichen Frameworks, die dich auf ein einzelnes Protokoll festlegen, lässt dich Go-Orb deine Dienste gleichzeitig über mehrere Protokolle verfügbar machen. Konfiguriere verschiedene Handler für verschiedene Protokolle in einem einzigen, zusammenhängenden Dienst.

### Leistungsstarkes, mehrschichtiges Konfigurationssystem

Das Konfigurationssystem von Go-Orb bietet unglaubliche Flexibilität:

- **Intelligentes Konfigurationsladen**: Führt automatisch Konfigurationen aus mehreren Quellen in einer Prioritätsreihenfolge zusammen:
  1. Vordefinierte Standardwerte
  2. Benutzerkonfigurationsdateien (lokal oder entfernt)
  3. Umgebungsvariablen
  4. Kommandozeilenargumente
- **Abschnittsbasiertes Laden**: Lade nur die Konfigurationsabschnitte, die du brauchst, wenn du sie brauchst
- **Keine riesigen Structs nötig**: Definiere Konfigurationstypen, die genau zu den Anforderungen deiner Komponente passen
- **Format-agnostisch**: Unterstützung für YAML, TOML, JSON und mehr

### Erweiterte Konfiguration

Definiere deine gesamte Service-Architektur in einfachem YAML, TOML oder JSON:

```yaml
service1:
  server:
    logging:
        plugin: lumberjack
        level: INFO
    handlers:
      - UserInfo
    middlewares:
      - middleware-1
      - middleware-2
    entrypoints:
      - name: hertzhttp
        plugin: hertz
        http2: false
        insecure: true

      - name: grpc
        plugin: grpc
        insecure: true
        reflection: false

      - name: http
        plugin: http
        insecure: true

      - name: drpc
        plugin: drpc
  client:
    middlewares:
      - name: log
      - name: retry
    logging:
        level: TRACE
  registry:
    plugin: kvstore
    kvstore:
        plugin: natsjs
        servers:
        - nats://localhost:9222
```

### Unsere Vision: Ein Unix-artiges Service-Ökosystem

Wir bauen Go-Orb mit einer Vision, die von der Unix-Philosophie "Eine Sache machen und sie gut machen" inspiriert ist:

- **Fokussiertes Service-Design**: Jeder Dienst sollte eine klare, einzelne Verantwortung haben
- **Komponierbare Architektur**: Dienste arbeiten über klar definierte Schnittstellen zusammen
- **Sofort einsetzbare Komponenten**: Wir arbeiten an einer Bibliothek vorgefertigter Dienste wie API-Gateways und Auth-Systeme
- **Einfache Integration**: Wire-basierte Dependency Injection macht Service-Komposition natürlich
- **Baue dein Ökosystem**: Erstelle deine eigene maßgeschneiderte Plattform, indem du genau das mischst und kombinierst, was du brauchst

Während wir Go-Orb weiterentwickeln, leitet diese Vision unseren Fahrplan und unsere Architekturentscheidungen.

Schau dir [services](https://github.com/go-orb/services) für mehr Details an.

### Protokollkonforme Handler

Schreibe Handler einmal und stelle sie über jedes Protokoll bereit:

```go
// Einfache, typsichere Client-Aufrufe
resp, err := client.Call[HelloResponse](
    context.Background(), 
    clientDi, 
    "org.orb.svc.hello", 
    "Say.Hello", 
    &req
)
```

Oder typsichere generierte Handler:

```go
cli := authproto.NewAuthClient(clientFromWire)
req := &authproto.Req{Token: "someToken"}
resp, err := cli.Authenticate(
    ctx,
    serverName,
    req,
    opts...,
)
```

### Strukturiertes Logging

Go-Orb kommt mit eingebauter Unterstützung für modernes strukturiertes Logging basierend auf Go's Standard-Bibliothek slog, was Debugging und Monitoring zum Kinderspiel macht.

## Anwendungsfälle

- **Microservices-Architektur**: Baue, deploye und skaliere einzelne Dienste unabhängig
- **API-Gateways**: Erstelle einheitliche Eingangspunkte für dein Microservices-Ökosystem
- **Ereignisgesteuerte Systeme**: Implementiere Pub/Sub-Muster für asynchrone Verarbeitung
- **Edge Computing**: Deploye leichtgewichtige Dienste näher an deinen Benutzern
- **Cloud-Native-Anwendungen**: Perfekt für containerisierte Umgebungen und Kubernetes

## Erste Schritte

Der beste Weg, mit Go-Orb zu starten, ist unser Examples-Repository unter [github.com/go-orb/examples](https://github.com/go-orb/examples). Es ist vollgepackt mit Beispieldiensten, die Go-Orbs Fähigkeiten zeigen:

- Einfache Dienste mit verschiedenen Protokollen
- Ereignisgesteuerte Architekturen
- API-Gateways und Proxies
- Authentifizierung und Autorisierung
- Performance-Benchmarks

Schau im Repository vorbei, um Schritt-für-Schritt-Anleitungen zum Ausführen der Beispiele zu erhalten und deine eigenen Dienste mit Go-Orb zu erstellen.

Besuche [github.com/go-orb/go-orb](https://github.com/go-orb/go-orb), um mehr über das Core-Framework zu erfahren.

## Community und Support

Tritt unserer freundlichen Community bei:

- **Matrix**: [https://matrix.to/#go-orb:jochum.dev](https://matrix.to/#/#go-orb:jochum.dev) - Echtzeit-Chat und Support
- **Discord**: [https://discord.gg/go-orb](https://discord.gg/go-orb) - Eine weitere Möglichkeit, dich zu vernetzen
- **GitHub**: [https://github.com/go-orb](https://github.com/go-orb) - Open-Source-Entwicklung und Issue-Tracking