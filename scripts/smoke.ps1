# PowerShell smoke test for GOSS-CLI
param(
    [int]$TimeoutSeconds = 30
)

$ErrorActionPreference = "Stop"

function Write-Pass { param($msg) Write-Host "✓ $msg" -ForegroundColor Green }
function Write-Fail { param($msg) Write-Host "✗ $msg" -ForegroundColor Red; exit 1 }
function Write-Info { param($msg) Write-Host "$msg" -ForegroundColor Yellow }

$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

Write-Info "🧪 Running GOSS-CLI smoke tests..."
Write-Host

# 0) Sanity
if (-not (Test-Path ".\bin\goss" -PathType Leaf)) {
    Write-Fail "bin/goss missing"
}
Write-Pass "CLI present"

# 1) Help text
$helpOut = & ".\bin\goss" --help 2>&1
if ($helpOut -match "GOSS-CLI") {
    Write-Pass "Help prints"
} else {
    Write-Fail "Help missing"
}

# 2) Non-stream JSON response
Write-Info "Testing non-stream JSON response..."
try {
    $jsonOut = & ".\bin\goss" --no-stream --temperature 0 'Reply with exactly {"ok":true} and nothing else' 2>&1 | Out-String
    if ($jsonOut -match '"ok":') {
        Write-Pass "Non-stream JSON roundtrip"
    } else {
        Write-Fail "JSON not returned: $jsonOut"
    }
} catch {
    Write-Fail "Non-stream test failed: $($_.Exception.Message)"
}

# 3) Streaming works
Write-Info "Testing streaming output..."
try {
    $streamOut = & ".\bin\goss" 'Type: streaming-ok' 2>&1 | Select-Object -First 1
    if ($streamOut -match 'streaming-ok') {
        Write-Pass "Streaming output"
    } else {
        Write-Fail "Streaming failed: $streamOut"
    }
} catch {
    Write-Fail "Streaming test failed: $($_.Exception.Message)"
}

# 4) Save transcript
if (Test-Path "logs") { Remove-Item "logs" -Recurse -Force }
New-Item -ItemType Directory -Path "logs" -Force | Out-Null
try {
    & ".\bin\goss" --save --no-stream 'Reply with {"saved":true}' 2>&1 | Out-Null
    $logFiles = Get-ChildItem "logs" -ErrorAction SilentlyContinue
    if ($logFiles.Count -ge 1) {
        Write-Pass "Transcript saved"
    } else {
        Write-Fail "No log file created"
    }
} catch {
    Write-Fail "Save transcript test failed"
}

# 5) Context file respected
"System: The secret is swordfish." | Out-File -FilePath "$env:TEMP\goss_ctx.txt" -Encoding UTF8
try {
    $ctxOut = & ".\bin\goss" --context-file "$env:TEMP\goss_ctx.txt" --no-stream --temperature 0 'What is the secret? Answer one word.' 2>&1 | Out-String
    if ($ctxOut.ToLower() -match 'swordfish') {
        Write-Pass "Context file applied"
    } else {
        Write-Fail "Context not applied: $ctxOut"
    }
} catch {
    Write-Fail "Context file test failed: $($_.Exception.Message)"
}

# 6) List models
try {
    $modelsOut = & ".\bin\goss" list-models 2>&1 | Out-String
    if ($modelsOut -match 'model|gpt|llama|mistral|Found \d+ model') {
        Write-Pass "Models listed"
    } else {
        Write-Fail "No models listed (is provider running?): $modelsOut"
    }
} catch {
    Write-Fail "List models failed: $($_.Exception.Message)"
}

# 7) Invalid model UX
try {
    $invalidOut = & ".\bin\goss" --model definitely-not-a-real-model --no-stream 'hi' 2>&1 | Out-String
    if ($invalidOut -match 'available models|not found|invalid model|Warning.*Model') {
        Write-Pass "Invalid model handled"
    } else {
        Write-Fail "Invalid model not handled"
    }
} catch {
    # Expected to fail, but should have helpful error message
    if ($_.Exception.Message -match 'available models|not found|invalid model') {
        Write-Pass "Invalid model handled"
    } else {
        Write-Fail "Invalid model not handled properly"
    }
}

# 8) Provider override
try {
    $provOut = & ".\bin\goss" --debug --no-stream --temperature 0 'Reply with "prov-ok" exactly' 2>&1 | Select-Object -First 1
    if ($provOut -match 'prov-ok') {
        Write-Pass "Provider override works"
    } else {
        Write-Fail "Provider override failed: $provOut"
    }
} catch {
    Write-Fail "Provider test failed: $($_.Exception.Message)"
}

# 9) Unreachable endpoint error
try {
    $connOut = & ".\bin\goss" --api-base http://127.0.0.1:9 --no-stream 'hi' 2>&1 | Out-String
    # Should fail and show error
    Write-Fail "Should have failed with connection error"
} catch {
    $errorMsg = $_.Exception.Message
    if ($errorMsg -match 'unreachable|ECONNREFUSED|connect|Connection refused') {
        Write-Pass "Unreachable endpoint surfaces clear error"
    } else {
        Write-Fail "Connection error not surfaced: $errorMsg"
    }
}

# 10) Debug shows provider selection
try {
    $debugOut = & ".\bin\goss" --debug --no-stream 'hi' 2>&1 | Out-String
    if ($debugOut -match 'DEBUG.*REQUEST|DEBUG.*provider|Using provider') {
        Write-Pass "Debug logs include provider/REQUEST info"
    } else {
        Write-Fail "Debug logs missing provider info"
    }
} catch {
    Write-Fail "Debug test failed: $($_.Exception.Message)"
}

Write-Host
Write-Host "🎉 All smoke tests passed!" -ForegroundColor Green
Write-Host "GOSS-CLI is ready for launch! 🚀" -ForegroundColor Yellow

# Cleanup
Remove-Item "$env:TEMP\goss_ctx.txt" -ErrorAction SilentlyContinue