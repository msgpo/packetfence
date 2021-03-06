package rewrite

import (
	"github.com/inverse-inc/packetfence/go/coredns/core/dnsserver"
	"github.com/inverse-inc/packetfence/go/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("rewrite", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	rewrites, err := rewriteParse(c)
	if err != nil {
		return plugin.Error("rewrite", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Rewrite{Next: next, Rules: rewrites}
	})

	return nil
}

func rewriteParse(c *caddy.Controller) ([]Rule, error) {
	var rules []Rule

	for c.Next() {
		args := c.RemainingArgs()
		rule, err := newRule(args...)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
