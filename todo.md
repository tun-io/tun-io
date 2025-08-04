# TODO:

## Client
- [ ] Implement config file
- [ ] Implement flags for CLI
- [ ] Make a pretty TUI
using: `github.com/charmbracelet/bubbletea` `github.com/NimbleMarkets/ntcharts` `github.com/skip2/go-qrcode` 

## Server
- [ ] Implement authentication middleware
- [ ] Implement authentication database (gorm or bun)
- [x] Implement prometheus metrics
- [x] Implement actual logging

## Readme:
- [ ] Use VHS to record a demo of the client and server in action.
- [ ] Add a section on how to run the server and client.


## Extras:
- [x] Implement a docker image for the server
- [ ] Implement a docker image for the client (optional)
- [ ] Implement a docker-compose file to run both the server, database and maybe a prometheus/grafana stack.

# REQUIRED
- [ ] Implement https support
- [ ] Allow sending errors respones via HTTP/WS to the server to make sure the requester errors out correctly
- [ ] Handle websockets possibly crashing
- [ ] Allow the server to gracefully time out / retry packets after X amount of time.
