@echo off
 
set PROTO_DIR=..\proto
set OUT_DIR=..\proto
 
for /D %%S in (%PROTO_DIR%\*) do (
    for %%F in (%%S\*.proto) do (
        protoc ^
            --proto_path=%PROTO_DIR% ^
            --go_out=%OUT_DIR% ^
            --go_opt=paths=source_relative ^
            --go-grpc_out=%OUT_DIR% ^
            --go-grpc_opt=paths=source_relative ^
            %%F
        echo Generated: %%F
    )
)
 
echo Done.