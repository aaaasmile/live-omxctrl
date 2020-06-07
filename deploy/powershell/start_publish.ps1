# Start build process
Start-Process powershell.exe -ArgumentList '-NoExit', '-Command', "cd '$pwd'; & '.\publish.ps1';"