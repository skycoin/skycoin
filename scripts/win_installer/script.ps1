
$ErrorActionPreference = "Stop"
$version = $args[0]
$arch = $args[1]

function CleanStage
{
    if (Test-Path ".\archive.zip") { Remove-Item ".\archive.zip" -Recurse -Force}
    if (Test-Path ".\archive") { Remove-Item ".\archive" -Recurse -Force}
    if (Test-Path ".\scripts\win_installer\build") { Remove-Item ".\scripts\win_installer\build" -Recurse -Force}
    if (Test-Path ".\scripts\win_installer\wix.zip") { Remove-Item ".\scripts\win_installer\wix.zip" -Recurse -Force}
    if (Test-Path ".\scripts\win_installer\wix") { Remove-Item ".\scripts\win_installer\wix" -Recurse -Force}
    if (Test-Path ".\scripts\win_installer\Product.wixobj") { Remove-Item ".\scripts\win_installer\Product.wixobj" -Recurse -Force}
    if (Test-Path ".\scripts\win_installer\skycoin.msi") { Remove-Item ".\scripts\win_installer\skycoin.msi" -Recurse -Force}
}

function InstallWix
{
    Set-Location .\scripts\win_installer
    Invoke-WebRequest "https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip" -OutFile wix.zip
    Expand-Archive wix.zip
    Set-Location ../../
}

function BuildInstaller()
{
    if ($arch -eq "386") {
        $arch_title="386  "
        $wix_arch="x86"
    } else {
        $arch_title="amd64"
        $wix_arch="x64"
    }

    Write-Output "#                                                        #"
    Write-Output "#    => Create Installer for $arch_title                       #"
    Write-Output "#       0. Preparing Stage...                            #"
    CleanStage
    Write-Output "#       1. Installing Wix...                             #"
    InstallWix
    Write-Output "#       2. Fetching Archive from GitHub...               #"
    
    $fileName = "skycoin-$version-windows-$arch"
    $msiName = "skycoin-installer-$version-windows-$arch"
    $downloadURL = "https://github.com/skycoin/skycoin/releases/download/$version/$filename.zip"
    Invoke-WebRequest $downloadURL -OutFile archive.zip

    Write-Output "#       3. Extracting Downloaded Archive File...         #"
    Expand-Archive -Path archive.zip -DestinationPath archive

    Write-Output "#       4. Preparing Environment for Wix...              #"
    Set-Location .\scripts\win_installer
    mkdir -p ".\build" > $null
    Copy-Item ..\..\archive\skycoin.exe .\build\skycoin.exe
    # Only copy skyhw for amd64
    if ($arch -eq "amd64" -and (Test-Path "..\..\archive\skyhw.exe")) {
        Copy-Item ..\..\archive\skyhw.exe .\build\skyhw.exe
    }
    
    $installerVersion = $version -replace '(^v|-.+$)', ''
    $productWxs = Get-Content -Path Product.wxs
    $newContent = $productWxs -replace "skycoinVersion", $installerVersion
    Set-Content -Path Product.wxs -Value $newContent

    Write-Output "#       5. Building MSI Installer...                     #"
    .\wix\candle.exe Product.wxs -arch $wix_arch  > $null
    .\wix\light.exe -ext WixUIExtension -ext WixUtilExtension -sacl -spdb -out skycoin.msi Product.wixobj  > $null
    Move-Item skycoin.msi ../../$msiName.msi -Force

    Write-Output "#          ==> Build Completed for $arch_title!                #"
    
    Write-Output "#       6. Cleaning Stage...                             #"
    $productWxs = Get-Content -Path Product.wxs
    $newContent = $productWxs -replace $installerVersion, "skycoinVersion"
    Set-Content -Path Product.wxs -Value $newContent
    Set-Location ../../
    CleanStage

    Write-Output "#       7. Done!                                         #"
} 

Write-Output "`n##########################################################"
Write-Output "#                                                        #"
Write-Output "#        .:::: Create MSI Installer Package ::::.        #"
BuildInstaller
Write-Output "#                                                        #"
Write-Output "##########################################################`n"
