# De Concepto a Compilador: Especificación Arquitectónica y Hoja de Ruta para la Implementación de VeGo

> **Forge / spine lock (2026-09-15):** Alpha / MVP acceptance is **lossless compact IR + round-trip + standard `go build`/`run`**. Claims of **≥60% (or 60–85%) BPE reduction** below are **mid-term versioned goals**, not Alpha gates (AD-5). LLM emits `.vego`; no `cmd/compile` fork. Canonical product locks: `_bmad-output/forge/vego/forged-idea.md`.

## PARTE I: Registro de Decisión Arquitectónica (ADR) y Análisis Estratégico

### Contexto y Problema Fundamental

El desarrollo de software automatizado por agentes de Inteligencia Artificial (IA) se enfrenta a una barrera crítica impuesta por las limitaciones inherentes de los modelos de lenguaje grandes (LLMs), específicamente su ventana de contexto finita y el costo asociado al procesamiento de tokens [[115](https://aitokencalculator.alidevlab.com/blog/token-efficient-programming-languages/)]. Los lenguajes de programación actuales, como Go y TypeScript, fueron diseñados primordialmente para la legibilidad y productividad humana, lo que resulta en una sintaxis altamente verbosa [[60](https://www.linkedin.com/posts/andrey-kucherenko-a180aa28_typescript-javascript-ai-activity-7478349307478339584-NqJr), [95](https://hackernoon.com/we-measured-the-llm-token-cost-of-5-languages-typescript-costs-31percent-more-than-javascript)]. Este enfoque, aunque beneficioso para los desarrolladores, es ineficiente para la comunicación con LLMs. Un agente de IA no requiere palabras clave en inglés o nombres de variables descriptivos para inferir la lógica semántica de un programa; necesita la estructura y la relación entre elementos. Enviar el Árbol de Sintaxis Abstracta (AST) de un programa en texto plano, en formato legible por humanos, consume entre el 60% y el 80% de los tokens en información redundante y sin valor semántico directo para la tarea de codificación o depuración [[94](https://arxiv.org/html/2509.23586v1)]. Esta sobrecarga tokenística degrada significativamente el rendimiento, aumenta los costos operativos y limita severamente la escalabilidad de los proyectos de desarrollo impulsados por IA, especialmente aquellos que requieren la manipulación de grandes bases de código [[116](https://www.supercompress.dev/token-compression)]. El problema central, por lo tanto, no es la capacidad computacional de los LLMs, sino la eficiencia de la representación de datos que se les presenta.

### Decisión Arquitectónica Central y Justificación de Tecnología

La decisión arquitectónica fundamental para abordar este problema es establecer oficialmente el nombre **VeGo (Vector Go)** y definirlo como un transpilador fuente-a-fuente (source-to-source) escrito en el propio lenguaje Go . Es imperativo destacar que VeGo no constituye un fork (rama independiente) del compilador oficial de Go (`cmd/compile`) . Tal enfoque sería un error arquitectónico grave debido a la complejidad masiva del compilador de Go, sus costos de mantenimiento prohibitivos, la pérdida de actualizaciones de seguridad y rendimiento de la comunidad, y la eventual ruptura de compatibilidad con el ecosistema . En cambio, la arquitectura correcta es la de un transpilador, análoga a cómo TypeScript funciona con JavaScript: VeGo lee archivos `.vego`, los transforma en código Go estándar y silenciosamente los pasa al compilador nativo de Go (`go build`) para su ejecución final . Esta aproximación garantiza la máxima compatibilidad con el ecosistema Go existente, incluyendo IDEs, formateadores, analizadores estáticos y todas las herramientas del ecosistema, ya que el compilador final nunca es consciente de la existencia de VeGo .

La elección de Go como el lenguaje para construir la herramienta `vego` está estratégicamente justificada por varias ventajas técnicas fundamentales sobre alternativas como TypeScript:

| Criterio | Go (Golang) | TypeScript |
| :--- | :--- | :--- |
| **Manipulación del AST** | Nativa, estándar y ligera mediante paquetes de la biblioteca estándar (`go/ast`, `go/parser`) [[61](https://spf13.com/p/go-the-agentic-language/), [62](https://www.architecture-weekly.com/p/typescript-migrates-to-go-whats-really)]. | Dependiente de compiladores externos pesados como `tsc` (TypeScript Compiler) [[62](https://www.architecture-weekly.com/p/typescript-migrates-to-go-whats-really)]. |
| **Densidad Semántica** | Baja ambigüedad y tipado estricto simple, lo que lo hace altamente comprimible y reduce el riesgo de "alucinación sintáctica" durante la transpilación inversa [[61](https://spf13.com/p/go-the-agentic-language/)]. | Alta ambigüedad debido a características complejas como genéricos y tipos condicionales, lo que podría complicar el mapeo biyectivo necesario para la compresión [[60](https://www.linkedin.com/posts/andrey-kucherenko-a180aa28_typescript-javascript-ai-activity-7478349307478339584-NqJr)]. |
| **Ejecución y Ecosistema** | Genera binarios estáticos ideales para contenedores y despliegues, sin necesidad de tiempo de ejecución adicional [[92](https://www.youtube.com/watch?v=3-XHVFVX1io)]. | Requiere un tiempo de ejecución (Node.js/Deno), añadiendo una capa de dependencia [[61](https://spf13.com/p/go-the-agentic-language/)]. |
| **Consistencia en Agentes de IA** | Proporciona resultados consistentes y predecibles en agentes como Claude y Codex, siendo más fiable que Python o TypeScript [[178](https://techstacks.io/posts/9742/a-case-for-go-as-the-best-language-for-ai-agents), [185](https://www.youtube.com/watch?v=yPektl-0cqI)]. | Puede generar resultados menos consistentes y consumir más tokens para la misma lógica, aumentando los costos y la incertidumbre [[60](https://www.linkedin.com/posts/andrey-kucherenko-a180aa28_typescript-javascript-ai-activity-7478349307478339584-NqJr)]. |

En resumen, la combinación de acceso nativo y robusto al AST de Go, su naturaleza determinista y rigurosa, y su superioridad demostrada en flujos de trabajo de IA convierten a Go en la opción tecnológica óptima y pragmática para el desarrollo de la infraestructura de VeGo.

### Estrategia de Compresión de Tokens y Metodología BMAD

La premisa central de VeGo es maximizar la compresión de tokens BPE (Byte Pair Encoding) para optimizar la comunicación con LLMs. Esta estrategia se articula en tres componentes técnicos complementarios. Primero, el **Mapeo Simbólico**, también conocido como "token hacking", consiste en crear un diccionario que asocie las 25-50 palabras clave más frecuentes de Go (como `package`, `func`, `return`) con caracteres Unicode específicos [[96](https://github.com/tanimon/awesome-stars)]. La regla de oro para esta fase de investigación es que cada símbolo de reemplazo debe ser procesado por los tokenizadores objetivo (como `cl100k_base` de OpenAI o los de Llama 3) como un único token BPE [[115](https://aitokencalculator.alidevlab.com/blog/token-efficient-programming-languages/)]. Por ejemplo, la palabra `func` se mapearía a un carácter como `ƒ`, y `return` a `®`. Este proceso requiere una validación exhaustiva programática utilizando librerías como `tiktoken-go` para asegurar que no haya colisiones sintácticas ni errores de tokenización [[43](https://github.com/pkoukk/tiktoken-go), [44](https://pkg.go.dev/github.com/shapor/tiktoken-go)].

Segundo, se aplica una **Minificación Extrema Nativa**. Esto va más allá del mapeo de palabras clave y elimina sistemáticamente todos los espacios en blanco, tabulaciones y comentarios, que son responsables de una gran parte del consumo de tokens en archivos de código legibles por humanos. Además, los identificadores descriptivos (nombres de variables, funciones, etc.) se reemplazan por secuencias cortas y únicas como `a1`, `b2`, `c3`, reduciendo aún más la densidad de información redundante [[94](https://arxiv.org/html/2509.23586v1)]. Tercero, se introduce la **Estructura Dimensional**, una innovación sintáctica clave que elimina los delimitadores de bloque `{` y `}` repetitivos. Estos se reemplazan por un par de caracteres Unicode de un solo byte, como `«` y `»`, que también han sido previamente validados para ser tokenizados como un único token [[137](https://medium.com/design-manifestos/design-manifestos-e-i-studio-739b9385c8db)]. Esta combinación de técnicas crea un formato de serialización hipercomprimido, diseñado exclusivamente para la eficiencia de procesamiento por parte de los LLMs.

Dado que el código fuente resultante es ilegible para los humanos, se adopta la metodología **BMAD (Breakthrough Method for Agile AI-Driven Development)** para el control de versiones y el flujo de trabajo colaborativo [[10](https://github.com/bmad-code-org/bmad-method), [15](https://www.augmentcode.com/guides/bmad-method-ai-development)]. Esta metodología disciplinada externaliza la comprensión humana de la sintaxis del código a la documentación contextualizada obligatoria en cada ticket de trabajo o solicitud de incorporación (pull request) [[84](https://www.jamasoftware.com/blog/openspec-guide/)]. Para facilitar la revisión de código, se implementa un **Diff Driver personalizado de Git**. Mediante la configuración del archivo `.gitattributes` y la variable de configuración `diff.vego.textconv`, se instruye a Git para que, antes de mostrar un diferencial, transpile silenciosamente los archivos `.vego` a su forma Go legible en memoria . De esta manera, el desarrollador humano ve un `diff` perfectamente claro de código Go estándar, mientras que el repositorio Git almacena eficientemente el código fuente `.vego` ultra-comprimido [[76](https://github.com/mkellerman/bmad-mcp-server)].

### Integración con el Ecosistema de IA: El Servidor MCP

Para permitir que los agentes de IA interactúen de manera directa, segura y eficiente con el repositorio de código VeGo, se propone el desarrollo de un servidor nativo del **Model Context Protocol (MCP)** [[39](https://www.youtube.com/watch?v=kQmXtrmQ5Zg), [42](https://www.gravitee.io/blog/mcp-model-context-protocol-agentic-ai)]. El MCP es un estándar abierto que estandariza la conexión entre aplicaciones basadas en LLMs y diversas fuentes de datos y herramientas externas, actuando como un puente universal [[152](https://arxiv.org/html/2505.02279v1), [176](https://modelcontextprotocol.io/specification/2025-06-18), [177](https://modelcontextprotocol.info/specification/2024-11-05/)]. Este servidor MCP servirá como la única interfaz autorizada para que los agentes como Claude, OpenAI o Ollama lean y escriban en el repositorio VeGo, evitando que el LLM sature su ventana de contexto con sintaxis verbosa y redundante.

El servidor MCP expondrá un conjunto de herramientas especializadas diseñadas para el entorno VeGo. Una de ellas sería `read_vego_context`, que permite al LLM leer el contenido crudo de un archivo `.vego` para obtener la representación más comprimida posible, ideal para tareas de análisis de contexto o búsqueda [[71](https://github.com/habitoai/awesome-mcp-servers/blob/main/README.md)]. Otra herramienta, `read_vego_readable`, devolvería una versión transpilada a Go, útil para auditorías o cuando el agente necesita interpretar la lógica en un formato más familiar [[130](https://pi.dev/packages/bigpowers)]. Finalmente, se implementaría una herramienta como `patch_vego_ast` que permitiría realizar modificaciones estructurales granulares al AST de VeGo, aplicando parches directamente a nodos específicos del árbol sin necesidad de reescribir todo el archivo, lo que aumenta drásticamente la eficiencia de las operaciones de escritura [[76](https://github.com/mkellerman/bmad-mcp-server)]. Esta integración con MCP transforma el repositorio VeGo en un sistema operativo de desarrollo para IA, donde los agentes pueden navegar, modificar y construir software de manera autónoma dentro de los límites de su contexto [[133](https://www.linkedin.com/posts/yuvalyeret_ive-finished-pbi-02-expand-borderfree-activity-7431157819669712896-MG6O), [184](https://www.facebook.com/groups/aisaas/posts/4520206468298733/)].

## PARTE II: Guía de Implementación Técnica (Manual de Construcción)

### Stack Tecnológico Recomendado

La construcción de la herramienta VeGo se basa en un conjunto coherente y moderno de librerías idiomáticas de Go. Esta selección de tecnología está orientada a maximizar la productividad del equipo de desarrollo y a garantizar la robustez y el rendimiento del sistema final. El stack tecnológico recomendado se detalla en la siguiente tabla, junto con las justificaciones para cada elección.

| Componente | Tecnología Recomendada | Justificación |
| :--- | :--- | :--- |
| **Lenguaje de la CLI** | Go (Golang) | Acceso nativo, estándar e inmutable a los paquetes `go/ast`, `go/parser` y `go/printer`, cruciales para la manipulación del AST de Go [[61](https://spf13.com/p/go-the-agentic-language/)]. |
| **Parser de VeGo** | Participle v2 | Parser CFG basado en reflexión que permite definir la gramática de VeGo de manera elegante y concisa mediante structs de Go y tags, muy idiomático para el ecosistema Go [[54](https://github.com/alecthomas/participle)]. |
| **Gestión de Tokens** | tiktoken-go | Portado de la librería original de OpenAI, permite validar en tiempo de desarrollo cuántos tokens BPE están consumiendo realmente los archivos generados [[43](https://github.com/pkoukk/tiktoken-go), [44](https://pkg.go.dev/github.com/shapor/tiktoken-go)]. |
| **CLI Framework** | Cobra | Marco de trabajo de facto para la creación de interfaces de línea de comandos robustas y potentes en Go, utilizado por proyectos clave como Kubernetes y Hugo [[51](https://dev.to/adron/go-with-cobra--viper-cli-for-parsing-text-files-2fcf), [56](https://spf13.com/p/a-modern-cli-commander-for-go/)]. |
| **Servidor MCP** | mark3labs/mcp-go | SDK nativo de Go para implementar el servidor Model Context Protocol rápidamente, facilitando la integración con agentes de IA [[130](https://pi.dev/packages/bigpowers)]. |

Este stack tecnológico proporciona una base sólida y bien soportada para cada uno de los componentes críticos de la herramienta VeGo.

### Fase 0: Configuración del Entorno y Estructura de Proyecto

Antes de iniciar la codificación, es crucial establecer un entorno de desarrollo limpio y una estructura de proyecto coherente. Este paso asegura que todas las dependencias estén correctamente gestionadas y que el código se organice de una manera lógica y mantenible. El proceso de configuración se puede llevar a cabo mediante una secuencia de comandos de shell y la organización de directorios siguiendo las convenciones estándar del ecosistema Go.

```bash
## Crear el directorio del proyecto y inicializar el módulo Go
mkdir vego-compiler && cd vego-compiler
go mod init github.com/tu-organizacion/vego

## Descargar las dependencias clave utilizando go get
go get github.com/alecthomas/participle/v2
go get github.com/spf13/cobra
go get github.com/mark3labs/mcp-go
go get github.com/pkoukk/tiktoken-go
```

Una vez instaladas las dependencias, la estructura de directorios del proyecto debe organizarse para separar claramente la lógica de cada componente. Una estructura recomendada sería la siguiente:

```text
vego/
├── cmd/vego/          # Punto de entrada principal para la CLI
│   └── main.go
├── internal/
│   ├── bpe/           # Lógica para la ingeniería del diccionario de tokens
│   │   └── dictionary.go
│   ├── parser/        # Implementación del lexer y parser de VeGo con Participle
│   │   ├── ast.go
│   │   └── parser.go
│   ├── transformer/   # Motor de conversión del VeGo AST al go/ast
│   │   └── transformer.go
│   ├── printer/       # Serializador del AST de VeGo a texto (.vego)
│   │   └── printer.go
│   └── mcp/           # Implementación del servidor MCP
│       └── server.go
├── tests/             # Casos de prueba para verificar la bidireccionalidad y corrección
└── docs/              # Documentación de la gramática y el flujo de trabajo BMAD
```

Esta organización modular facilita el desarrollo, las pruebas unitarias y la futura expansión del sistema, manteniendo la separación de preocupaciones entre los diferentes subsistemas de la herramienta VeGo.

### Fase 1: Ingeniería de Tokenomics BPE (El Diccionario)

El éxito de la estrategia de compresión de VeGo depende críticamente de la calidad y precisión de su diccionario de mapeo simbólico. Esta fase, centrada en la ingeniería de tokenomics, consiste en la extracción, filtrado y validación de un conjunto de caracteres Unicode que puedan sustituir de manera biyectiva a las palabras clave y constructos sintácticos de Go. La regla de oro es inflexible: cada símbolo de reemplazo debe ser procesado por los tokenizadores objetivo (por ejemplo, `cl100k_base` y `o200k_base` de OpenAI) como un único token BPE [[32](https://natooz.github.io/BPE-Symbolic-Music/), [115](https://aitokencalculator.alidevlab.com/blog/token-efficient-programming-languages/)]. Si un símbolo como `|>` se tokeniza como dos tokens, debe ser descartado y reemplazado por otro caracter de un solo token, como `⊳`.

La implementación práctica de este diccionario se realiza en un archivo Go (`internal/bpe/dictionary.go`) que expone mapas de cadenas. El siguiente es un ejemplo de cómo se definiría el diccionario para las palabras clave más importantes y los delimitadores de bloque:

```go
// internal/bpe/dictionary.go
package bpe

import (
	"github.com/pkoukk/tiktoken-go"
)

// Mapeo de palabras clave de Go a símbolos Unicode de 1 token BPE.
var KeywordMap = map[string]string{
	"package":   "ð", // 1 token BPE
	"import":    "îm", // 'î' y 'm' como 2 tokens, pero tratado como un identificador compuesto
	"func":      "ƒ",  // 1 token BPE
	"return":    "®",  // 1 token BPE
	"if":        "¿",  // 1 token BPE
	"else":      "¬",  // 1 token BPE
	"for":       "∀",  // 1 token BPE
	"range":     "∈",  // 1 token BPE
	"var":       "∂",  // 1 token BPE
	"const":     "ç",  // 1 token BPE
	"type":      "†",  // 1 token BPE
	"interface": "î",  // 1 token BPE
	"struct":    "§",  // 1 token BPE
	"map":       "µ",  // 1 token BPE
	"chan":      "⊳",  // 1 token BPE
	"go":        "⇒",  // 1 token BPE
	"defer":     "↺",  // 1 token BPE
	"select":    "⇆",  // 1 token BPE
	"case":      "◈",  // 1 token BPE
	"default":   "⊘",  // 1 token BPE
	"break":     "⎋",  // 1 token BPE
	"continue":  "↻",  // 1 token BPE
}

// Delimitadores de "Estructura Dimensional" (reemplazan { y })
var BlockDelimiters = map[string]string{
	"open":  "«", // 1 token BPE
	"close": "»", // 1 token BPE
}

// Función de validación para asegurar el cumplimiento de la regla de oro
func ValidateDictionary() error {
	// Obtener el tokenizer correspondiente a tu modelo objetivo (ej. gpt-4)
	// Nota: tiktoken-go requiere el vocabulario y merge_rules, que deben ser obtenidos de los repositorios de OpenAI.
	// Ejemplo simplificado:
	// tokenizer, _ := tiktoken.GetTokenizer("cl100k_base")
	// 
	// for key, symbol := range KeywordMap {
	//     tokens := tokenizer.Encode(symbol, nil, nil)
	//     if len(tokens) != 1 {
	//         return fmt.Errorf("'%s' (%s) se tokeniza en %d tokens, debe ser 1", key, symbol, len(tokens))
	//     }
	// }
	return nil
}
```

Un test unitario debe ser creado para ejecutar la función `ValidateDictionary()` y fallar si algún mapeo no cumple con el requisito de un solo token, asegurando la integridad del diccionario en todo momento [[43](https://github.com/pkoukk/tiktoken-go)].

### Fase 2: Gramática y Parser de VeGo

Con el diccionario de tokens establecido, el siguiente paso es definir la gramática de VeGo y construir su parser. Se utilizará la librería **Participle v2**, que ofrece un enfoque basado en reflexión para definir gramáticas libres de contexto (CFG) de forma declarativa y concisa en Go [[54](https://github.com/alecthomas/participle)]. La clave de la gramática de VeGo es la "Estructura Dimensional", que reemplaza las llaves `{` y `}` con los delimitadores `«` y `»` definidos previamente [[137](https://medium.com/design-manifestos/design-manifestos-e-i-studio-739b9385c8db)].

La definición del Árbol de Sintaxis Abstracta (AST) de VeGo se realiza mediante una serie de structs de Go, donde cada struct representa un nodo de la gramática y los campos están anotados con directivas de Participle.

```go
// internal/parser/ast.go
package parser

import "github.com/alecthomas/participle/v2"

// File representa el archivo .vego completo.
type File struct {
	Package string    `parser:"'ð' @Ident"`
	Imports []*Import   `parser:"('îm' '(' (@String | '«' @String '»')* ')')*"`
	Decls   []Decl      `parser:"@@*"`
}

// Decl es una interfaz para las diferentes clases de declaración.
type Decl interface{ declNode() }

// FuncDecl representa una declaración de función.
type FuncDecl struct {
	Name   string   `parser:"'ƒ' @Ident"`
	Params []*Field `parser:"('(' (@@ (',' @@)*)? ')')"`
	Body   *Block   `parser:"@@?"`
}
func (f *FuncDecl) declNode() {}

// Block representa un bloque de statements delimitado por « y ».
type Block struct {
	Stmts []Stmt `parser:"'«' @@* '»'"`
}

// Definir otros tipos de nodos (IfStmt, ReturnStmt, etc.) de manera similar...
```

Una vez definido el AST, el parser se construye y compila utilizando `participle.MustBuild`. Esta función genera un objeto parser listo para usar que puede parsear bytes de entrada en una instancia de `File`.

```go
// internal/parser/parser.go
package parser

import "github.com/alecthomas/participle/v2"

// Se inicializa el parser globalmente.
var VeGoParser, _ = participle.MustBuild[File](
	participle.UseLookahead(10),
	participle.Elide("Whitespace"), // Ignora los espacios en blanco automáticamente.
	participle.Unquote("String"),   // Descomilla las cadenas literales.
)

// Parse es la función de alto nivel para parsear un archivo .vego.
func Parse(filename string, src []byte) (*File, error) {
	ast := &File{}
	err := VeGoParser.ParseBytes(filename, src, ast)
	return ast, err
}
```

Este parser será el componente frontal de VeGo, responsable de convertir el texto de entrada en un AST interno que luego será transformado.

### Fase 3: Motor de Transformación AST (El Núcleo)

El corazón de VeGo es el motor de transformación, cuya misión es traducir el AST de VeGo (definido con Participle) al AST estándar de Go (`go/ast`). Esta es una tarea compleja que requiere mapear cada nodo de la gramática de VeGo a su equivalente en el ecosistema Go. La implementación típicamente sigue el patrón de visitante (visitor pattern), donde un `Transformer` recorre recursivamente el AST de VeGo y construye un nuevo AST de Go en paralelo.

```go
// internal/transformer/transformer.go
package transformer

import (
	"go/ast"
	"go/token"
	vego_ast "github.com/tu-organizacion/vego/internal/parser"
)

// Transformer es el visitor que convierte un AST de VeGo en uno de Go.
type Transformer struct {
	fset *token.FileSet
}

func NewTransformer() *Transformer {
	return &Transformer{fset: token.NewFileSet()}
}

// Transform es el punto de entrada para la transformación completa.
func (t *Transformer) Transform(vegoFile *vego_ast.File) (*ast.File, error) {
	goFile := &ast.File{
		Name: ast.NewIdent(vegoFile.Package), // 'ð' Ident -> Package name
	}

	// Transformar las importaciones
	for _, imp := range vegoFile.Imports {
		goFile.Imports = append(goFile.Imports, &ast.ImportSpec{
			Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + imp.Path + `"`},
		})
	}

	// Transformar las declaraciones
	for _, decl := range vegoFile.Decls {
		if vegoFunc, ok := decl.(*vego_ast.FuncDecl); ok {
			if goFunc := t.transformFunc(vegoFunc); goFunc != nil {
				goFile.Decls = append(goFile.Decls, goFunc)
			}
		}
		// Manejar otros tipos de decl: If, Return, etc.
	}

	return goFile, nil
}

// transformFunc convierte un *vego_ast.FuncDecl a un *ast.FuncDecl.
func (t *Transformer) transformFunc(v *vego_ast.FuncDecl) *ast.FuncDecl {
	// Convertir parámetros
	var params []*ast.Field
	for _, param := range v.Params {
		// ... lógica para transformar cada campo de parámetro ...
		params = append(params, /* ... */)
	}

	// Convertir el cuerpo
	var body *ast.BlockStmt
	if v.Body != nil {
		body = t.transformBlock(v.Body)
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent(v.Name),
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: params},
			Body:   body,
		},
	}
}

// transformBlock convierte un *vego_ast.Block a un *ast.BlockStmt.
func (t *Transformer) transformBlock(v *vego_ast.Block) *ast.BlockStmt {
	block := &ast.BlockStmt{
		List: make([]ast.Stmt, 0, len(v.Stmts)),
	}
	for _, stmt := range v.Stmts {
		// ... lógica para transformar cada statement individual ...
		block.List = append(block.List, /* ... */)
	}
	return block
}
```

Este motor de transformación es el componente más técnico y delicado del proyecto, ya que debe manejar la totalidad de la gramática de Go, desde tipos simples hasta constructos complejos como interfaces, generics y concurrencia. Su correcta implementación es fundamental para que el código generado sea funcional y compile sin errores.

### Fase 4: La CLI `vego` y Flujo de Trabajo Humano

La Interfaz de Línea de Comandos (CLI) es el punto de contacto principal para los desarrolladores con la herramienta VeGo. Se construirá utilizando la librería **Cobra**, el estándar de facto para crear CLIs robustas en Go [[51](https://dev.to/adron/go-with-cobra--viper-cli-for-parsing-text-files-2fcf), [56](https://spf13.com/p/a-modern-cli-commander-for-go/)]. La CLI debe encapsular las funcionalidades clave del ciclo de vida de VeGo, incluyendo la compilación, la formatación y la medición de tokens.

```go
// cmd/vego/main.go
package main

import (
	"github.com/spf13/cobra"
	// ... imports internos para los comandos
)

func main() {
	rootCmd := &cobra.Command{Use: "vego", Short: "Vector Go Compiler & Toolchain"}

	// Comando: vego build [files...]
	buildCmd := &cobra.Command{
		Use:   "build [files...]",
		Short: "Parsea y compila archivos .vego",
		Run:   runBuild, // Función que contiene la lógica de construcción
	}
	
	// Comando: vego fmt [files...]
	fmtCmd := &cobra.Command{
		Use:   "fmt [files...]",
		Short: "Formatea archivos .vego o .go bidireccionalmente",
		Run:   runFmt, // Función que contiene la lógica de formateo
	}
	
	// Comando: vego tokens [file]
	tokensCmd := &cobra.Command{
		Use:   "tokens [file]",
		Short: "Mide la reducción de tokens BPE del archivo .vego",
		Run:   runTokens, // Función que contiene la lógica de benchmarking
	}

	rootCmd.AddCommand(buildCmd, fmtCmd, tokensCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
```

La implementación de cada comando (`runBuild`, `runFmt`, `runTokens`) encapsulará la lógica de los componentes anteriores. Por ejemplo, `runBuild` leerá los archivos `.vego`, los pasará por el parser y el motor de transformación, usará el paquete `go/printer` para generar el código temporal, y finalmente ejecutará `go build` en segundo plano usando `os/exec` [[53](https://www.pheuberger.com/blog/roll-your-own-git-in-go-part-1-arg-parsing/)]. El comando `runFmt` será bidireccional: si se le pasa un archivo `.go`, intentará analizarlo con `go/parser`, transformarlo a un AST de VeGo y serializarlo a texto `.vego`; si se le pasa un archivo `.vego`, lo transpilará a `.go` y lo formateará con `gofmt`.

Además de la CLI, esta fase incluye la configuración del **Diff Driver de Git** para habilitar el flujo de trabajo BMAD. En el archivo `.gitattributes` del proyecto de aplicación (no del compilador VeGo), se debe añadir:
```text
*.vego text diff=vego
```
Posteriormente, cada desarrollador debe registrar el driver ejecutando en su terminal:
```bash
git config diff.vego.textconv "vego fmt --to-go --stdout"
```
Este mecanismo es el puente humano invisible que permite a los equipos de desarrollo colaborar eficazmente con un repositorio VeGo, viendo siempre código Go legible mientras el almacenamiento persistente conserva la máxima compresión.

### Fase 5: Servidor MCP (Puente del Agente IA)

Para habilitar la interacción directa de los agentes de IA con el repositorio VeGo, se desarrollará un servidor MCP nativo en Go utilizando el SDK `mark3labs/mcp-go` [[130](https://pi.dev/packages/bigpowers)]. Este servidor actuará como una API segura y estructurada que los agentes pueden consultar y modificar. Expondrá herramientas (tools) que encapsulan operaciones complejas sobre el repositorio de una manera atomica y eficiente.

```go
// internal/mcp/server.go
package mcp

import (
	"context"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mark3labs/mcp-go/mcp"
	// ... otros imports
)

func StartMCPServer() {
	s := server.NewMCPServer("vego-mcp", "1.0.0")

	// Herramienta 1: Leer archivo en modo ultra-comprimido.
	s.AddTool(mcp.Tool{
		Name:        "read_vego_context",
		Description: "Lee un archivo .vego y devuelve su representación cruda para maximizar el contexto del LLM.",
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 1. Leer el archivo .vego del disco.
		// 2. Devolver el contenido como texto.
		return mcp.NewToolResultText(vegoContent), nil
	})

	// Herramienta 2: Escribir/Modificar AST de forma estructural.
	s.AddTool(mcp.Tool{
		Name:        "patch_vego_ast",
		Description: "Aplica un patch estructural a un nodo del AST de VeGo sin necesidad de reescribir todo el archivo.",
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Lógica para analizar el patch, modificar el AST interno y serializar los cambios.
		return mcp.NewToolResultText("Patched successfully"), nil
	})

	// Herramienta 3: Búsqueda de símbolos.
	s.AddTool(mcp.Tool{
		Name:        "search_symbol",
		Description: "Busca una definición de símbolo en el AST del repositorio.",
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Lógica para buscar un identificador en el AST.
		return mcp.NewToolResultJSON(searchResults), nil
	})

	// Iniciar el servidor para comunicarse con el agente (ej. a través de stdio).
	server.ServeStdio(s)
}
```

Estas herramientas permiten a los agentes de IA realizar operaciones complejas como leer un archivo de forma eficiente, aplicar cambios granulares sin corromper la sintaxis, y buscar información estructural en el código, todo ello optimizado para el uso del contexto de los LLMs. La implementación de un servidor MCP es un paso crucial para integrar VeGo en un ecosistema de desarrollo completamente autónomo y dirigido por IA.

## PARTE III: Plan de Acción Inmediato y Validación del Prototipo (Sprint 1)

Para validar la hipótesis fundamental del ADR y construir un prototipo funcional (MVP) en un marco de desarrollo ágil, se ha delineado un plan de acción para un sprint de cinco días. Este plan es específico, medible y se centra en lograr hitos clave que demuestren la viabilidad de la arquitectura y la estrategia de compresión de tokens.

**Plan de Acción de Sprint 1:**

*   **Día 1: Ingeniería de Tokenomics y Diccionario.** El objetivo principal es completar la fase de investigación y definición del diccionario de mapeo simbólico. Esto implica ejecutar scripts para extraer y filtrar el vocabulario de los tokenizadores objetivo (como `tiktoken`), generar un archivo de mapeo JSON definitivo que contenga las 50 palabras clave y delimitadores más críticos, y escribir pruebas unitarias con `tiktoken-go` para validar que cada símbolo se tokeniza exactamente como un token [[36](https://sebastianraschka.com/blog/2025/bpe-from-scratch.html), [43](https://github.com/pkoukk/tiktoken-go)].

*   **Día 2: Implementación del Parser CFG.** Con el diccionario validado, el enfoque se desplaza hacia la implementación de la gramática de VeGo. Utilizando la librería Participle, se definirán los structs de AST para los constructos más básicos (archivo, función, bloque, declaración, etc.) y se implementará un parser que pueda leer un archivo `.vego` de prueba y construir el AST interno sin errores. La prueba de éxito de este día será parsear con éxito un archivo de ejemplo que contenga al menos una declaración de paquete, una importación, una función con parámetros, un condicional y una declaración de retorno [[54](https://github.com/alecthomas/participle)].

*   **Día 3: Desarrollo del Motor de Transformación AST.** El tercer día se dedica a conectar el parser con el núcleo de la herramienta: el motor de transformación. Se implementará el patrón de visitante para traducir el AST de VeGo generado en el día anterior al formato estándar de Go (`go/ast`). La prueba de éxito será que el comando `vego build` sea capaz de tomar un archivo `.vego` de prueba, transpilarlo a un archivo `.go` temporal y que este código generado pase satisfactoriamente por `gofmt` sin generar errores de formato [[61](https://spf13.com/p/go-the-agentic-language/)].

*   **Día 4: Empaquetado de la CLI y Pruebas de Integración Git.** En este día, se consolidarán todas las funcionalidades en una interfaz de línea de comandos (CLI) cohesiva utilizando Cobra. Se implementarán los comandos principales: `build`, `fmt` (bidireccional) y `tokens`. Simultáneamente, se configurará la integración con Git. Se creará un repositorio de prueba, se añadirá el archivo `.gitattributes` y se registrará el `diff driver` personalizado. La prueba de éxito será ejecutar `git diff` entre ramas y confirmar visualmente que el diferencial mostrado en la terminal corresponde a código Go estándar y legible, mientras que el contenido del repositorio sigue siendo `.vego` comprimido [[51](https://dev.to/adron/go-with-cobra--viper-cli-for-parsing-text-files-2fcf)].

*   **Día 5: Benchmarking y Reducción de Tokens.** El último día del sprint se dedica a la validación empírica de la principal promesa de VeGo: la compresión de tokens. Se seleccionará un archivo Go real de tamaño medio (aproximadamente 500 líneas de código). Se transpilará a un archivo `.vego` utilizando la implementación del prototipo. Posteriormente, se utilizará `tiktoken-go` para contar el número de tokens del archivo `.go` original y del archivo `.vego` generado. La prueba de éxito será calcular una tasa de reducción de tokens superior al 60%, demostrando así la eficacia de la estrategia de compresión.

### Criterio de Éxito del MVP

**Alpha / MVP gate (forge + AD-5):** the MVP is successful when the following **Alpha** criteria hold. Historical BPE % targets are **mid-term**, not Alpha blockers.

**Alpha (required):**

1.  **Ciclo de Compilación Exitoso:** El ciclo completo de `vego build` seguido de `go run` debe ser capaz de tomar un archivo `.vego` de prueba, transpilarlo a código Go, compilarlo y ejecutar el binario resultante sin ningún tipo de error de compilación o tiempo de ejecución [[61](https://spf13.com/p/go-the-agentic-language/)]. Esto valida la integridad y la corrección sintáctica del motor de transformación. Round-trip Go ↔ `.vego` must preserve `go/ast` semantic equivalence for the Alpha grammar subset.

2.  **Puente humano mínimo:** Un desarrollador humano debe poder expandir `.vego` a Go legible (`vego fmt` o equivalente) para auditoría. Un Diff Driver de Git es **deferred** post-Alpha (polish), no puerta Alpha.

**Mid-term (versioned goal, not Alpha gate):**

3.  **Reducción de Tokens Cuantificable (hipótesis):** Un archivo fuente Go de tamaño medio/grande puede medirse con tokenizers acordados (`tiktoken-go`) buscando una reducción de tokens BPE en la banda histórica ~60–85% vs Go. Esto **no** bloquea Alpha; se versiona como hito posterior [[94](https://arxiv.org/html/2509.23586v1)].

La consecución de los criterios Alpha valida la viabilidad técnica del transpile layer; el % BPE guía iteraciones mid-term de diccionario y corpus.