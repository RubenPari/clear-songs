# Controlla se il sistema operativo è Windows o Linux
if ($env:OS -match "Windows") {
    $src = ".\src.exe"
}
else {
    $src = "./src"
}

# Elimina i file "src" e "src.exe" nella directory corrente
Remove-Item -Path ".\src" -ErrorAction SilentlyContinue
Remove-Item -Path ".\src.exe" -ErrorAction SilentlyContinue

# Compila il codice sorgente con il comando "go build"
go build .

# Avvia l'eseguibile "src" o "src.exe"
& $src
