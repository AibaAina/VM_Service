# Define the path to the JSON file
$jsonFilePath = "D:\Documents\VScode\VM_Service\VM_Service\appinfo.json"

# Read the content of the JSON file
$jsonContent = Get-Content -Path $jsonFilePath -Raw

# Parse the JSON content into a PowerShell object
$jsonObject = ConvertFrom-Json -InputObject $jsonContent

# Extract the value of the "version" property and save it in a variable
$appVersion = $jsonObject.version

# Display the app version
Write-Host "App Version: $appVersion"
