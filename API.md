# Geep protocol specification

## Architecture

Geep is module application. It can be used two ways. As client/server and Go libriary. In client/server way, Client communicate with server with simple GDATA protol, which allow list of service commands. In general response has format:

## API

### SET
REQUEST:  SET <key> <ttl> <value>
RESPONSE: [201]
Set new <key> to <value> with expire in <ttl> sec or return error.
Note: key can't be set twice. Use 'UPDATE' instead.

### GET
REQUEST:  GET <key>
RESPONSE: <value>
GET <value> for <key> if exists and not expired.

### UPDATE
REQUEST:  UPDATE <key> <ttl> <value>
RESPONSE: [204]
UPDATE <value> for <key> if exists or return error.

### UPDATE
REQUEST:  DELETE <key>
RESPONSE: [204]
DELETE <key> if exists.

TODO: add list and dict documentation
