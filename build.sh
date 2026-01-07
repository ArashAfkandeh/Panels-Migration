cd /root/Panels-Migration-main
go mod tidy
go build -o Panels_Migration ./cmd/main
chmod +x Panels_Migration
