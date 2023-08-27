# Ottieni la cartella corrente
$currentDir = Get-Location

# Estrai l'ultimo segmento del percorso
$lastSegment = Split-Path $currentDir -Leaf

# Controlla se ti trovi nella cartella 'src'
if ($lastSegment -ne "src") {
    # Se non sei in 'src', prova a spostarti in 'src'
    if (Test-Path ".\src") {
        Set-Location -Path .\src
    } else {
        Write-Output "Cartella 'src' non trovata!"
        exit
    }
}

# Esegui la build del progetto
& go build .

#  Rimouvo il file eseguibile "src" o "src.exe"
Remove-Item -Path .\src -Force

# Controlla se il comando è stato eseguito correttamente
if ($LASTEXITCODE -ne 0) {
    Write-Output "Build fallita!"
    exit
}

# Esegui il comando specificato
& infisical run --env=dev -- go run .
