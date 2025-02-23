---
title: Comprendre la différence entre log.Fatal() et panic() en Golang
type: blog
date: 2025-02-11
comments: true
---

## "Dans ce cas, il semble préférable d'utiliser `panic()` plutôt que `log.Fatal()`"
Récemment, j'ai reçu ce retour en utilisant `log.Fatal()`.

Oh ? `log.Fatal()` n'est-il pas simplement censé afficher les logs de manière plus précise ? C'est ce que je pensais.

C'est gênant, mais c'est récemment que j'ai vraiment compris la différence entre `log.Fatal()` et `panic()` en Golang. Je vais donc profiter de cette occasion pour clarifier cela.

## Différence entre log.Fatal() et panic()
`log.Fatal()` et `panic()` sont tous deux des fonctions qui arrêtent le programme. Examinons le fonctionnement à travers du code.

```go
package main

import (
    "log"
	"log/slog"
)

func RunWithFatal() {
	log.Fatal("Ceci est une erreur fatale")
}

func main() {
	RunWithFatal()

	slog.Info("Ceci ne s'exécute pas")
}
```

En exécutant le code ci-dessus, vous obtiendrez le résultat suivant.

```shell
2025/02/11 20:02:31 Ceci est une erreur fatale
```

Maintenant, examinons le code utilisant `panic()`.

```go
package main

import (
    "log/slog"
)

func RunWithPanic() {
    panic("Ceci est une erreur de panique")
}

func main() {
    RunWithPanic()

    slog.Info("Ceci ne s'exécute pas")
}
```

En exécutant le code ci-dessus, vous obtiendrez le résultat suivant.

```shell
panic: Ceci est une erreur de panique

goroutine 1 [running]:
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:8
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:12 +0x30
```

En observant ces codes, on peut dire que `log.Fatal()` affiche une erreur et arrête le programme, tandis que `panic()` affiche une erreur, arrête le programme de manière similaire mais affiche également une trace de la pile.

### Récupérer avec recover()
Lorsqu'on utilise `panic()`, le programme s'arrête, mais on peut le récupérer en utilisant `recover()`.

En réalité, il est préférable d'éviter les paniques, mais les développeurs étant humains, ils peuvent faire des erreurs. Par conséquent, dans les applications telles que les serveurs API, on crée souvent un middleware pour gérer les paniques et permettre ainsi d'éviter que le serveur ne se bloque de manière inattendue.

Pour mieux comprendre la différence, commençons par récupérer log.Fatal() à l'aide de recover().

```go
package main

import (
    "log"
    "log/slog"
)

func RunWithFatal() {
    log.Fatal("Ceci est une erreur fatale")
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            slog.Info("Récupération depuis", "erreur", r)
        }
    }()

    RunWithFatal()

    slog.Info("Ceci ne s'exécute pas")
}
```

En exécutant le code ci-dessus, vous obtiendrez le résultat suivant.

```shell
2025/02/11 20:07:49 Ceci est une erreur fatale
```

Il n'y a pas eu de récupération. Essayons maintenant de récupérer `panic()` avec recover().

```go
package main

import (
	"log/slog"
	"runtime/debug"
)

func RunWithPanic() {
	panic("Ceci est une erreur de panique")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("Récupération depuis", "erreur", r)
			debug.PrintStack()
		}
	}()

	RunWithPanic()
}
```

En exécutant le code ci-dessus, vous obtiendrez le résultat suivant.

```shell
2025/02/11 20:09:51 INFO Récupération depuis erreur="Ceci est une erreur de panique"
goroutine 1 [running]:
runtime/debug.Stack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:26 +0x64
runtime/debug.PrintStack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:18 +0x1c
main.main.func1()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:16 +0x8c
panic({0x1004d6560?, 0x1004f4190?})
	/opt/homebrew/opt/go/libexec/src/runtime/panic.go:785 +0x124
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:9
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:20 +0x4c
```

Bien que cela ne soit pas nécessaire, j'ai utilisé `debug.PrintStack()` pour afficher une trace de la pile similaire à celle de tout à l'heure.

## Utilisation
`log.Fatal()` appelle `os.Exit(1)` de manière interne.

Étant conçu pour arrêter immédiatement le programme avec un code d'erreur, il n'est pas récupérable avec `recover()`.

D'autre part, `panic()` peut être récupéré avec `recover()`.

Si vous pouvez éviter les erreurs, c'est bien, mais si vous vous trouvez dans une situation où vous souhaitez récupérer une erreur survenue, il est préférable d'utiliser `panic()`.

En général, il est recommandé d'utiliser `panic()` pour les fonctions de bibliothèque ou de certains packages. (Si à cause d'une bibliothèque le serveur plante mais qu'aucune récupération n'est possible, les conséquences pourraient être désastreuses).

Dans le cas de `log.Fatal()`, il est préférable que cela soit utilisé dans des parties comme `main()` où la gestion finale des erreurs est effectuée.

Par exemple, si une erreur se produit lors de l'initialisation des dépendances et que le module d'initialisation des dépendances retourne une erreur, c'est `main()` qui doit appeler `log.Fatal()`.

Si l'on regarde la structure, elle pourrait ressembler à ceci.
> Ceci est juste un exemple minimaliste. À titre indicatif seulement.
```go
package main

type Dependencies struct {
	DB *sql.DB
	redis *redis.Client
	...
}

func NewDependencies() (*Dependencies, error) {
    db, err := sql.Open("mysql", "utilisateur:motdepasse@/basededonnees")
    if err != nil {
        return nil, err
    }

    redis := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
	
    if err = redis.Ping().Err(); err != nil {
        return nil, err
    }

    return &Dependencies{
        DB: db,
        redis: redis,
    }, nil
    
}


func main() {
    deps, err := NewDependencies()
    if err != nil {
        log.Fatal(err)
    }
	
	// ...
}
```

## + `log.Panic()`
Dans le package `log`, il existe une fonction appelée `log.Panic()`.

Il s'agit d'une extension à `panic()` qui ajoute des fonctionnalités de logging. En exécutant le code ci-dessous :

```go
package main

import (
	"log"
	"log/slog"
	"runtime/debug"
)

func RunWithPanic() {
	log.Panic("Ceci est une erreur de panique")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("Récupération depuis", "erreur", r)
			debug.PrintStack()
		}
	}()

	RunWithPanic()
}

```

On obtiendrait la sortie suivante.

```shell
2025/02/11 20:23:17 Ceci est une erreur de panique
2025/02/11 20:23:17 INFO Récupération depuis erreur="Ceci est une erreur de panique"
goroutine 1 [running]:
runtime/debug.Stack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:26 +0x64
runtime/debug.PrintStack()
	/opt/homebrew/opt/go/libexec/src/runtime/debug/stack.go:18 +0x1c
main.main.func1()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:17 +0x8c
panic({0x100ad24e0?, 0x140000100a0?})
	/opt/homebrew/opt/go/libexec/src/runtime/panic.go:785 +0x124
log.Panic({0x1400010af20?, 0x0?, 0x68?})
	/opt/homebrew/opt/go/libexec/src/log/log.go:432 +0x60
main.RunWithPanic(...)
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:10
main.main()
	/Users/code_kirin/dev/personal/awesomeProject6/main.go:21 +0x60
```

Comparé à `panic()`, cette fonctionnalité de logging est la principale différence. Cette fonction `log.Panic()` déclenche une panique de la même manière que `panic()`, mais elle enregistre également des logs.

## Référence
- https://pkg.go.dev/log#Fatal
- Remarques de relecture fournies par le directeur général
