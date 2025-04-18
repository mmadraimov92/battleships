# Battleships

Battleships terminal UI multiplayer game. Creatd as learning project.

## TODO

* Implement networking
  * Use mDNS to find other player.
    * hashicorp/mdns can be used.
    * client can lookup specific service. If does not find then start own server and advertise dns.
    * if finds then connects with tcp or websocket.
    * as server, if other player connects, then shuts down mdns server.
  * discard messages based on who's turn it is.
