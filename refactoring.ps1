# Move to the src directory
Set-Location .\src

# Run gofmt on all Go files in the project
Get-ChildItem -Recurse -Filter "*.go" | ForEach-Object {
    gofmt -s -w $_.FullName
}

# Run goimports on all Go files in the project
Get-ChildItem -Recurse -Filter "*.go" | ForEach-Object {
    goimports -w $_.FullName
}
