$outarm = "live-omxctrl.bin"
$confirmation = Read-Host "Do you want to build a new version target raspberry pi4 (y/n)?"
if ($confirmation -ne 'y') {
    Write-Host "Nothing to do. All the best."
    return
}

#create the exe
Write-Host "Build the $outarm in $pwd"
$env:GOOS = "linux"
$env:GOARCH="arm"
$env:GOARM="5"
go build -o $outarm

# Create the Zip package
Write-Host "Create a deploy package"
cd ./deploy
.\deploy.exe -target pi4

Write-Host "Done. Now the process continue with the WLC using the batch copy_app_to_pi4.sh"

cd ../../Deployed
Write-Host "current dir is now $pwd"
Write-Host "Bash script copy and update the service"
bash -c /mnt/d/scratch/go-lang/live-omxctrl/Deployed/copy_app_to_pi4.sh 

Write-Host "Done!"

