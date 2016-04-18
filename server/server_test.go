package server

import (
    "fmt"
    "time"
	"testing"
    "github.com/stretchr/testify/assert"

	"github.com/avsolo/gache/lib"
    "github.com/avsolo/gache/server"
)

var _s = fmt.Sprintf
var log *lib.Logger
var addr = lib.CliParams.ServerAddr
var srv = server.NewServer(addr)
var cln *server.Client

func init() {
	log = lib.NewLogger("server_test")
    go func() {
        srv.ListenTCP()
        log.Info("Server stopped")
    }()
    // Sometimest client trying connect before server start.
    // So, we just wait a bit
    time.Sleep(time.Duration(100 * time.Microsecond))
    cln = server.NewClient(addr)
}

// checkGet is shortcut for GET request and asserting expected value
func checkGet(t *testing.T, key, expected string) {
    res, err := cln.Sendf("GET %s", key)
    assert.Nil(t, err)
    assert.Equal(t, expected, res)
}

func TestSET(t *testing.T) {
    cln.Send("OPT flush")
    k1, v1, t1 := "k1", "some v alue", 100

    // Set first time
    res, err := cln.Sendf("SET %s %s %d", k1, v1, t1)
    assert.Nil(t, err)
    assert.Regexp(t, `\[201\]\s*`, res)

    // Check GET
    checkGet(t, k1, v1)

    // Set second time same key denied
    res, err = cln.Sendf("SET %s %s %d", k1, v1, t1)
    assert.Nil(t, err)
    assert.Regexp(t, `\[400\]\s*`, res)

    // Check GET again
    checkGet(t, k1, v1)
}

func TestUPDATE(t *testing.T) {
    cln.Send("OPT flush")
    k1, v1, t1 := "k1", "some v alue", 100

    _, _ = cln.Sendf("SET %s %s %d", k1, v1, t1)

    v2, t2 := "new value", 50
    r, err := cln.Sendf("UPD %s %s %d", k1, v2, t2)
    assert.Nil(t, err)
    assert.Regexp(t, `\[204\]\s*`, r)

    // Check updated value again
    checkGet(t, k1, v2)
}

func TestDELETE(t *testing.T) {
    cln.Send("OPT flush")
    k1, v1, t1 := "k1", "some v alue", 100

    _, _ = cln.Sendf("SET %s %s %d", k1, v1, t1)

    r, err := cln.Sendf("DEL %s", k1)
    assert.Nil(t, err)
    assert.Regexp(t, `\[204\]\s*`, r)

    // Check updated value again
    res, err := cln.Sendf("GET %s", k1)
    assert.Nil(t, err)
    assert.Regexp(t, `\[400\]\s*`, res)
}

func checkLPOP(t *testing.T, key string, vals []string) {
    for _, exp := range vals {
        v, err := cln.Sendf("LPOP %s", key)
        assert.Nil(t, err)
        assert.Equal(t, exp, v)
    }
}

// LISTS
func TestLSET(t *testing.T) {
    cln.Send("OPT flush")
    v1, v2, v3 := "val1", "val2", "val3"
    rkey, vs, t1 := "rkey", _s("%s %s %s", v1, v2, v3), 100

    // Set first time
    res, err := cln.Sendf("LSET %s %s %d", rkey, vs, t1)
    assert.Nil(t, err)
    assert.Regexp(t, `\[201\]\s*`, res)

    // Check LPOP
    checkLPOP(t, rkey, []string{v3, v2, v1})
}

func TestLPUSH(t *testing.T) {
    cln.Send("OPT flush")
    v1, v2, v3 := "val1", "val2", "val3"
    rkey, vs, t1 := "rkey", _s("%s %s %s", v1, v2, v3), 100

    // Set first time
	res, err := cln.Sendf("LSET %s %s %d", rkey, vs, t1)
    assert.Nil(t, err)
    assert.Regexp(t, `\[201\]\s*`, res)

    // Push one more key
    v4 := "val4"
    res, err = cln.Sendf("LPUSH %s %s", rkey, v4)
    assert.Nil(t, err)
    assert.Regexp(t, `\[204\]\s*`, res)

    v5 := "val5"
    res, err = cln.Sendf("LPUSH %s %s", rkey, v5)
    assert.Nil(t, err)
    assert.Regexp(t, `\[204\]\s*`, res)

	checkLPOP(t, rkey, []string{v5, v4, v3, v2, v1})
}

func checkDGET(t *testing.T, key string, expMap map[string]string) {
    for k, v := range expMap {
        res, err := cln.Sendf("DGET %s %s", key, k)
        // log.Debugf("DGET %s %s result: %s", key, k, res)
        assert.Nil(t, err)
        assert.Equal(t, v, res)
    }
}

// DICTS
func TestDSET(t *testing.T) {
    cln.Send("OPT flush")
    dictKey := "dictKey"
    k1, k2, k3 := "key1", "key2", "key3"
    v1, v2, v3 := "val1", "val2", "val3"
    dMap := map[string]string{k1:v1, k2:v2, k3:v3}
    t1 := 10

    // Set first time
    kvs := _s("%s %s %s %s %s %s", k1,v1, k2,v2, k3,v3)
    res, err := cln.Sendf("DSET %s %s %d", dictKey, kvs, t1)
    assert.Nil(t, err)
    assert.Regexp(t, `\[201\]\s*`, res)

    // Check DGET
    checkDGET(t, dictKey, dMap)

    // Try to DADD
    k4, v4 := "key4", "val4"
    res, err = cln.Sendf("DADD %s %s %s", dictKey, k4, v4)
    assert.Nil(t, err)
    assert.Regexp(t, `\[204\]\s*`, res)

    // Check DGET for new val
    dMap[k4] = v4
    checkDGET(t, dictKey, dMap)

    // Check DDEL
    delete(dMap, k3)
    res, err = cln.Sendf("DGET %s %s %s", dictKey, k3)
    assert.Nil(t, err)
    assert.Regexp(t, `\[400\]\s*`, res)
}
