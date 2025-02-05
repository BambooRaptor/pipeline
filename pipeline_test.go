package pipeline

import (
	"fmt"
	"testing"
)

type testConfig struct {
	maxConn uint
	id      string
	name    string
}

var defaultConfig = testConfig{10, "default", "Default Config"}

func buildConfig(opts ...Pipe[testConfig]) testConfig {
	return New(opts...).CloseWith(defaultConfig)
}

func setName(name string) Pipe[testConfig] {
	return func(c testConfig) testConfig {
		c.name = name
		return c
	}
}

func setId(id string) Pipe[testConfig] {
	return func(c testConfig) testConfig {
		c.id = id
		return c
	}
}

func addConns(amount uint) Pipe[testConfig] {
	return func(c testConfig) testConfig {
		c.maxConn += amount
		return c
	}
}

func (c testConfig) matchAgainst(config testConfig) bool {
	return c.id == config.id && c.name == config.name && c.maxConn == config.maxConn
}

func (c testConfig) String() string {
	return fmt.Sprintf("[%s] %q => %v", c.id, c.name, c.maxConn)
}

func TestConfigBuilder(t *testing.T) {
	config := buildConfig(
		setName("What what"),
		setName("Testing Config"),
		setId("testing"),
		addConns(5),
		addConns(5),
	)

	want := testConfig{
		name:    "Testing Config",
		id:      "testing",
		maxConn: 20,
	}

	if !want.matchAgainst(config) {
		t.Fatalf("Test Failed\n WANT [%s] <=/=> [%s] GOT", want, config)
	}
}

type testInterface interface {
	Use(num int)
}

type intUser func(num int)

func (i intUser) Use(num int) {
	i(num)
}

func add(amount int) Pipe[testInterface] {
	return func(p testInterface) testInterface {
		return intUser(func(num int) {
			p.Use(amount + num)
		})
	}
}

func mult(amount int) Pipe[testInterface] {
	return func(p testInterface) testInterface {
		return intUser(func(num int) {
			p.Use(amount * num)
		})
	}
}

type testIntChecker struct {
	t    *testing.T
	want int
}

func checkFinalInt(t *testing.T, want int) testIntChecker {
	return testIntChecker{t, want}
}

func (i testIntChecker) Use(num int) {
	if i.want != num {
		i.t.Fatalf("[WANT] %v <=/=> %v [GOT]", i.want, num)
	}
}

func TestInterfacePipeline(t *testing.T) {
	line1 := New(
		add(5),
		mult(10),
	).EmbedFinal(checkFinalInt(t, 70))
	line1.Use(2)

	line2 := New(
		mult(10),
		add(5),
	).EmbedFinal(checkFinalInt(t, 25))
	line2.Use(2)
}
