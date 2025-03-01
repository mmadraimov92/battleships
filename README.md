# Basic TUI app

Basic terminal UI app written in Go as a learning project.

## TODO

* Implement networking
  * Use mDNS to find other player.
    * hashicorp/mdns can be used.
    * client can lookup specific service. If does not find then start own server and advertise dns.
    * if finds then connects with tcp or websocket.
    * as server, if other player connects, then shuts down mdns server.
  * Use websockets for communication? or just basic tcp streaming? coder/websocket can be used.
  * Design message protocol.
  * discard messages based on who's turn it is.
* Cleanup project
  * Remove tui menu? Go directly to battleships?
