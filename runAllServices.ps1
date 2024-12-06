# set Execution Policy to Unrestricted if needed
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy Unrestricted

Write-Host "---------------------------------------------------------"
Write-Host "Starting All Services..."
Write-Host "---------------------------------------------------------"

Set-Location .\server-side\services\

## BILLING SERVICE
Write-Host "Starting Billing Service..."
Set-Location .\billing\
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\billing-service.go"
cd ..

## COMMON SERVICE
Write-Host "Starting Common Service..."
Set-Location .\common\
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\common-service.go"
cd ..

## USER SERVICE
Write-Host "Starting User Service..."
Set-Location .\user\
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\user-service.go"
cd ..

## VEHICLE SERVICE
Write-Host "Starting Vehicle Service..."
Set-Location .\vehicle\
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\vehicle-service.go"
cd ..

cd ..


Write-Host "---------------------------------------------------------"
Write-Host "All services are up and running"
Write-Host "---------------------------------------------------------"

