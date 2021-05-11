# Consumo dos endpoints através do pacote Dynamic Go Actions.

## Endpoints disponivéis:
1. **Utilidade** _(health check)_ :
    * Ping
    
1. **Ficheiros** _(file meta info)_ :
    * ApagarFicheiroMetaData
    * CriarFicheiroMetaData
    * BuscarMetaData

1. **Repositorios** _(repos)_ :
    * UpdateRepositorio
    * BuscarRepositorio
    * CriarRepositorio
    * DropRepositorio

<br>

## Utilidade:
### **Ping:**

    func PingServico(name string) (retorno map[string]interface{})

### Action a ser feita (ex.):

    action:
      funcs:
        "Ping":
          "<NOME>",
<hr>


## Ficheiros:
### **ApagarFicheiroMetaData:**

    func ApagarFicheiroMetaData(campos map[string]interface{}, token string) (retorno map[string]interface{})

### Action a ser feita (ex.):

    action:
      funcs:
        "ApagarFicheiroMetaData":
          {"nome":"teste", "autor":"teste", "repo_nome":"one_test", "path":["repo","one_test", "teste"]},
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA1NjkwMDAsImlzcyI6IlJvYm...",

<br>

### **CriarFicheiroMetaData:**

    func CriarFicheiroMetaData(ficheiroMetaData map[string]interface{}, token string) (retorno map[string]interface{})

### Action a ser feita (ex.):

    action:
      funcs:
        "CriarFicheiroMetaData":
          {"nome":"FILE0.txt", "autor":"tomas", "reponome":"teste03", "path":["repo","teste03", "FILE0.txt"]},
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA1NjkwMDAsImlzcyI6IlJvYm...",

<br>

### **BuscarMetaData:**

    func CriarFicheiroMetaData(ficheiroMetaData map[string]interface{}, token string) (retorno map[string]interface{})

### Action a ser feita (ex.):

    action:
      funcs:
        "CriarFicheiroMetaData":
          {"nome":"FILE02.txt","reponome":"teste01", "path":["repo", "teste01", "FILE02.txt"]},  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA1NjkwMDAsImlzcyI6IlJvYm...",
<hr>

## Repositorios:
### **UpdateRepositorio:**

    func UpdateRepositorio(campos map[string]interface{}, updateQuery map[string]interface{}, token string) (retorno map[string]interface{})

### Action a ser feita (ex.):

    action:
      funcs:
        "UpdateRepositorio":
          {"nome": "teste03"},
          {"tema": "Reparação de projetores"},
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA1NjkwMDAsImlzcyI6IlJvYm...",

<br>

### **BuscarRepositorio:**

    func BuscarRepositorio(campos map[string]interface{}, token string) (retorno map[string]interface{})

### Action a ser feita (ex.):

    action:
      funcs:
        "BuscarRepositorio":
          {"nome": "teste03"},
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA1NjkwMDAsImlzcyI6IlJvYm...",

<br>

### **DropRepositorio:**

    func DropRepositorio(campos map[string]interface{}, token string) (retorno map[string]interface{})

### Action a ser feita (ex.):

    action:
      funcs:
        "DropRepositorio":
          {"nome": "teste03"},
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA1NjkwMDAsImlzcyI6IlJvYm...",

<br>

### **CriarRepositorio:**

    func CriarRepositorio(repoInfo map[string]interface{}, token string) (retorno map[string]interface{})

### Action a ser feita (ex.):

    action:
      funcs:
        "CriarRepositorio":
          {"nome": "teste03", "tema": "testes para apagar meta data", "autor": "admin"},
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA1NjkwMDAsImlzcyI6IlJvYm...",
<br>          