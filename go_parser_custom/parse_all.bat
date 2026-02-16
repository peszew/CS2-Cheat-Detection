@echo off
echo Deleting old CSV...
del enhanced_kills.csv

for %%f in (..\demos\*.dem) do (
    echo Parsing %%f ...
    go run main.go %%f
)
