# matrix-server-chain-check

## Description

This check-plugin complies to the [monitoring plugin development guidelines](https://www.monitoring-plugins.org/doc/guidelines.html), and is therefore compatible with [Nagios](https://nagios.com), [Icinga](https://icinga.com), [Zabbix](https://zabbix.com), [CheckMK](https://checkmk.com), etc.

This is a chain check for testing connectivity between two [matrix servers](https://matrix.org), written in Go. The check sends a message to a (remote) room (on the other server) with one user (the sending user) and the other user (the receiving user) checks if the message contains the same message it should receive.



## Usage

Flags you have to set to use this check:

### Required-Flags

- -sending-homeserver     (here you have to specify the matrix server, which contains the sending user that you´ve registered.)
- -sending-username       (here you have to specify the username, that you have registered for your user, on the sending homeserver)
- -sending-password       (here you have to specify the password, that you have entered with the registration of your sending user)
- -receiving-homeserver   (here you have to specify the matrix server, which contains the receiving user that you´ve registered.)
- -receiving-username     (here you have to specify the username, that you have registered for your user, on the receiving homeserver)
- -receiving-password     (here you have to specify the password, that you have entered with the registration of your sending user) 
- -room-id                (here you have to specify the room-id where you want to send you message, that will be checked)



### Optional Flags

- -timeout                (here you can specify the timeout in seconds, when theres no matching message)

The output of this check-plugin will be handled by the [go-monitoringplugin](https://github.com/inexio/go-monitoringplugin). So you dont need to prettify or do something else to implement this in your monitoring system.



## How to Build

```cli
go build -o matrix-server-chain-check -o main.go
```



## Contribution

We are always looking forward for your ideas and suggestions.

If you want to help us, please make sure that your code is conform to our [coding-style](https://github.com/uber-go/guide/blob/master/style.md)

Happy Coding!
