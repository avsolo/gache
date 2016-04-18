# Gache protocol specification

## Architecture

Gache is module application. It can be used two ways. As client/server and Go libriary. In client/server way, Client communicate with server with simple GDATA protol, which allow list of service commands. In general response has format:

## API

### SET

    REQUEST:  SET key value ttl
    RESPONSE: [201]
    
Sets new <key> to <value> with expire in <ttl> sec or return error. Note: key can't be set twice. Use 'UPDATE' instead.

### GET

    REQUEST:  GET key
    RESPONSE: value
    
Returns <value> for <key> if exists and not expired.

### UPDATE

    REQUEST:  UPD key value ttl
    RESPONSE: [204]

Updates <value> for <key> if exists or return error.

### DELETE

    REQUEST:  DEL <key>
    RESPONSE: [204]

Deletes key if exists.

## Lists

### LSET

    REQUEST:  LSET key val1 [val2, val3,.. valN] ttl
    RESPONSE: [204]

Creates List (LIFO) named <key> with values val1...valN

### LPUSH

    REQUEST:  LADD key value
    RESPONSE: [204]

Pushes <value> to list <key>

### LPOP

    REQUEST:  LPOP key
    RESPONSE: [204]

Delete last pushed <value> for <key> and return it as result

## Maps

### DSET

    REQUEST: DSET key k1 v1 k2 v2.. kn vn ttl
    RESPONSE: [204]

Create map {k1:v1, k2:v2,.. kn:vn} and saves it as <key>

### DGET

    REQUEST: DSET rkey skey
    RESPONSE: value

Returns <value> for <skey> in map <rkey>

### DADD

    REQUEST: DADD rkey skey value
    RESPONSE: [204]

Add <skey>:<value> in map <rkey>

### DDEL

    REQUEST: DDEL rkey skey
    RESPONSE: [204]

Deletes <skey> in map <rkey>

TODO: add LSET, LGET... DSET.. documentation
